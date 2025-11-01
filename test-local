#!/bin/bash

# Local test script for Lambda handler
# Usage: ./test-local.sh [URL] [METRIC_NAME] [REGION]

set -e

URL=${1:-"https://www.google.com"}
METRIC_NAME=${2:-"test-local"}
REGION=${3:-"us-west-2"}

echo "Testing Lambda handler locally..."
echo "URL: $URL"
echo "METRIC_NAME: $METRIC_NAME"
echo "REGION: $REGION"
echo ""

# Set environment variables
export URL=$URL
export METRIC_NAME=$METRIC_NAME
export REGION=$REGION

# Skip CloudWatch unless AWS credentials are configured
if [ -z "$AWS_ACCESS_KEY_ID" ]; then
    echo "No AWS credentials found, skipping CloudWatch (set SKIP_CLOUDWATCH=true)"
    export SKIP_CLOUDWATCH=true
else
    echo "AWS credentials found, will send metrics to CloudWatch"
    unset SKIP_CLOUDWATCH
fi

echo ""
echo "Running handler..."
go run main.go
