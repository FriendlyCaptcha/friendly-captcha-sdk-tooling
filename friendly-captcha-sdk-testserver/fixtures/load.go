package fixtures

import (
	"encoding/json"
	"fmt"
	"os"

	_ "embed"

	"github.com/friendlycaptcha/friendly-captcha-sdk-tooling/friendly-captcha-sdk-testserver/model"
)

//go:embed captcha_siteverify_test_cases.json
var captchaSiteverifyTestCasesFileBytes []byte

//go:embed risk_intelligence_retrieve_test_cases.json
var riskIntelligenceRetrieveTestCasesFileBytes []byte

// LoadCaptchaSiteverify loads captcha siteverify test cases from the embedded JSON file or from the provided filepath.
func LoadCaptchaSiteverify(filepath string) (model.CaptchaSiteverifyTestCasesFile, error) {
	var err error
	b := captchaSiteverifyTestCasesFileBytes

	// Only read from disk if a filepath was provided.
	if filepath != "" {
		b, err = os.ReadFile(filepath)
		if err != nil {
			return model.CaptchaSiteverifyTestCasesFile{}, fmt.Errorf("failed to read test cases: %w", err)
		}
	}

	var testCases model.CaptchaSiteverifyTestCasesFile
	err = json.Unmarshal(b, &testCases)
	if err != nil {
		return model.CaptchaSiteverifyTestCasesFile{}, fmt.Errorf("failed to parse test cases as JSON: %w", err)
	}

	return testCases, nil
}

// LoadRiskIntelligenceRetrieve loads risk intelligence retrieve test cases from the embedded JSON file
// or from the provided filepath.
func LoadRiskIntelligenceRetrieve(filepath string) (model.RiskIntelligenceRetrieveTestCasesFile, error) {
	var err error
	b := riskIntelligenceRetrieveTestCasesFileBytes

	// Only read from disk if a filepath was provided.
	if filepath != "" {
		b, err = os.ReadFile(filepath)
		if err != nil {
			return model.RiskIntelligenceRetrieveTestCasesFile{}, fmt.Errorf("failed to read retrieve test cases: %w", err)
		}
	}

	var testCases model.RiskIntelligenceRetrieveTestCasesFile
	err = json.Unmarshal(b, &testCases)
	if err != nil {
		return model.RiskIntelligenceRetrieveTestCasesFile{}, fmt.Errorf("failed to parse retrieve test cases as JSON: %w", err)
	}

	return testCases, nil
}
