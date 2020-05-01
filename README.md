# Purslane
[WIP] more details to come soon...

# Getting started

1. You need a docker container to start with for testing. See [the included docker image readme](example_docker_image/README.md) for how to build and push that. We are using github packages for docker image hosting as it's not the default docker hub. We want to test out Purslane can connect to a repository that's private or not hub. Build the push the docker image. 

2. Run the `./script.sh` on your machine. Read the header on the top of the file. It will tell you what environment variables you must pass into the script. The script is not automated. It requires commenting out lines, running it, uncommenting and commenting lines, running, etc. It's full of notes in order to help make the real automated CLI tool. 