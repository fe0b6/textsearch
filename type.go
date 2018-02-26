package textsearch

// Query - Объект запроса
type Query struct {
	Words []string
}

type analys struct {
	Lex string  `json:"lex"`
	Wt  float64 `json:"wt"`
}

type answer struct {
	Text     string   `json:"text"`
	Analysis []analys `json:"analysis"`
}
