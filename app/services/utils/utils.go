package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

func MakePostRequest(url string, headers map[string]interface{}, query map[string]interface{}, timeout time.Duration, response interface{}) error {
	client := &http.Client{
		Timeout: timeout,
	}
	jsonStr, err := json.Marshal(query)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(jsonStr)))
	if err != nil {
		return err
	}

	for k, v := range headers {
		req.Header.Set(k, fmt.Sprintf("%v", v))
	}

	res, err := client.Do(req)
	if err != nil {
		return err
	}

	body := res.Body
	defer body.Close()

	bytes, err := ioutil.ReadAll(body)
	if err != nil {
		return err
	}

	if res.StatusCode >= 400 {
		return GetRequestHTTPError{
			StatusCode: res.StatusCode,
			Body:       bytes,
		}
	}
	return json.Unmarshal(bytes, response)
}

type GetRequestHTTPError struct {
	StatusCode int
	Body       []byte
}

func (a GetRequestHTTPError) Error() string {
	return fmt.Sprintf("status: %d, body: %s", a.StatusCode, string(a.Body))
}

func MakeGetRequest(url string, headers map[string]interface{}, query map[string]string, timeout time.Duration, response interface{}) error {
	client := &http.Client{
		Timeout: timeout,
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	for k, v := range headers {
		req.Header.Set(k, fmt.Sprintf("%v", v))
	}

	q := req.URL.Query()
	for key, value := range query {
		q.Add(key, value)
	}
	req.URL.RawQuery = q.Encode()

	res, err := client.Do(req)
	if err != nil {
		return err
	}

	body := res.Body
	defer body.Close()

	bytes, err := ioutil.ReadAll(body)
	if err != nil {
		return err
	}

	if res.StatusCode >= 400 {
		return fmt.Errorf("got unexpected statusCode %d %s", res.StatusCode, req.URL)
	}

	return json.Unmarshal(bytes, response)
}
