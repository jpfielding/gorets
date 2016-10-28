RETS Explorer
======

RETS client with a ReactJS front end, Gorilla RPC services on the backend.  

Find me at gophers.slack.com#gorets


[Example Explorer UI](cmds/explorer/main.go)
```
All hosted from go
cd explorer/client
npm i
npm run build
go run ../../cmds/explorer/main.go -port 8000 -react ./build

To run in dev mode
npm i
cd explorer/client
npm run start
go run ../../cmds/explorer/main.go -port 8080

```
