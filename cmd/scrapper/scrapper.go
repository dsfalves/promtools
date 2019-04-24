// Package to scrap metric data from Prometheus instance
package main

import (
	"os"
	"flag"
	"io/ioutil"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"github.com/dsfalves/promtools"
)

type TimePoint struct {
	time float64
	measurement string
}

type TimeSeries []TimePoint

type Response struct {
	metrics map[string] string
	values TimeSeries
}

func Decode(raw map[string] interface{}) *Response {
	response := Response{
		metrics: make(map[string] string),
	}
	response.metrics = raw["metric"].(map[string] string)

	return &response
}

func logErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func checkStatus(response map[string] interface{}) {
	status, ok := response["status"].(string)
	if !ok {
		log.Fatal("Response missing status field: ", response)
	}
	if status != "success" {
		log.Fatal("Received error: ", response["error"].(string))
	}
}

type Scrapper struct {
	root string
}

func NewScrapper(host string, port int) *Scrapper {
	root := fmt.Sprintf("http://%s:%d/api/v1", host, port)
	return &Scrapper{root: root}
}

func (s *Scrapper) Request(path string, v map[string]string) (*http.Response, error){
	fullpath := fmt.Sprintf("%s/%s", s.root, path)
	u, err := url.Parse(fullpath)
	logErr(err)
	if v != nil {
		q := make(url.Values)
		for k, m := range v {
			q.Add(k, m)
		}
		u.RawQuery = q.Encode()
	}
	response, err := http.Get(u.String())
	return response, err
}

func (s Scrapper) Measurements(metric string) {
	var v map[string] string
	v = make(map[string] string)
	v["query"] = metric + fmt.Sprintf("[%dh]", 24*7)
	response, err := s.Request("query", v)
	logErr(err)
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	logErr(err)
	ioutil.WriteFile(metric + ".json", body, 0644)
}

func main() {
	l := log.New(os.Stderr, "", 0)
	ip := flag.String("address", "", "address for Prometheus") //"129.114.108.78"
	port := flag.Int("port", 30900, "port for Prometheus")
	flag.Parse()
	if *ip == "" {
		l.Fatal("The --address argument is required")
	}
	path := fmt.Sprintf("http://%s:%d", *ip, *port)
	scrapper := NewScrapper(*ip, *port)
	metrics, err := promtools.MetricsRequest(path)
	if err != nil {
		log.Fatal(err)
	}
	for _, metric := range metrics {
		fmt.Println(metric)
		scrapper.Measurements(metric)
	}
}
