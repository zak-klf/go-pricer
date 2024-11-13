# go-pricer

`go-pricer` is a Go-based application for pricing financial instruments.
Currently bonds, more features like option pricing coming soon!


## Prerequisites

Before building and running the application, make sure you have at least either of the following installed:

- [Go 1.23.x](https://golang.org/dl/)
- [Docker](https://www.docker.com/get-started)

## Building the Application

### With Go

To build the `go-pricer` binary locally using Go, follow these steps:

1. Clone the repository (if you haven’t already):

   ```bash
   git clone https://github.com/zak-klf/go-pricer.git
   cd go-pricer
   ```

2. Install dependencies

    ```bash
    go mod tidy
    ```

3. Build the binary
    ```bash
    make build
    ```

### With Docker

Alternatively, you can build the Docker image directly.

1. Build the docker image

    ```bash
    docker build -t go-pricer . 
    ```


## Running the Application

### With Go

1. Compute some clean prices, ("today" or "now" are also acceptable as issue dates)
    ```bash
    ./bin/go-pricer bond --issue-date 2020-01-01 --maturity-date 2030-01-01 --coupon-rate 0.05 --yield 0.04 --face-value 1000 --frequency 2
    ```

    ```bash
    ./bin/go-pricer bond --issue-date today --maturity-date 2030-01-01 --coupon-rate 0.05 --yield 0.04 --face-value 1000 --frequency 2
    ```

2. Compute some dirty prices, (`today` or `now` are also acceptable as issue dates)
    ```bash
    ./bin/go-pricer bond --issue-date 2020-01-01 --maturity-date 2030-01-01 --coupon-rate 0.05 --yield 0.04 --face-value 1000 --frequency 2 --settlement-date 2024-12-01 --dirty-price
    ```

### With Docker

1. Compute some clean prices, (`today` or `now` are also acceptable as issue dates)
    ```bash
    docker run --rm go-pricer ./go-pricer bond --issue-date 2020-01-01 --maturity-date 2030-01-01 --coupon-rate 0.05 --yield 0.04 --face-value 1000 --frequency 2
    ```

    ```bash
    docker run --rm go-pricer ./go-pricer bond --issue-date today --maturity-date 2030-01-01 --coupon-rate 0.05 --yield 0.04 --face-value 1000 --frequency 2
    ```

2. Compute some dirty prices, ("`today`" or "`now`" are also acceptable as issue dates)
    ```bash
    docker run --rm go-pricer ./go-pricer bond --issue-date 2020-01-01 --maturity-date 2030-01-01 --coupon-rate 0.05 --yield 0.04 --face-value 1000 --frequency 2 --settlement-date 2024-12-01 --dirty-price
    ```

## Running tests

### With Go

Run:
```bash
make test