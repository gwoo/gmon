package handlers

import (
	"fmt"
)

type StdoutHandler struct{}

func (handler *StdoutHandler) Config(config []byte) {}

func (handler *StdoutHandler) Store(results []*Metric) bool {
	value := ""
	for _, m := range results {
		value += fmt.Sprintf("%s\t%g\t%s\t%s\t%s\t%v\n",
			m.Host+"/"+m.Name+"/"+m.Type, m.Value,
			m.Time, m.Duration, m.Message, m.Tags)
	}
	println(value)
	return true
}

func init() {
	RegisterHandler("stdout", &StdoutHandler{})
}
