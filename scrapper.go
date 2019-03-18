// Package to scrap metric data from Prometheus instance
package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"encoding/json"
)

type Scrapper struct {
	root string
}

func NewScrapper(host string, port int) *Scrapper {
	root := fmt.Sprintf("http://%s:%d/api/v1", host, port)
	return &Scrapper{root: root}
}

func (s *Scrapper) Request(path string) (*http.Response, error){
	fullpath := fmt.Sprintf("%s/%s", s.root, path)
	u, err := url.Parse(fullpath)
	if err != nil {
		log.Fatal(err)
	}
	response, err := http.Get(u.String())
	return response, err
}

func (s *Scrapper) Metrics() []string {
	var out []string
	res := make(map[string] interface{})
	response, err := s.Request("label/__name__/values")
	if err != nil {
		log.Fatal(err)
	}
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

func main() {
	ip := "129.114.108.78"
	port := 30900
	scrapper := NewScrapper(ip, port)
	for _, metric := range scrapper.Metrics() {
		fmt.Println(metric)
	}
}
