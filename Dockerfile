FROM golang:1.16-alpine AS production

# Override if necessary during build time
ENV DISCORD_TOKEN=""
ENV DISCORD_TOKEN_ID="discord/token"
ENV LOCAL_RUN="false"
ENV MEDIA_PATH="/media/"
ENV TRACKING_CONFIG_PATH="/pkg/configs/tracking/"
ENV WEBCRAWLER_HOST="webcrawler"
ENV WEBCRAWLER_PORT=9090

COPY app .

CMD ["./app"]