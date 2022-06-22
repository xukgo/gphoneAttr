package gphoneAttr

type Attribute struct {
	*CityAttribute
	Province *string `json:"province"`

	//运营商，0未初始化 1移动 2电信 3联通 4广电 -1错误
	Isp         int     `json:"isp,omitempty"`
	MainIspName *string `json:"mainIspName,omitempty"`
	SubIspName  string  `json:"subIspName,omitempty"`
}

type CityAttributeWithCounter struct {
	*CityAttribute
	MatchCount int
}

type CityAttribute struct {
	City     string `json:"city"`
	ZipCode  string `json:"zipCode"`
	ZoneCode int    `json:"zoneCode"`
	Cid      string `json:"cid"`
}

func (this *CityAttribute) Equals(target *CityAttribute) bool {
	if this.City != target.City {
		return false
	}
	if this.ZipCode != target.ZipCode {
		return false
	}
	if this.ZoneCode != target.ZoneCode {
		return false
	}
	if this.Cid != target.Cid {
		return false
	}
	return true
}
