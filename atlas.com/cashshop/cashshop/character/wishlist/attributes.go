package wishlist

import (
	"atlas-cashshop/rest/response"
	"strconv"
)

type Attributes struct {
	SerialNumber uint32 `json:"serial_number"`
}

func MakeAttribute(m Model) Attributes {
	return Attributes{
		SerialNumber: m.serialNumber,
	}
}

func MakeRelationshipData(ms []Model) []response.RelationshipData {
	result := make([]response.RelationshipData, 0)
	for _, m := range ms {
		result = append(result, response.RelationshipData{
			Type: "wishlist",
			Id:   strconv.Itoa(int(m.Id())),
		})
	}
	return result
}
