#!/bin/bash 

build() {
  rm app
  GOOS=linux go build -o app
  docker rmi discordbot
  docker build . -t discordbot
}

push() {
  send_command "docker system prune -f"
  aws ecr get-login-password --region us-east-1 | docker login --username AWS --password-stdin 732407346024.dkr.ecr.us-east-1.amazonaws.com
  docker tag discordbot:latest 732407346024.dkr.ecr.us-east-1.amazonaws.com/discordbot:latest
  docker push 732407346024.dkr.ecr.us-east-1.amazonaws.com/discordbot:latest
}

send_command() {
  aws ssm send-command \
    --document-name "AWS-RunShellScript" \
    --document-version "1" \
    --targets "Key=instanceids,Values=i-09006bb3612ab4df6" \
    --parameters "commands='$1'" \
    --timeout-seconds 600 \
    --max-concurrency "50" \
    --max-errors "0" \
    --region us-east-1
}

pull() {
  echo "Pulling docker image from ECR to EC2 Instance"
  send_command "aws ecr get-login-password --region us-east-1 | docker login --username AWS --password-stdin 732407346024.dkr.ecr.us-east-1.amazonaws.com"
  send_command "docker pull 732407346024.dkr.ecr.us-east-1.amazonaws.com/discordbot:latest"
  send_command "docker image tag 732407346024.dkr.ecr.us-east-1.amazonaws.com/discordbot:latest discordBot:latest"
  echo "Successfully pulled docker image from ECR to EC2 Instance"
}

run() {
  send_command "docker run -d -p 9090:9090 --name discordbot discordBot"
}

stop() {
  send_command "docker stop discordbot"
}

init_ec2(){
  send_command "sudo chmod 666 /var/run/docker.sock"
}

echo "Starting build"
case $1 in
  build_and_deploy)
    echo "Building go executable"
    build
    push
    pull
    run
    ;;
  init)
    echo "Init EC2"
    init_ec2
    ;;
  run)
    echo "Running docker container"
    run
    ;;
  stop)
    echo "Stoping EC2 instance"
    stop
    ;;
  *)
    echo "No flags passed, doing nothing"
    ;;
esac
echo "Finished build"
