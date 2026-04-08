package banking

const (
	InReview = Status("in_review")
	Active   = Status("active")
	Closed   = Status("closed")
)

type Status string
