package handlers

import (
	"encoding/json"
	"github.com/mattbaird/elastigo/api"
	"github.com/mattbaird/elastigo/core"
	"log"
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
	api.Domain = handler.Host
	api.Port = handler.Port
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
		_, err := core.Index(true, "gmon", m.Name, "", r)
		if err != nil {
			log.Println(err.Error())
			continue
		}
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
