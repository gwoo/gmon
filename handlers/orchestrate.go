package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type OrHandler struct {
	Orchestrate `json:"orchestrate"`
}

type Orchestrate struct {
	Url    string
	ApiKey string `json:"api_key"`
}

func (handler *OrHandler) Config(config []byte) {
	json.Unmarshal(config, handler)
}

func (handler *OrHandler) Store(results []*Metric) bool {
	if handler.Url == "" {
		log.Println("Orchestrate url is missing.")
		return false
	}
	if handler.ApiKey == "" {
		log.Println("Orchestrate api_key is missing.")
		return false
	}
	for _, m := range results {
		_, err := handler.send(m)
		if err != nil {
			log.Println(err.Error())
			continue
		}
	}
	return false
}

func (handler *OrHandler) send(m *Metric) ([]byte, error) {
	url := fmt.Sprintf(
		"%s/%s/events/%s?timestamp=%d",
		strings.TrimRight(handler.Url, "/"), m.Name, m.Type, m.Time.Unix(),
	)
	data := map[string]interface{}{
		"timestamp": m.Time,
		"host":      m.Host,
		"type":      m.Name,
		"name":      m.Type,
		"tags":      m.Tags,
		"message":   m.Message,
		"value":     m.Value,
		"duration":  m.Duration,
	}
	d, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	json := bytes.NewBuffer(d)
	req, err := http.NewRequest("PUT", url, json)
	if err != nil {
		return nil, err
	}
	req.Close = true
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(handler.ApiKey, "")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func init() {
	RegisterHandler("orchestrate", &OrHandler{})
}
