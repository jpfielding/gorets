gorets
======

RETS client in Go

[![Build Status](https://travis-ci.org/jpfielding/gorets.svg?branch=master)](https://travis-ci.org/jpfielding/gorets.

The attempt is to meet 1.8.0 compliance.

http://www.reso.org/assets/RETS/Specifications/rets_1_8.pdf.

Find me at gophers.slack.com#gorets


```
connect.json:
{
	"username": "user",
	"password": "pwd",
	"url":	  "http://www.rets.com/rets/login",
	"user-agent": "Company/1.0",
	"user-agent-pw": "",
	"rets-version": "RETS/1.7",
	"wirelog": "/tmp/gorets/wire.log"
}
```

[Example Search](cmds/gorets/search.go)
```
go run cmds/gorets/*.go search -c connect.json -params sp.json -o /tmp/

sp.json:
{
	"SearchType": "Property",
	"Class": "Residential",
	"Select": "",
	"Format": "COMPACT_DECODED",
	"Count": 1,
	"Offset": 1,
	"Limit": 2500,
	"QueryType": "dmql2",
	"Query": "(ModifiedDateTime=2016-08-01T00:00:00+)"
}

[Example GetObject](cmds/gorets/getobject.go)
```
go run cmds/getobject/*.go getobject -c connect.json -type "HiRes" -i "343234:*" -o /tmp/

[Example Metadata](cmds/gorets/search.go)
```
go run cmds/metadata/*.go metadata -c connect.json -output /tmp/

```
