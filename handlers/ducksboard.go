package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type DuckHandler struct {
	Ducksboard `json:"ducksboard"`
}

type Ducksboard struct {
	ApiKey string `json:"api_key"`
}

func (handler *DuckHandler) Config(config []byte) {
	json.Unmarshal(config, handler)
}

func (handler *DuckHandler) Store(results []*Metric) bool {
	if handler.ApiKey == "" {
		log.Println("Ducksboard api_key is missing.")
		return false
	}
	for _, m := range results {
		_, err := handler.send(m)
		if err != nil {
			log.Println(err.Error())
			return false
		}
	}
	return true
}

func (handler *DuckHandler) send(m *Metric) ([]byte, error) {
	url := fmt.Sprintf(
		"%s/v/%s-%s-%s",
		"https://push.ducksboard.com", m.Host, m.Name, m.Type,
	)
	log.Println(url)
	data := map[string]interface{}{
		"timestamp": m.Time.Unix(),
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
	req, err := http.NewRequest("POST", url, json)
	if err != nil {
		return nil, err
	}
	req.Close = true
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(handler.ApiKey, "unused")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, errors.New("Bad Response from API.\n" + resp.Status + "\n" + string(body))
	}
	return body, nil
}

func init() {
	RegisterHandler("ducksboard", &DuckHandler{})
}
