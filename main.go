package main

import (
	"context"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"
	"github.com/tcnksm/go-httpstat"
)

func handler() (string, error) {
	ctx := context.Background()

	url := os.Getenv("URL")
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}

	var result httpstat.Result
	httpCtx := httpstat.WithHTTPStat(req.Context(), &result)
	req = req.WithContext(httpCtx)

	client := http.DefaultClient
	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	if _, err := io.Copy(io.Discard, res.Body); err != nil {
		log.Fatal(err)
	}
	res.Body.Close()
	result.End(time.Now())

	metricName := os.Getenv("METRIC_NAME")
	namespace := "AppHealth"
	dimension := types.Dimension{
		Name:  aws.String("ServiceName"),
		Value: aws.String(metricName),
	}

	// Prepare metric data
	var metricData []types.MetricDatum

	if res.StatusCode == 200 {
		totalDuration := int(result.Total(time.Now()) / time.Millisecond)
		metricData = []types.MetricDatum{
			{
				MetricName: aws.String("is-up"),
				Value:      aws.Float64(1),
				Unit:       types.StandardUnitCount,
				Dimensions: []types.Dimension{dimension},
			},
			{
				MetricName: aws.String("dns-lookup"),
				Value:      aws.Float64(float64(result.DNSLookup / time.Millisecond)),
				Unit:       types.StandardUnitMilliseconds,
				Dimensions: []types.Dimension{dimension},
			},
			{
				MetricName: aws.String("tcp-connection"),
				Value:      aws.Float64(float64(result.TCPConnection / time.Millisecond)),
				Unit:       types.StandardUnitMilliseconds,
				Dimensions: []types.Dimension{dimension},
			},
			{
				MetricName: aws.String("tls-handshake"),
				Value:      aws.Float64(float64(result.TLSHandshake / time.Millisecond)),
				Unit:       types.StandardUnitMilliseconds,
				Dimensions: []types.Dimension{dimension},
			},
			{
				MetricName: aws.String("server-processing"),
				Value:      aws.Float64(float64(result.ServerProcessing / time.Millisecond)),
				Unit:       types.StandardUnitMilliseconds,
				Dimensions: []types.Dimension{dimension},
			},
			{
				MetricName: aws.String("total"),
				Value:      aws.Float64(float64(totalDuration)),
				Unit:       types.StandardUnitMilliseconds,
				Dimensions: []types.Dimension{dimension},
			},
		}
	} else {
		metricData = []types.MetricDatum{
			{
				MetricName: aws.String("is-up"),
				Value:      aws.Float64(0),
				Unit:       types.StandardUnitCount,
				Dimensions: []types.Dimension{dimension},
			},
		}
	}

	// Put metrics to CloudWatch
	// Skip CloudWatch calls in test mode (when AWS credentials aren't configured)
	if os.Getenv("SKIP_CLOUDWATCH") == "" {
		// Get region from REGION env var, fallback to AWS_REGION
		region := os.Getenv("REGION")
		if region == "" {
			region = os.Getenv("AWS_REGION")
		}
		if region == "" {
			region = "us-west-2" // default fallback
		}

		// Initialize CloudWatch client
		cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
		if err != nil {
			log.Fatalf("unable to load AWS config: %v", err)
		}
		cwClient := cloudwatch.NewFromConfig(cfg)

		_, err = cwClient.PutMetricData(ctx, &cloudwatch.PutMetricDataInput{
			Namespace:  aws.String(namespace),
			MetricData: metricData,
		})
		if err != nil {
			log.Fatalf("failed to put metric data: %v", err)
		}
	}

	return "", nil
}

func main() {
	lambda.Start(handler)
}
