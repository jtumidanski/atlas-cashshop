package character

type attributes struct {
	Credit  uint32 `json:"credit"`
	Points  uint32 `json:"points"`
	Prepaid uint32 `json:"prepaid"`
}

func makeAttribute(c Model) attributes {
	return attributes{
		Credit:  c.credit,
		Points:  c.points,
		Prepaid: c.prepaid,
	}
}
