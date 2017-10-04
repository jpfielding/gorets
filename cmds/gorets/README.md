gorets
======

RETS command line client in Go 

these are mostly intended to be examples for how a Go client might use the rets client. 


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
