package main

import (
	"os"
	"testing"
)

func TestLambdaHandler(t *testing.T) {
	// mock environment variables for METRIC_NAME
	os.Setenv("METRIC_NAME", "test-metric")
	os.Setenv("URL", "https://www.google.com")
	_, err := handler()
	if err != nil {
		t.Errorf("LambdaHandler failed: %v", err)
	}

}
