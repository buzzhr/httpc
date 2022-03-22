package httpc

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
)

func Request(url string, headers map[string]string, method string, body interface{}) (int, []byte, error) {
	if url == "" || method == "" {
		return 0, nil, fmt.Errorf("url and method are required")
	}

	transport := &http.Transport{DisableKeepAlives: true, TLSClientConfig: &tls.Config{MinVersion: tls.VersionTLS12}}
	var err error
	var bodyReader io.Reader
	if body != nil {
		tp := reflect.TypeOf(body)
		tps := tp.String()
		if tps == "string" {
			bodyReader = strings.NewReader(body.(string))
		} else {
			var bodyBytes []byte
			bodyBytes, err = json.Marshal(body)
			if err != nil {
				return 0, nil, err
			}
			bodyReader = strings.NewReader(string(bodyBytes))
		}
	}
	var req *http.Request
	var res *http.Response
	req, err = http.NewRequest(strings.ToUpper(method), url, bodyReader)
	if err != nil {
		return 0, nil, err
	}
	setRequestHeader(headers, req)
	res, err = transport.RoundTrip(req)
	if err != nil {
		return 0, nil, err
	}

	if transport != nil {
		defer transport.CloseIdleConnections()
	}
	if res != nil && res.Body != nil {
		defer res.Body.Close()
	}
	var resData []byte
	resData, err = ioutil.ReadAll(res.Body)
	return res.StatusCode, resData, err
}
func setRequestHeader(headers map[string]string, req *http.Request) {
	if headers == nil {
		return
	}
	for k, v := range headers {
		if k == "BasicAuth" {
			usernamePwdPair := strings.Split(v, "::")
			if len(usernamePwdPair) == 2 {
				u := usernamePwdPair[0]
				p := usernamePwdPair[1]
				req.SetBasicAuth(u, p)
			}
			continue
		}
		req.Header.Set(k, v)
	}
}
