// Copyright 2013 GWoo. All rights reserved.
// The BSD License http://opensource.org/licenses/bsd-license.php.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/mattbaird/elastigo/api"
	"github.com/mattbaird/elastigo/core"
)

var wd, err1 = os.Getwd()
var path = flag.String("path", wd+"/scripts/", "Path to scripts directory.")
var oFormat = flag.String("output", "json", "Comma seperate list of output formats. ex: graphite,json,string")

type Check struct {
	host     string
	name     string
	tag      string
	value    float64
	time     time.Time
	duration time.Duration
	message  string
}

// Scripts should return `value message\n`
func Exec(pub chan []*Check, host string, name string, script string) {
	start := time.Now()
	defer func(name string) {
		if x := recover(); x != nil {
			log.Printf("%s %s\n", name, x)
		}
	}(name)
	c := exec.Command(script)
	output, err := c.CombinedOutput()
	if err != nil {
		log.Printf("Error running %s: %s", name, err)
		return
	}
	if string(output) == "" {
		log.Printf("Error running %s: %s", name, "no response")
		return
	}
	end := time.Now()
	duration := end.Sub(start)
	results := strings.Split(strings.Trim(
		strings.NewReplacer("\r", "").Replace(string(output)), "\n"), "\n")
	responses := make([]*Check, 0)
	for _, v := range results {
		matches := strings.SplitAfterN(v, " ", 3)
		message := name + " is running."
		if len(matches) > 2 {
			message = strings.TrimSpace(matches[2])
		}
		value, _ := strconv.ParseFloat(strings.TrimSpace(matches[1]), 64)
		responses = append(responses, &Check{
			host:     host,
			name:     name,
			tag:      strings.TrimSpace(matches[0]),
			value:    value,
			time:     end,
			message:  message,
			duration: duration,
		})
	}
	pub <- responses
}

func toString(results []*Check) string {
	value := ""
	for _, m := range results {
		value += fmt.Sprintf("%s>%s>%s>%s>%s\n",
			m.host+"/"+m.name+"/"+m.tag, m.value,
			m.time, m.duration, m.message)
	}
	return value
}

//metric_path value timestamp\n
func toGraphite(results []*Check) string {
	value := ""
	for _, m := range results {
		value += fmt.Sprintf("%s %s %d\n",
			m.host+"/"+m.name+"/"+m.tag, m.value, m.time.Unix())
	}
	return value
}

// Check should be in logstash format.
// {
//   "@timestamp": "2012-12-18T01:01:46.092538Z".
//   "tags": [ "kernel", "dmesg" ]
//   "type": "syslog"
//   "message": "usb 3-1.2: USB disconnect, device number 4",
//   "path": "/var/log/messages",
//   "host": "pork.home"
// }
func toJson(results []*Check) string {
	api.Domain = "localhost"
	records := make([]interface{}, 0)
	for _, m := range results {
		r := map[string]interface{}{
			"@timestamp": m.time,
			"type":       m.name,
			"name":       m.tag,
			"message":    m.message,
			"host":       m.host,
			"value":      m.value,
			"duration":   m.duration,
		}
		_, err := core.Index(true, "gmon", m.name, "", r)
		if err != nil {
			log.Println(err.Error())
		}
		records = append(records, r)
	}
	b, _ := json.Marshal(records)
	return string(b)
}

func Send(sub chan []*Check, oFormat *string) {
	results := <-sub

	formats := strings.Split(*oFormat, ",")
	for _, form := range formats {
		if form == "string" {
			println(toString(results))
		}
		if form == "graphite" {
			println(toGraphite(results))
		}
		if form == "json" {
			println(toJson(results))
		}
	}
}

func main() {
	c := exec.Command("hostname")
	output, err := c.CombinedOutput()
	if err != nil {
		log.Printf("Could not get hostname")
		return
	}
	host := strings.TrimSpace(string(output))
	log.Printf("Host: %s", host)
	log.Printf("Running scripts in %s", *path)
	scripts, err := filepath.Glob(*path + "/*")
	if err != nil {
		panic(err)
		return
	}
	cs := make(chan []*Check)
	for {
		for _, script := range scripts {
			name := strings.Replace(script, *path, "", 1)
			if strings.Index(name, ".") == 0 {
				continue
			}
			script, err := filepath.EvalSymlinks(script)
			if err != nil {
				log.Println(err.Error())
				continue
			}
			name = strings.Replace(name, filepath.Ext(script), "", 1)
			go func(cs chan []*Check, host string, name string, script string) {
				Exec(cs, host, name, script)
			}(cs, host, name, script)

			go func(cs chan []*Check, oFormat *string) {
				Send(cs, oFormat)
			}(cs, oFormat)
		}
		time.Sleep(10 * 1e9)
	}
}
