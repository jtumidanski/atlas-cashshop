package attributes

type LocationInputDataContainer struct {
	Data LocationData `json:"data"`
}

type LocationDataContainer struct {
	Data LocationData `json:"data"`
}

type LocationDataListContainer struct {
	Data []LocationData `json:"data"`
}

type LocationData struct {
	Id         string             `json:"id"`
	Type       string             `json:"type"`
	Attributes LocationAttributes `json:"attributes"`
}

type LocationAttributes struct {
	Type     string `json:"type"`
	MapId    uint32 `json:"mapId"`
	PortalId uint32 `json:"portalId"`
}
