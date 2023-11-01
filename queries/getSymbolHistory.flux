bucketName = "{bucket}"

// Definition des `last_values` Flusses, der den Durchschnitt des "last"-Feldes 
// von den Börsen "bitstamp" und "kraken" für Pair für jedes 1-Stunden-Fenster berechnet.
last_values =
    from(bucket: bucketName)  // Datenquelle spezifizieren
        |> range(start: -7d)  // Zeitbereich für die Abfrage festlegen
        |> filter(fn: (r) => r["_measurement"] == "crypto_data")  // Nur Messwerte namens "crypto_data" auswählen
        |> filter(fn: (r) => r["_field"] == "last")  // Nur Daten für das Feld "last" auswählen
        // Nur Daten von den Börsen "bitstamp" oder "kraken" auswählen
        |> filter(fn: (r) => r["exchange"] == "bitstamp" or r["exchange"] == "kraken")
        // AggregateWindow berechnet den Durchschnitt für jedes 1-Stunden-Zeitfenster, um glattere Daten zu erzeugen
        |> aggregateWindow(every: 5m, fn: mean, createEmpty: false)
        |> filter(fn: (r) => r["pair"] == "%s")  // Nur Pair-Paardaten auswählen
        |> sort(columns: ["_time"])  // Ergebnisse nach Zeit sortieren
        // Nicht benötigte Spalten entfernen, um die Ausgabe zu vereinfachen
        |> drop(columns: ["_measurement", "_field", "_start", "_stop"])

// Definition des `average_value` Flusses, der den Durchschnitt des "last"-Feldes 
// berechnet, indem er die zuvor berechneten stündlichen Durchschnittswerte nutzt.
average_value =
    last_values
        // Gruppieren nach "pair" und "_time", um einen Durchschnitt über alle Börsen zu berechnen
        |> group(columns: ["_time", "pair"], mode: "by")
        |> mean()  // Durchschnitt über alle Börsen berechnen
        // Einen festen Wert "average" zur "exchange" Spalte hinzufügen, um diesen Datenstrom zu kennzeichnen
        |> map(fn: (r) => ({r with exchange: "average"}))

// Die beiden Datenflüsse `last_values` und `average_value` werden vereinigt und als kombinierte Ausgabe bereitgestellt
union(tables: [last_values, average_value])
    |> map(fn: (r) => ({ r with _time: int(v: r._time)}))
    |> yield(name: "combined")
