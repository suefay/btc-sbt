package ordinals

import (
	"errors"
	"strings"

	"github.com/tidwall/gjson"
)

// GetInscriptionsResponse represents the response of GetInscriptions
type GetInscriptionsResponse struct {
	Inscriptions []gjson.Result `json:"inscriptions"`
}

// GetInscriptions gets the inscriptions from the response
func (r GetInscriptionsResponse) GetInscriptions() []string {
	result := make([]string, 0)

	for _, e := range r.Inscriptions {
		inscriptionRef := e.Get("href").String()
		result = append(result, strings.Split(inscriptionRef, "/")[2])
	}

	return result
}

// UnmarshalJSON unmarshals the given data to the GetInscriptionsResponse struct
func (r *GetInscriptionsResponse) UnmarshalJSON(data []byte) error {
	if !gjson.ValidBytes(data) {
		return errors.New("invalid JSON")
	}

	json := gjson.ParseBytes(data)

	r.Inscriptions = json.Get("inscriptions").Array()

	return nil
}
