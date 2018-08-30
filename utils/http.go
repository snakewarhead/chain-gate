package utils

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/snakewarhead/chain-gate/errors"
)

func HTTPGet(_url string, _params map[string]string) ([]byte, *errors.AppError) {

	baseUrl, err := url.Parse(_url)

	if err != nil {
		return nil, errors.NewAppError(err, "cannot parse url", -1, nil)
	}

	params := url.Values{}

	for k, v := range _params {

		u := &url.URL{Path: k}
		k = u.String()

		u = &url.URL{Path: v}
		v = u.String()

		params.Add(k, v)
	}

	baseUrl.RawQuery = params.Encode()

	req, _ := http.NewRequest("GET", baseUrl.String(), strings.NewReader(""))

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		return nil, errors.NewAppError(err, "error trying to reach", -1, nil)
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return nil, errors.NewAppError(err, "error reading body response", -1, nil)
	}

	return body, nil
}

func HTTPPost(url string, keyValues map[string]interface{}, bytes []byte) ([]byte, *errors.AppError) {

	var err error

	if keyValues != nil {
		bytes, err = json.Marshal(&keyValues)
	}

	if err != nil {
		return nil, errors.NewAppError(err, "error marshalling params", -1, nil)
	}

	req, _ := http.NewRequest("POST", url, strings.NewReader(string(bytes)))

	req.Header.Add("content-type", "application/json")

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		return nil, errors.NewAppError(err, "error trying to reach", -1, nil)
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return nil, errors.NewAppError(err, "error reading body response", -1, nil)
	}

	return body, nil
}

func HTTPPostRawData(url string, raw string) ([]byte, *errors.AppError) {

	req, _ := http.NewRequest("POST", url, strings.NewReader(raw))

	req.Header.Add("content-type", "application/json")

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		return nil, errors.NewAppError(err, "error trying to reach ", -1, nil)
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return nil, errors.NewAppError(err, "error reading body response", -1, nil)
	}

	return body, nil
}
