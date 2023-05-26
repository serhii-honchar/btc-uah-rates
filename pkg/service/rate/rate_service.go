package rate

import (
	"btc-uah-rates/pkg/utils"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type CryptoCompareResponse struct {
	Rate float64 `json:"UAH"`
}

type RateService interface {
	GetRate() (float64, error)
}

type rateServiceImpl struct {
	rateProviderUrl string
}

func NewRateService() RateService {
	rateProviderUrl := utils.GetEnvOrDefault("RATE_PROVIDER_URL", "https://min-api.cryptocompare.com/data/price?fsym=BTC&tsyms=UAH")
	return &rateServiceImpl{
		rateProviderUrl: rateProviderUrl,
	}
}

func (r rateServiceImpl) GetRate() (float64, error) {
	resp, err := http.Get(r.rateProviderUrl)
	if err != nil {
		return 0, fmt.Errorf("error making API request: %w", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("error reading response body: %w", err)
	}

	var data CryptoCompareResponse
	err = json.Unmarshal(body, &data)
	if err != nil {
		return 0, fmt.Errorf("error parsing API response: %w", err)
	}

	return data.Rate, nil
}
