package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/alecthomas/kong"
	"github.com/friendlycaptcha/friendly-captcha-sdk-tooling/friendly-captcha-sdk-testserver/buildinfo"
	"github.com/friendlycaptcha/friendly-captcha-sdk-tooling/friendly-captcha-sdk-testserver/fixtures"
	"github.com/friendlycaptcha/friendly-captcha-sdk-tooling/friendly-captcha-sdk-testserver/model"
	"github.com/friendlycaptcha/friendly-captcha-sdk-tooling/friendly-captcha-sdk-testserver/wire"
)

// Usage: go run main.go serve --port 1090 --tests ./test_cases.json
// Or just use the defaults: go run main.go serve

const defaultSiteverifyEndpoint = "/api/v2/captcha/siteverify"
const defaultTestsJSONEndpoint = "/api/v1/tests"

var CLI struct {
	Serve struct {
		Port  int    `help:"Port to listen on." default:"1090"`
		Tests string `help:"Path to the test cases (JSON) file, leave empty to use the embedded tests." default:""`
	} `cmd:"" help:"Start the SDK test server."`
	Version struct{} `cmd:"" help:"Print version information and exit."`
}

func main() {
	ctx := kong.Parse(&CLI)
	switch ctx.Command() {
	case "serve":
		serve(CLI.Serve.Port, CLI.Serve.Tests)
	case "version":
		fmt.Println(buildinfo.FullVersion())
	default:
		panic(ctx.Command())
	}

}

func serve(port int, testsPath string) {
	tf, err := fixtures.Load(testsPath)
	if err != nil {
		panic(err)
	}

	if tf.Version != 1 {
		panic("Unsupported test file version")
	}

	if len(tf.Tests) == 0 {
		panic("No tests found")
	}

	responseToCase := make(map[string]model.TestCase)

	for _, test := range tf.Tests {
		responseToCase[test.Response] = test
	}

	mux := http.NewServeMux()

	handler := func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		if r.Header.Get("X-Api-Key") == "" {
			http.Error(w, "Missing X-Api-Key header", http.StatusBadRequest)
			fmt.Println("Missing X-Api-Key header")
			return
		}

		ct := r.Header.Get("Content-Type")

		if ct != "application/json" {
			fmt.Printf("Invalid content-type header %s\n", ct)
			http.Error(w, "Invalid content type header "+ct, http.StatusBadRequest)
			return
		}

		// Check that `Frc-Sdk` header is present.
		if r.Header.Get("Frc-Sdk") == "" {
			fmt.Println("Missing Frc-Sdk header")
			http.Error(w, "Missing Frc-Sdk header", http.StatusBadRequest)
			return
		}

		var req wire.SiteverifyRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to decode request: %s", err), http.StatusBadRequest)
			return
		}

		testCase, ok := responseToCase[req.Response]
		if !ok {
			http.Error(w, fmt.Sprintf("No test case found for response: %s", req.Response), http.StatusBadRequest)
			return
		}
		w.Header().Add("Content-Type", "application/json")

		w.WriteHeader(testCase.SiteverifyStatusCode)
		_, err = w.Write([]byte(testCase.SiteverifyResponse))
		if err != nil {
			fmt.Println("Failed to write response: ", err)
			return
		}
	}

	mux.HandleFunc(defaultSiteverifyEndpoint, handler)
	mux.HandleFunc(defaultTestsJSONEndpoint, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		j, err := json.Marshal(tf)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to marshal tests: %v", err), http.StatusInternalServerError)
			return
		}
		_, err = w.Write(j)
		if err != nil {
			fmt.Println("Failed to write response: ", err)
			return
		}
	})

	fmt.Printf("Serving test cases version %d on port %d", tf.Version, port)
	err = http.ListenAndServe(fmt.Sprintf(":%d", port), mux)
	if err != nil {
		if err == http.ErrServerClosed {
			return
		}
		panic(err)
	}
}
