package model

import "encoding/json"

// What the SDK should conclude from the API response.
type TestCaseExpectation struct {
	ShouldAccept    bool `json:"should_accept"`
	WasAbleToVerify bool `json:"was_able_to_verify"`
	IsClientError   bool `json:"is_client_error"`
}

type TestCasesFile struct {
	Version int        `json:"version"`
	Tests   []TestCase `json:"tests"`
}

type TestCase struct {
	Name     string `json:"name"`
	Response string `json:"response"`
	Strict   bool   `json:"strict"`

	SiteverifyResponse   json.RawMessage     `json:"siteverify_response"`
	SiteverifyStatusCode int                 `json:"siteverify_status_code"`
	Expectation          TestCaseExpectation `json:"expectation"`
}
