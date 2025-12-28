package plugins

import (
    "bytes"
    "encoding/json"
    "fmt"
    "io/ioutil"
    "net/http"
    "time"

    "github.com/pkg/errors"
)

type HttpExecutor struct {
    client *http.Client
}

func NewHttpExecutor() *HttpExecutor {
    return &HttpExecutor{
        client: &http.Client{
            Timeout: 30 * time.Second,
        },
    }
}

func (h *HttpExecutor) Execute(url string, method string, headers map[string]string, body interface{}) (int, []byte, error) {
    var requestBody []byte
    var err error

    if body != nil {
        requestBody, err = json.Marshal(body)
        if err != nil {
            return 0, nil, errors.Wrap(err, "failed to marshal request body")
        }
    }

    req, err := http.NewRequest(method, url, bytes.NewBuffer(requestBody))
    if err != nil {
        return 0, nil, errors.Wrap(err, "failed to create HTTP request")
    }

    for key, value := range headers {
        req.Header.Set(key, value)
    }

    resp, err := h.client.Do(req)
    if err != nil {
        return 0, nil, errors.Wrap(err, "failed to execute HTTP request")
    }
    defer resp.Body.Close()

    responseBody, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return 0, nil, errors.Wrap(err, "failed to read response body")
    }

    return resp.StatusCode, responseBody, nil
}