
#!/bin/sh

set -xe

go generate tomate/gen.go

go run main.go &

CURL="curl -s -D -"
sleep 1

$CURL http://localhost:8080/GetByID?id=0
$CURL http://localhost:8080/GetByID?id=10
$CURL --data "color=blue" http://localhost:8080/Create
$CURL --data "color=blue" http://localhost:8080/Create
$CURL --data "color=" http://localhost:8080/Create
$CURL --data "color=green" http://localhost:8080/Create
$CURL -H "Content-Type: application/json" -X POST -d '{"color":"yellow"}' http://localhost:8080/write/1
$CURL http://localhost:8080/GetByID?id=1
$CURL -H "Content-Type: application/json" -X POST -d '{"color":"yellow"}' http://localhost:8080/write/0
$CURL -H "Content-Type: application/json" -X POST -d '{"color":"yellow"}' http://localhost:8080/write/1
$CURL http://localhost:8080/GetByID?id=0
killall main
