package db

type Mfc struct {
	Id                    string
	Name                  string
	OrganizationFullName  string
	OrganizationAddress   string
	CompletedTicketsCount int
	PendingTicketsCount   int
	Lat                   float64
	Lan                   float64
}

func (m Mfc) GetId() string {
	return m.Id
}
