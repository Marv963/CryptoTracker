import "join"

bucketName = "{bucket}"

// Funktion zur Abfrage des letzten Preises zu einem bestimmten Zeitpunkt für ein bestimmtes Paar
getLastPrice =
    () =>
        {
            return
                from(bucket: bucketName)
                    |> range(start: -1h)
                    |> filter(fn: (r) => r._measurement == "crypto_data" and r._field == "last")
                    |> last()
                    |> keep(columns: ["_value", "pair"])
                    |> group(columns: ["pair", "_stop"])
                    //    // Gruppierung nach Paar und Stopp-Zeit, um den durchschnittspreis auf den exchanges zu berechnen
                    |> mean()
        }

// Funktion zur Abfrage des letzten Preises zu einem bestimmten Zeitpunkt für ein bestimmtes Paar
getLastPriceFromTime =
    (time) =>
        {
            return
                from(bucket: bucketName)
                    |> range(start: time)
                    |> filter(fn: (r) => r._measurement == "crypto_data" and r._field == "last")
                    |> first()
                    |> keep(columns: ["_value", "pair"])
                    |> group(columns: ["pair", "_stop"])
                    //    // Gruppierung nach Paar und Stopp-Zeit, um den durchschnittspreis auf den exchanges zu berechnen
                    |> mean()
        }

// Abrufen des letzten Preises für BTC/USD
lastPrice = getLastPrice()

price1hAgo = getLastPriceFromTime(time: -1h)

price1dAgo = getLastPriceFromTime(time: -1d)

price7dAgo = getLastPriceFromTime(time: -7d)

// Die Preistabellen zusammenführen
joinFirst =
    join.left(
        left: lastPrice,
        right: price1hAgo,
        on: (l, r) => l.pair == r.pair,
        as: (l, r) => ({pair: l.pair, priceLast: l._value, price1hAgo: r._value}),
    )

joinSecond =
    join.left(
        left: joinFirst,
        right: price1dAgo,
        on: (l, r) => l.pair == r.pair,
        as: (l, r) =>
            ({
                pair: l.pair,
                priceLast: l.priceLast,
                price1hAgo: l.price1hAgo,
                price1dAgo: r._value,
            }),
    )

joinThird =
    join.left(
        left: joinSecond,
        right: price7dAgo,
        on: (l, r) => l.pair == r.pair,
        as: (l, r) =>
            ({
                pair: l.pair,
                priceLast: l.priceLast,
                price1hAgo: l.price1hAgo,
                price1dAgo: l.price1dAgo,
                price7dAgo: r._value,
            }),
    )
        |> map(
            fn: (r) =>
                ({
                    pair: r.pair,
                    priceLast: r.priceLast,
                    price1hAgo: r.price1hAgo,
                    price1dAgo: r.price1dAgo,
                    price7dAgo: r.price7dAgo,
                    priceChange1h: (r.priceLast - r.price1hAgo) / r.price1hAgo * 100.0,
                    priceChange1d: (r.priceLast - r.price1dAgo) / r.price1dAgo * 100.0,
                    priceChange7d: (r.priceLast - r.price7dAgo) / r.price7dAgo * 100.0,
                }),
        )
        |> group()

// Abrufen der letzten Preise von den Exchanges
lastPricesFromExchanges =
    from(bucket: bucketName)
        |> range(start: -1h)
        |> filter(
            fn: (r) =>
                r._measurement == "crypto_data" and (r._field == "last" or r._field == "ask"
                        or
                        r._field == "bid" or r._field == "ask"),
        )
        |> last()
        |> group()
        // remove grouping keys
        |> pivot(rowKey: ["pair"], columnKey: ["exchange", "_field"], valueColumn: "_value")

// Die Preistabellen zusammenführen
join.left(
    left: joinThird,
    right: lastPricesFromExchanges,
    on: (l, r) => l.pair == r.pair,
    as: (l, r) =>
        ({
            pair: l.pair,
            priceLast: l.priceLast,
            price1hAgo: l.price1hAgo,
            price1dAgo: l.price1dAgo,
            price7dAgo: l.price7dAgo,
            priceChange1h: l.priceChange1h,
            priceChange1d: l.priceChange1d,
            priceChange7d: l.priceChange7d,
            priceBitstamp: r.bitstamp_last,
            priceKraken: r.kraken_last,
            bidKraken: r.kraken_bid,
            askKraken: r.kraken_ask,
            bidBitstamp: r.bitstamp_bid,
            askBitstamp: r.bitstamp_ask,
        }),
)
    |> filter(fn: (r) => r.priceLast != 0)
    |> yield(name: "result")

