# Build 

```bash
# build the test image. Run from root directory of project. 
docker build -t docker.pkg.github.com/levibostian/purslane/test-docker-image:latest -f example_docker_image/Dockerfile example_docker_image/

# Test out the build image 
docker run --rm docker.pkg.github.com/levibostian/purslane/test-docker-image:latest

# Push to github repository. We want to test a docker repository that's not docker hub to test purslane works with unique docker repo. 
docker login docker.pkg.github.com --username levibostian
docker push docker.pkg.github.com/levibostian/purslane/test-docker-image:latest
```