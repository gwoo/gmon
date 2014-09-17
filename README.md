GMon
====

GMon is a Go program to monitor and distribute metrics associated with any
system. A preferred use case is the ability to capture cpu and memory usage
in Elasticsearch to be used with Kibana.

![GMon Dashboard](http://gwoo.github.io/gmon/img/gmon-dashboard.png)


Install
-------
	go install github.com/gwoo/gmon

Build
-----
	go get github.com/gwoo/gmon
	cd $GOPATH/github.com/gwoo/gmon
	go get
	go build

Running
-------
./gmon

	Usage of ./gmon:
	  -conf="gmon.json": Path to config file.
	  -handlers="stdout": Comma seperate list of handlers. ex: elasticseach,stdout.
	  -path="scripts": Path to scripts directory.
	  -interval="5m": Time between each check. Examples: 10s, 5m, 1h


Example Metrics
---------------
Symlink [gmon-scripts](https://github.com/gwoo/gmon-scripts) into a `scripts`
directory in the current directory where `gmon` is installed. Or pass
an absolute path to the scripts directory.


Creating Metrics
----------------
A metric is a script in any language that outputs in a standardized format
that GMon can understand. The format is a simple pipe seperated list of values.

	<name>|<value>|<message>|<tag1> <tag2> <tagN>

You can have any number of tags after the message. If you do not want a
message you can use two spaces. Both the message and tags are optional.


Creating Handlers
-----------------
A Handler is a Go interface. The interface has a `Store()` method that receives
a slice of Metrics. Have a look at handles/main.go. GMon ships with
handlers for stdout, Elasticsearch, Ducksboard, and Orchestrate.io.
You can add your own or contribute them back to the project.
