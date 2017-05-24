# ggt
ggt's generator toolbox

Check the demo [here](https://github.com/mh-cbon/ggt/tree/master/ademo)

# notes

This repository will welcome the refactoring and rationalization of those repos,

- https://github.com/mh-cbon/http-clienter
- https://github.com/mh-cbon/goriller
- https://github.com/mh-cbon/httper
- https://github.com/mh-cbon/jsoner
- https://github.com/mh-cbon/channeler
- https://github.com/mh-cbon/mutexer
- https://github.com/mh-cbon/lister
- https://github.com/mh-cbon/astutil

Please come back later, watch it if you are interested.

```sh
go install && go generate ./ademo/*go
go run ademo/*go &
curl http://localhost:8080/GetById?id=2
curl http://localhost:8080/GetById?id=0
curl --data "color=blue" http://localhost:8080/Create
curl -H "Content-Type: application/json" -X POST -d '{"color":""}' http://localhost:8080/Create
curl -H "Content-Type: application/json" -X POST -d '{"color":"blue"}' http://localhost:8080/Create
```
