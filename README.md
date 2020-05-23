# Purslane

Scalable, flexible, affordable, and safe way to perform periodic tasks. All from a convenient CLI binary!

*Note:* Purslane is in development. It is being used in production for 1 or 2 projects to test it out but that's the extent of the testing thus far. 

![logo](misc/logo.jpg)
> credits: PaoloBis / Getty Images

## Why use Purslane? 

Let's use an example. You have a production database that you want to perform daily backups on. This seems like a simple task, but there are some issues that you could run into:

* **What machine should I run my backup on?...** Performing database backups takes resources (CPU, memory) and you don't want a daily backup to slow down your servers. 
* **Disk storage...** As the size of your database grows, you will need larger disk storage to store the dump on. 
* **Pain to scale...** After you setup your server for database dumps, your app gets popular as time goes on, you will need to spend time upgrading your backup system to handle the new load.
* **Higher costs...** If I run backups daily, that means that 23 1/2 hours of each day that server is doing nothing and you still need to pay for it.

Purslane solves all of these problems for you! Let's read about how Purslane works to understand how all of these problems are gone.

## How does Purslane work?

Purslane is quite simple, really. Each time Purslane is run, it will:
1. Create a new cloud server and optionally attach disk storage to it. 
2. Run a Docker container in the newly created cloud server. 
3. When the Docker container exits, Purslane will destroy all of the resources it created such as the created server and disk storage. 

For you, this means:
* **Your application has no performance hit...** Because your task runs on a brand new server, your infrastructure does not take a hit. 
* **Scale with 1 line of code...** As the size of your application grows and you need to scale your periodic tasks, it only takes 1 quick change to your config file and Purslane will scale your server and disk storage the next time you run your task. 
* **Affordable...** You only pay for the time your Docker container runs. 
* **Flexible...** Run any task you need! You are the only that makes the Docker container, so you can run anything you need to run! All language, any software, any task. 

All you need to do is give Purslane authentication details to login to your cloud provider to create the server, set the size of the server and optional disk storage you want, and provide the Docker image to run. Purslane will take care of the rest. 

# Getting started

Getting started docs coming soon... Purslane needs to be deployed first, then it can be installed and executed. 

# Configuration 

...coming soon... in the mean time, check out the `example-config.yaml` file. 

## Development 

Purslane is a Go lang program. To start developing Purslane is as simple as (1) cloning the repo and (2) running `go run main.go`.

## Author

* Levi Bostian - [GitHub](https://github.com/levibostian), [Twitter](https://twitter.com/levibostian), [Website/blog](http://levibostian.com)

![Levi Bostian image](https://gravatar.com/avatar/22355580305146b21508c74ff6b44bc5?s=250)

## Contribute

Purslane is open for pull requests. Check out the [list of issues](https://github.com/levibostian/purslane/issues) for tasks planned to be worked on. Check them out if you wish to contribute in that way.

**Want to add features to Purslane?** Before you decide to take a bunch of time and add functionality to the CLI, please, [create an issue](https://github.com/levibostian/Purslane/issues/new) stating what you wish to add. This might save you some time in case your purpose does not fit well in the use cases of Purslane.
