# golang-waterjug-api
 
A simple REST API to solve the classic water jug riddle.

Check challenge's requirements [here](doc/waterJugChallenge.pdf).

## requirement
- [Docker](docker.com)

## algorithm

Check [here](algorithm.md) for a detailed explanation about the decisions taken.

## running it

```
make run PORT=<port>
```

### example using cURL

request:

```
curl --location 'http://localhost:8080/v1/measure' \
--header 'Content-Type: application/json' \
--data '{
  "x_capacity": 2,
  "y_capacity": 100,
  "z_amount_wanted": 96
}'
```

response:

```
{
  "solution": [
    {
      "step": 1,
      "bucketX": 0,
      "bucketY": 100,
      "action": "Fill bucket Y"
    },
    {
      "step": 2,
      "bucketX": 2,
      "bucketY": 98,
      "action": "Transfer from bucket Y to X"
    },
    {
      "step": 3,
      "bucketX": 0,
      "bucketY": 98,
      "action": "Empty bucket X"
    },
    {
      "step": 4,
      "bucketX": 2,
      "bucketY": 96,
      "action": "Transfer from bucket Y to X",
      "status": "Solved"
    }
  ]
}
```

## running tests

```
make test
```

## running integration tests

```
make int-test
```

## coverage report

```
make coverage
```

## api documentation

Two files hold api's documentation: [doc.go](doc/doc.go) and [api.go](doc/api.go).

To re-generate [doc/swagger.json](doc/swagger.json),

```
make swagger
```

To view it on a browser,

```
make swaggger-ui
```

then visit `localhost`.

## available `Makefile` targets

To generate the basic `Makefile` I've used [go-makefile-gen](https://github.com/tiagomelo/go-makefile-gen), a tool that I've written.

```
$ make help

Usage: make [target]

  help                        shows this help message
  test                        run unit tests
  int-test                    run integration tests
  coverage                    run unit tests and generate coverage report in html format
  swagger                     generates api's documentation
  swagger-ui                  launches swagger ui
  redis-cache                 launch redis cache docker container
  redis-cache-test-instance   launch redis cache docker container for tests
  run                         runs the API
```

## rest api structure

I've started with [example-rest-api](https://github.com/tiagomelo/go-templates/tree/main/example-rest-api) golang template that I've written.

```
gonew github.com/tiagomelo/go-templates/example-rest-api github.com/tiagomelo/golang-waterjug-api
```

## related articles of mine

- [Go: a clean and neat way for managing configuration data from environment variables](https://tiagomelo.info/quicktip/go/envconfig/2024/04/08/golang-envconfig-pdf-post.html)
- [Golang: declarative validation made similar to Ruby on Rails](https://tiagomelo.info/go/validation/2021/04/09/golang-declarative-validation-made-similar-ruby-rails-tiago-melo.html)