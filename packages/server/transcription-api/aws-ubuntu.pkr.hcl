variable "region" {
	type=string
	default="eu-west-2"
	sensitive=false
}

variable "environment" {
	type=string
	default="uat-ninja"
	sensitive=false
}


data "amazon-parameterstore" "customer_os_api_url" {
  name = "/config/transcription_api_${var.environment}/customer_os_api_url"
  with_decryption = false
}

data "amazon-parameterstore" "customer_os_api_key" {
  name = "/config/transcription_api_${var.environment}/customer_os_api_key"
  with_decryption = true
}

data "amazon-parameterstore" "transcription_key" {
  name = "/config/transcription_api_${var.environment}/transcription_key"
  with_decryption = true
}

data "amazon-parameterstore" "vcon_api_key" {
  name = "/config/transcription_api_${var.environment}/vcon_api_key"
  with_decryption = true
}

data "amazon-parameterstore" "vcon_api_url" {
  name = "/config/transcription_api_${var.environment}/vcon_api_url"
  with_decryption = false
}

data "amazon-parameterstore" "file_store_api_key" {
  name = "/config/transcription_api_${var.environment}/file_store_api_key"
  with_decryption = true
}

data "amazon-parameterstore" "file_store_api_url" {
  name = "/config/transcription_api_${var.environment}/file_store_api_url"
  with_decryption = false
}

data "amazon-parameterstore" "replicate_api_key" {
  name = "/config/transcription_api_${var.environment}/replicate_api_key"
  with_decryption = true
}

# usage example of the data source output
locals {
  customer_os_api_url  = data.amazon-parameterstore.customer_os_api_url.value
  customer_os_api_key  = data.amazon-parameterstore.customer_os_api_key.value
  transcription_key    = data.amazon-parameterstore.transcription_key.value
  vcon_api_key         = data.amazon-parameterstore.vcon_api_key.value
  vcon_api_url         = data.amazon-parameterstore.vcon_api_url.value
  file_store_api_key   = data.amazon-parameterstore.file_store_api_key.value
  file_store_api_url   = data.amazon-parameterstore.file_store_api_url.value
  replicate_api_key   = data.amazon-parameterstore.replicate_api_key.value
}
packer {
  required_plugins {
    amazon = {
      version = ">= 0.0.2"
      source  = "github.com/hashicorp/amazon"
    }
  }
}

source "amazon-ebs" "ubuntu" {
  ami_name      = "transcription-api-ami_${var.environment}"
  instance_type = "t2.micro"
  region        = "${var.region}"
  source_ami_filter {
    filters = {
      name                = "ubuntu/images/*ubuntu-jammy-22.04-amd64-server-*"
      root-device-type    = "ebs"
      virtualization-type = "hvm"
    }
    most_recent = true
    owners      = ["099720109477"]
  }
  ssh_username = "ubuntu"
}

build {
  name    = "build-transcription-api-image"
  sources = [
    "source.amazon-ebs.ubuntu"
  ]
  provisioner "shell" {
    inline = [
      "sudo sh -c 'add-apt-repository universe && apt-get update'",
      "sudo sh -c 'apt-get install -y python3-pip python2 ffmpeg'",
      "mkdir /tmp/transcribe/",

    ]
  }

  provisioner "file" { 
    source = "awslogs"
    destination = "/tmp/transcribe/"
  }

  provisioner "file" { 
    source = "model"
    destination = "/tmp/transcribe/"
  }

  provisioner "file" { 
    source = "routes"
    destination = "/tmp/transcribe/"
  }

  provisioner "file" { 
    source = "scripts"
    destination = "/tmp/transcribe/"
  }

  provisioner "file" { 
    source = "service"
    destination = "/tmp/transcribe/"
  }

  provisioner "file" { 
    source = "transcribe"
    destination = "/tmp/transcribe/"
  }

  provisioner "file" { 
    source = "requirements.txt"
    destination = "/tmp/transcribe/"
  }

  provisioner "file" { 
    source =  "main.py"
    destination = "/tmp/transcribe/"
  }
  
  provisioner "shell" {
    inline = [
      "sudo sh -c 'mkdir -p /etc/transcription'",
      "sudo sh -c 'CUSTOMER_OS_API_KEY=\"${local.customer_os_api_key}\" CUSTOMER_OS_API_URL=\"${local.customer_os_api_url}\" TRANSCRIPTION_KEY=\"${local.transcription_key}\" VCON_API_KEY=\"${local.vcon_api_key}\" VCON_API_URL=\"${local.vcon_api_url}\" FILE_STORE_API_KEY=\"${local.file_store_api_key}\" FILE_STORE_API_URL=\"${local.file_store_api_url}\" REPLICATE_API_TOKEN=\"${local.replicate_api_key}\" /tmp/transcribe/scripts/build_env.sh'",
      "sudo sh -c 'pip3 install --no-cache-dir -r /tmp/transcribe/requirements.txt'",
      "sudo sh -c 'mkdir -p /usr/local/transcribe'",
      "sudo sh -c 'useradd -m -r -U transcribe'",
      "sudo sh -c 'cp -a /tmp/transcribe/* /usr/local/transcribe'",
      "sudo sh -c 'find /usr/local/transcribe -type f ! -name \\*.py  -delete'",
      "sudo sh -c 'chown -R transcribe:transcribe /usr/local/transcribe'",
      "sudo sh -c 'chown -R transcribe:transcribe /etc/transcription'",
      "sudo sh -c 'mv /tmp/transcribe/scripts/transcription-api.service /etc/systemd/system'",
      "sudo sh -c 'chmod 644 /etc/systemd/system/transcription-api.service'",
      "sudo sh -c 'systemctl enable transcription-api.service'",
      "sudo sh -c 'cd /tmp/; curl https://s3.amazonaws.com/aws-cloudwatch/downloads/latest/awslogs-agent-setup.py -O; chmod a+x awslogs-agent-setup.py'",
      "sudo sh -c 'cd /tmp/; python2 ./awslogs-agent-setup.py -r ${var.region} -n -c /tmp/transcribe/awslogs/awslogs.conf'",
    ]
  }
}

