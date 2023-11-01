from(bucket: "CryptoTrackerBucket")
    |> range(start: -1m)
    |> filter(
        fn: (r) =>
            r._measurement == "crypto_data" and (r._field == "last" ),
    )
    |> last()
    |> group(columns: ["pair"])
    |> mean()
    |> yield(name:"prices")
