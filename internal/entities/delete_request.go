package entities

// UserAndShort - pair of user and shortened URL
type UserAndShort struct {
	URL string
	UID UserID
}

// DeleteRequest - a request to delete entries
type DeleteRequest []UserAndShort
