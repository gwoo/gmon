// Copyright 2013 GWoo. All rights reserved.
// The BSD License http://opensource.org/licenses/bsd-license.php.
package handlers

import (
	"io/ioutil"
	"strconv"
	"strings"
	"time"
)

// Read conf from path
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

//Parse the ouput from a script
func (m *Metric) Parse(output string) {
	message := ""
	parts := strings.SplitN(output, "|", 4)
	tags := make([]string, 0)
	if len(parts) > 2 {
		message = strings.TrimSpace(parts[2])
		tags = strings.SplitAfter(parts[3], " ")
	}
	value, _ := strconv.ParseFloat(strings.TrimSpace(parts[1]), 64)
	m.Type = strings.TrimSpace(parts[0])
	m.Value = value
	m.Message = message
	m.Tags = tags
}

type Handler interface {
	Config(config []byte)
	Store([]*Metric) bool
}

type Map map[string]Handler

// Convert a map of handlers to csv
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

//Register a new handler
//Should be called from init method of the handler.
func RegisterHandler(name string, handler Handler) {
	if handler == nil {
		panic("Handler is nil.")
	}
	Handlers[name] = handler
}
