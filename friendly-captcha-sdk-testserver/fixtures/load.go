package fixtures

import (
	"encoding/json"
	"fmt"
	"os"

	_ "embed"

	"github.com/friendlycaptcha/friendly-captcha-sdk-tooling/friendly-captcha-sdk-testserver/model"
)

//go:embed test_cases.json
var testCasesFileBytes []byte

// Load loads the test cases from the embedded JSON file or from the provided filepath.
func Load(filepath string) (model.TestCasesFile, error) {
	var err error
	b := testCasesFileBytes

	// Only read from disk if a filepath was provided.
	if filepath != "" {
		b, err = os.ReadFile(filepath)
		if err != nil {
			return model.TestCasesFile{}, fmt.Errorf("failed to read test cases: %w", err)
		}
	}

	var testCases model.TestCasesFile
	err = json.Unmarshal(b, &testCases)
	if err != nil {
		return model.TestCasesFile{}, fmt.Errorf("failed to parse test cases as JSON: %w", err)
	}

	return testCases, nil
}
