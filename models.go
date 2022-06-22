package gphoneAttr

type Attribute struct {
	*CityAttribute

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
	Province string `json:"province"`
	City     string `json:"city"`
	ZipCode  string `json:"zipCode"`  //邮政编码
	ZoneCode int    `json:"zoneCode"` //区号
	Cid      string `json:"cid"`      //身份证前6位
}

func (this *CityAttribute) Equals(target *CityAttribute) bool {
	if this.Province != target.Province {
		return false
	}
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
