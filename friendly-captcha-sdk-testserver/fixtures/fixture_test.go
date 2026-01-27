package fixtures

import (
	"encoding/json"
	"testing"

	"github.com/guregu/null/v6"
)

type SiteverifyResponse struct {
	Success bool                          `json:"success"`
	Data    SiteverifyResponseSuccessData `json:"data,omitempty"`
}

type SiteverifyResponseSuccessData struct {
	EventID string `json:"event_id"`

	Challenge SiteverifyResponseSuccessDataChallenge `json:"challenge"`

	RiskIntelligence null.Value[RiskIntelligenceData] `json:"risk_intelligence"`
}

type SiteverifyResponseSuccessDataChallenge struct {
	Timestamp string `json:"timestamp"`
	Origin    string `json:"origin"`
}

type RiskIntelligenceData struct {
	// Note: this only contains a subset of the actual API response values.

	RiskScores RiskIntelligenceDataRiskScores `json:"risk_scores"`

	Client RiskIntelligenceDataClient `json:"client"`
}

type RiskIntelligenceDataRiskScores struct {
	Overall uint8 `json:"overall"`
	Network uint8 `json:"network"`
	Browser uint8 `json:"browser"`
}

type RiskIntelligenceDataClient struct {
	// Note: this only contains a subset of the actual API response values.

	HeaderUserAgent string                                        `json:"header_user_agent"`
	Browser         null.Value[RiskIntelligenceDataClientBrowser] `json:"browser"`
}

type RiskIntelligenceDataClientBrowser struct {
	// Note: this only contains a subset of the actual API response values.

	ID string `json:"id"`
}

func TestFixtures(t *testing.T) {
	t.Parallel()

	testCases, err := Load("")
	if err != nil {
		t.Fatalf("Failed to load embedded test cases: %v", err)
	}

	// For each fixture make sure it has a name, response and expectation.
	for i, tc := range testCases.Tests {
		if tc.Name == "" {
			t.Errorf("Test case %d has an empty name", i)
		}
		if tc.Response == "" {
			t.Errorf("Test case %d (%s) has an empty response", i, tc.Name)
		}

		// Parse as SiteverifyResponse to ensure it's valid JSON.
		var svr SiteverifyResponse
		err := json.Unmarshal(tc.SiteverifyResponse, &svr)
		if err != nil {
			// Some responses are completely invalid JSON (e.g., HTML error pages).
			// We only validate the ones that are supposed to be valid JSON.

			// We hardcode those that we expect to fail to parse:
			invalidJSONCases := map[string]bool{
				"bad_response_200":                 true,
				"bad_response_200_strict":          true,
				"bad_response_500":                 true,
				"bad_response_400_strict":          true,
				"empty_string_response_200":        true,
				"empty_string_response_200_strict": true,
			}
			if !invalidJSONCases[tc.Name] {
				t.Errorf("Test case %d (%s) has invalid siteverify_response JSON: %v", i, tc.Name, err)
			}
			return
		}

		if !svr.Success { // For non-success cases we don't validate further.
			continue
		}

		ri := svr.Data.RiskIntelligence

		// If risk intelligence data is present, ensure it has expected fields.
		if ri.Valid {
			if ri.V.Client.HeaderUserAgent == "" {
				t.Errorf("Test case %d (%s) has risk intelligence data with empty client header_user_agent", i, tc.Name)
			}

			if ri.V.Client.Browser.Valid {
				if ri.V.Client.Browser.V.ID == "" {
					t.Errorf("Test case %d (%s) has risk intelligence data with empty browser id", i, tc.Name)
				}
			}
		}

	}
}
