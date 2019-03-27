package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
  // "github.com/aws/aws-sdk-go/aws"
  // "github.com/aws/aws-sdk-go/aws/session"
	"github.com/tcnksm/go-httpstat"
)

func main() {
	url := os.Getenv("URL")
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
  fmt.Printf("%+v\n", result.DNSLookup)
  fmt.Printf("%+v\n", result.TCPConnection)
  fmt.Printf("%+v\n", result.TLSHandshake)
  fmt.Printf("%+v\n", result.ServerProcessing)
  fmt.Printf("%+v\n", result.Total(time.Now()))
	// fmt.Printf("%+v\n", result)
}
