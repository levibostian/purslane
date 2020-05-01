#!/bin/bash

#----------------
# Required environment variables that you need to pass in
# DO_TOKEN - digitalocean api token
# DO_SSH_FINGERPRINT - fingerprint of ssh key registered to DO account to use. Go to https://cloud.digitalocean.com/account/security to add a key and to find the fingerprint. 
#----------------

# Notes found during dev to be helpful later. 
# * DO names for volumes and stuff requires names must be lowercase and alphanumeric. hyphens are ok. 
# * We have to use a ssh key to login to the droplet.
# * Can we use a user script to set environment variables on the system? if so, we sould try to set the docker login password so it's ready. 

# Create block storage volume 
echo "Creating volume"
CREATE_VOLUME_BODY='
{"size_gigabytes":1, "name": "purslane-testing-volume", "description": "Volume for project Purslane. Ok to destroy. Its just for testing.", "region": "nyc1", "filesystem_type": "ext4", "filesystem_label": "purslane"}
'
#curl -X POST -H "Content-Type: application/json" -H "Authorization: Bearer $DO_TOKEN" -d "$CREATE_VOLUME_BODY" "https://api.digitalocean.com/v2/volumes"
# Example response: {"volume":{"id":"ce45863d-8b50-11ea-ad6b-0a58ac14477a","name":"purslane-testing-volume","created_at":"2020-05-01T02:09:36Z","description":"Volume for project Purslane. Ok to destroy. Its just for testing.","droplet_ids":null,"region":{"features":["private_networking","backups","ipv6","metadata","install_agent","storage","image_transfer"],"name":"New York 1","slug":"nyc1","sizes":["s-1vcpu-1gb","512mb","s-1vcpu-2gb","1gb","s-3vcpu-1gb","s-2vcpu-2gb","s-1vcpu-3gb","s-2vcpu-4gb","2gb","s-4vcpu-8gb","m-1vcpu-8gb","c-2","4gb","g-2vcpu-8gb","m-16gb","s-6vcpu-16gb","c-4","8gb","m-2vcpu-16gb","m3-2vcpu-16gb","g-4vcpu-16gb","m6-2vcpu-16gb","m-32gb","s-8vcpu-32gb","c-8","16gb","m-4vcpu-32gb","m3-4vcpu-32gb","g-8vcpu-32gb","s-12vcpu-48gb","m6-4vcpu-32gb","m-64gb","s-16vcpu-64gb","c-16","32gb","m-8vcpu-64gb","m3-8vcpu-64gb","g-16vcpu-64gb","s-20vcpu-96gb","48gb","m6-8vcpu-64gb","m-128gb","s-24vcpu-128gb","64gb","g-32vcpu-128gb","s-32vcpu-192gb","m-24vcpu-192gb"],"available":true},"size_gigabytes":1,"filesystem_type":"ext4","filesystem_label":"purslane","tags":null}}

# You need to get the returned ID from the response because that's the ID you put into the next request to create a droplet with that volume attached. 

# Create droplet 
echo "Create droplet"
CREATE_DROPLET_BODY='
{"name":"purslane-testing-droplet","region":"nyc1","size":"s-1vcpu-1gb","image":"docker-18-04","backups":false,"ipv6":true,"user_data":null,"private_networking":null, "volumes": ["dd633e6a-8bac-11ea-ad6b-0a58ac14477a"]
'
CREATE_DROPLET_BODY="$CREATE_DROPLET_BODY, \"ssh_keys\": [\"$DO_SSH_FINGERPRINT\"] }"

# curl -X POST -H "Content-Type: application/json" -H "Authorization: Bearer $DO_TOKEN" -d "$CREATE_DROPLET_BODY" "https://api.digitalocean.com/v2/droplets"
# Example response: {"droplet":{"id":190584016,"name":"purslane-testing-droplet","memory":1024,"vcpus":1,"disk":25,"locked":false,"status":"active","kernel":null,"created_at":"2020-05-01T02:16:52Z","features":["ipv6"],"backup_ids":[],"next_backup_window":null,"snapshot_ids":[],"image":{"id":50944795,"name":"Docker 5:19.03.1~3 on 18.04","distribution":"Ubuntu","slug":"docker-18-04","public":true,"regions":["nyc3","nyc1","sfo1","nyc2","ams2","sgp1","lon1","nyc3","ams3","fra1","tor1","sfo2","blr1","sfo3"],"created_at":"2019-08-15T16:01:20Z","min_disk_size":20,"type":"application","size_gigabytes":0.8,"description":"Docker 5:19.03.1~3 on 18.04 20190815","tags":[],"status":"available"},"volume_ids":["265b43b2-8b51-11ea-ad6b-0a58ac14477a"],"size":{"slug":"s-1vcpu-1gb","memory":1024,"vcpus":1,"disk":25,"transfer":1.0,"price_monthly":5.0,"price_hourly":0.00744,"regions":["ams2","ams3","blr1","fra1","lon1","nyc1","nyc2","nyc3","sfo1","sfo2","sfo3","sgp1","tor1"],"available":true},"size_slug":"s-1vcpu-1gb","networks":{"v4":[{"ip_address":"142.93.124.171","netmask":"255.255.240.0","gateway":"142.93.112.1","type":"public"}],"v6":[{"ip_address":"2604:a880:400:d1::a7b:1001","netmask":64,"gateway":"2604:a880:400:d1::1","type":"public"}]},"region":{"name":"New York 1","slug":"nyc1","features":["private_networking","backups","ipv6","metadata","install_agent","storage","image_transfer"],"available":true,"sizes":["s-1vcpu-1gb","512mb","s-1vcpu-2gb","1gb","s-3vcpu-1gb","s-2vcpu-2gb","s-1vcpu-3gb","s-2vcpu-4gb","2gb","s-4vcpu-8gb","m-1vcpu-8gb","c-2","4gb","g-2vcpu-8gb","m-16gb","s-6vcpu-16gb","c-4","8gb","m-2vcpu-16gb","m3-2vcpu-16gb","g-4vcpu-16gb","m6-2vcpu-16gb","m-32gb","s-8vcpu-32gb","c-8","16gb","m-4vcpu-32gb","m3-4vcpu-32gb","g-8vcpu-32gb","s-12vcpu-48gb","m6-4vcpu-32gb","m-64gb","s-16vcpu-64gb","c-16","32gb","m-8vcpu-64gb","m3-8vcpu-64gb","g-16vcpu-64gb","s-20vcpu-96gb","48gb","m6-8vcpu-64gb","m-128gb","s-24vcpu-128gb","64gb","g-32vcpu-128gb","s-32vcpu-192gb","m-24vcpu-192gb"]},"tags":[]}}

# You want to ping the endpoint below and check the "status" key until it goes to "active" from "new" which means it's ready to connect to. 
# curl -X GET -H "Content-Type: application/json" -H "Authorization: Bearer $DO_TOKEN" "https://api.digitalocean.com/v2/droplets/190634206"
# example response: {"droplet":{"id":190634206,"name":"purslane-testing-droplet","memory":1024,"vcpus":1,"disk":25,"locked":false,"status":"active","kernel":null,"created_at":"2020-05-01T13:10:08Z","features":["ipv6"],"backup_ids":[],"next_backup_window":null,"snapshot_ids":[],"image":{"id":50944795,"name":"Docker 5:19.03.1~3 on 18.04","distribution":"Ubuntu","slug":"docker-18-04","public":true,"regions":["nyc3","nyc1","sfo1","nyc2","ams2","sgp1","lon1","nyc3","ams3","fra1","tor1","sfo2","blr1","sfo3"],"created_at":"2019-08-15T16:01:20Z","min_disk_size":20,"type":"application","size_gigabytes":0.8,"description":"Docker 5:19.03.1~3 on 18.04 20190815","tags":[],"status":"available"},"volume_ids":["dd633e6a-8bac-11ea-ad6b-0a58ac14477a"],"size":{"slug":"s-1vcpu-1gb","memory":1024,"vcpus":1,"disk":25,"transfer":1.0,"price_monthly":5.0,"price_hourly":0.00744,"regions":["ams2","ams3","blr1","fra1","lon1","nyc1","nyc2","nyc3","sfo1","sfo2","sfo3","sgp1","tor1"],"available":true},"size_slug":"s-1vcpu-1gb","networks":{"v4":[{"ip_address":"178.128.149.127","netmask":"255.255.240.0","gateway":"178.128.144.1","type":"public"}],"v6":[{"ip_address":"2604:a880:400:d1::97a:3001","netmask":64,"gateway":"2604:a880:400:d1::1","type":"public"}]},"region":{"name":"New York 1","slug":"nyc1","features":["private_networking","backups","ipv6","metadata","install_agent","storage","image_transfer"],"available":true,"sizes":["s-1vcpu-1gb","512mb","s-1vcpu-2gb","1gb","s-3vcpu-1gb","s-2vcpu-2gb","s-1vcpu-3gb","s-2vcpu-4gb","2gb","s-4vcpu-8gb","m-1vcpu-8gb","c-2","4gb","g-2vcpu-8gb","m-16gb","s-6vcpu-16gb","c-4","8gb","m-2vcpu-16gb","m3-2vcpu-16gb","g-4vcpu-16gb","m6-2vcpu-16gb","m-32gb","s-8vcpu-32gb","c-8","16gb","m-4vcpu-32gb","m3-4vcpu-32gb","g-8vcpu-32gb","s-12vcpu-48gb","m6-4vcpu-32gb","m-64gb","s-16vcpu-64gb","c-16","32gb","m-8vcpu-64gb","m3-8vcpu-64gb","g-16vcpu-64gb","s-20vcpu-96gb","48gb","m6-8vcpu-64gb","m-128gb","s-24vcpu-128gb","64gb","g-32vcpu-128gb","s-32vcpu-192gb","m-24vcpu-192gb"]},"tags":[]}}

# You get the ip address from the response. and because you added your ssh key, you're ready to perform the ssh connection once the status is active. 
# ssh into the droplet with `ssh root@ip-address` but make sure you run ssh with flags to accept the fingerprint yes/no question. 
# login to registry `echo $DOCKER_LOGIN_PASSWORD | docker login docker.pkg.github.com --username levibostian --password-stdin` this works without a password prompt. make sure to allow the user to configure this line because they might use a differnet registry or might use a public one. 

# Run docker container `docker run --rm -v /home/app/extra_stuff:/mnt/purslane_testing_volume/ docker.pkg.github.com/levibostian/purslane/test-docker-image:latest`
#   I will want to add more stuff to this like environment variables and port forwarding so the container can talk to a DB via ssh tunnel. this is somethign i still need to test. my docker container used for testing was too basic and we didn't test ssh tunnels and port binding. 
#   I will want to capture the exit code and stdout/stderr of the docker container to map that stderr/stdout to the CLI's stdout/stderr and make sure the exit code of the CLI is the exit code of the docker container running. 

# Time to delete the droplet and it's resources. This endpoint used below takes an array of volumes, snapshots, etc to delete along with droplet. There is another endpoint where you can delete everything attached to the droplet but we will use this one instead as it's safer in case we decide to attach to other stuff in the future. Only delete the volume we created for you. 
echo "Deleting resources"

# Path of the url is the droplet ID. the body contains the volume ID created above. 
# curl -X DELETE -H "Content-Type: application/json" -H "Authorization: Bearer $DO_TOKEN" -d '{"volumes": ["dd633e6a-8bac-11ea-ad6b-0a58ac14477a"]}' "https://api.digitalocean.com/v2/droplets/190634206/destroy_with_associated_resources/selective"



