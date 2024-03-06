#!/bin/zsh
zip function.zip emailhider.go go.mod go.sum
gcloud auth application-default login
terraform init -upgrade
terraform apply -var-file="secrets.tfvars"