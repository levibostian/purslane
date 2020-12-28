# Deployment 

Purslane can be deployed in many different ways for your tech stack. Because Purslane is just a CLI compiled for Windows, macOS, and Linux the possibilities are endless. 

No matter how you deploy Purslane, you need to follow these steps:

1. Install Purslane on the machine. See the [README](https://github.com/levibostian/purslane#readme) to learn how to install on your operating system. 
2. Configure Purslane so it knows what to do. You can configure Purslane with a configuration file, environment variables, or command line arguments to the CLI. 

> Tip: If you use a configuration file on a server, it's recommended to edit the permissions of this file to be readable *only by the user executing the CLI*. The configuration file may contain secrets such as API keys in it. 

3. Run the CLI program. Run it on your local machine as a CLI in your terminal, run it on a server via Crontab, run it on a Kubernetes cluster as a Cronjob, all up to you. Below in this document we cover some of the more complex examples to give you a boost on getting Purslane up and running. 


# Examples

### Docker

There are [official Docker images](https://hub.docker.com/r/levibostian/purslane/tags?page=1&ordering=last_updated) deployed on each Purslane release. 

You can run Purslane inside of the docker image easily:
```
docker run -v ${PWD}/purslane_config.yml:/.purslane.yml levibostian/purslane:latest
```
> Note: The command above assumes that you have a Purslane config file called `purslane_config.yml` in the current directory. You can also pass in environment variables and command line arguments to the Docker image. See [`docker run` documentation](https://docs.docker.com/engine/reference/commandline/run/) to learn more. 

The Docker images for Purslane are all setup to automatically run `purslane` when the container runs. You can easily pass arguments to the CLI by adding them to the end:
```
docker run ... levibostian/purslane:latest --help
```
This Docker container will run `purslane --help` inside of the Docker container. 

##### Run more commands then just Purslane 

If all you need to do is execute the Purslane CLI, you are good to use the default Docker image, `levibostian/purslane:X`. This image does *not* come with a shell. That means this image is the smallest size and the most secure. However, if you need to run more commands (for example a curl command), then you may need a more robust Docker image. 

That is where the `levibostian/purslane:X-alpine` image comes in. The alpine image comes with a shell, the ability to install programs like `curl`, and more. It is still a small image but allows you to do more then just run the Purslane CLI. 

> Tip: It might be easier to simply execute a shell script inside of the Docker image to run all of your commands:
```
docker run -v ${PWD}/db_backups.sh:/db_backups.sh --entrypoint /bin/sh levibostian/purslane:latest-alpine -c /db_backups.sh
```

### Kubernetes cronjob 

Kubernetes has [cronjob functionality](https://kubernetes.io/docs/concepts/workloads/controllers/cron-jobs/) allowing you to execute Purslane periodically. 

To do this, you need to create a Kubernetes manifest file. 

```yml
apiVersion: batch/v1beta1
kind: CronJob
metadata:
  name: purslane-cronjob
spec:
  schedule: "30 5 * * *" # every day at 5:30am
  concurrencyPolicy: Replace
  jobTemplate:
    spec:
      template:
        spec:
          containers:
          - name: purslane-job
            image: levibostian/purslane:1 
            resources: # Setting limits is optional, but recommended. The limits below were defaults found in docs and may/may not work for you. Running purslane requires minimal resources. 
              limits: 
                memory: "64Mi"
                cpu: "60m"
              requests:
                memory: "32Mi"
                cpu: "40m"
          restartPolicy: OnFailure
```

> Note: It's recommended to view the [official docs on Kubernetes cronjobs](https://kubernetes.io/docs/concepts/workloads/controllers/cron-jobs/) as this syntax below might be out-of-date. 

Now, we need to provide configuration to Purslane. It's *highly* recommended to do this with [Kubernetes secrets](https://kubernetes.io/docs/concepts/configuration/secret/). 

// TODO talk about secrets 