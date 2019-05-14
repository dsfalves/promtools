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

type SeriesResponse struct {
	*Response
	Data struct {
		ResultType string
		Result []Measurements `json:"result"`
	}
}

type Measurements struct {
	Metric map[string] string
	Values []Entry `json:"values"`
}

type SeriesData struct {
	ResultType string
	Result []Measurements
}

type MetricsResponse struct {
	*Response
	Data []string
}

type Response struct {
	Status string
	ErrorType string
	Error string
	Warnings []string
}

func (e *Entry) UnmarshalJSON(buf []byte) error {
	tmp := []interface{}{&e.Time, &e.Value}
	wantLen := len(tmp)
	if err := json.Unmarshal(buf, &tmp); err != nil {
		return err
	}
	if l := len(tmp); l != wantLen {
		return fmt.Errorf("wrong number of fields in Entry: %d != %d", l, wantLen)
	}
	return nil
}

func SeriesRequest(address, metric, span string) ([]Measurements, error) {
	var res SeriesResponse
	fullpath := fmt.Sprintf("%s/api/v1/query", address)
	u, err := url.Parse(fullpath)
	if err != nil {
		return nil, err
	}
	v := u.Query()
	v.Set("query", fmt.Sprintf("%s[%s]", metric, span))
	u.RawQuery = v.Encode()
	response, err := http.Get(u.String())
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	decoder := json.NewDecoder(response.Body)
	if err = decoder.Decode(&res); err != nil {
		return nil, err
	} else if res.Status != "success" {
		return nil, errors.New(fmt.Sprintf("Error in response: %s", res.Error))
	}
	return res.Data.Result, nil
}

func MetricsRequest(address string) ([]string, error) {
	var res MetricsResponse
	fullpath := fmt.Sprintf("%s/api/v1/label/__name__/values", address)
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
	if err = decoder.Decode(&res); err != nil {
		return nil, err
	} else if res.Status != "success" {
		return nil, errors.New(fmt.Sprintf("Error in response: %s", res.Error))
	}
	return res.Data, nil
}
