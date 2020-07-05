# Purslane

Scalable, flexible, affordable, and safe way to perform periodic tasks. All from a convenient CLI binary!

*Note:* 
1. Purslane is being used in production for 1 or 2 projects to test it out but that's the extent of the testing thus far. 
2. The project is young so breaking changes may occur. However, this project will do it's best to always be safe to use by following semantic versioning. 
3. Purslane has only been tested through use and is not unit tested. It is planned, but has not yet been completed. 

![logo](misc/logo.jpg)
> credits: PaoloBis / Getty Images

## Why use Purslane? 

Let's use an example. You have a production database that you want to perform daily backups on. This seems like a simple task, but there are some issues that you could run into:

* **What machine should I run my backup on?...** Performing database backups takes resources (CPU, memory) and you don't want a daily backup to slow down your servers so it's best to run backups on a separate machine with resources that can handle the backup. 
* **Disk storage...** As the size of your database grows, you will need larger disk storage to store the dump on. 
* **Pain to scale...** After you setup your server for database dumps, your app gets popular. As time goes on you will need to spend time upgrading your backup system to handle the new load.
* **Higher costs...** If I run backups daily, that means that 23 1/2 hours of each day that server is doing nothing and you still need to pay for it.

Purslane solves all of these problems for you! Let's read about how Purslane works to understand how all of these problems are gone.

## How does Purslane work?

Purslane is quite simple, really. Each time Purslane is run, it will:
1. Create a new cloud server and optionally attach disk storage to it. (Purslane works with [DigitalOcean](https://www.digitalocean.com/) at this time but may include other cloud providers in the future)
2. Run a Docker container of your choice in the newly created cloud server. 
3. When the Docker container exits, Purslane will destroy all of the resources it created such as the created server and disk storage. 

For you, this means:
* **Your application has no performance hit...** Because your task runs on a brand new server, your infrastructure does not take a hit. 
* **Scale with 1 line of code...** As the size of your application grows and you need to scale your periodic tasks, it only takes 1 quick change to your Purslane config file and Purslane will scale the created server and disk storage the next time you run your task. 
* **Affordable...** You only pay for the time your Docker container runs. 
* **Flexible...** Run any task you need! You are the only that makes the Docker container, so you can run anything you need to run! Any language, any software, any task. 

All you need to do is give Purslane authentication details to login to your cloud provider to create the server, set the size of the server and optional disk storage you want, and provide the Docker image to run. Purslane will take care of the rest. 

# Getting started

* The first thing you will need is to create a Docker image for Purslane to execute. You can host this Docker image in a public or private Docker image registry. 

* Create the Purslane config file that tells Purslane how to execute your Docker image. You can create it in the default location of `~/.purslane` or create it anywhere and point to it with the `--config` CLI argument. 

```yaml
cloud: # Required. Must set a cloud provider. At this time, Purslane only works with DigitalOcean
  digitalocean:
    # Create an API key with write access as we are creating resources. 
    # https://www.digitalocean.com/docs/apis-clis/api/create-personal-access-token/
    api_key: "123" 

volume: # Optional, but required if you want to create a volume that attaches to the server. 
  gigs: 1 # How many GBs you want the volume to be. 
  container_mount_path: "/home" # The path inside of your Docker container you want the volume to attach to. 

server: # Optional - sets default values for you if left out. 
  # this format is universal with all providers but some providers may not have the combination you specify. This format will be converted to the string specific to the cloud provider for you. 
  # the format of this string: `s-1cpu-1gb`
  # s - standard. Can also be "m" for high memory or "c" for dedicated CPU. 
  # 1cpu - how many CPUs you want. 
  # 1gb - how many gigs of memory you want. 
  size: "s-1cpu-1gb" 
  # private network to add server to when created. (optional)
  #  - DigitalOcean: Enter UUID of already created private network. Find UUID in DigitalOcean's website > select a VPC > See UUID in URL of webpage. 
  private_network: "XXXXXXXX-XXXX-XXXX-XXXX-XXXXXXXXXXXX" 

# Required. These keys are used to run commands against the server once it's created. The machine you run the CLI from will SSH into the created server and send commands to it. 
# **Note:** At this time, SSH keys must not have a passphrase on them 
public_ssh_key_path: "~/.ssh/id_rsa.pub"
private_ssh_key_path: "~/.ssh/id_rsa"

docker_run: # Optional. Purslane runs a default set that will run the docker container just fine. 
  extra_args: "-p 5000:5000" # Append arguments to the `docker run` command. Great place to add ports bindings or environment variables, for example. 

docker: # Required. 
  image: "image-name:tag" # Image with tag of docker image to pull. Can be from private or public repo. 
  registry: # optional. Only needed if image in private repo you need to authenticate with. 
    name: "name-of-registry"
    username: "docker-username"
    password: "docker-password-to-registry"
```

> Note: This file contains secret information such as API keys. It's recommended that you make this config file read-only by select users on your machine. 

* Time to execute Purslane. You decide how you would like to execute the CLI. [Run the CLI on a server](#run-cli-on-a-server) or [run via Docker](#run-cli-via-docker). 

### Run CLI on a server 

Since Purslane is a very low resources CLI application, you can run it directly on a server of yours. 

If you create your Purslane config file on your server, you can install the CLI easily from [GitHub releases](https://github.com/levibostian/purslane/releases) with 1 command:

```
curl -sf https://gobinaries.com/levibostian/Purslane | sh
```

You now have the CLI installed on your machine. Run `purslane --help` to learn about the CLI. 

### Run CLI via Docker

You can run the Purslane CLI via a Docker image. There is [a public Purslane Docker image](https://hub.docker.com/levibostian/purslane) for you to use. When you want to use the Docker image, you need to provide to the Docker image your config file, ssh public/private key files, and maybe you need to provide other files. 

Here is an example `docker run` statement to run the Docker purslane CLI:
```
docker run --rm -v $(pwd)/purslane.yaml:/config.yaml -v /home/.ssh/:/ssh levibostian/purslane:latest run --config /config.yaml
```

> Note: The home directory of the Docker image is set to `/root`. You can use `~/` in the config file to refer to this path. 

[The Dockerfile](Dockerfile) is setup quite easily. It just runs the purslane CLI and everything that you add to the end of the `docker run` command gets passed as arguments to the CLI. So the example above with `docker run ... run --config /config.yaml` means that you are running `/purslane run --config /config.yaml` inside of the Docker container. 

## Development 

Purslane is a Go lang program. To start developing Purslane is as simple as (1) cloning the repo and (2) running `go run main.go`.

## Contribute

Purslane is open for pull requests. Check out the [list of issues](https://github.com/levibostian/purslane/issues) for tasks planned to be worked on. Check them out if you wish to contribute in that way.

**Want to add features to Purslane?** Before you decide to take a bunch of time and add functionality to the CLI, please, [create an issue](https://github.com/levibostian/Purslane/issues/new) stating what you wish to add. This might save you some time in case your purpose does not fit well in the use cases of Purslane.
