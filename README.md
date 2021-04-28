# Go REST Template
Go REST Template is a starting point for writing REST APIs in Go. This minimal example shows how to set up an API complete with:
- SQL database and Redis for storing and caching data
- Health checks
- Error handling
- Access and server logs
- Tests (of course)

## Application Structure
This structure of this application is inspired by [Mat Ryer's](https://github.com/matryer) GopherCon 2019 presentation,
["How I Write HTTP Web Services After 8 Years"](https://www.youtube.com/watch?v=8TLiGHJTlig&ab_channel=GopherConEurope).

TODO more explanation about what's going on here and why

## Technologies
We use [go-chi/chi](https://github.com/go-chi/chi) as the HTTP routing technology. chi is a lightweight, performant package that makes it
easier to define the routes of your API while staying true to the standard library. With chi, witing routes is just like writing a sentence.
It also provides nice middleware support that runs for all requests, unlike popular alternative
[gorilla/mux](https://github.com/gorilla/mux/issues/416).

[go-redis/redis] is used for all Redis communication. go-redis is the only officially-recommended Redis client for Go that employes
type-safety when executing commands. It supports Redis clusters, pipelining, and pub/sub and has a large, active community.

## Usage
Set up Redis:
```zsh
docker run --name redis -p 6379:6379 -d redis:6
```

Build application:
```zsh
go build
```

Load data:
```zsh
curl -X POST localhost:8080/dummy \
    -H 'Accept: application/json' \
    -d '{
        "name": "something"
    }'
```

Retrieve data:
```zsh
curl localhost:8080/dummy/6ed70dd5-4c34-4a24-b11f-1cdf1c9f69a0
```
