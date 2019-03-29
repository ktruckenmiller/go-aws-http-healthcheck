package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
  "github.com/aws/aws-sdk-go/aws"
  "github.com/aws/aws-sdk-go/aws/session"
  "github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/tcnksm/go-httpstat"
)

func serializeDatum(metricname string, value float64) *cloudwatch.MetricDatum {
  return &cloudwatch.MetricDatum{
      MetricName: aws.String(os.Getenv("METRIC_NAME")),
      Timestamp:  aws.Time(time.Now()),
      Value:      aws.Float64(value),
      Unit:       aws.String("Milliseconds"),
      Dimensions: []*cloudwatch.Dimension{
        &cloudwatch.Dimension{
          Name:  aws.String("Connection"),
          Value: aws.String(metricname),
        },
      },
  }
}


func main() {
  var datums []*cloudwatch.MetricDatum
	url := os.Getenv("URL")
  region := os.Getenv("REGION")
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)

	}

	var result httpstat.Result
	ctx := httpstat.WithHTTPStat(req.Context(), &result)
	req = req.WithContext(ctx)

	client := http.DefaultClient
	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	if _, err := io.Copy(ioutil.Discard, res.Body); err != nil {
		log.Fatal(err)
	}
	res.Body.Close()
	result.End(time.Now())

  // datums = append(datums, serializeDatum("is-up", 1))

  if res.StatusCode == 200 {
    datums = append(datums, serializeDatum("is-up", 1))
    datums = append(datums, serializeDatum("dns-lookup", float64(result.DNSLookup / time.Millisecond)))
    datums = append(datums, serializeDatum("tcp-connection", float64(result.TCPConnection / time.Millisecond)))
    datums = append(datums, serializeDatum("tls-handshake", float64(result.TLSHandshake / time.Millisecond)))
    datums = append(datums, serializeDatum("server-processing", float64(result.ServerProcessing / time.Millisecond)))
    datums = append(datums, serializeDatum("total", float64(result.Total(time.Now()) / time.Millisecond )))
  } else {
    datums = append(datums, serializeDatum("is-up", 0))
  }

  sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(region),
	}))

  svc := cloudwatch.New(sess)

  fmt.Printf("%v", datums)
  _, err = svc.PutMetricData(&cloudwatch.PutMetricDataInput{
    MetricData: datums,
    Namespace: aws.String("AppHealth"),
  })

  if err != nil {
		log.Fatal(err)
	}
}
