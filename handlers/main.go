package handlers

import (
	"io/ioutil"
	"time"
)

func Config(path string) (conf []byte, err error) {
	conf, err = ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return conf, nil
}

type Metric struct {
	Host     string
	Name     string
	Script   string
	Type     string
	Tags     []string
	Value    float64
	Time     time.Time
	Duration time.Duration
	Message  string
}

type Handler interface {
	Config(config []byte)
	Store([]*Metric) bool
}

type Map map[string]Handler

func (m *Map) String() string {
	keys := ""
	i := 0
	for k, _ := range *m {
		keys = keys + k
		if i++; i < len(*m) {
			keys = keys + ","
		}
	}
	return keys
}

var Handlers = make(Map)

func RegisterHandler(name string, handler Handler) {
	if handler == nil {
		panic("Handler is nil.")
	}
	Handlers[name] = handler
}
