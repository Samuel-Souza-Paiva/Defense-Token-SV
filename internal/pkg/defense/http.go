package defense

import (
	"crypto/tls"
	"errors"
	"io"
	"net/http"
	"strings"
	"time"
)

type api struct {
  Host               string
  Port               string
  Scheme             string
  InsecureSkipVerify bool
  Token              string
  headers            map[string]string
}

func (defense *api) Get(path string) (uint, []byte, error) {
  code, resData, err := defense.req(http.MethodGet, path, nil)
  if err != nil {
    return 0, resData, errors.New("error on GET request to defense\n" + err.Error())
  }

  return uint(code), resData, nil
}

func (defense *api) Post(path string, data any) (uint, []byte, error) {
  code, resData, err := defense.req(http.MethodPost, path, data)
  if err != nil {
    return 0, resData, errors.New("error on POST request to defense\n" + err.Error())
  }

  return uint(code), resData, nil
}

func (defense *api) Put(path string, data any) (uint, []byte, error) {
  code, resData, err := defense.req(http.MethodPut, path, data)
  if err != nil {
    return 0, resData, errors.New("error on PUT request to defense\n" + err.Error())
  }

  return uint(code), resData, nil
}

func (defense *api) Patch(path string, data any) (uint, []byte, error) {
  code, resData, err := defense.req(http.MethodPatch, path, data)
  if err != nil {
    return 0, resData, errors.New("error on PATCH request to defense\n" + err.Error())
  }

  return uint(code), resData, nil
}

func (defense *api) Delete(path string, data any) (uint, []byte, error) {
  code, resData, err := defense.req(http.MethodDelete, path, data)
  if err != nil {
    return 0, resData, errors.New("error on DELETE request to defense\n" + err.Error())
  }

  return uint(code), resData, nil
}

func (defense *api) createHeaders(token string) {
  headers := make(map[string]string)

  headers["Content-Type"] = "application/json"

  if token != "" {
    headers["X-Subject-Token"] = token
  }

  defense.headers = headers
}

func (defense *api) createUrl(scheme, host, port, path string) string {
  scheme = strings.TrimSpace(scheme)
  if scheme == "" {
    scheme = "https"
  }
  return scheme + "://" + host + ":" + port + path
}

func (defense *api) req(method, path string, data any) (int, []byte, error) {
  url := defense.createUrl(defense.Scheme, defense.Host, defense.Port, path)
  defense.createHeaders(defense.Token)

  req, err := http.NewRequest(method, url, ToBuffer(data))
  if err != nil {
    return 0, nil, err
  }

  for key, value := range defense.headers {
    req.Header.Add(key, value)
  }

  client := &http.Client{
    Timeout: 8 * time.Second,
  }
  if strings.EqualFold(defense.Scheme, "https") {
    client.Transport = &http.Transport{
      TLSClientConfig: &tls.Config{InsecureSkipVerify: defense.InsecureSkipVerify},
    }
  }

  res, err := client.Do(req)
  if err != nil {
    return 0, nil, err
  }
  defer res.Body.Close()

  resData, err := io.ReadAll(res.Body)
  if err != nil {
    return 0, nil, err
  }

  return res.StatusCode, resData, nil
}
