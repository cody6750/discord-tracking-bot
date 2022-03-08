package trackingbot

import (
	"fmt"
	"os"
	"strconv"
)

func (t *TrackingBot) getEnvVariables() {
	var err error
	t.logger.Info("Getting environment variables")
	if os.Getenv("AWS_MAX_RETRIES") != "" {
		t.options.AWSMaxRetries, err = strconv.Atoi(os.Getenv("AWS_MAX_RETRIES"))
		if err != nil {
			t.logger.WithError(err).Fatal("Failed to convert AWS_MAX_RETRIES from string to int")
		}
		t.logger.WithField("AWS_MAX_RETRIES: ", t.options.AWSMaxRetries).Info("Successfully got environment variable")
	}

	if os.Getenv("AWS_REGION") != "" {
		t.options.AWSRegion = os.Getenv("AWS_REGION")
		t.logger.WithField("AWS_REGION: ", t.options.AWSRegion).Info("Successfully got environment variable")
	}

	if os.Getenv("DISCORD_TOKEN_AWS_SECRET_NAME") != "" {
		t.options.DiscordTokenAWSSecretName = os.Getenv("DISCORD_TOKEN_AWS_SECRET_NAME")
		t.logger.WithField("DISCORD_TOKEN_AWS_SECRET_NAME: ", t.options.DiscordTokenAWSSecretName).Info("Successfully got environment variable")
	}

	if os.Getenv("DISCORD_TOKEN") != "" {
		t.options.DiscordToken = os.Getenv("DISCORD_TOKEN")
	}

	if os.Getenv("LOCAL_RUN") != "" {
		t.options.LocalRun, err = getEnvBool("LOCAL_RUN")
		if err != nil {
			t.logger.WithError(err).Fatal("Failed to convert LOCAL_RUN from string to bool")
		}
		t.logger.WithField("LOCAL_RUN: ", t.options.LocalRun).Info("Successfully got environment variable")
	}

	if os.Getenv("LOG_TO_DISCORD") != "" {
		t.options.LogToDiscord, err = getEnvBool("LOG_TO_DISCORD")
		if err != nil {
			t.logger.WithError(err).Fatal("Failed to convert LOG_TO_DISCORD from string to bool")
		}
		t.logger.WithField("LOG_TO_DISCORD: ", t.options.LogToDiscord).Info("Successfully got environment variable")

	}

	if os.Getenv("MEDIA_PATH") != "" {
		t.options.MediaPath = os.Getenv("MEDIA_PATH")
		t.logger.WithField("MEDIA_PATH: ", t.options.MediaPath).Info("Successfully got environment variable")
	}

	if os.Getenv("METRICS_TO_DISCORD") != "" {
		t.options.MetricsToDiscord, err = getEnvBool("METRICS_TO_DISCORD")
		if err != nil {
			t.logger.WithError(err).Fatal("Failed to convert METRICS_TO_DISCORD from string to bool")
		}
		t.logger.WithField("METRICS_TO_DISCORD: ", t.options.MetricsToDiscord).Info("Successfully got environment variable")
	}

	if os.Getenv("TRACKING_CONFIG_PATH") != "" {
		t.options.TrackingConfigPath = os.Getenv("TRACKING_CONFIG_PATH")
		t.logger.WithField("TRACKING_CONFIG_PATH: ", t.options.TrackingConfigPath).Info("Successfully got environment variable")
	}

	if os.Getenv("TRACKING_CHANNELS_DELAY") != "" {
		t.options.TrackingChannelsDelay, err = getEnvInt("TRACKING_CHANNELS_DELAY")
		if err != nil {
			t.logger.WithError(err).Fatal("Failed to convert TRACKING_CHANNELS_DELAY from string to int")
		}
		t.logger.WithField("TRACKING_CHANNELS_DELAY: ", t.options.TrackingChannelsDelay).Info("Successfully got environment variable")
	}

	if os.Getenv("WEBCRAWLER_HOST") != "" {
		t.options.WebcrawlerHost = os.Getenv("WEBCRAWLER_HOST")
		t.logger.WithField("WEBCRAWLER_HOST: ", t.options.WebcrawlerHost).Info("Successfully got environment variable")
	}

	if os.Getenv("WEBCRAWLER_PORT") != "" {
		t.options.WebcrawlerPort, err = getEnvInt("WEBCRAWLER_PORT")
		if err != nil {
			t.logger.WithError(err).Fatal("Failed to convert WEBCRAWLER_PORT from string to int")
		}
		t.logger.WithField("WEBCRAWLER_PORT: ", t.options.WebcrawlerPort).Info("Successfully got environment variable")
	}
	t.logger.Info("Successfully got environment variables")

}

func (t *TrackingBot) checkEnvVariables() {
	t.logger.Info("Checking environment variables")
	if t.options.DiscordToken == "" {
		t.logger.Fatalf("Discord token is not set, please set DISCORD_TOKEN environment variable or check AWS")
	}
	t.logger.Info("Successfully checked environment variables")

}

func getEnvBool(envVar string) (bool, error) {
	s := os.Getenv(envVar)
	if s == "" {
		return false, fmt.Errorf("")
	}
	v, err := strconv.ParseBool(s)
	if err != nil {
		return false, err
	}
	return v, nil
}

func getEnvInt(envVar string) (int, error) {
	s := os.Getenv(envVar)
	if s == "" {
		return 0, fmt.Errorf("")
	}
	strconv.Atoi(s)
	v, err := strconv.Atoi(s)
	if err != nil {
		return v, err
	}
	return v, nil
}
