package model

import (
	"encoding/json"
)

// What the SDK should conclude from the API response.
type CaptchaSiteverifyTestCaseExpectation struct {
	ShouldAccept    bool `json:"should_accept"`
	WasAbleToVerify bool `json:"was_able_to_verify"`
	IsClientError   bool `json:"is_client_error"`
}

type CaptchaSiteverifyTestCasesFile struct {
	Version int                         `json:"version"`
	Tests   []CaptchaSiteverifyTestCase `json:"tests"`
}

type CaptchaSiteverifyTestCase struct {
	Name     string `json:"name"`
	Response string `json:"response"`
	Strict   bool   `json:"strict"`

	SiteverifyResponse   json.RawMessage                      `json:"siteverify_response"`
	SiteverifyStatusCode int                                  `json:"siteverify_status_code"`
	Expectation          CaptchaSiteverifyTestCaseExpectation `json:"expectation"`
}

type RiskIntelligenceRetrieveTestCaseExpectation struct {
	WasAbleToRetrieve bool `json:"was_able_to_retrieve"`
	IsClientError     bool `json:"is_client_error"`
}

type RiskIntelligenceRetrieveTestCasesFile struct {
	Version int                                `json:"version"`
	Tests   []RiskIntelligenceRetrieveTestCase `json:"tests"`
}

type RiskIntelligenceRetrieveTestCase struct {
	Name  string `json:"name"`
	Token string `json:"token"`

	RiskIntelligenceRetrieveResponse   json.RawMessage                             `json:"retrieve_response"`
	RiskIntelligenceRetrieveStatusCode int                                         `json:"retrieve_status_code"`
	Expectation                        RiskIntelligenceRetrieveTestCaseExpectation `json:"expectation"`
}
