# Go REST Template
⚠️ This is a work in progress and is by no means ready for consumption

Go REST Template is a starting point for writing REST APIs in Go. This minimal example shows how to set up an API complete with:
- SQL database and Redis for storing and caching data
- Health checks
- Error handling
- Access and server logs
- Tests (of course)

## Application Structure
This structure of this application is inspired by [Mat Ryer's](https://github.com/matryer) GopherCon 2019 presentation,
["How I Write HTTP Web Services After 8 Years"](https://www.youtube.com/watch?v=rWBSMsLG8po).

TODO more explanation about what's going on here and why

## Technologies
We use [go-chi/chi](https://github.com/go-chi/chi) as the HTTP routing technology. chi is a lightweight, performant package that makes it
easier to define the routes of your API while staying true to the standard library. With chi, witing routes is just like writing a sentence.
It also provides nice middleware support that runs for all requests, unlike popular alternative
[gorilla/mux](https://github.com/gorilla/mux/issues/416).

[go-redis/redis](https://github.com/go-redis/redis) is used for all Redis communication. go-redis is the only officially-recommended Redis
client for Go that employes type-safety when executing commands. It supports Redis clusters, pipelining, and pub/sub and has a large,
active community.

## Techniques at Play
We are defining our own interfaces and wrapping *some* third party libraries for a couple of reasons:

1. Interfaces are critical in Go for stubbing/mocking dependencies and injecting those replacements for unit testing purposes
(["accept interfaces, return structs"](https://medium.com/@cep21/what-accept-interfaces-return-structs-means-in-go-2fe879e25ee8)). However,
idiomatic Go libraries [generally should not export interfaces](https://github.com/golang/go/wiki/CodeReviewComments#interfaces). library users are
instead advised to define their own abbreviated interfaces that only specify the code that is actually being used from the library and
[implement those interfaces by wrapping the library code](https://rakyll.org/interface-pollution/). This is primarily because of how
interfaces in Go are designed; interfaces are implicit rather than explicit and adding new functions to an interface breaks all instances
of that interface because they must now implement that new function. This is particularly problematic in a library setting where you have
many users depending on an exported interface and creating their own implementations as needed. If the library
user instead creates their own minimal interface, a breakage would not occur.

2. Wrapping a third-party library protects your project from future breaking changes and enables you to switch the underlying library or
technology safely without modifying a lot of code everywhere. Modifications can be made under the hood in the wrapper without the callers
needing to be aware of the change. This is a somewhat debated
[practice outside of the Go realm](https://softwareengineering.stackexchange.com/questions/107338/) that the Go language
designers effectively enshrined into the language when the typing and interface systems were built.

TODO more explanation about testing techniques

We use the short flag to skip anything that can take a bit to set up -- namely integration tests.

## Build & Develop
Make is used to script all of the necessary development actions into a single easy command.

Build application:
```
make build
```

Execute tests:
```
make test
```

Or do both at once:
```
make all
```

## Usage
Set up Redis:
```zsh
# TODO set up docker compose
docker run --name redis -p 6379:6379 -d redis:6
```

Build and run application:
```zsh
make run
```

Load data:
```zsh
curl -vX POST localhost:8080/tasks \
    -H 'Accept: application/json' \
    -d '{
        "description": "buy socks"
    }'
```

Retrieve data:
```zsh
curl -v localhost:8080/tasks/<ID>
```

Get health:
```zsh
curl -v localhost:8080/health
```
