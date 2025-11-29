# Massive Go Client Documentation

## Introduction

The Massive Go Client (formerly Polygon.io) provides official Go bindings for accessing real-time and historical financial market data through both REST and WebSocket APIs. This library serves as a comprehensive interface to Massive's extensive market data infrastructure, covering stocks, options, forex, crypto, indices, and futures markets. Built with Go 1.18+ and utilizing generics, the client offers type-safe access to market data with automatic pagination, connection management, and error handling.

The library is designed for financial applications requiring robust market data integration, from algorithmic trading systems to market analysis tools and portfolio management platforms. It abstracts the complexity of API authentication, rate limiting, and data streaming while providing idiomatic Go interfaces for both synchronous REST operations and asynchronous WebSocket data feeds. The client handles automatic reconnection, subscription management, and provides structured data models for all market data types.

## REST API Client Initialization

```go
package main

import (
    "context"
    "log"
    "os"

    massive "github.com/massive-com/client-go/v2/rest"
    "github.com/massive-com/client-go/v2/rest/models"
)

func main() {
    // Create client with API key from environment
    c := massive.New(os.Getenv("MASSIVE_API_KEY"))

    // Alternative: Create client with custom HTTP client
    // import "net/http"
    // httpClient := &http.Client{Timeout: 30 * time.Second}
    // c := massive.NewWithClient(os.Getenv("MASSIVE_API_KEY"), httpClient)

    // Client is ready to make API calls
    _ = c
}
```

## List Aggregate Bars with Pagination

```go
package main

import (
    "context"
    "log"
    "time"

    massive "github.com/massive-com/client-go/v2/rest"
    "github.com/massive-com/client-go/v2/rest/models"
)

func main() {
    c := massive.New("YOUR_API_KEY")

    // Configure parameters for aggregate bars (OHLCV data)
    params := models.ListAggsParams{
        Ticker:     "AAPL",
        Multiplier: 1,
        Timespan:   "day",
        From:       models.Millis(time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)),
        To:         models.Millis(time.Date(2023, 12, 31, 0, 0, 0, 0, time.UTC)),
    }.WithOrder(models.Asc).WithLimit(50000).WithAdjusted(true)

    // Get iterator - handles pagination automatically
    iter := c.ListAggs(context.Background(), params)

    // Iterate through all results across all pages
    count := 0
    for iter.Next() {
        agg := iter.Item()
        log.Printf("Date: %v, Open: %.2f, High: %.2f, Low: %.2f, Close: %.2f, Volume: %d",
            time.Unix(0, agg.Timestamp*1e6), agg.Open, agg.High, agg.Low, agg.Close, agg.Volume)
        count++
    }

    // Check for errors after iteration
    if iter.Err() != nil {
        log.Fatal(iter.Err())
    }

    log.Printf("Retrieved %d aggregate bars", count)
}
```

## Get Last Trade for a Ticker

```go
package main

import (
    "context"
    "log"

    massive "github.com/massive-com/client-go/v2/rest"
    "github.com/massive-com/client-go/v2/rest/models"
)

func main() {
    c := massive.New("YOUR_API_KEY")

    // Set parameters for last trade request
    params := &models.GetLastTradeParams{
        Ticker: "AAPL",
    }

    // Make synchronous request
    res, err := c.GetLastTrade(context.Background(), params)
    if err != nil {
        log.Fatal(err)
    }

    // Access trade data
    log.Printf("Last Trade - Ticker: %s", res.Results.Ticker)
    log.Printf("Price: %.2f", res.Results.Price)
    log.Printf("Size: %d", res.Results.Size)
    log.Printf("Exchange: %d", res.Results.Exchange)
    log.Printf("Timestamp: %d", res.Results.SipTimestamp)
    log.Printf("Conditions: %v", res.Results.Conditions)
}
```

## List Trades with Time Filtering

```go
package main

import (
    "context"
    "log"
    "time"

    massive "github.com/massive-com/client-go/v2/rest"
    "github.com/massive-com/client-go/v2/rest/models"
)

func main() {
    c := massive.New("YOUR_API_KEY")

    // Get all trades for a specific day using helper method
    params := models.ListTradesParams{Ticker: "AAPL"}.
        WithDay(2023, 7, 22).
        WithOrder(models.Asc).
        WithLimit(50000)

    // Alternative: Use timestamp range
    // params := models.ListTradesParams{Ticker: "AAPL"}.
    //     WithTimestamp(models.GTE, models.Nanos(time.Date(2023, 7, 22, 0, 0, 0, 0, time.UTC))).
    //     WithTimestamp(models.LTE, models.Nanos(time.Date(2023, 7, 22, 23, 59, 59, 0, time.UTC))).
    //     WithOrder(models.Asc)

    iter := c.ListTrades(context.Background(), params)

    tradeCount := 0
    var totalVolume int64

    for iter.Next() {
        trade := iter.Item()
        totalVolume += int64(trade.Size)
        tradeCount++

        if tradeCount <= 5 {
            log.Printf("Trade #%d - Price: %.2f, Size: %d, Time: %d",
                tradeCount, trade.Price, trade.Size, trade.SipTimestamp)
        }
    }

    if iter.Err() != nil {
        log.Fatal(iter.Err())
    }

    log.Printf("Total trades: %d, Total volume: %d", tradeCount, totalVolume)
}
```

## List Quotes with Pagination

```go
package main

import (
    "context"
    "log"
    "time"

    massive "github.com/massive-com/client-go/v2/rest"
    "github.com/massive-com/client-go/v2/rest/models"
)

func main() {
    c := massive.New("YOUR_API_KEY")

    // Get all quotes for a specific day
    params := models.ListQuotesParams{Ticker: "AAPL"}.
        WithDay(2023, 7, 22).
        WithOrder(models.Asc).
        WithLimit(50000)

    iter := c.ListQuotes(context.Background(), params)

    quoteCount := 0
    for iter.Next() {
        quote := iter.Item()
        quoteCount++

        if quoteCount <= 5 {
            log.Printf("Quote #%d - Bid: %.2f (size: %d), Ask: %.2f (size: %d)",
                quoteCount, quote.BidPrice, quote.BidSize, quote.AskPrice, quote.AskSize)
        }
    }

    if iter.Err() != nil {
        log.Fatal(iter.Err())
    }

    log.Printf("Total quotes retrieved: %d", quoteCount)
}
```

## Get Ticker Details

```go
package main

import (
    "context"
    "log"
    "time"

    massive "github.com/massive-com/client-go/v2/rest"
    "github.com/massive-com/client-go/v2/rest/models"
)

func main() {
    c := massive.New("YOUR_API_KEY")

    // Get ticker details with optional date parameter
    params := models.GetTickerDetailsParams{
        Ticker: "AAPL",
    }.WithDate(models.Date(time.Date(2023, 7, 22, 0, 0, 0, 0, time.Local)))

    res, err := c.GetTickerDetails(context.Background(), params)
    if err != nil {
        log.Fatal(err)
    }

    // Access ticker information
    log.Printf("Ticker: %s", res.Results.Ticker)
    log.Printf("Name: %s", res.Results.Name)
    log.Printf("Market: %s", res.Results.Market)
    log.Printf("Locale: %s", res.Results.Locale)
    log.Printf("Primary Exchange: %s", res.Results.PrimaryExchange)
    log.Printf("Type: %s", res.Results.Type)
    log.Printf("Currency: %s", res.Results.CurrencyName)
    log.Printf("Market Cap: %.0f", res.Results.MarketCap)
    log.Printf("Outstanding Shares: %.0f", res.Results.ShareClassSharesOutstanding)
}
```

## List Ticker News

```go
package main

import (
    "context"
    "log"

    massive "github.com/massive-com/client-go/v2/rest"
    "github.com/massive-com/client-go/v2/rest/models"
)

func main() {
    c := massive.New("YOUR_API_KEY")

    // Configure news parameters with filters
    params := models.ListTickerNewsParams{}.
        WithTicker("AAPL").
        WithLimit(10).
        WithOrder(models.Desc)

    iter := c.ListTickerNews(context.Background(), params)

    for iter.Next() {
        news := iter.Item()
        log.Printf("Title: %s", news.Title)
        log.Printf("Author: %s", news.Author)
        log.Printf("Published: %s", news.PublishedUtc)
        log.Printf("URL: %s", news.ArticleURL)
        log.Printf("Description: %s", news.Description)
        log.Printf("Tickers: %v", news.Tickers)
        log.Println("---")
    }

    if iter.Err() != nil {
        log.Fatal(iter.Err())
    }
}
```

## Get Grouped Daily Aggregates

```go
package main

import (
    "context"
    "log"
    "time"

    massive "github.com/massive-com/client-go/v2/rest"
    "github.com/massive-com/client-go/v2/rest/models"
)

func main() {
    c := massive.New("YOUR_API_KEY")

    // Get grouped daily bars for entire market
    params := &models.GetGroupedDailyAggsParams{
        Locale:     models.US,
        MarketType: models.Stocks,
        Date:       models.Date(time.Date(2023, 7, 22, 0, 0, 0, 0, time.Local)),
    }

    res, err := c.GetGroupedDailyAggs(context.Background(), params)
    if err != nil {
        log.Fatal(err)
    }

    log.Printf("Market date: %s", params.Date)
    log.Printf("Total tickers: %d", res.ResultsCount)

    // Display top 10 by volume
    for i, agg := range res.Results {
        if i >= 10 {
            break
        }
        log.Printf("%d. %s - Close: %.2f, Volume: %d, Change: %.2f%%",
            i+1, agg.Ticker, agg.Close, agg.Volume,
            ((agg.Close - agg.Open) / agg.Open * 100))
    }
}
```

## Get Market Status and Holidays

```go
package main

import (
    "context"
    "log"

    massive "github.com/massive-com/client-go/v2/rest"
)

func main() {
    c := massive.New("YOUR_API_KEY")

    // Get current market status
    status, err := c.GetMarketStatus(context.Background())
    if err != nil {
        log.Fatal(err)
    }

    log.Printf("Market: %s", status.Market)
    log.Printf("Server Time: %s", status.ServerTime)
    log.Printf("Exchanges:")
    log.Printf("  NYSE: %s", status.Exchanges.NYSE)
    log.Printf("  NASDAQ: %s", status.Exchanges.Nasdaq)
    log.Printf("  OTC: %s", status.Exchanges.Otc)

    // Get upcoming market holidays
    holidays, err := c.GetMarketHolidays(context.Background())
    if err != nil {
        log.Fatal(err)
    }

    log.Printf("\nUpcoming Market Holidays:")
    for _, holiday := range *holidays {
        log.Printf("  %s - %s (Status: %s)", holiday.Date, holiday.Name, holiday.Status)
        if holiday.Open != nil {
            log.Printf("    Open: %s, Close: %s", *holiday.Open, *holiday.Close)
        }
    }
}
```

## Get Snapshot Data

```go
package main

import (
    "context"
    "log"

    massive "github.com/massive-com/client-go/v2/rest"
    "github.com/massive-com/client-go/v2/rest/models"
)

func main() {
    c := massive.New("YOUR_API_KEY")

    // Get snapshot for a single ticker
    params := &models.GetTickerSnapshotParams{
        Ticker:     "AAPL",
        Locale:     models.US,
        MarketType: models.Stocks,
    }

    res, err := c.GetTickerSnapshot(context.Background(), params)
    if err != nil {
        log.Fatal(err)
    }

    snapshot := res.Ticker
    log.Printf("Ticker: %s", snapshot.Ticker)
    log.Printf("Today's Change: %.2f (%.2f%%)", snapshot.TodaysChange, snapshot.TodaysChangePerc)
    log.Printf("Updated: %d", snapshot.Updated)

    // Current day aggregate
    if snapshot.Day != nil {
        log.Printf("Day - Open: %.2f, High: %.2f, Low: %.2f, Close: %.2f, Volume: %d",
            snapshot.Day.Open, snapshot.Day.High, snapshot.Day.Low, snapshot.Day.Close, snapshot.Day.Volume)
    }

    // Previous day aggregate
    if snapshot.PrevDay != nil {
        log.Printf("Previous Day - Close: %.2f, Volume: %d", snapshot.PrevDay.Close, snapshot.PrevDay.Volume)
    }

    // Last trade
    if snapshot.LastTrade != nil {
        log.Printf("Last Trade - Price: %.2f, Size: %d", snapshot.LastTrade.Price, snapshot.LastTrade.Size)
    }

    // Last quote
    if snapshot.LastQuote != nil {
        log.Printf("Last Quote - Bid: %.2f, Ask: %.2f", snapshot.LastQuote.BidPrice, snapshot.LastQuote.AskPrice)
    }
}
```

## Get Gainers and Losers

```go
package main

import (
    "context"
    "log"

    massive "github.com/massive-com/client-go/v2/rest"
    "github.com/massive-com/client-go/v2/rest/models"
)

func main() {
    c := massive.New("YOUR_API_KEY")

    // Get top gainers
    gainersParams := &models.GetGainersLosersSnapshotParams{
        Locale:     models.US,
        MarketType: models.Stocks,
        Direction:  models.Gainers,
    }

    gainers, err := c.GetGainersLosersSnapshot(context.Background(), gainersParams)
    if err != nil {
        log.Fatal(err)
    }

    log.Println("Top Gainers:")
    for i, ticker := range gainers.Tickers {
        log.Printf("%d. %s - Change: %.2f%%", i+1, ticker.Ticker, ticker.TodaysChangePerc)
    }

    // Get top losers
    losersParams := &models.GetGainersLosersSnapshotParams{
        Locale:     models.US,
        MarketType: models.Stocks,
        Direction:  models.Losers,
    }

    losers, err := c.GetGainersLosersSnapshot(context.Background(), losersParams)
    if err != nil {
        log.Fatal(err)
    }

    log.Println("\nTop Losers:")
    for i, ticker := range losers.Tickers {
        log.Printf("%d. %s - Change: %.2f%%", i+1, ticker.Ticker, ticker.TodaysChangePerc)
    }
}
```

## Get Technical Indicators - SMA

```go
package main

import (
    "context"
    "log"
    "time"

    massive "github.com/massive-com/client-go/v2/rest"
    "github.com/massive-com/client-go/v2/rest/models"
)

func main() {
    c := massive.New("YOUR_API_KEY")

    // Calculate Simple Moving Average
    params := &models.GetSMAParams{
        Ticker:    "AAPL",
        Timespan:  "day",
        Adjusted:  true,
        Window:    50, // 50-day SMA
        SeriesType: "close",
        Timestamp: models.GTE,
        TimestampValue: models.Millis(time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)),
    }

    res, err := c.GetSMA(context.Background(), params)
    if err != nil {
        log.Fatal(err)
    }

    log.Printf("SMA Results for %s", params.Ticker)
    log.Printf("Window: %d %s", params.Window, params.Timespan)

    if len(res.Results.Values) > 0 {
        // Show last 5 values
        start := len(res.Results.Values) - 5
        if start < 0 {
            start = 0
        }

        for _, val := range res.Results.Values[start:] {
            log.Printf("Timestamp: %d, SMA Value: %.2f", val.Timestamp, val.Value)
        }
    }
}
```

## Get Technical Indicators - RSI

```go
package main

import (
    "context"
    "log"
    "time"

    massive "github.com/massive-com/client-go/v2/rest"
    "github.com/massive-com/client-go/v2/rest/models"
)

func main() {
    c := massive.New("YOUR_API_KEY")

    // Calculate Relative Strength Index
    params := &models.GetRSIParams{
        Ticker:    "AAPL",
        Timespan:  "day",
        Adjusted:  true,
        Window:    14, // 14-day RSI
        SeriesType: "close",
        Timestamp: models.GTE,
        TimestampValue: models.Millis(time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)),
    }

    res, err := c.GetRSI(context.Background(), params)
    if err != nil {
        log.Fatal(err)
    }

    log.Printf("RSI Results for %s", params.Ticker)

    if len(res.Results.Values) > 0 {
        latest := res.Results.Values[len(res.Results.Values)-1]
        log.Printf("Latest RSI: %.2f", latest.Value)

        if latest.Value > 70 {
            log.Println("Status: Overbought")
        } else if latest.Value < 30 {
            log.Println("Status: Oversold")
        } else {
            log.Println("Status: Neutral")
        }
    }
}
```

## Get Crypto Data

```go
package main

import (
    "context"
    "log"
    "time"

    massive "github.com/massive-com/client-go/v2/rest"
    "github.com/massive-com/client-go/v2/rest/models"
)

func main() {
    c := massive.New("YOUR_API_KEY")

    // Get crypto aggregate bars
    params := models.ListAggsParams{
        Ticker:     "X:BTCUSD",
        Multiplier: 1,
        Timespan:   "day",
        From:       models.Millis(time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)),
        To:         models.Millis(time.Date(2023, 12, 31, 0, 0, 0, 0, time.UTC)),
    }.WithOrder(models.Desc).WithLimit(365)

    iter := c.ListAggs(context.Background(), params)

    count := 0
    for iter.Next() {
        agg := iter.Item()
        count++

        if count <= 5 {
            log.Printf("BTC/USD - Close: %.2f, High: %.2f, Low: %.2f, Volume: %.2f",
                agg.Close, agg.High, agg.Low, agg.Volume)
        }
    }

    if iter.Err() != nil {
        log.Fatal(iter.Err())
    }

    log.Printf("Total bars retrieved: %d", count)
}
```

## Get Forex Currency Conversion

```go
package main

import (
    "context"
    "log"

    massive "github.com/massive-com/client-go/v2/rest"
    "github.com/massive-com/client-go/v2/rest/models"
)

func main() {
    c := massive.New("YOUR_API_KEY")

    // Real-time currency conversion
    params := &models.GetRealTimeCurrencyConversionParams{
        From:      "USD",
        To:        "EUR",
        Amount:    100.0,
        Precision: 2,
    }

    res, err := c.GetRealTimeCurrencyConversion(context.Background(), params)
    if err != nil {
        log.Fatal(err)
    }

    log.Printf("Converting %.2f %s to %s", params.Amount, params.From, params.To)
    log.Printf("Converted Amount: %.2f %s", res.Converted, res.To)
    log.Printf("Exchange Rate: %.4f", res.Last.Ask)
    log.Printf("Timestamp: %d", res.Last.Timestamp)
}
```

## WebSocket Client - Real-time Stock Data

```go
package main

import (
    "os"
    "os/signal"
    "log"

    massivews "github.com/massive-com/client-go/v2/websocket"
    "github.com/massive-com/client-go/v2/websocket/models"
)

func main() {
    // Create WebSocket client
    c, err := massivews.New(massivews.Config{
        APIKey: os.Getenv("MASSIVE_API_KEY"),
        Feed:   massivews.RealTime,
        Market: massivews.Stocks,
    })
    if err != nil {
        log.Fatal(err)
    }
    defer c.Close()

    // Connect to server
    if err := c.Connect(); err != nil {
        log.Fatal(err)
    }

    // Subscribe to trades and quotes for specific tickers
    if err := c.Subscribe(massivews.StocksTrades, "AAPL", "MSFT", "GOOGL"); err != nil {
        log.Fatal(err)
    }

    if err := c.Subscribe(massivews.StocksQuotes, "AAPL", "MSFT", "GOOGL"); err != nil {
        log.Fatal(err)
    }

    // Handle graceful shutdown
    sigint := make(chan os.Signal, 1)
    signal.Notify(sigint, os.Interrupt)

    // Process incoming data
    for {
        select {
        case <-sigint:
            log.Println("Shutting down...")
            return
        case err := <-c.Error():
            log.Printf("Error: %v", err)
            return
        case out, more := <-c.Output():
            if !more {
                return
            }

            switch data := out.(type) {
            case models.EquityTrade:
                log.Printf("Trade - %s: Price=%.2f, Size=%d, Time=%d",
                    data.Symbol, data.Price, data.Size, data.Timestamp)
            case models.EquityQuote:
                log.Printf("Quote - %s: Bid=%.2f, Ask=%.2f, BidSize=%d, AskSize=%d",
                    data.Symbol, data.BidPrice, data.AskPrice, data.BidSize, data.AskSize)
            case models.EquityAgg:
                log.Printf("Agg - %s: Open=%.2f, High=%.2f, Low=%.2f, Close=%.2f, Volume=%d",
                    data.Symbol, data.Open, data.High, data.Low, data.Close, data.Volume)
            }
        }
    }
}
```

## WebSocket Client - Multiple Markets

```go
package main

import (
    "os"
    "os/signal"
    "log"

    massivews "github.com/massive-com/client-go/v2/websocket"
    "github.com/massive-com/client-go/v2/websocket/models"
)

func main() {
    // Create crypto WebSocket client
    cryptoClient, err := massivews.New(massivews.Config{
        APIKey: os.Getenv("MASSIVE_API_KEY"),
        Feed:   massivews.RealTime,
        Market: massivews.Crypto,
    })
    if err != nil {
        log.Fatal(err)
    }
    defer cryptoClient.Close()

    if err := cryptoClient.Connect(); err != nil {
        log.Fatal(err)
    }

    // Subscribe to crypto trades
    if err := cryptoClient.Subscribe(massivews.CryptoTrades, "X:BTCUSD", "X:ETHUSD"); err != nil {
        log.Fatal(err)
    }

    // Subscribe to all crypto aggregates
    if err := cryptoClient.Subscribe(massivews.CryptoMinAggs); err != nil {
        log.Fatal(err)
    }

    sigint := make(chan os.Signal, 1)
    signal.Notify(sigint, os.Interrupt)

    for {
        select {
        case <-sigint:
            return
        case err := <-cryptoClient.Error():
            log.Printf("Error: %v", err)
            return
        case out, more := <-cryptoClient.Output():
            if !more {
                return
            }

            switch data := out.(type) {
            case models.CryptoTrade:
                log.Printf("Crypto Trade - %s: Price=%.2f, Size=%.6f",
                    data.Pair, data.Price, data.Size)
            case models.CurrencyAgg:
                log.Printf("Crypto Agg - %s: Close=%.2f, Volume=%.2f",
                    data.Pair, data.Close, data.Volume)
            }
        }
    }
}
```

## WebSocket Client - With Custom Logger

```go
package main

import (
    "os"
    "os/signal"

    massivews "github.com/massive-com/client-go/v2/websocket"
    "github.com/massive-com/client-go/v2/websocket/models"
    "github.com/sirupsen/logrus"
)

func main() {
    // Configure custom logger
    log := logrus.New()
    log.SetLevel(logrus.DebugLevel)
    log.SetFormatter(&logrus.JSONFormatter{})

    // Create client with logger
    c, err := massivews.New(massivews.Config{
        APIKey: os.Getenv("MASSIVE_API_KEY"),
        Feed:   massivews.RealTime,
        Market: massivews.Stocks,
        Log:    log,
    })
    if err != nil {
        log.Fatal(err)
    }
    defer c.Close()

    if err := c.Connect(); err != nil {
        log.Error(err)
        return
    }

    // Subscribe to specific data streams
    _ = c.Subscribe(massivews.StocksTrades, "SPY", "QQQ")
    _ = c.Subscribe(massivews.StocksQuotes, "SPY", "QQQ")

    sigint := make(chan os.Signal, 1)
    signal.Notify(sigint, os.Interrupt)

    for {
        select {
        case <-sigint:
            return
        case <-c.Error():
            return
        case out, more := <-c.Output():
            if !more {
                return
            }

            switch data := out.(type) {
            case models.EquityAgg:
                log.WithFields(logrus.Fields{"aggregate": data}).Info("Received aggregate")
            case models.EquityTrade:
                log.WithFields(logrus.Fields{"trade": data}).Info("Received trade")
            case models.EquityQuote:
                log.WithFields(logrus.Fields{"quote": data}).Info("Received quote")
            }
        }
    }
}
```

## Advanced Request Options

```go
package main

import (
    "context"
    "log"
    "strconv"
    "time"

    massive "github.com/massive-com/client-go/v2/rest"
    "github.com/massive-com/client-go/v2/rest/models"
)

func main() {
    c := massive.New("YOUR_API_KEY")

    // Use custom request options
    params := &models.GetGroupedDailyAggsParams{
        Locale:     models.US,
        MarketType: models.Stocks,
        Date:       models.Date(time.Date(2023, 7, 22, 0, 0, 0, 0, time.Local)),
    }

    // Add custom headers and query parameters
    res, err := c.GetGroupedDailyAggs(context.Background(), params,
        models.APIKey("ALTERNATE_API_KEY"),
        models.Header("X-Custom-Header", "CustomValue"),
        models.QueryParam("adjusted", strconv.FormatBool(true)),
    )

    if err != nil {
        log.Fatal(err)
    }

    log.Printf("Retrieved %d tickers", res.ResultsCount)
}
```

## Debug Mode for Troubleshooting

```go
package main

import (
    "context"
    "log"
    "time"

    massive "github.com/massive-com/client-go/v2/rest"
    "github.com/massive-com/client-go/v2/rest/models"
)

func main() {
    c := massive.New("YOUR_API_KEY")

    // Enable trace mode to see request/response details
    params := models.ListAggsParams{
        Ticker:     "AAPL",
        Multiplier: 1,
        Timespan:   "day",
        From:       models.Millis(time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)),
        To:         models.Millis(time.Date(2023, 1, 31, 0, 0, 0, 0, time.UTC)),
    }

    // WithTrace(true) prints request URL, headers, and response headers
    iter := c.ListAggs(context.Background(), params, models.WithTrace(true))

    // Output will include:
    // Request URL: /v2/aggs/ticker/AAPL/range/1/day/...
    // Request Headers: map[...]
    // Response Headers: map[...]

    for iter.Next() {
        log.Print(iter.Item())
    }

    if iter.Err() != nil {
        log.Fatal(iter.Err())
    }
}
```

## Summary

The Massive Go Client provides a production-ready solution for integrating financial market data into Go applications. The REST API client handles all synchronous data retrieval operations with automatic pagination through iterators, making it simple to fetch historical trades, quotes, aggregates, and reference data without manually managing cursor-based pagination. All list operations return iterators that abstract pagination complexity while providing error handling and type-safe access to results.

The WebSocket client delivers real-time streaming data with automatic reconnection, subscription management, and structured data models. It supports multiple concurrent connections to different markets (stocks, options, forex, crypto, indices, futures) and handles connection lifecycle management transparently. Both clients integrate seamlessly with standard Go patterns using context for cancellation, proper error handling through returned errors and channels, and support for custom HTTP clients and loggers. The library's comprehensive type system ensures compile-time safety while the extensive example collection demonstrates real-world usage patterns for algorithmic trading, market analysis, portfolio management, and financial application development.
