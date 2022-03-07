package options

const (
	// Default value for DiscordTokenID. Overrideable via environment variables.
	defaultDiscordTokenID = "discord/token"
	// Default value for DiscordToken. Overrideable via environment variables.
	defaultDiscordToken = ""
	// Default value for MaxRetries. Overrideable via environment variables.
	defaultMaxRetries = 5
	// Default valur for Region. Overrideable via environment variables.
	defaultRegion = "us-east-1"
	// Default value for LocalRun. Overrideable via environment variables.
	defaultLocalRun = false
	// Default value for MediaPath. Overridaeble via environment variables.
	defaultMediaPath = "media/"
	// Default value for MediaPath. Overridaeble via environment variables.
	defaultTrackingConfigPath = "pkg/configs/tracking/"
	// Default value for LogToDiscord. Overridaeble via environment variables.
	defaultLogToDiscord = true
	// Default value for MetricsToDiscord. Overridaeble via environment variables.
	defaultMetricsToDiscord = true
	// Default value for WebcrawlerHost. Overridaeble via environment variables.
	defaultWebcrawlerHost = "localhost"
	// Default value for WebcrawlerProt. Overridaeble via environment variables.
	defaultWebcrawlerPort = 9090
)

type Options struct {
	LogToDiscord       bool
	LocalRun           bool
	MetricsToDiscord   bool
	MaxRetries         int
	WebcrawlerPort     int
	DiscordTokenID     string
	DiscordToken       string
	MediaPath          string
	TrackingConfigPath string
	Region             string
	WebcrawlerHost     string
}

func New() *Options {
	return &Options{
		DiscordTokenID:     defaultDiscordTokenID,
		DiscordToken:       defaultDiscordToken,
		MaxRetries:         defaultMaxRetries,
		MediaPath:          defaultMediaPath,
		TrackingConfigPath: defaultTrackingConfigPath,
		Region:             defaultRegion,
		LocalRun:           defaultLocalRun,
		LogToDiscord:       defaultLogToDiscord,
		MetricsToDiscord:   defaultMetricsToDiscord,
		WebcrawlerHost:     defaultWebcrawlerHost,
		WebcrawlerPort:     defaultWebcrawlerPort,
	}
}
