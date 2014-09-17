// Copyright 2013 GWoo. All rights reserved.
// The BSD License http://opensource.org/licenses/bsd-license.php.
package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/gwoo/greq"
)

type EsHandler struct {
	Elasticsearch `json:"elasticsearch"`
}

type Elasticsearch struct {
	Host string
	Port string
}

func (handler *EsHandler) Config(config []byte) {
	json.Unmarshal(config, handler)
}

func (es *Elasticsearch) Addr() string {
	return fmt.Sprintf("%s:%s", es.Host, es.Port)
}

// Uses something similar to logstash.
// {
//   "@timestamp": "2012-12-18T01:01:46.092538Z".
//   "tags": [ "kernel", "dmesg" ]
//   "type": "syslog"
//   "message": "usb 3-1.2: USB disconnect, device number 4",
//   "path": "/var/log/messages",
//   "host": "pork.home"
// }
func (handler *EsHandler) Store(results []*Metric) bool {
	es := greq.New(handler.Addr(), true)
	records := make([]interface{}, 0)
	for _, m := range results {
		r := map[string]interface{}{
			"@timestamp": m.Time,
			"host":       m.Host,
			"type":       m.Name,
			"name":       m.Type,
			"tags":       m.Tags,
			"message":    m.Message,
			"value":      m.Value,
			"duration":   m.Duration,
		}
		year, month, day := time.Now().Date()
		index := fmt.Sprintf("gmon-%d.%02d.%02d", year, month, day)
		url := fmt.Sprintf("/%s/%s", index, m.Name)
		body, _, err := es.Post(url, r)
		if err != nil {
			log.Println(err.Error())
			continue
		}
		log.Println(string(body))
		records = append(records, r)
	}
	b, _ := json.Marshal(records)
	if b != nil {
		return true
	}
	return false
}

func init() {
	RegisterHandler("elasticsearch", &EsHandler{})
}
