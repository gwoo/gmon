// Copyright 2013 GWoo. All rights reserved.
// The BSD License http://opensource.org/licenses/bsd-license.php.
package main

import (
	"flag"
	h "github.com/gwoo/gmon/handlers"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

var wd, _ = os.Getwd()
var conf = flag.String("conf", wd+"/gmon.json", "Path to config file.")
var handlers = flag.String("handlers", "stdout", "Comma seperate list of handlers. ex: elasticseach,stdout.")
var interval = flag.String("interval", "5m", "Time between each check. Examples: 10s, 5m, 1h")
var path = flag.String("path", wd+"/scripts/", "Path to scripts directory.")

func main() {
	flag.Parse()
	config, err := h.Config(*conf)
	if err != nil {
		log.Printf("Config Error: %s", err.Error())
		return
	}
	host := hostname()
	log.Printf("Host: %s", host)
	log.Printf("Running scripts in %s", *path)
	log.Printf("Available Handlers: %s", h.Handlers.String())
	log.Printf("Using handlers %s", *handlers)

	scripts := scripts()
	cs := make(chan []*h.Metric)
	for {
		for _, s := range scripts {
			script, name := scriptname(s)
			if script == "" {
				continue
			}
			m := h.Metric{Host: host, Name: name, Script: script}
			go func(cs chan []*h.Metric, m h.Metric) {
				Exec(cs, m)
			}(cs, m)

			go func(cs chan []*h.Metric, config []byte, handlers *string) {
				Send(cs, config, handlers)
			}(cs, config, handlers)
		}
		t, _ := time.ParseDuration(*interval)
		time.Sleep(t)
	}
}

// Scripts should return `name value message\n`
func Exec(pub chan []*h.Metric, m h.Metric) {
	start := time.Now()
	defer func(name string) {
		if x := recover(); x != nil {
			log.Printf("%s %s\n", name, x)
		}
	}(m.Name)
	c := exec.Command(m.Script)
	output, err := c.CombinedOutput()
	if err != nil {
		log.Printf("Error running %s: %s", m.Name, err)
		return
	}
	if string(output) == "" {
		log.Printf("Error running %s: %s", m.Name, "no response")
		return
	}
	end := time.Now()
	duration := end.Sub(start)
	results := strings.Split(strings.Trim(
		strings.NewReplacer("\r", "").Replace(string(output)), "\n"), "\n")
	responses := make([]*h.Metric, 0)
	for _, r := range results {
		m.Time = end
		m.Duration = duration
		responses = append(responses, response(r, m))
	}
	pub <- responses
}

// Send the collected metrics to registered handlers
func Send(sub chan []*h.Metric, config []byte, handlers *string) {
	results := <-sub
	hs := strings.Split(*handlers, ",")
	for _, name := range hs {
		if _, ok := h.Handlers[name]; ok {
			h.Handlers[name].Config(config)
			if !h.Handlers[name].Store(results) {
				log.Printf("%s could not store results.", name)
			}
		}
	}
}

// Get the hostname to add to Metric.
func hostname() string {
	c := exec.Command("hostname")
	output, err := c.CombinedOutput()
	if err != nil {
		log.Printf("Could not get hostname.")
		return ""
	}
	host := strings.TrimSpace(string(output))
	return host
}

// Get the metrics that should be run based on the "path" flag.
func scripts() []string {
	scripts, err := filepath.Glob(*path + "/*")
	if err != nil {
		log.Panic(err.Error())
	}
	return scripts
}

// Get the full script path and short name of the metric.
func scriptname(s string) (script string, name string) {
	name = strings.Replace(s, *path, "", 1)
	if strings.Index(name, ".") == 0 {
		return "", ""
	}
	script, err := filepath.EvalSymlinks(s)
	if err != nil {
		log.Print(err.Error())
		return "", ""
	}
	name = strings.Replace(name, filepath.Ext(script), "", 1)
	return script, name
}

//Convert string response to a Metric
func response(r string, m h.Metric) *h.Metric {
	parts := strings.SplitN(r, "|", 4)
	message := m.Name + " is running."
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
	return &m
}
