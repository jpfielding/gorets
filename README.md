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

[Example Search](cmds/search/example.go)
```
go run cmds/search/*.go -connect ~/.gorets/config.json -search ~/gorets/search.json -output /tmp/

search.json:
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

```
[Example GetObject](cmds/getobject/example.go)
```
go run cmds/getobject/*.go -connect ~/.gorets/connect.json -objects ~/gorets/getobjects.json -output /tmp/

getobjects.json:
{
	"resource": "Property",
	"type": "Photo",
	"id": "1330918:*,1555397:*"
}
```
[Example Metadata](cmds/metadata/example.go)
```
go run cmds/metadata/*.go -connect ~/.gorets/connect.json -metadata-options ~/gorets/metadata.json -output /tmp

metadata.json
{
        "metadata-type": "METADATA-SYSTEM",
        "format":       "COMPACT",
        "id":           "*"
}
```
