#!/bin/zsh
zip function.zip emailhider.go go.mod go.sum
terraform apply -var-file="secrets.tfvars"