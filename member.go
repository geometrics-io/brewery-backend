package brewery

type Member struct {
	Id      string
	Number  int `db:"membernumber`
	Names   []byte
	Contact []byte
}

type Contact struct {
	id   string
	data []byte
}

type MemberStorage interface {
	MemberByID(id string) (*Member, error)
	Add(m *Member) error
	RemoveByID(m string) error
	Update(m *Member) error
}
