package util

import (
	"bytes"
	"encoding/json"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
	u "net/url"
)

func PostUrl(url string, params map[string]string, body interface{}, headers map[string]string) ([]byte, error) {
	var (
		bodyJson []byte
		req      *http.Request
		err      error
	)

	if body != nil {
		bodyJson, err = json.Marshal(body)
		if err != nil {
			return nil, errors.Wrap(err, "json marshal request body error")
		}
	}

	req, err = http.NewRequest(http.MethodPost, url, bytes.NewBuffer(bodyJson))
	if err != nil {
		return nil, errors.Wrap(err, "NewRequest error")
	}

	contentType := "Content-type"
	req.Header.Set(contentType, headers[contentType])

	q := req.URL.Query()
	if params != nil {
		for key, val := range params {
			q.Add(key, val)
		}
		req.URL.RawQuery = q.Encode()
	}

	if headers != nil {
		for key, val := range headers {
			req.Header.Add(key, val)
		}
	}
	client := &http.Client{}

	response, err := client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "client.Do error")
	}
	defer response.Body.Close()
	d, err_ := ioutil.ReadAll(response.Body)
	if err_ != nil {
		return nil, errors.Wrap(err, "ioutil.ReadAll error")
	}

	return d, nil
}

func GetUrl(url string, data map[string]string) ([]byte, error) {
	params := u.Values{}
	ur, err := u.Parse(url)
	if err != nil {
		return nil,errors.Wrap(err,"u.Parse error")
	}

	for key, value := range data {
		params.Set(key, value)
	}

	ur.RawQuery = params.Encode()
	urlPath := ur.String()

	res, err := http.Get(urlPath)
	if err != nil {
		return nil,errors.Wrap(err,"http.Get error")
	}
	defer res.Body.Close()

	d, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil,errors.Wrap(err,"ioutil.ReadAll error")
	}
	return d,nil
}
