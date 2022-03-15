#!/bin/bash 

build_docker() {
  build_linux
  docker rmi discordbot
  docker build -t discordbot .
}

build_docker_and_run() {
  docker rm -vf $(docker ps -aq)
  docker rmi -f $(docker images -aq)
  build_linux
  docker rmi discordbot
  docker build -t discordbot .
  ACCESS_KEY=$(cat  ~/.aws/credentials | grep aws_access_key_id | cut -d "=" -f2)
  SECRET_KEY=$(cat  ~/.aws/credentials | grep aws_secret_access_key | cut -d "=" -f2)
  docker run -e AWS_ACCESS_KEY_ID=$ACCESS_KEY -e AWS_SECRET_ACCESS_KEY=$SECRET_KEY discordbot  
}

build_windows() {
  rm app.exe
  GOOS=windows go build -o app.exe
}

build_linux() {
  rm app
  GOOS=linux go build -o app
}

upload_and_pull_config(){
  upload_to_s3 pkg/configs s3://${S3_BUCKET}/pkg/configs/
  upload_to_s3 media/ s3://${S3_BUCKET}/media
  send_command "sudo aws s3 cp s3://${S3_BUCKET}/pkg/ /pkg --recursive"
  send_command "sudo aws s3 cp s3://${S3_BUCKET}/media/ /media --recursive"
}

upload_to_s3() {
  file_to_upload=$1
  bucket_name=$2
  aws s3 cp $1 $2 --recursive
}

push() {
  aws ecr get-login-password --region ${REGION} | docker login --username AWS --password-stdin ${ACCOUNT_ID}.dkr.ecr.${REGION}.amazonaws.com
  docker tag discordbot:latest ${ACCOUNT_ID}.dkr.ecr.${REGION}.amazonaws.com/discordbot:latest 
  docker push ${ACCOUNT_ID}.dkr.ecr.${REGION}.amazonaws.com/discordbot:latest

}

send_command() {
  aws ssm send-command \
    --document-name "AWS-RunShellScript" \
    --document-version "1" \
    --targets "Key=instanceids,Values=${INSTANCE_ID}" \
    --parameters "commands='$1'" \
    --timeout-seconds 600 \
    --max-concurrency "50" \
    --max-errors "0" \
    --region ${REGION}
}

pull() {
  echo "Pulling docker image from ECR to EC2 Instance"
  send_command "aws ecr get-login-password --region ${REGION} | docker login --username AWS --password-stdin ${ACCOUNT_ID}.dkr.ecr.${REGION}.amazonaws.com"
  send_command "docker pull ${ACCOUNT_ID}.dkr.ecr.${REGION}.amazonaws.com/discordbot:latest"
  send_command "docker image tag ${ACCOUNT_ID}.dkr.ecr.${REGION}.amazonaws.com/discordbot:latest discordbot:latest"
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

echo "Getting flags"
while getopts :a:r:i:s:b: opt; do
  case $opt in
    a)
      ACCOUNT_ID=${OPTARG}
      ;;
    r)
      REGION=${OPTARG}
      ;;
    i)
      INSTANCE_ID=${OPTARG}
      ;;
    s)
      S3_BUCKET=${OPTARG}
      ;;
    b)
      BUILD_STEP=${OPTARG}
      ;;
  *)
    echo "$opt Flag not supported"
    ;;
 esac
done
echo "Successfully got all flags"

echo "Starting build step: $BUILD_STEP"
case $BUILD_STEP in
  r)
    build_docker_and_run
    ;;
  build_docker)
    build_docker
    ;;
  build_windows)
    build_windows
    ;;
  build_linux)
    build_linux
    ;;
  upload_config)
    upload_and_pull_config
    ;;
  build_and_deploy)
    echo "Building go executable"
    build_docker
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
