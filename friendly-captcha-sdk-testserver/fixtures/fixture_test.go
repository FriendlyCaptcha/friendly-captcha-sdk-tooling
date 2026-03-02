package fixtures

import (
	"encoding/json"
	"testing"

	"github.com/guregu/null/v6"
)

type CaptchaSiteverifyResponse struct {
	Success bool                                 `json:"success"`
	Data    CaptchaSiteverifyResponseSuccessData `json:"data,omitempty"`
}

type CaptchaSiteverifyResponseSuccessData struct {
	EventID string `json:"event_id"`

	Challenge CaptchaSiteverifyResponseSuccessDataChallenge `json:"challenge"`

	RiskIntelligence null.Value[RiskIntelligenceData] `json:"risk_intelligence"`
}

type CaptchaSiteverifyResponseSuccessDataChallenge struct {
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

type RiskIntelligenceRetrieveResponse struct {
	Success bool                                        `json:"success"`
	Data    RiskIntelligenceRetrieveResponseSuccessData `json:"data,omitempty"`
}

type RiskIntelligenceRetrieveResponseSuccessData struct {
	RiskIntelligence null.Value[RiskIntelligenceData]        `json:"risk_intelligence"`
	Details          RiskIntelligenceRetrieveResponseDetails `json:"details"`
}

type RiskIntelligenceRetrieveResponseDetails struct {
	Timestamp string `json:"timestamp"`
	ExpiresAt string `json:"expires_at"`
	NumUses   int64  `json:"num_uses"`
}

func TestCaptchaSiteverifyFixtures(t *testing.T) {
	t.Parallel()

	testCases, err := LoadCaptchaSiteverify("")
	if err != nil {
		t.Fatalf("Failed to load embedded siteverify test cases: %v", err)
	}

	invalidJSONCases := map[string]bool{
		"bad_response_200":                 true,
		"bad_response_200_strict":          true,
		"bad_response_500":                 true,
		"bad_response_400_strict":          true,
		"empty_string_response_200":        true,
		"empty_string_response_200_strict": true,
	}

	for i, tc := range testCases.Tests {
		if tc.Name == "" {
			t.Errorf("Test case %d has an empty name", i)
		}
		if tc.Response == "" {
			t.Errorf("Test case %d (%s) has an empty response", i, tc.Name)
		}

		var svr CaptchaSiteverifyResponse
		err := json.Unmarshal(tc.SiteverifyResponse, &svr)
		if err != nil {
			if !invalidJSONCases[tc.Name] {
				t.Errorf("Test case %d (%s) has invalid siteverify_response JSON: %v", i, tc.Name, err)
			}
			continue
		}

		if !svr.Success {
			continue
		}
		validateRiskIntelligenceData(t, i, tc.Name, svr.Data.RiskIntelligence)
	}
}

func TestRiskIntelligenceRetrieveFixtures(t *testing.T) {
	t.Parallel()

	testCases, err := LoadRiskIntelligenceRetrieve("")
	if err != nil {
		t.Fatalf("Failed to load embedded retrieve test cases: %v", err)
	}

	invalidJSONCases := map[string]bool{
		"bad_response_200":          true,
		"bad_response_500":          true,
		"empty_string_response_200": true,
	}

	for i, tc := range testCases.Tests {
		if tc.Name == "" {
			t.Errorf("Test case %d has an empty name", i)
		}
		if tc.Token == "" {
			t.Errorf("Test case %d (%s) has an empty token", i, tc.Name)
		}

		var rr RiskIntelligenceRetrieveResponse
		err := json.Unmarshal(tc.RiskIntelligenceRetrieveResponse, &rr)
		if err != nil {
			if !invalidJSONCases[tc.Name] {
				t.Errorf("Test case %d (%s) has invalid retrieve_response JSON: %v", i, tc.Name, err)
			}
			continue
		}

		if !rr.Success {
			continue
		}

		if rr.Data.Details.Timestamp == "" {
			t.Errorf("Test case %d (%s) has empty details.timestamp", i, tc.Name)
		}
		if rr.Data.Details.ExpiresAt == "" {
			t.Errorf("Test case %d (%s) has empty details.expires_at", i, tc.Name)
		}
		validateRiskIntelligenceData(t, i, tc.Name, rr.Data.RiskIntelligence)
	}
}

func validateRiskIntelligenceData(
	t *testing.T,
	testIndex int,
	testName string,
	ri null.Value[RiskIntelligenceData],
) {
	t.Helper()

	if !ri.Valid {
		return
	}

	if ri.V.Client.HeaderUserAgent == "" {
		t.Errorf("Test case %d (%s) has risk intelligence data with empty client header_user_agent", testIndex, testName)
	}

	if ri.V.Client.Browser.Valid && ri.V.Client.Browser.V.ID == "" {
		t.Errorf("Test case %d (%s) has risk intelligence data with empty browser id", testIndex, testName)
	}
}
