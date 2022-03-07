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

const (
	// Default value for discordTokenID. Overrideable via environment variables.
	defaultDiscordTokenID = "discord/token"
	// Default value for discordToken. Overrideable via environment variables.
	defaultDiscordToken = ""
	// Default value for maxRetries. Overrideable via environment variables.
	defaultMaxRetries = 5
	// Default valur for region. Overrideable via environment variables.
	defaultRegion = "us-east-1"
	// Default value for localRun. Overrideable via environment variables.
	defaultLocalRun = false
)

type channels struct {
	stopTracking chan struct{}
}

var (
	mediaPath *string
)

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
		stopTracking: make(chan struct{}),
	}
	bot.logger.SetFormatter(&logrus.TextFormatter{ForceColors: true, FullTimestamp: true})
	return bot
}

func (t *TrackingBot) initBot() {
	var err error
	t.logger.Info("Discord bot initializing......")
	t.getEnvVariables()

	if !t.options.LocalRun {
		t.initAWS(t.options.AWSMaxRetries, t.options.AWSRegion)
		t.options.DiscordToken, err = services.GetSecret(t.secretsSvc, t.options.DiscordTokenAWSSecretName)
		if err != nil {
			t.logger.WithError(err).Fatalf("failed to get discord token secret from AWS")
		}
	}

	t.checkEnvVariables()

	// Create a new Discord session using the provided bot token.
	t.discordSession, err = discordgo.New("Bot " + t.options.DiscordToken)
	if err != nil {
		t.logger.WithError(err).Fatalf("failed to create discord session")
	}

	// Register the messageCreate func as a callback for MessageCreate events.
	handlers.EnableHandlerMedia(t.options.MediaPath)
	t.discordSession.AddHandler(handlers.Ready)
	t.discordSession.AddHandler(handlers.MessageCreate)
	t.discordSession.AddHandler(handlers.Disconnect)
	t.discordSession.AddHandler(handlers.SlashCommands)

	t.discordSession.Identify.Intents = discordgo.IntentsGuildMessages

	// Open a websocket connection to Discord and begin listening.
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
	functions.LogToDiscordAndStdOut(t.logger, t.discordSession, t.discordLogChannel, t.logger.Info, "Bot is now running.  Press CTRL-C to exit.")
	go t.startTracking("tracking", t.options.TrackingConfigPath)

	// Wait here until CTRL-C or other term signal is received.
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
	functions.LogToDiscordAndStdOut(t.logger, t.discordSession, t.discordLogChannel, t.logger.Info, "Bot has shut down")
	t.discordSession.Close()
}

func (t *TrackingBot) startTracking(channel, trackingConfigPath string) {
	functions.LogToDiscordAndStdOut(t.logger, t.discordSession, t.discordLogChannel, t.logger.Info, "Starting to track channels")
	handlers.EnableStatus("Running, tracking Channels")
	defer handlers.EnableStatus("Idle")
	go func() {
		t.TrackItemChannels(t.discordSession, functions.GetChannels(t.discordSession, channel), trackingConfigPath)
		t.channels.stopTracking <- struct{}{}
	}()
	select {
	case <-t.channels.stopTracking:
		functions.LogToDiscordAndStdOut(t.logger, t.discordSession, t.discordLogChannel, t.logger.Info, "Tracking has stopped")
		return
	}
}

func (t *TrackingBot) stopTracking() {
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
