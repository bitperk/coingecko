# coingecko

Simple CoinGecko API client. Get market value of cryptos

## Usage

```go
import "github.com/bitperk/coingecko"

func main(){
  coingecko.Init()

  // Pass ID to MarketValue method
  // Refer to coingecko.com/api for list of ids
  coingecko.Instance().MarketValue("ripple")
}
```
