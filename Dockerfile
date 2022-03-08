FROM golang:1.16-alpine AS production

# Override if necessary during build time
ENV AWS_REGION="us-east-1"
ENV AWS_MAX_RETIRES="5"
ENV DISCORD_TOKEN_AWS_SECRET_NAME="discord/token"
ENV DISCORD_TOKEN=""
ENV LOCAL_RUN="false"
ENV LOG_TO_DISCORD="true"
ENV MEDIA_PATH="/media/"
ENV METRICS_TO_DISCORD="true"
ENV TRACKING_CONFIG_PATH="/pkg/configs/tracking/"
ENV TRACKING_CHANNELS_DELAY="21600"
ENV WEBCRAWLER_HOST="webcrawler"
ENV WEBCRAWLER_PORT="9090"
COPY app .

CMD ["./app"]