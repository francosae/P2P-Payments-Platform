package enrollment

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strings"

	"github.com/google/uuid"

	"github.com/Sharefunds/auth-service/pkg/galileo"
)

type AccountResponse struct {
	StatusCode      int            `json:"status_code"`
	Status          string         `json:"status"`
	ProcessingTime  float64        `json:"processing_time"`
	ResponseData    []ResponseData `json:"response_data,omitempty"`
	Echo            EchoData       `json:"echo"`
	SystemTimestamp string         `json:"system_timestamp"`
	RToken          string         `json:"rtoken"`
}

type ResponseData struct {
	PMTRefNo             string `json:"pmt_ref_no,omitempty"`
	ProductID            string `json:"product_id,omitempty"`
	GalileoAccountNumber string `json:"galileo_account_number,omitempty"`
	CIP                  string `json:"cip,omitempty"`
	CardID               string `json:"card_id,omitempty"`
	CardNumber           string `json:"card_number,omitempty"`
	ExpiryDate           string `json:"expiry_date,omitempty"`
	CardSecurityCode     string `json:"card_security_code,omitempty"`
	EmbossLine2          string `json:"emboss_line_2,omitempty"`
	BillingCycleDay      string `json:"billing_cycle_day,omitempty"`
	NewEmbossUUID        string `json:"new_emboss_uuid,omitempty"`
}

type EchoData struct {
	ProviderTransactionID string `json:"provider_transaction_id"`
	ProviderTimestamp     string `json:"provider_timestamp"`
	TransactionID         string `json:"transaction_id"`
}

func CreateAccount(c *galileo.GalileoClient) (*AccountResponse, error) {
	u, err := url.Parse(c.GalileoUrl)
	if err != nil {
		return nil, fmt.Errorf("error parsing base URL: %v", err)
	}

	u.Path = path.Join(u.Path, "/createAccount")
	apiUrl := u.String()

	transactionId := uuid.New().String()

	data := c.PrepareRequestData(transactionId)

	req, err := http.NewRequest("POST", apiUrl, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	c.AddCommonHeaders(req)

	res, err := c.SendRequest(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var response AccountResponse
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	fmt.Println("response", response)

	return &AccountResponse{
		StatusCode:      response.StatusCode,
		Status:          response.Status,
		ProcessingTime:  response.ProcessingTime,
		ResponseData:    response.ResponseData,
		Echo:            response.Echo,
		SystemTimestamp: response.SystemTimestamp,
		RToken:          response.RToken,
	}, nil
}
