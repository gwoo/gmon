// Copyright 2013 GWoo. All rights reserved.
// The BSD License http://opensource.org/licenses/bsd-license.php.
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

var wd, err1 = os.Getwd()
var path = flag.String("path", wd+"/scripts/", "Path to scripts directory.")

// struct Message {
// 	status int
// 	name, message, time string
// }

func Exec(pub chan string, command string) {
	t0 := time.Now()
	name := strings.Replace(command, *path, "", 1)
	defer func(name string) {
		if x := recover(); x != nil {
			log.Printf("%s %s\n", name, x)
		}
	}(name)
	c := exec.Command(command)
	output, err := c.CombinedOutput()
	if err != nil {
		log.Printf("Error running %s: %s", name, err)
		return
	}
	t := time.Now()
	if string(output) == "" {
		log.Printf("Error running %s: %s", name, "no response")
		return
	}
	response := strings.SplitAfterN(string(output), " ", 2)
	duration := t.Sub(t0)
	result := fmt.Sprintf("%s>%s>%s>%s>%s",
		strings.TrimSpace(response[0]), name, t, duration, response[1])
	pub <- result
}

func Send(sub chan string) {
	result := <-sub
	log.Print(result)
}

func main() {
	log.Printf("Running scripts in %s", *path)
	scripts, err := filepath.Glob(*path + "/*")
	if err != nil {
		panic(err)
		return
	}
	cs := make(chan string)
	for {
		for _, script := range scripts {
			go func(cs chan string, script string) {
				Exec(cs, script)
			}(cs, script)

			go func(cs chan string) {
				Send(cs)
			}(cs)
		}
		time.Sleep(3 * 1e9)
	}
}
