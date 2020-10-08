package coingecko

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

var api = "https://api.coingecko.com/api/v3"

var simplePriceEndpoint = "/simple/price"

var defaultCoinIDs = "ripple,ethereum,tron,neo"
var defaultVSCurrencies = "eur,usd"

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

// MarketValue returns EUR and USD values for passed cryptocurrecny
func MarketValue(cryptocurrency string) Value {
	if _, ok := cache.values[cryptocurrency]; !ok {
		updateCache(defaultCoinIDs+","+cryptocurrency, defaultVSCurrencies)
		return cache.values[cryptocurrency]
	}
	if timeDelta := time.Now().Sub(cache.Timestamp); timeDelta >= time.Minute*5 {
		updateCache(defaultCoinIDs+","+cryptocurrency, defaultVSCurrencies)
	}
	return cache.values[cryptocurrency]
}

func updateCache(ids, vsCurrencies string) {
	req, err := http.NewRequest("GET", api+simplePriceEndpoint, nil)
	if err != nil {
		log.Errorln(err)
		return
	}

	q := req.URL.Query()
	q.Add("ids", ids)
	q.Add("vs_currencies", vsCurrencies)
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
}

// Init populates cache with default cryptocurrency values
//
// defaultCoinIDs = ripple, ethereum, tron, neo
//
// defaultVSCurrencies = eur, usd
func Init() {
	updateCache(defaultCoinIDs, defaultVSCurrencies)
}
