# Run Invoker In A Docker Container

The purpose of this code snippet is to determine if its possible the run the invoker portion of the
application inside a Docker container. This isnt totally straight forward because the invoker
interacts with a Docker daemon to pull and spawn Docker containers. In order to be able to interact
with a Docker daemon from within a Docker container we can either:
1. Actually run an instance of the Docker program inside your container using something like [Docker
   in Docker](https://hub.docker.com/_/docker)
2. Mount the hosts docker.sock to your container so that the process in the container can access the
   host machines Docker daemon

In this example Im going to test out option 2 because if this works, this option is simpler and can
avoid some of the overhead and hidden issues of Docker in Docker


See `run.sh` for the docker commands


RESULT: Option 2 works. A program in a Docker container can interact with the host machines Docker
daemon if you mount the hosts docker.sock to the container
