package entities

type ShortLongPair struct {
	Short string
	Long  string
}

type BatchRequest struct {
	ID      string
	LongURL string
}

type BatchResponse struct {
	ID       string
	ShortURL string
}
