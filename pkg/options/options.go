package options

const (
	// Default value for AWSRegion. Overrideable via environment variables.
	defaultAWSRegion = "us-east-1"

	// Default value for AWSAWSMaxRetries. Overrideable via environment variables.
	defaultAWSMaxRetries = 5

	// Default value for DiscordTokenAWSSecretName. Overrideable via environment variables.
	defaultDiscordTokenAWSSecretName = "discord/token"

	// Default value for DiscordToken. Overrideable via environment variables.
	defaultDiscordToken = ""

	// Default value for LocalRun. Overrideable via environment variables.
	defaultLocalRun = false

	// Default value for LogToDiscord. Overridaeble via environment variables.
	defaultLogToDiscord = true

	// Default value for MediaPath. Overridaeble via environment variables.
	defaultMediaPath = "media/"

	// Default value for MetricsToDiscord. Overridaeble via environment variables.
	defaultMetricsToDiscord = true

	// Default value for TrackingConfigPath. Overridaeble via environment variables.
	defaultTrackingChannelsDelay = 21600

	// Default value for TrackingConfigPath. Overridaeble via environment variables.
	defaultTrackingConfigPath = "pkg/configs/tracking/"

	// Default value for WebcrawlerHost. Overridaeble via environment variables.
	defaultWebcrawlerHost = "localhost"

	// Default value for WebcrawlerProt. Overridaeble via environment variables.
	defaultWebcrawlerPort = 9090
)

//Options represents all avlaiable configurable options for the tracking bot
type Options struct {
	// LogToDiscord flag used to determine whether to send logs to discord.
	LogToDiscord bool

	// LocalRun flag used to determine whether the runtime platform is on AWS or not. If set to false, the
	// tracking bot will set up the connection to AWS and retrive a secret value.
	LocalRun bool

	// LogToDiscord flag used to determine whether to send metrics logs to discord.
	MetricsToDiscord bool

	// AWSMaxRetries used to determine the maxmimum retries when establishes a session with AWS.
	AWSMaxRetries int

	// WebcrawlerPort port that is used to establish the connection to the webcrawler.
	WebcrawlerPort int

	// WebcrawlerPort delay between each channel tracking.
	TrackingChannelsDelay int

	// AWSMaxRetries region the applications is deployed in.
	AWSRegion string

	// DiscordTokenAWSSecretName AWS secret name that stores discord token.
	DiscordTokenAWSSecretName string

	// DiscordToken token used to establish a session with discord.
	DiscordToken string

	// MediaPath path within the filesystem that contains the media directory.
	MediaPath string

	//TrackingConfigPath path within the filesystem that contains the config tracking directory
	TrackingConfigPath string

	//WebcrawlerHost host name set by the webcrawler. Used to connect to the webcrawler alongisde the WebcrawlerPort.
	WebcrawlerHost string
}

//New returns Options with default values.
func New(discordToken ...string) *Options {
	options := &Options{
		LogToDiscord:              defaultLogToDiscord,
		LocalRun:                  defaultLocalRun,
		MetricsToDiscord:          defaultMetricsToDiscord,
		WebcrawlerPort:            defaultWebcrawlerPort,
		TrackingChannelsDelay:     defaultTrackingChannelsDelay,
		AWSRegion:                 defaultAWSRegion,
		DiscordTokenAWSSecretName: defaultDiscordTokenAWSSecretName,
		DiscordToken:              defaultDiscordToken,
		AWSMaxRetries:             defaultAWSMaxRetries,
		MediaPath:                 defaultMediaPath,
		TrackingConfigPath:        defaultTrackingConfigPath,
		WebcrawlerHost:            defaultWebcrawlerHost,
	}
	if len(discordToken) > 0 {
		options.DiscordToken = discordToken[0]
	}
	return options
}
