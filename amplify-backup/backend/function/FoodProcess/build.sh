#!/bin/bash
cd src
GOOS=linux GOARCH=amd64 go build -o ../bin/main
