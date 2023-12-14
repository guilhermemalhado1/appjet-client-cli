#!/bin/bash

# This script builds the appjet-client-cli executable

# Get today's date and time in the format MMDDYYYYHHMM
datetime=$(date +"%m%d%Y%H%M")

# Create a build folder with the appended date and time
go build -o "client-cli-build-$datetime/appjet-client-cli" ./cmd
