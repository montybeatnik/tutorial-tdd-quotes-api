# Tutorial - TDD Quotes API
This repo illustrates the usage of TDD to build an API dealing with Quotes. 

## Run the server
```bash
go run . 
# or 
go run main.go
# or build it and then run it. 
go build 
./tdd-quotes-api
```

## Testing 
```bash
# all tests
go test ./... # this could just be go test
# all tests with verbosity 
go test -v 
# all specific sub tests
 go test -run TestHandleQuotes/post -v 
# a specific sub tests
 go test -run TestHandleQuotes/post_no_body -v 
```

## Benchmarks
```bash
go test -bench=. -run Bench
```

### Sample from InMem store
```bash
➜  tdd-quotes-api git:(bench) ✗ go test -bench=. -run Bench                         
goos: darwin
goarch: arm64
pkg: github.com/montybeatnik/tutorials/tdd-quotes-api
BenchmarkCreateQuote-10          3380576               324.2 ns/op
PASS
ok      github.com/montybeatnik/tutorials/tdd-quotes-api        1.569s
```

## Curl Examples
```bash
# POST
curl -X POST -d '{
    "author": "Sandi Metz",
    "message": "Design is the art of arranging code that needs to work today, and to be easy to change forever."
}' http://localhost:8000
# GET all quotes
curl -X GET http://localhost:8000
# GET a specific quote
curl -X GET http://localhost:8000/1
```