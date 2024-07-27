#!/bin/bash

DOCKER_IMAGE_NAME="thrashraf/tender-scraper"
DOCKER_TAG="latest"
VPS_USER="root"
VPS_IP="178.128.97.167"

# Build the Docker image
echo "Building Docker image..."
docker build -t $DOCKER_IMAGE_NAME:$DOCKER_TAG .

# Push the Docker image to DockerHub
echo "Pushing image to Docker Hub..."
docker push $DOCKER_IMAGE_NAME:$DOCKER_TAG

# SSH to VPS and pull/run the Docker container
echo "Deploying on VPS..."
ssh $VPS_USER@$VPS_IP << EOF

echo "Checking for any process using port 8080..."
PID=\$(sudo netstat -tulpn | grep :8080 | awk '{print \$7}' | cut -d'/' -f1)

if [ -n "\$PID" ]; then
  echo "Port 8080 is busy, trying to kill process ID \$PID..."
  sudo kill -9 \$PID
fi

echo "Pulling Docker image from Docker Hub..."
docker pull $DOCKER_IMAGE_NAME:$DOCKER_TAG

echo "Stopping any running container..."
docker stop tender-scraper || true
docker rm tender-scraper || true

echo "Starting Docker container..."
# docker run -d --name tender-scraper -p 8080:8080 --restart always $DOCKER_IMAGE_NAME:$DOCKER_TAG
EOF

echo "Deployment completed successfully."
