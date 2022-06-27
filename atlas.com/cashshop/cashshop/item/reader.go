package item

import (
	"atlas-cashshop/wz"
	"atlas-cashshop/xml"
	"errors"
)

func readItems() ([]Model, error) {
	ci, err := getCommodityInfo()
	if err != nil {
		return nil, err
	}

	results := make([]Model, 0)
	for _, c := range ci.Children() {
		i, err := createItem(c)
		if err != nil {
			return results, err
		}
		results = append(results, i)
	}
	return results, nil
}

func createItem(node xml.Noder) (Model, error) {
	in, ok := node.(xml.Parent)
	if !ok {
		return Model{}, errors.New("invalid xml structure")
	}

	serialNumber, err := xml.GetInteger(in, "SN")
	if err != nil {
		return Model{}, err
	}
	itemId, err := xml.GetInteger(in, "ItemId")
	if err != nil {
		return Model{}, err
	}
	price := xml.GetIntegerWithDefault(in, "Price", 0)

	period := xml.GetLongWithDefault(in, "Period", 1)

	count := xml.GetShortWithDefault(in, "Count", 1)

	onSale := xml.GetBooleanWithDefault(in, "OnSale", false)

	return Model{
		serialNumber: uint32(serialNumber),
		itemId:       uint32(itemId),
		price:        uint32(price),
		period:       uint64(period),
		count:        uint16(count),
		onSale:       onSale,
	}, nil
}

func getCommodityInfo() (xml.Parent, error) {
	fe, err := wz.GetFileCache().GetFile("Commodity.img.xml")
	if err != nil {
		return nil, err
	}
	ci, err := xml.Read(fe.Path())
	if err != nil {
		return nil, err
	}
	return ci, nil
}
