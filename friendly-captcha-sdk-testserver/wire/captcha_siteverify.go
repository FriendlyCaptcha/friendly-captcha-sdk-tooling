package wire

type CaptchaSiteverifyRequest struct {
	// Optional, the sitekey that you want to make sure the puzzle was generated from.
	//
	// Not really used in this mock server.
	Sitekey  string `json:"sitekey,omitempty"`
	Response string `json:"response"`
}
