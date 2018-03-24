gorets
======

RETS command line client in Go 

These are mostly intended to be examples for how a Go client might use the RETS client. 


```json
// connect.json:
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

[Example Search](search.go):
```sh
go run cmds/gorets/*.go search -c connect.json -params sp.json -o /tmp/
```
```json
// sp.json:
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

[Example GetObject](getobject.go):
```sh
go run cmds/getobject/*.go getobject -c connect.json -type "HiRes" -i "343234:*" -o /tmp/
```

[Example Metadata](search.go):
```sh
go run cmds/metadata/*.go metadata -c connect.json -output /tmp/
```
