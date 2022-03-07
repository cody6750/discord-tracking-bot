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
	// Default value for MediaPath. Overridaeble via environment variables.
	defaultTrackingConfigPath = "pkg/configs/tracking/"
	// Default value for WebcrawlerHost. Overridaeble via environment variables.
	defaultWebcrawlerHost = "localhost"
	// Default value for WebcrawlerProt. Overridaeble via environment variables.
	defaultWebcrawlerPort = 9090
)

type Options struct {
	LogToDiscord              bool
	LocalRun                  bool
	MetricsToDiscord          bool
	AWSMaxRetries             int
	WebcrawlerPort            int
	DiscordTokenAWSSecretName string
	DiscordToken              string
	MediaPath                 string
	TrackingConfigPath        string
	AWSRegion                 string
	WebcrawlerHost            string
}

func New() *Options {
	return &Options{
		DiscordTokenAWSSecretName: defaultDiscordTokenAWSSecretName,
		DiscordToken:              defaultDiscordToken,
		AWSMaxRetries:             defaultAWSMaxRetries,
		MediaPath:                 defaultMediaPath,
		TrackingConfigPath:        defaultTrackingConfigPath,
		AWSRegion:                 defaultAWSRegion,
		LocalRun:                  defaultLocalRun,
		LogToDiscord:              defaultLogToDiscord,
		MetricsToDiscord:          defaultMetricsToDiscord,
		WebcrawlerHost:            defaultWebcrawlerHost,
		WebcrawlerPort:            defaultWebcrawlerPort,
	}
}
