// Copyright 2013 GWoo. All rights reserved.
// The BSD License http://opensource.org/licenses/bsd-license.php.

package main

import (
	"fmt"
	h "github.com/gwoo/gmon/handlers"
	"testing"
)

type MockHandler struct{}

var TestResults []*h.Metric

func (m *MockHandler) Config(config []byte) {}
func (m *MockHandler) Store(results []*h.Metric) bool {
	TestResults = results
	return true
}

func init() {
	h.RegisterHandler("mock", &MockHandler{})
}

func Test_response(t *testing.T) {
	result := response("test|100\n", h.Metric{Name: "foo"})

	if result.Name != "foo" {
		t.Error("response `name` failed.")
	}
	if result.Type != "test" {
		t.Error("response `type` failed.")
	}
	if result.Value != 100 {
		t.Error("response `value` failed.")
	}
}

func Test_scriptname(t *testing.T) {
	*path = "/usr/bin/"
	script, name := scriptname("/usr/bin/sar")
	if script != "/usr/bin/sar" {
		t.Error(fmt.Sprintf(
			"scriptname `script` should be %s, but was %s",
			"/usr/bin/sar", script))
	}
	if name != "sar" {
		t.Error(fmt.Sprintf(
			"scriptname `name` should be %s, but was %s",
			"sar", name))
	}
}

func Test_Send(t *testing.T) {
	cs := make(chan []*h.Metric)
	*handlers = "mock"
	go func(chan []*h.Metric) {
		responses := []*h.Metric{&h.Metric{Name: "foo", Value: 100, Type: "test"}}
		cs <- responses
	}(cs)

	Send(cs, []byte{}, handlers)
	if len(TestResults) <= 0 {
		t.Error("No results from Send.")
		return
	}
	if TestResults[0].Name != "foo" {
		t.Error("Send failed with wrong name.")
	}

}
