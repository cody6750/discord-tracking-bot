package options

const (
	// DefaultAWSRegion default value for AWSRegion. Overrideable via environment variables.
	DefaultAWSRegion = "us-east-1"

	// DefaultAWSMaxRetries default value for AWSAWSMaxRetries. Overrideable via environment variables.
	DefaultAWSMaxRetries = 5

	// DefaultDiscordTokenAWSSecretName default value for DiscordTokenAWSSecretName. Overrideable via environment variables.
	DefaultDiscordTokenAWSSecretName = "discord/token"

	// DefaultDiscordToken default value for DiscordToken. Overrideable via environment variables.
	DefaultDiscordToken = ""

	// DefaultLocalRun default value for LocalRun. Overrideable via environment variables.
	DefaultLocalRun = false

	// DefaultLogToDiscord default value for LogToDiscord. Overridaeble via environment variables.
	DefaultLogToDiscord = true

	// DefaultMediaPath default value for MediaPath. Overridaeble via environment variables.
	DefaultMediaPath = "media/"

	// DefaultMetricsToDiscord default value for MetricsToDiscord. Overridaeble via environment variables.
	DefaultMetricsToDiscord = true

	// DefaultTrackingChannelsDelay default value for TrackingConfigPath. Overridaeble via environment variables.
	DefaultTrackingChannelsDelay = 21600

	// DefaultTrackingConfigPath default value for TrackingConfigPath. Overridaeble via environment variables.
	DefaultTrackingConfigPath = "pkg/configs/tracking/"

	// DefaultWebcrawlerHost default value for WebcrawlerHost. Overridaeble via environment variables.
	DefaultWebcrawlerHost = "localhost"

	// DefaultWebcrawlerPort default value for WebcrawlerProt. Overridaeble via environment variables.
	DefaultWebcrawlerPort = 9090
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

//New returns Options with Default values.
func New(discordToken ...string) *Options {
	options := &Options{
		LogToDiscord:              DefaultLogToDiscord,
		LocalRun:                  DefaultLocalRun,
		MetricsToDiscord:          DefaultMetricsToDiscord,
		WebcrawlerPort:            DefaultWebcrawlerPort,
		TrackingChannelsDelay:     DefaultTrackingChannelsDelay,
		AWSRegion:                 DefaultAWSRegion,
		DiscordTokenAWSSecretName: DefaultDiscordTokenAWSSecretName,
		DiscordToken:              DefaultDiscordToken,
		AWSMaxRetries:             DefaultAWSMaxRetries,
		MediaPath:                 DefaultMediaPath,
		TrackingConfigPath:        DefaultTrackingConfigPath,
		WebcrawlerHost:            DefaultWebcrawlerHost,
	}
	if len(discordToken) > 0 {
		options.DiscordToken = discordToken[0]
	}
	return options
}
