package coingecko

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

var api = "https://api.coingecko.com/api/v3"

var simplePriceEndpoint = "/simple/price"

var coinIDs = []string{"ripple", "ethereum", "tron", "neo"}
var vsCurrencies = []string{"eur", "usd"}

var cache = struct {
	Timestamp time.Time
	values    map[string]Value
}{
	time.Now(),
	make(map[string]Value),
}

// Value contains market value for cryptocurrency
type Value struct {
	EUR float32 `json:"eur"`
	USD float32 `json:"usd"`
}

// MarketValue returns EUR and USD values for passed cryptocurrency
func MarketValue(cryptocurrency string) Value {
	if _, ok := cache.values[cryptocurrency]; !ok {
		AddCoinID(cryptocurrency)
		updateCache()
		return cache.values[cryptocurrency]
	}
	if timeDelta := time.Now().Sub(cache.Timestamp); timeDelta >= time.Minute*5 {
		updateCache()
	}
	return cache.values[cryptocurrency]
}

// AddCoinID to list of coins for getting market value
func AddCoinID(id string) {
	coinIDs = append(coinIDs, id)
}

// AddVsCurrency to list of vs_currencies for getting market value
// func AddVsCurrency(vsCurrency string) {
// 	vsCurrencies = append(vsCurrencies, vsCurrency)
// }

// RemoveCoinID removes passed crypto id from list of coinIDs
func RemoveCoinID(id string) {
	index := indexOf(id, coinIDs)
	coinIDs = append(coinIDs[:index], coinIDs[index+1:]...)
}

// RemoveVsCurrency removes passed vs_currency id from list of vs_currencies
// func RemoveVsCurrency(vsCurrency string) {
// 	index := indexOf(vsCurrency, coinIDs)
// 	coinIDs = append(vsCurrencies[:index], vsCurrencies[index+1:]...)
// }

func updateCache() {
	req, err := http.NewRequest("GET", api+simplePriceEndpoint, nil)
	if err != nil {
		log.Errorln(err)
		return
	}

	q := req.URL.Query()
	q.Add("ids", strings.Join(coinIDs, ","))
	q.Add("vs_currencies", strings.Join(vsCurrencies, ","))
	req.URL.RawQuery = q.Encode()

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Errorln(err)
		return
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Errorln(err)
		return
	}
	err = json.Unmarshal(data, &cache.values)
	if err != nil {
		log.Errorln(err)
		return
	}

	// set missing coins market value to 0
	for _, coin := range coinIDs {
		if _, ok := cache.values[coin]; !ok {
			cache.values[coin] = Value{
				EUR: 0,
				USD: 0,
			}
		}
	}

	cache.Timestamp = time.Now()
}

// Init populates cache with default cryptocurrency values
//
// defaultCoinIDs = ripple, ethereum, tron, neo
//
// defaultVSCurrencies = eur, usd
func Init() {
	updateCache()
}

func indexOf(element string, data []string) int {
	for k, v := range data {
		if element == v {
			return k
		}
	}
	return -1
}
