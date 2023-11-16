package entities

// ShortLongPair - pair of long and shortned URLs
type ShortLongPair struct {
	Short string `json:"short_url"`
	Long  string `json:"original_url"`
}

// BatchRequest - request to batch process URLs
type BatchRequest struct {
	ID      string `json:"correlation_id"`
	LongURL string `json:"original_url"`
}

// BatchResponse - response of URLs batch processing
type BatchResponse struct {
	ID       string `json:"correlation_id"`
	ShortURL string `json:"short_url"`
}

// PairArray - array of short-long pairs
type PairArray []ShortLongPair
