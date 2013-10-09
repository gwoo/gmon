GMon
====

GMon is a Go program to monitor and distribute metrics associated with any
system. A preferred use case is the ability to capture cpu and memory usage
in elastic search to be used with Kibana.

Running
-------
./gmon


Creating Metrics
----------------
A metric is a script in any language that outputs in a standardized format
that GMon can understand. The format is a simple space seperated list of values.

	<name>|<value>|<message>|<tag1> <tag2> <tagN>

You can have any number of tags after the message. If you do not want a
message you can use two spaces. Both the message and tags are optional.


Creating Handlers
-----------------
A Handler is a Go interface. The interface has a store method that receives
a slice of Metrics. Have a look at handles/main.go. GMon ships with two
handlers, Elasticsearch and Stdout. You can add your own or contribute them
back to the project.

#### Dependencies
go get github.com/mattbaird/elastigo