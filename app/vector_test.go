package main

import "testing"
import "encoding/json"
import "reflect"

func TestTransform(t *testing.T) {
	var example = `
{
      "id": "tx-3330991687",
      "transaction":      { "amount": 9505.97, "installments": 10, "requested_at": "2026-03-14T05:15:12Z" },
      "customer":         { "avg_amount": 81.28, "tx_count_24h": 20, "known_merchants": ["MERC-008", "MERC-007", "MERC-005"] },
      "merchant":         { "id": "MERC-068", "mcc": "7802", "avg_amount": 54.86 },
      "terminal":         { "is_online": false, "card_present": true, "km_from_home": 952.27 },
      "last_transaction": null
}
`
	var p Payload
	err := json.Unmarshal([]byte(example), &p)
	if err != nil {
		t.Error(err)
		return
	}
	result := Transform(p)
	if !reflect.DeepEqual(result, []float32{0.9506, 0.8333, 1.0, 0.2174, 0.8333, -1, -1, 0.9523, 1.0, 0, 1, 1, 0.75, 0.0055}) {
		t.Errorf("Wrong result: %v", result)
	}
}
