package db

type Mfc struct {
	Id                    string  `json:"id"`
	Name                  string  `json:"name"`
	OrganizationFullName  string  `json:"organizationFullName"`
	OrganizationAddress   string  `json:"organizationAddress"`
	CompletedTicketsCount int     `json:"completedTicketsCount"`
	PendingTicketsCount   int     `json:"pendingTicketsCount"`
	Lat                   float64 `json:"lat"`
	Lan                   float64 `json:"lan"`
}

func (m Mfc) GetId() string {
	return m.Id
}
