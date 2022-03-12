package trackingbot

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/bwmarrin/discordgo"
	"github.com/cody6750/discordbot/pkg/functions"
	"github.com/cody6750/discordbot/pkg/handlers"
	"github.com/cody6750/discordbot/pkg/options"
	services "github.com/cody6750/discordbot/pkg/services/aws"
	"github.com/sirupsen/logrus"
)

// channels contains all of the channels used by the tracking bot. These are initialized
// upon the creation of the tracking bot object
type channels struct {
	stopTracking  chan struct{}
	startTracking chan struct{}
}

// TrackingBot represents all of the necessary dependencies required to run the tracking
// bot. The tracking bot itself is designed to run as a microservice on AWS with feature
// flags that allow for local testing using docker or the go exectuable.
type TrackingBot struct {
	// discordSession established a session with Discord. Requires an API token
	// which is set using options.DiscordToken or DISCORD_TOKEN environment variable.
	discordSession *discordgo.Session

	// session established a session with AWS. Requires AWS to be configured on the
	// machine. The session is created through initAWS which is set using options.AWSMaxRetries
	// and options.AWSRegion or AWS_MAX_RETRIES and AWS_REGION environent variables.
	session *session.Session

	// secretsSvc establishes a session with AWS Secrets manager using the AWS session.
	// Allows us to get the discord token from AWS secrets.
	secretsSvc *secretsmanager.SecretsManager

	// logger represnts the structured logger used within the tracking bot and discord handlers
	logger *logrus.Logger

	// options represents the configurable tracking bot options. This sturct allows users to
	// override the default options. Users can tailor the tracking bot to their needs based on
	// the running platform.
	options *options.Options

	// channels includes all of the channels used within the tracking bot. These have different functionalities.
	channels *channels

	// discordLogChannel represents the discord channel that is used to output logs.
	discordLogChannel *discordgo.Channel

	// discordMetricsChannel represnts the discord channel that is use to output metrics.
	discordMetricsChannel *discordgo.Channel
}

// NewTrackingBot creates an instance of the tracking bot with the default settings.
// Environment variables can be used to override the default settings.
func NewTrackingBot(discordToken ...string) *TrackingBot {
	return NewTrackingBotWithOptions(options.New(discordToken...))
}

// NewTrackingBotWithOptions creates an instance of the tracking bot with custom settings.
// Environment variables will override the custom settings.
func NewTrackingBotWithOptions(option *options.Options) *TrackingBot {
	bot := &TrackingBot{}
	bot.options = option
	bot.logger = logrus.New()
	bot.channels = &channels{
		stopTracking:  make(chan struct{}),
		startTracking: make(chan struct{}),
	}
	handlers.SetChannel("startTracking", bot.channels.startTracking)
	handlers.SetChannel("stopTracking", bot.channels.startTracking)
	bot.logger.SetFormatter(&logrus.TextFormatter{ForceColors: true, FullTimestamp: true})
	return bot
}

// initBot is used to bootstrap the bot. It initializes all of the dependencies and creates the necessary connections
// with 3rd party services such as Discord and AWS. The environment variables are gathered and
// set. All handlers are added to the discord session, and a websocket connection is created.
func (t *TrackingBot) initBot() {
	var err error
	t.logger.Info("Discord bot initializing")
	handlers.SetStatus("Discord bot initializing")
	t.getEnvVariables()

	if !t.options.LocalRun {
		t.initAWS(t.options.AWSMaxRetries, t.options.AWSRegion)
		t.options.DiscordToken, err = services.GetSecret(t.secretsSvc, t.options.DiscordTokenAWSSecretName)
		if err != nil {
			t.logger.WithError(err).Fatalf("failed to get discord token secret from AWS")
		}
	}

	t.checkEnvVariables()

	t.discordSession, err = discordgo.New("Bot " + t.options.DiscordToken)
	if err != nil {
		t.logger.WithError(err).Fatalf("failed to create discord session")
	}

	handlers.SetHandlerMediaPath(t.options.MediaPath)
	t.discordSession.AddHandler(handlers.Ready)
	t.discordSession.AddHandler(handlers.MessageCreate)
	t.discordSession.AddHandler(handlers.Disconnect)
	t.discordSession.AddHandler(handlers.SlashCommands)

	t.discordSession.Identify.Intents = discordgo.IntentsGuildMessages

	err = t.discordSession.Open()
	if err != nil {
		t.logger.WithError(err).Fatalf("failed to open websocket connection")
	}

	if t.options.LogToDiscord {
		t.discordLogChannel = functions.GetChannel(t.discordSession, "logs")
	}

	if t.options.MetricsToDiscord {
		t.discordMetricsChannel = functions.GetChannel(t.discordSession, "metrics")
	}
	functions.LogToDiscordAndStdOut(t.logger, t.discordSession, t.discordLogChannel, t.logger.Info, "Succesfully initialized bot")

}

// Run serves as the entrypoint to the tracking bot. It needs only to be called once per session.
// Immediatley after the bot is initialized, the bot will start tracking. Once the bot is live,
// commands will be sent to the bot from the handlers, and channels. If there are any signal interrupts,
// the bot will shutdown.
func (t *TrackingBot) Run() {
	t.initBot()
	handlers.SetStatus("Idle")
	functions.LogToDiscordAndStdOut(t.logger, t.discordSession, t.discordLogChannel, t.logger.Info, "Bot is now running.  Press CTRL-C to exit.")

	go t.startTracking("tracking", t.options.TrackingConfigPath)

	go t.monitoringChannelSignals()

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
	handlers.SetStatus("Bot shutting down")
	functions.LogToDiscordAndStdOut(t.logger, t.discordSession, t.discordLogChannel, t.logger.Info, "Bot has shut down")
	t.discordSession.Close()
}

// startTracking begins the tracking process within the tracking bot. All of the tracking channels within
// the discord server will be intialized. The TrackItemChannels function is wrapped inside of a goroutine
// anonymous function which allows the execution of the code to be asynchronous. This allows signals to be sent
// to the stopTracking channel, which is used to stop TrackItemChannels at any moment.
//
//  parameters:
//
//	channelsToTrack string : The prefix of discord channels to track.
//
//  trackingConfigPath string : The path of the tracking file configs within the file system. These files are used
//  to define what to track on a given channel.
func (t *TrackingBot) startTracking(channelsToTrack, trackingConfigPath string) {
	functions.LogToDiscordAndStdOut(t.logger, t.discordSession, t.discordLogChannel, t.logger.Info, "Starting to track channels")
	handlers.SetStatus("Running, currently tracking channels")
	go func() {
		t.TrackItemChannels(t.discordSession, functions.GetChannels(t.discordSession, channelsToTrack), trackingConfigPath, t.options.TrackingChannelsDelay)
		t.channels.stopTracking <- struct{}{}
		handlers.SetStatus("Idle")
	}()

	select {
	case <-t.channels.stopTracking:
		functions.LogToDiscordAndStdOut(t.logger, t.discordSession, t.discordLogChannel, t.logger.Info, "Tracking has stopped")
		return
	}
}

// monitoringChannelSignals is used to act upon a signal that is recieved from the handlers.
func (t *TrackingBot) monitoringChannelSignals() {
	for {
		select {
		case <-t.channels.startTracking:
			go t.startTracking("tracking", t.options.TrackingConfigPath)
		}
	}
}

// StopTracking is used to stop the tracking bot via a signal on the stopTracking channel.
func (t *TrackingBot) StopTracking() {
	t.channels.stopTracking <- struct{}{}
}

// initAWS creates the required AWS session and services.
func (t *TrackingBot) initAWS(maxRetries int, region string) {
	configs := aws.Config{
		Region:     aws.String(region),
		MaxRetries: aws.Int(maxRetries),
	}
	t.session = session.Must(session.NewSession(&configs))
	t.secretsSvc = secretsmanager.New(t.session)
}
