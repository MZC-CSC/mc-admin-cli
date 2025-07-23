# M-CMP ADMIN CLI (mcc)
This repository provides a Multi-Cloud ADMIN CLI.    
The name of this tool is mcc(Multi-Cloud admin CLI).    
A sub-system of [M-CMP platform](https://github.com/m-cmp) to deploy and manage Multi-Cloud Infrastructures.    


```
[NOTE]
mcc is currently under development.
So, we do not recommend using the current release in production.
Please note that the functionalities of mcc are not stable and secure yet.
If you have any difficulties in using mcc, please let us know.
(Open an issue or Join the M-CMP Slack)
```

## mcc Overview
- Management tool that supports the installation, execution, status information provision, termination, and API calls of the M-CMP system.
- Currently, infra subcommand is only support docker compose base infra install and management.
  - [infra subcommand](./docs/mc-admin-cli-infra.md)
- If you want to checkout how to run the whole subsystem on the single instance on CSP Instance, see [this document](./docs/mc-admin-cli-infra.md).

## Development & Test Environment
- Go 1.23
- Docker version 27.3.1
- Docker Compose version v2.29

## Install Docker & Docker Compose V2

- [Install Docker Engine on Ubuntu](https://docs.docker.com/engine/install/ubuntu/)

checkout the commands down below.

```shell
sudo apt-get install -y apt-transport-https ca-certificates curl gnupg-agent software-properties-common

curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo apt-key add -

sudo add-apt-repository "deb [arch=amd64] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable"

sudo apt-get update

sudo apt-get install -y docker-ce docker-ce-cli docker-compose-plugin
```

# Quick Guide
 Required: Domain, Email
  - SSL and IDP configuration requires certificates. (If you already have certificates in use, substitute them into the nginx configuration values of the iam manager)

1. Certificate Issuance (docker execution)
 Execution location: ./conf/docker
```shell 
sudo docker compose -f docker-compose.cert.yaml up
```
When issuance is successful, pem files are created in the folder with the corresponding domain name.
ex)
 ./container-volume/mc-iam-manager/certs/live/<my domain>fullchain.pem
 ./container-volume/mc-iam-manager/certs/live/<my domain>/privkey.pem

After domain issuance is complete, run docker down (to be used again for future renewals)
```shell 
sudo docker compose -f docker-compose.cert.yaml up
```

2. Execute nginx configuration script
 (If the owner of container-volume is root, change ownership to the appropriate user)
```shell 
  sudo chown -R ubuntu:ubuntu ./container-volume
  cd mc-iam-manager
  ./0_preset_create_nginx_conf.sh
```
 Creates files in container-volume/mc-iam-manager/nginx using template files and configuration information under the mc-iam-manager folder.

3. After completing admin cli docker configuration, start the service with docker compose.
```shell 
sudo docker compose -f docker-compose.yaml
```


---

# Command to build the operator from souce code
```Shell
$ git clone https://github.com/m-cmp/mc-admin-cli.git
$ cd mc-admin-cli/src

(Setup dependencies)
mc-admin-cli/src$ go get -u

(Build a binary for mcc)
mc-admin-cli/src$ go build -o mcc

**Build a binary for mcc using Makerfile depends on your machine\'s os type**
mc-admin-cli/src$ make
mc-admin-cli/src$ make win
mc-admin-cli/src$ make mac
mc-admin-cli/src$ make linux-arm
mc-admin-cli/src$ make win86
mc-admin-cli/src$ make mac-arm
```

# How to use the mcc

```
mc-admin-cli/bin$ ./mcc -h

The mcc is a tool to operate Cloud-Barista system. 
  
Usage:
  mcc [command]

Available Commands:
  api         Call the M-CMP system's Open APIs as services and actions
  infra       A tool to operate M-CMP system
  help        Help about any command
  rest        rest api call

Flags:
  -h, --help   help for mcc

Use "mcc [command] --help" for more information about a command.
```

For more detailed explanations, see the articles below.   
- [infra sub-command guide](./docs/mc-admin-cli-infra.md)
- [rest sub-command guide](./docs/mc-admin-cli-rest.md)

## docker-compose.yaml
```
The necessary service information for the M-CMP System configuration is defined in the mc-admin-cli/docker-compose-mode-files/docker-compose.yaml file.(By default, it is set to build the desired configuration and data volume in the docker-compose-mode-files folder.)

If you want to change the information for each container you want to deploy, modify the mc-admin-cli/docker-compose-mode-files/docker-compose.yaml file or use the -f option.
```

## infra subcommand
For more information, check out [the infra subcommand document.](./docs/mc-admin-cli-infra.md)


For now, it supports infra's run/stop/info/pull/remove commands.

Use the -h option at the end of the sub-command requiring assistance, or executing 'mcc' without any options will display the help manual.


```
Usage:
  mcc infra [flags]
  mcc infra [command]

Available Commands:
  info        Get information of M-CMP System
  pull        Pull images of M-CMP System containers
  remove      Stop and Remove M-CMP System
  run         Setup and Run M-CMP System
  stop        Stop M-CMP System

Flags:
  -h, --help   help for infra

Use "mcc infra [command] --help" for more information about a command.
```

## infra subcommand examples
Simple usage examples for infra subcommand
```
- ./mcc infra pull [-f ../conf/docker/docker-compose.yaml]
- ./mcc infra run [-f ../conf/docker/docker-compose.yaml]  -d
- ./mcc infra info
- ./mcc infra stop [-f ../conf/docker/docker-compose.yaml]
- ./mcc infra remove [-f ../conf/docker/docker-compose.yaml] -v -i

```
## k8s subcommand
**K8S is not currently supported and will be supported in the near future.**

## rest subcommand
The rest subcommands are developed around the basic features of REST to make it easy to use the open APIs of M-CMP-related frameworks from the CLI.
For now, it supports get/post/delete/put/patch commands.

For more information, check out [the rest subcommand document.](./docs/mc-admin-cli-rest.md)

```
rest api call

Usage:
  mcc rest [flags]
  mcc rest [command]

Available Commands:
  delete      REST API calls with DELETE methods
  get         REST API calls with GET methods
  patch       REST API calls with PATCH methods
  post        REST API calls with POST methods
  put         REST API calls with PUT methods

Flags:
      --authScheme string   sets the auth scheme type in the HTTP request.(Exam. OAuth)(The default auth scheme is Bearer)
      --authToken string    sets the auth token of the 'Authorization' header for all HTTP requests.(The default auth scheme is 'Bearer')
  -d, --data string         Data to send to the server
  -f, --file string         Data to send to the server from file
  -I, --head                Show response headers only
  -H, --header strings      Pass custom header(s) to server
  -h, --help                help for rest
  -p, --password string     Password for basic authentication
  -u, --user string         Username for basic authentication
  -v, --verbose             Show more detail information

Use "mcc rest [command] --help" for more information about a command.
```

## rest command examples
Simple usage examples for rest commands

```
./mcc rest get -u default -p default http://localhost:1323/tumblebug/health
./mcc rest post https://reqres.in/api/users -d '{
                "name": "morpheus",
                "job": "leader"
        }'
```

## api subcommand
For more information, check out [the infra subcommand document.](./docs/mc-admin-cli-api.md)
The api subcommands are developed to make it easy to use the open APIs of M-CMP-related frameworks from the CLI.

```
Call the action of the service defined in api.yaml. 

Usage:
  mcc api [flags]
  mcc api [command]

Available Commands:
  tool        Swagger JSON parsing tool to assist in writing api.yaml files

Flags:
  -a, --action string        Action to perform
  -c, --config string        config file (default "../conf/api.yaml")
  -d, --data string          Data to send to the server
  -f, --file string          Data to send to the server from file
  -h, --help                 help for api
  -l, --list                 Show Service or Action list
  -m, --method string        HTTP Method
  -p, --pathParam string     Variable path info set "key1:value1 key2:value2" for URIs
  -q, --queryString string   Use if you have a query string to add to URIs
  -s, --service string       Service to perform
  -v, --verbose              Show more detail information

Use "mcc api [command] --help" for more information about a command.
```

## api subcommand examples
Simple usage examples for api subcommand.

```
./mcc api --help
./mcc api --list
./mcc api --service spider --list
./mcc api --service spider --action ListCloudOS
./mcc api --service spider --action GetCloudDriver --pathParam driver_name:AWS
./mcc api --service spider --action GetRegionZone --pathParam region_name:ap-northeast-3 --queryString ConnectionName:aws-config01
```