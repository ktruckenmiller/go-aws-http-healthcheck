package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/prozz/aws-embedded-metrics-golang/emf"
	"github.com/tcnksm/go-httpstat"
)

func handler() (string, error) {

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

	m := emf.New(
		emf.WithLogGroup("lambda-emf-metrics"),
	).Namespace("AppHealth").Dimension("ServiceName", os.Getenv("METRIC_NAME"))
	defer m.Log()
	// datums = append(datums, serializeDatum("is-up", 1))

	if res.StatusCode == 200 {
		m.MetricAs("is-up", 1, emf.Count)
		m.MetricAs("dns-lookup", int(result.DNSLookup/time.Millisecond), emf.Milliseconds)
		m.MetricAs("tcp-connection", int(result.TCPConnection/time.Millisecond), emf.Milliseconds)
		m.MetricAs("tls-handshake", int(result.TLSHandshake/time.Millisecond), emf.Milliseconds)
		m.MetricAs("server-processing", int(result.ServerProcessing/time.Millisecond), emf.Milliseconds)
		m.MetricAs("total", int(result.Total(time.Now())/time.Millisecond), emf.Milliseconds)
	} else {
		m.MetricAs("is-up", 0, emf.Count)
	}

	return "", nil
}

func main() {
	lambda.Start(handler)
}
