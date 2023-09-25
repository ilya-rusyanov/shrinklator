package entities

type UserAndShort struct {
	URL string
	UID UserID
}

type DeleteRequest []UserAndShort
