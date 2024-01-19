package ordinals

import (
	"encoding/json"
	"errors"
	"strings"

	"github.com/tidwall/gjson"
)

const INSCRIPTION_ID_SEPARATOR = "i"
const INSCRIPTION_SATPOINT_SEPARATOR = ":"

// Inscription represents the inscription details
type Inscription struct {
	Id            string `json:"inscription_id"`
	Number        int64  `json:"inscription_number"`
	GenesisHeight int64  `json:"genesis_height"`
	GenesisFee    int64  `json:"genesis_fee"`
	OutputValue   int64  `json:"output_value"`
	Address       string `json:"address"`
	Sat           int64  `json:"sat"`
	SatPoint      string `json:"satpoint"`
	ContentType   string `json:"content_type"`
	ContentLength int    `json:"content_length"`
	Timestamp     uint32 `json:"timestamp"`
}

// GetOwner gets the owner of the inscription
func (i Inscription) GetOwner() string {
	return i.Address
}

// GetGenesisHeight gets the genesis height of the inscription
func (i Inscription) GetGenesisHeight() int64 {
	return i.GenesisHeight
}

// GetGenesisTransaction gets the genesis transaction of the inscription
func (i Inscription) GetGenesisTransaction() string {
	genesisOutput := strings.Split(i.Id, INSCRIPTION_ID_SEPARATOR)

	return genesisOutput[0]
}

// GetOutput gets the output of the inscription
func (i Inscription) GetOutput() string {
	index := strings.LastIndex(i.SatPoint, INSCRIPTION_SATPOINT_SEPARATOR)
	if index == -1 {
		return ""
	}

	return i.SatPoint[0:index]
}

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

// GetInscriptionsByOutputResponse represents the response of GetInscriptionsByOutput
type GetInscriptionsByOutputResponse struct {
	Inscriptions []string `json:"inscriptions"`
}

// GetInscriptions gets the inscriptions from the response
func (r GetInscriptionsByOutputResponse) GetInscriptions() []string {
	return r.Inscriptions
}

// UnmarshalJSON unmarshals the given data to the GetInscriptionsByOutputResponse struct
func (r *GetInscriptionsByOutputResponse) Unmarshal(data []byte) error {
	return json.Unmarshal(data, r)
}
