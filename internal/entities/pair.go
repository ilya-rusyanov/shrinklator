package entities

type ShortLongPair struct {
	Short string `json:"short_url"`
	Long  string `json:"original_url"`
}

type BatchRequest struct {
	ID      string `json:"correlation_id"`
	LongURL string `json:"original_url"`
}

type BatchResponse struct {
	ID       string `json:"correlation_id"`
	ShortURL string `json:"short_url"`
}

type PairArray []ShortLongPair
