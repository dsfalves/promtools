// Package to scrap metric data from Prometheus instance
package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"encoding/json"
)

func logErr(err error) {
	if err != nil {
		log.Fatal(err)
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

func (s *Scrapper) Metrics() []string {
	var out []string
	res := make(map[string] interface{})
	response, err := s.Request("label/__name__/values", nil)
	logErr(err)
	defer response.Body.Close()
	decoder := json.NewDecoder(response.Body)
	err = decoder.Decode(&res)
	if err != nil || res["status"] != "success"{
		log.Fatal(err)
	}

	switch v := res["data"].(type) {
	case []interface{}:
		for _, m := range v {
			switch name := m.(type) {
			case string:
				out = append(out, name)
			default:
				log.Fatal("Stranger element: ", name)
			}
		}
	default:
		log.Fatal("Data not a string: ", res["data"])
	}
	return out
}

func (s Scrapper) Measurements(metric string) {
	var v map[string] string
	v = make(map[string] string)
	v["query"] = metric + fmt.Sprintf("[%ds]", 60*60*24*7)
	fmt.Println(v["query"])
	response, err := s.Request("query", v)
	logErr(err)
	defer response.Body.Close()

	data := make(map[string] interface{})
	decoder := json.NewDecoder(response.Body)
	err = decoder.Decode(&data)
	if status := data["status"].(string); status != "success" {
		fmt.Println(data["error"].(string))
	}
	fmt.Println(response)
}

func main() {
	ip := "129.114.108.78"
	port := 30900
	scrapper := NewScrapper(ip, port)
	metrics := scrapper.Metrics()
	for _, metric := range metrics {
		fmt.Println(metric)
	}
	scrapper.Measurements(metrics[0])
}
