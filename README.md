gorets
======

RETS client in Go

[![Build Status](https://travis-ci.org/jpfielding/gorets.svg?branch=master)](https://travis-ci.org/jpfielding/gorets.

The attempt is to meet 1.8.0 compliance.

http://www.reso.org/assets/RETS/Specifications/rets_1_8.pdf.

Find me at gophers.slack.com#gorets

```
config.json:
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

[Example Search](cmds/search/main.go)
```
go run cmds/search/*.go -config-file ~/gorets/config.json -search-options ~/gorets/search.json -output=/tmp/

search.json:
{
	"resource": "Property",
	"class": "Residential",
	"select": "",
	"format": "COMPACT_DECODED",
	"query-type": "dmql2",
	"count-type": 1,
	"limit": 2500,
	"query": "(ModifiedDateTime=2016-08-01T00:00:00+)"
}

```
[Example GetObject](cmds/getobject/main.go)
```
go run cmds/getobject/*.go -config-file ~/gorets/config.json -object-options ~/gorets/getobject.json -output=/tmp/

getobject.json:
{
	"resource": "Property",
	"type": "Photo",
	"id": "1330918:*,1555397:*"
}
```
[Example Metadata](cmds/metadata/main.go)
