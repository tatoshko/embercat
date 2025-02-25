package huggingface

import (
    "bytes"
    "crypto/tls"
    "encoding/json"
    "errors"
    "fmt"
    "io"
    "net/http"
    "net/url"
)

var (
    client *HuggingFaceClient
)

type HuggingFaceClient struct {
    config Config
    client *http.Client
}

type Request struct {
    Inputs string `json:"inputs"`
}

type Response struct {
    Answer string `json:"answer"`
}

func Init(config Config) {
    purl := url.URL{}
    url_proxy, _ := purl.Parse(config.Proxy)

    transport := http.Transport{}
    transport.Proxy = http.ProxyURL(url_proxy)
    transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
    transport.ProxyConnectHeader.Set("Authorization", "Bearer "+config.Token)
    client = &HuggingFaceClient{client: &http.Client{Transport: &transport}, config: config}
}

func GetClient() (*HuggingFaceClient, error) {
    if client == nil {
        return nil, errors.New("huggingface client need to be initialized")
    }
    return client, nil
}

func (hf *HuggingFaceClient) Ask(text string) (result string, err error) {
    reqBody := Request{Inputs: text}

    var jsonBody []byte
    if jsonBody, err = json.Marshal(reqBody); err != nil {
        return
    }

    var resp *http.Response
    var req *http.Request

    if req, err = http.NewRequest("POST", hf.config.Api, bytes.NewBuffer(jsonBody)); err != nil {
        return
    }

    req.Header.Set("Authorization", "Bearer "+hf.config.Token)
    req.Header.Set("Content-Type", "application/json ")
    req.Header.Set("Accept", "application/json")

    if resp, err = hf.client.Do(req); err != nil {
        return
    }
    defer resp.Body.Close()

    if resp.StatusCode == http.StatusOK {
        if err = json.NewDecoder(resp.Body).Decode(&result); err != nil {
            var bodyBytes []byte
            if bodyBytes, err = io.ReadAll(resp.Body); err != nil {
                return
            }

            return "", errors.New(fmt.Sprintf("Unexpected response body: %s. Actual body: %s", err.Error(), string(bodyBytes)))
        }
    } else {
        return "", errors.New(fmt.Sprintf("Ask http code error. ACTUAL CODE [%d]", resp.StatusCode))
    }

    return
}
