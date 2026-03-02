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

// load loads test cases from the embedded bytes or from the provided filepath.
// If filepath is non-empty, the file is read from disk; otherwise embedded is used.
func load[T any](filepath string, embedded []byte) (T, error) {
	var zero T
	var err error
	b := embedded

	if filepath != "" {
		b, err = os.ReadFile(filepath)
		if err != nil {
			return zero, fmt.Errorf("failed to read test cases: %w", err)
		}
	}

	var out T
	err = json.Unmarshal(b, &out)
	if err != nil {
		return zero, fmt.Errorf("failed to parse test cases as JSON: %w", err)
	}

	return out, nil
}

// LoadCaptchaSiteverify loads captcha siteverify test cases from the embedded JSON file or from the provided filepath.
func LoadCaptchaSiteverify(filepath string) (model.CaptchaSiteverifyTestCasesFile, error) {
	return load[model.CaptchaSiteverifyTestCasesFile](filepath, captchaSiteverifyTestCasesFileBytes)
}

// LoadRiskIntelligenceRetrieve loads risk intelligence retrieve test cases from the embedded JSON file or from the provided filepath.
func LoadRiskIntelligenceRetrieve(filepath string) (model.RiskIntelligenceRetrieveTestCasesFile, error) {
	return load[model.RiskIntelligenceRetrieveTestCasesFile](filepath, riskIntelligenceRetrieveTestCasesFileBytes)
}
