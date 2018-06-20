package wit

// wit response structure
type WitResponseStructMap struct {
	Text string `json:"_text"`
	MsgId string `json:"msg_id"`
	Entities map[string][] WitEntity `json:"entities"`
}

type WitEntity struct {
	Confidence float64 `json:"confidence"`
	Value string `json:"value"`
	Type string `json:"type"`
}