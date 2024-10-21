# Friendly Captcha SDK Test Server
This module is dedicated to testing the integration of Friendly Captcha SDKs with a mocked server environment.

The goal is to unify the behavior of server side SDK implementations, especially when the API returns errors.

## Installation
There are different ways to use this tool.

* Download and unpack the binary from the latest version on the [**releases**](https://github.com/FriendlyCaptcha/friendly-captcha-sdk-tooling/releases) page.
* Install using Go
  ```
  go install github.com/friendlycaptcha/friendly-captcha-sdk-tooling/friendly-captcha-sdk-testserver@latest
  ```
* Run it using Docker
  ```
  docker run -p 1090:1090 friendlycaptcha/sdk-testserver:latest
  ```
* Build it yourself (clone this repo, use `go run main.go`).

## Running the server
```shell
# With a binary you downloaded from Github Releases
friendly-captcha-sdk-testserver serve

# Alternative: build and run it locally
go run main.go serve

# With Docker
docker run -p 1090:1090 friendlycaptcha/sdk-testserver:latest
```

This starts the SDK test server on port `1090`.

Next, run the tests that talk to this test server in the SDK implementation.

### Command-line options
You can pass some optional settings:

* `--port 1234` run the server on a custom port.
* `--tests some/path/my_test_cases_file.json` serve the test cases in a custom fixtures file.

## Adding new sdk tests
The expected behavior of the SDK is defined in the [test_cases.json](./fixtures/test_cases.json) file.

## Development

### Minting a release

Releases are built using Github Actions. To mint a local release, install `goreleaser` and run

```shell
goreleaser --snapshot --skip=publish --clean
```

### Publishing to docker

```shell
docker login
docker buildx create --use
docker buildx build --platform=linux/amd64,linux/arm64 -t friendlycaptcha/sdk-testserver:latest . --push
```