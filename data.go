package promtools

import (
	"errors"
	"fmt"
	"encoding/json"
	"net/http"
	"net/url"
)

type Entry struct {
	Time float64
	Value string
}

type Series []Entry

type SeriesResponse struct {
	Status string
	Data SeriesData
	ErrorType string
	Error string
	Warnings []string
}

type Measurements struct {
	Metric map[string] string
	Values Series
}

type SeriesData struct {
	ResultType string
	Result []Measurements
}

type MetricsResponse struct {
	Status string
	Data []string
	ErrorType string
	Error string
	Warnings []string
}

func MetricsRequest(address string) ([]string, error) {
	res := new(MetricsResponse)
	fullpath := fmt.Sprintf("http://%s/api/v1/label/__name__/values", address)
	u, err := url.Parse(fullpath)
	if err != nil {
		return nil, err
	}
	response, err := http.Get(u.String())
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	decoder := json.NewDecoder(response.Body)
	if err = decoder.Decode(res); err != nil {
		return nil, err
	} else if res.Status != "success" {
		return nil, errors.New(fmt.Sprintf("Error in response: %s", res.Error))
	}
	return res.Data, nil
}
