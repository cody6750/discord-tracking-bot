FROM golang:1.16-alpine AS production

ENV DISCORD_TOKEN = ''
ENV DISCORD_TOKEN_ID = 'discord/token'
ENV RUN_LOCALLY = 'falses'
COPY app .
CMD ["./app"]