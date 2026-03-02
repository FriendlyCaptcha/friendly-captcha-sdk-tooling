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

// Usage: go run main.go serve --port 1090 --siteverify-tests <siteverify_test_cases_file_path> --retrieve-tests <retrieve_test_cases_file_path>
// Or just use the defaults: go run main.go serve

const (
	defaultCaptchaSiteverifyEndpoint        = "/api/v2/captcha/siteverify"
	defaultRiskIntelligenceRetrieveEndpoint = "/api/v2/riskIntelligence/retrieve"

	// Legacy endpoint for backwards compatibility with SDKs that haven't been updated yet.
	defaultLegacyTestsJSONEndpoint                   = "/api/v1/tests"
	defaultCaptchaSiteverifyTestsJSONEndpoint        = "/api/v1/captcha/siteverifyTests"
	defaultRiskIntelligenceRetrieveTestsJSONEndpoint = "/api/v1/riskIntelligence/retrieveTests"

	expectedCaptchaSiteverifyTestsFileVersion        = 2
	expectedRiskIntelligenceRetrieveTestsFileVersion = 1
)

var CLI struct {
	Serve struct {
		Port          int    `help:"Port to listen on." default:"1090"`
		Tests         string `name:"siteverify-tests" help:"Path to captcha siteverify test cases (JSON), leave empty to use embedded tests." default:""`
		RetrieveTests string `name:"retrieve-tests" help:"Path to risk intelligence retrieve test cases (JSON), leave empty to use embedded tests." default:""`
	} `cmd:"" help:"Start the SDK test server."`
	Version struct{} `cmd:"" help:"Print version information and exit."`
}

func main() {
	ctx := kong.Parse(&CLI)
	switch ctx.Command() {
	case "serve":
		serve(CLI.Serve.Port, CLI.Serve.Tests, CLI.Serve.RetrieveTests)
	case "version":
		fmt.Println(buildinfo.FullVersion())
	default:
		panic(ctx.Command())
	}
}

func serve(port int, siteverifytestsPath string, retrieveTestsPath string) {
	captchaSiteverifyTestsFile, err := fixtures.LoadCaptchaSiteverify(siteverifytestsPath)
	if err != nil {
		panic("failed to load captcha siteverify test cases: " + err.Error())
	}
	riskIntelligenceRetrieveTestsFile, err := fixtures.LoadRiskIntelligenceRetrieve(retrieveTestsPath)
	if err != nil {
		panic("failed to load risk intelligence retrieve test cases: " + err.Error())
	}

	if captchaSiteverifyTestsFile.Version != expectedCaptchaSiteverifyTestsFileVersion {
		panic("Unsupported captcha siteverify test file version")
	}
	if riskIntelligenceRetrieveTestsFile.Version != expectedRiskIntelligenceRetrieveTestsFileVersion {
		panic("Unsupported risk intelligence retrieve test file version")
	}

	if len(captchaSiteverifyTestsFile.Tests) == 0 {
		panic("No captcha siteverify tests found")
	}
	if len(riskIntelligenceRetrieveTestsFile.Tests) == 0 {
		panic("No risk intelligence retrieve tests found")
	}

	responseToCaptchaSiteverifyCase := make(map[string]model.CaptchaSiteverifyTestCase)
	for _, test := range captchaSiteverifyTestsFile.Tests {
		responseToCaptchaSiteverifyCase[test.Response] = test
	}

	responseToRiskIntelligenceRetrieveCase := make(map[string]model.RiskIntelligenceRetrieveTestCase)
	for _, test := range riskIntelligenceRetrieveTestsFile.Tests {
		responseToRiskIntelligenceRetrieveCase[test.Token] = test
	}

	mux := http.NewServeMux()

	validateHeaders := func(w http.ResponseWriter, r *http.Request) bool {
		if r.Header.Get("X-Api-Key") == "" {
			http.Error(w, "Missing X-Api-Key header", http.StatusBadRequest)
			fmt.Println("Missing X-Api-Key header")
			return false
		}

		ct := r.Header.Get("Content-Type")
		if ct != "application/json" {
			fmt.Printf("Invalid content-type header %s\n", ct)
			http.Error(w, "Invalid content type header "+ct, http.StatusBadRequest)
			return false
		}

		if r.Header.Get("Frc-Sdk") == "" {
			fmt.Println("Missing Frc-Sdk header")
			http.Error(w, "Missing Frc-Sdk header", http.StatusBadRequest)
			return false
		}
		return true
	}

	captchaSiteverifyHandler := func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		if !validateHeaders(w, r) {
			return
		}

		var req wire.CaptchaSiteverifyRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to decode request: %s", err), http.StatusBadRequest)
			return
		}

		testCase, ok := responseToCaptchaSiteverifyCase[req.Response]
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

	riskIntelligenceRetrieveHandler := func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		if !validateHeaders(w, r) {
			return
		}

		var req wire.RiskIntelligenceRetrieveRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to decode request: %s", err), http.StatusBadRequest)
			return
		}

		testCase, ok := responseToRiskIntelligenceRetrieveCase[req.Token]
		if !ok {
			http.Error(w, fmt.Sprintf("No test case found for token: %s", req.Token), http.StatusBadRequest)
			return
		}
		w.Header().Add("Content-Type", "application/json")

		w.WriteHeader(testCase.RiskIntelligenceRetrieveStatusCode)
		_, err = w.Write([]byte(testCase.RiskIntelligenceRetrieveResponse))
		if err != nil {
			fmt.Println("Failed to write response: ", err)
			return
		}
	}

	serveJSON := func(w http.ResponseWriter, value any) {
		w.Header().Add("Content-Type", "application/json")
		j, err := json.Marshal(value)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to marshal tests: %v", err), http.StatusInternalServerError)
			return
		}
		_, err = w.Write(j)
		if err != nil {
			fmt.Println("Failed to write response: ", err)
			return
		}
	}

	mux.HandleFunc(defaultCaptchaSiteverifyEndpoint, captchaSiteverifyHandler)
	mux.HandleFunc(defaultRiskIntelligenceRetrieveEndpoint, riskIntelligenceRetrieveHandler)

	mux.HandleFunc(defaultLegacyTestsJSONEndpoint, func(w http.ResponseWriter, r *http.Request) {
		serveJSON(w, captchaSiteverifyTestsFile)
	})
	mux.HandleFunc(defaultCaptchaSiteverifyTestsJSONEndpoint, func(w http.ResponseWriter, r *http.Request) {
		serveJSON(w, captchaSiteverifyTestsFile)
	})
	mux.HandleFunc(defaultRiskIntelligenceRetrieveTestsJSONEndpoint, func(w http.ResponseWriter, r *http.Request) {
		serveJSON(w, riskIntelligenceRetrieveTestsFile)
	})

	fmt.Printf("Serving test cases (captcha siteverify v%d, risk intelligence retrieve v%d) on port %d\n", captchaSiteverifyTestsFile.Version, riskIntelligenceRetrieveTestsFile.Version, port)
	err = http.ListenAndServe(fmt.Sprintf(":%d", port), mux)
	if err != nil {
		if err == http.ErrServerClosed {
			return
		}
		panic(err)
	}
}
