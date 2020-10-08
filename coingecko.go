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

// Service .
type Service interface {
	MarketValue(id string) Value
}

// Coingecko .
type Coingecko struct{}

// MarketValue .
func (*Coingecko) MarketValue(id string) Value {
	if _, ok := cache.values[id]; !ok {
		updateCache(defaultCoinIDs+","+id, defaultVSCurrencies)
		return cache.values[id]
	}
	if timeDelta := time.Now().Sub(cache.Timestamp); timeDelta >= time.Minute*5 {
		updateCache(defaultCoinIDs+","+id, defaultVSCurrencies)
	}
	return cache.values[id]
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

var c *Coingecko

// Instance returns coingecko service instance
func Instance() Service {
	return c
}

// Init .
func Init() {
	c = &Coingecko{}
	updateCache(defaultCoinIDs, defaultVSCurrencies)
}
