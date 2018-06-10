package brewery

// A Bartender
type Bartender struct {
	id   string
	name struct {
		first string
		last  string
	}
}

// An Authenticater authenticates a bartender to record transactions and manage memebers
type Authenticater interface {
	AuthenticateBartender(bart *Bartender) (string, error)
}
