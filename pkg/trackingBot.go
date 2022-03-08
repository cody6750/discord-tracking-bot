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

// channels
type channels struct {
	stopTracking  chan struct{}
	startTracking chan struct{}
}

// TrackingBot ...
type TrackingBot struct {
	discordSession        *discordgo.Session
	session               *session.Session
	secretsSvc            *secretsmanager.SecretsManager
	logger                *logrus.Logger
	options               *options.Options
	channels              *channels
	discordLogChannel     *discordgo.Channel
	discordMetricsChannel *discordgo.Channel
}

//NewTrackingBot ...
func NewTrackingBot() *TrackingBot {
	return NewTrackingBotWithOptions(options.New())
}

//NewTrackingBotWithOptions ...
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

//Run initializes and runs the discord bot
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

func (t *TrackingBot) monitoringChannelSignals() {
	for {
		select {
		case <-t.channels.startTracking:
			go t.startTracking("tracking", t.options.TrackingConfigPath)
		}
	}
}

func (t *TrackingBot) StartTracking() {
	t.channels.startTracking <- struct{}{}
}

func (t *TrackingBot) StopTracking() {
	t.channels.stopTracking <- struct{}{}
}

func (t *TrackingBot) initAWS(maxRetries int, region string) {
	configs := aws.Config{
		Region:     aws.String(region),
		MaxRetries: aws.Int(maxRetries),
	}
	t.session = session.Must(session.NewSession(&configs))
	t.secretsSvc = secretsmanager.New(t.session)
}
