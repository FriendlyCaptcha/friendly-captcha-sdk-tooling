# Friendly Captcha SDK Test Server
This module is dedicated to testing the integration of Friendly Captcha SDKs with a mocked server environment.

The goal is to unify the behavior of server side SDK implementations, especially when the API returns errors.

## How to run the server
```shell
# With a binary you downloaded from Github Releases
friendly-captcha-sdk-testserver serve

# Alternative: build and run it locally
go run main.go serve
```

This starts the SDK test server on port `1090`.

Next, run the tests that talk to this test server in the SDK implementation.

### Command-line options

* `--port 1234` run the server on a custom por
* `--tests some/path/my_test_cases_file.json` serve the test cases in a custom fixtures file.

## Adding new sdk tests
The expected behavior of the SDK is defined in the [test_cases.json](./fixtures/test_cases.json) file.