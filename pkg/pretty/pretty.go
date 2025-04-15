package pretty

import "encoding/json"

func JSON(body []byte) []byte {
	var m map[string]interface{}
	json.Unmarshal(body, &m)
	body, _ = json.MarshalIndent(m, "", "  ")
	return body
}
