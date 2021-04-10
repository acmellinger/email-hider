#!/bin/zsh
GOOS=linux go build main.go
zip function.zip main
terraform apply -var-file="secrets.tfvars"