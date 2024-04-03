package galileo

import (
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type GalileoClient struct {
	httpClient        *http.Client
	GalileoUrl        string
	GalileoLogin      string
	GalileoTranskey   string
	GalileoProviderId string
	GalileoProductId  string
}

func InitGalileoClient(galileoUrl string, galileoLogin string, galileoTranskey string, galileoProviderId string, galileoProductId string) *GalileoClient {
	fmt.Println("Initializing Galileo Client")
	return &GalileoClient{
		httpClient: &http.Client{
			Timeout: time.Second * 10,
		},
		GalileoUrl:        galileoUrl,
		GalileoLogin:      galileoLogin,
		GalileoTranskey:   galileoTranskey,
		GalileoProviderId: galileoProviderId,
		GalileoProductId:  galileoProductId,
	}

}

func (c *GalileoClient) AddCommonHeaders(req *http.Request) {
	req.Header.Add("accept", "application/json")
	req.Header.Add("response-content-type", "json")
	req.Header.Add("content-type", "application/x-www-form-urlencoded")
}

func (c *GalileoClient) PrepareRequestData(transactionId string) url.Values {
	data := url.Values{}
	data.Set("apiLogin", c.GalileoLogin)
	data.Set("apiTransKey", c.GalileoTranskey)
	data.Set("providerId", c.GalileoProviderId)
	data.Set("transactionId", transactionId)
	data.Set("prodId", c.GalileoProductId)
	return data
}

func (c *GalileoClient) SendRequest(req *http.Request) (*http.Response, error) {
	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %v", err)
	}
	return res, nil
}
