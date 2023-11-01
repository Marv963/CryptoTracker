import "strings"
import "date"

bucketName = "{bucket}"

// Abrufen der aktuellen Werte für 'last' und 'quote_volume' für alle Währungspaare, die auf USD enden.
from(bucket: bucketName)
    |> range(start: -20s)
    // Filter auf die letzten 20 Sekunden
    |> filter(
        fn: (r) =>
            r._measurement == "crypto_data" and (r._field == "last" or r._field == "quote_volume")
                and
                r.pair =~ /\/USD$|\/EUR$/,
    )
    |> group(columns: ["exchange", "_field", "pair"])
    // Gruppierung nach Exchange, Field und Pair, um eindeutige Zeitreihen zu erstellen
    |> last()
    // Extrahieren des letzten Wertes
    |> group(columns: ["_field", "pair"])
    // Neue Gruppierung vor der Berechnung des Durchschnitts
    |> mean()
    // Berechnen des Durchschnittswertes von beiden Börsen
    |> pivot(rowKey: ["pair"], columnKey: ["_field"], valueColumn: "_value")
    |> group()
    // Entfernen aller Gruppierungen
    |> sort(columns: ["quote_volume"], desc: true)
    // Sortieren nach 'quote_volume' in absteigender Reihenfolge
    |> yield(name: "getSymbolQuotes")// Ergebnisausgabe
