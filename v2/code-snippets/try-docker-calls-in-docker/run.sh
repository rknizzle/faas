#!/usr/bin/env bash

# build the image
docker build -t try-docker-calls-in-docker-image .

# run without having access to the hosts docker daemon socket. we are expecting this container to
# fail because it wont have access to a Docker daemon to run the docker pull command
docker run --rm try-docker-calls-in-docker-image
echo ""
echo ""
echo "##############################################"
echo "##############################################"
echo "This failed to pull the image because this container doesnt have access to a Docker daemon"
echo "##############################################"
echo "##############################################"
echo ""

docker run --rm -v /var/run/docker.sock:/var/run/docker.sock try-docker-calls-in-docker-image

echo ""
echo "##############################################"
echo "##############################################"
echo "This container should have successfully pulled the Docker image because it has access to the host machines Docker daemon"
echo "##############################################"
echo "##############################################"

# cleanup the image
docker rmi try-docker-calls-in-docker-image >/dev/null
