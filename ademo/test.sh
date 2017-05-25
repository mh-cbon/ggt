
#!/bin/sh

set -xe

go generate tomate/gen.go

go run main.go &

CURL="curl -s -D -"
sleep 1

echo "GetByID";
echo ""
$CURL http://localhost:8080/GetByID?id=0
echo ""
$CURL http://localhost:8080/GetByID?id=10
echo ""
echo "Create";
echo ""
$CURL --data "color=blue" http://localhost:8080/Create
echo ""
$CURL --data "color=blue" http://localhost:8080/Create
echo ""
$CURL --data "color=" http://localhost:8080/Create
echo ""
$CURL --data "color=green" http://localhost:8080/Create
echo ""
echo "Update";
echo ""
$CURL -H "Content-Type: application/json" -X POST -d '{"color":"yellow"}' http://localhost:8080/write/1
echo ""
$CURL http://localhost:8080/GetByID?id=1
echo ""
$CURL -H "Content-Type: application/json" -X POST -d '{"color":"yellow"}' http://localhost:8080/write/0
echo ""
$CURL -H "Content-Type: application/json" -X POST -d '{"color":"yellow"}' http://localhost:8080/write/1
echo ""
$CURL http://localhost:8080/GetByID?id=0
killall main
