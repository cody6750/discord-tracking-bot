#!/bin/bash 

build() {
  rm app
  GOOS=linux go build -o app
  docker rmi discordbot
  docker build -t discordbot .

}


upload_and_pull_config(){
  upload_to_s3 pkg/configs s3://discordbott/pkg/configs/
  upload_to_s3 media/ s3://discordbott/media
  send_command "sudo aws s3 cp s3://discordbott/pkg/ /pkg --recursive"
  send_command "sudo aws s3 cp s3://discordbott/media/ /media --recursive"
}
upload_to_s3() {
  file_to_upload=$1
  bucket_name=$2
  aws s3 cp $1 $2 --recursive
}
push() {
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
  send_command "docker image tag 732407346024.dkr.ecr.us-east-1.amazonaws.com/discordbot:latest discordbot:latest"
  echo "Successfully pulled docker image from ECR to EC2 Instance"
}

run() {
  send_command "docker run -v /pkg:/pkg -v /media:/media --network=discord --name discordbot discordbot"
}

stop() {
  send_command "docker stop discordbot"
  send_command "docker container prune -f"
}

init_ec2(){
  send_command "sudo chmod 666 /var/run/docker.sock"
}

echo "Starting build"
case $1 in
  upload_config)
    upload_and_pull_config
    ;;
  build_and_deploy)
    echo "Building go executable"
    build
    upload_and_pull_config
    push
    pull
    stop
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
