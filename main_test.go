package main

import (
	"os"
	"testing"
)

func TestLambdaHandler(t *testing.T) {
	// mock environment variables
	os.Setenv("METRIC_NAME", "test-metric")
	os.Setenv("URL", "https://www.google.com")
	os.Setenv("REGION", "us-west-2")
	// Skip CloudWatch calls in tests
	os.Setenv("SKIP_CLOUDWATCH", "true")
	
	_, err := handler()
	if err != nil {
		t.Errorf("LambdaHandler failed: %v", err)
	}
}
