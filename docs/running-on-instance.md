# Running on Single Instance Guide

This document provides step-by-step instructions to set up and execute the `mc-admin-cli` for the MCMP (Multi-Cloud Management Platform) environment on an instance. Please follow the steps closely to ensure proper installation and configuration.

This guide covers the necessary preparations for deploying the MCMP platform on a single virtual machine (VM).

---

## ✅ Prerequisites

Ensure you have **sudo** privileges and access to the **VM instance** where you intend to set up the MCMP platform. The guide covers installing Docker, setting up necessary directories, cloning repositories, and initializing credentials and IAM (Identity and Access Management).

To enable full functionality, open your firewall or security group to allow all traffic.

### 1) Provisioning the VM

To run the entire platform on a single instance using mc-admin-cli, the following minimum requirements must be met:

- OS: Ubuntu 20.04 or higher (22.04 recommended)
- Instance Specifications: t2.xlarge (4vCPU, 16GiB RAM)
- EBS Storage Size: 15 GB
  - **31 images size**: approximately 10GB
  - **additional space for vm**: up to you!

### 2) Configuring Firewall Rules (Security Groups)

- Allow all TCP communication across all ports and from all IP addresses. (Note: This configuration will be improved in the future to address security concerns.)
- Ensure port 22 is open for SSH access.

---

<br>
<br>
<br>

## 🚀 Running MCMP on Instance Guide

The following instructions **should be executed on the provisioned VM.**

> If the key pair is correctly stored on your local host, you can connect to the instance via SSH. Otherwise, you may use the web terminal provided by AWS or other cloud consoles to establish an SSH connection and access the instance’s terminal before proceeding with the next steps.

```bash
ssh -i <YOUR_KEY_PAIR_DIRECTORY> <VM_USER_NAME>@<VM_PUBLIC_IP>
```

<br>

## Step 1: Install Docker

Install Docker and necessary dependencies to facilitate containerized operations:

```bash
sudo apt-get install -y apt-transport-https ca-certificates curl gnupg-agent software-properties-common
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo apt-key add -
sudo add-apt-repository "deb [arch=amd64] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable"
sudo apt-get update
sudo apt-get install -y docker-ce docker-ce-cli docker-compose-plugin
```

### Add Docker Permissions for Current User

```bash
sudo usermod -aG docker $(whoami)
newgrp docker
docker ps  # Verifies Docker installation by listing running containers
```

<br>

## Step 2: Create Required Directories and Credentials File

Set up a working directory and initialize the necessary credentials configuration:

```bash
mkdir -p ~/workspace
mkdir -p ~/.cloud-barista
```

<br>

## Step 3: Clone Required Git Repositories

Navigate to the workspace directory and clone the necessary repositories:

```bash
cd ~/workspace
git clone --branch v0.10.0 https://github.com/cloud-barista/cb-tumblebug.git
git clone --branch v0.3.2 https://github.com/m-cmp/mc-admin-cli.git
git clone --branch v0.3.0 https://github.com/m-cmp/mc-iam-manager.git
```

## Step 4: Run mc-admin-cli

Execute `mc-admin-cli` to initialize the MCMP infrastructure:

```bash
cd ~/workspace/mc-admin-cli/bin
./mcc infra run -d
```

Wait for Services to Initialize
Allow some time for all services to start and reach a healthy state. You may verify health checks for each service if required.
**It will take approximately 5 mins.**

> ❗Following Two steps (step 5 and 6) is crutial steps for using MCMP normally.

## ❗ Step 5: Initialize Credentials ⭐

If you're setting up a new instance of Tumblebug, follow these initialization steps (otherwise, skip this section if you already have Tumblebug set up).

For more information, refer to the [Tumblebug initialization guide.](https://github.com/cloud-barista/cb-tumblebug?tab=readme-ov-file#3-initialize-cb-tumblebug-to-configure-multi-cloud-info)

## ❗ Step 6: Initialize MC-IAM-MANAGER ⭐

Install jq, a lightweight JSON processor, and set up IAM configurations:

```bash
sudo apt-get install -y jq
```

### Configure MC-IAM-MANAGER Environment Variables

Edit .env to configure IAM service properties:

```bash
cd ~/workspace/mc-iam-manager/scripts/init
cp .env.initsample .env
sed -i 's|MCIAMMANAGER_HOST=https://MCIAMMANAGER_HOST|MCIAMMANAGER_HOST=http://127.0.0.1:5005|' .env
sed -i 's|MCIAMMANAGER_PLATFORMADMIN_ID=|MCIAMMANAGER_PLATFORMADMIN_ID=mcmpadmin|' .env
sed -i 's|MCIAMMANAGER_PLATFORMADMIN_PASSWORD=|MCIAMMANAGER_PLATFORMADMIN_PASSWORD=mcmpAdminPassword#@!|' .env
```

### Finalize MC-IAM-MANAGER Initialization

Execute the IAM auto-initialization script:

```bash
./initauto.sh -f


# Login successful.
# Role created successfully: admin
# First Role ID saved as ROLE_ID: d291b29a-2d36-41e7-b50d-cac7df88dde3
# Role created successfully: operator
# Role created successfully: viewer
# Role created successfully: billadmin
# Role created successfully: billviewer
# Downloaded mcwebconsoleMenu.yaml successfully.
# Uploaded mcwebconsoleMenu.yaml successfully.
# Downloaded permission.csv successfully.
# Uploaded permission.csv successfully.
# {"id":"95d97e4e-5882-4782-b6ce-d17bf949a43f","name":"workspace1","description":"workspace1 desc","created_at":"2024-10-31T07:27:36.533011Z","updated_at":"2024-10-31T07:27:36.533011Z"} 200
# Workspace created successfully. ID: 95d97e4e-5882-4782-b6ce-d17bf949a43f
# {"id":"ed17c0a7-2317-40c3-b4f1-6b7f224fc681","ns_id":"project1","name":"project1","description":"project1 desc","created_at":"2024-10-31T07:27:36.572527Z","updated_at":"2024-10-31T07:27:36.572527Z"} 200
# Project created successfully. ID: ed17c0a7-2317-40c3-b4f1-6b7f224fc681
# {"workspace":{"id":"95d97e4e-5882-4782-b6ce-d17bf949a43f","name":"workspace1","description":"workspace1 desc","created_at":"2024-10-31T07:27:36.533011Z","updated_at":"2024-10-31T07:27:36.533011Z"},"projects":[{"id":"ed17c0a7-2317-40c3-b4f1-6b7f224fc681","ns_id":"project1","name":"project1","description":"project1 desc","created_at":"2024-10-31T07:27:36.572527Z","updated_at":"2024-10-31T07:27:36.572527Z"}]} 200
# Project Worksapce mapping created successfully.
# User role assigned to workspace successfully
```

### Add user for Console user

Execute the IAM auto-add-user script:

```bash
./add_demo_user.sh -f

# Login successful.
# {"id":"eric","password":"changeMe!","firstName":"Eric","lastName":"Schmidt","email":"eric@mcmpemail.com","description":"ericDesc"} {"id":"elon","password":"changeMe!","firstName":"Elon","lastName":"Musk","email":"elon@mcmpemail.com","description":"elonDesc"} {"id":"jeffrey","password":"changeMe!","firstName":"Jeffrey","lastName":"PrestonBezos","email":"jeffrey@mcmpemail.com","description":"jeffreyDesc"} {"id":"gates","password":"changeMe!","firstName":"Bill","lastName":"Gates","email":"gates@mcmpemail.com","description":"gatesDesc"}
# {"id":"eric","password":"changeMe!","firstName":"Eric","lastName":"Schmidt","email":"eric@mcmpemail.com","description":"ericDesc"}
# User created successfully: eric
# User activated successfully: eric
# {"id":"elon","password":"changeMe!","firstName":"Elon","lastName":"Musk","email":"elon@mcmpemail.com","description":"elonDesc"}
# User created successfully: elon
# User activated successfully: elon
# {"id":"jeffrey","password":"changeMe!","firstName":"Jeffrey","lastName":"PrestonBezos","email":"jeffrey@mcmpemail.com","description":"jeffreyDesc"}
# User created successfully: jeffrey
# User activated successfully: jeffrey
# {"id":"gates","password":"changeMe!","firstName":"Bill","lastName":"Gates","email":"gates@mcmpemail.com","description":"gatesDesc"}
# User created successfully: gates
# User activated successfully: gates
```

## Step 7: Access the MCMP Platform

Upon successful initialization, access the MCMP platform via:

```bash
http://{vm-public-ip}:3001
```

#### - initial id: mcmpadmin

#### - initial password: mcmpAdminPassword#@!

Replace {vm-public-ip} with the actual public IP of your VM instance.

This completes the setup. You are now ready to manage multi-cloud services using MCMP on your instance. Happy managing!
