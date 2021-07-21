#!/usr/bin/env bash
export TF_LOG=DEBUG
terraform init
TF_VAR_api_token=${LOGZIO_API_TOKEN} terraform plan -out terraform.plan
terraform apply terraform.plan