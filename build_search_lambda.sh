#!/bin/bash
GOOS=linux GOARCH=arm64 go build -tags lambda.norpc -o bootstrap cmd/search_lambda/main.go
zip bin/search_lambda.zip bootstrap
rm bootstrap
