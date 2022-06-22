package gphoneAttr

import (
	"fmt"
)

//400开头的固定10位，这里不查
//95001这种的一般5/6位，还有10086这种可以带区号
//服务类最少三位，比如110,114,但是一般都要加区号
//手机号固定11位
//我国的固定电话的区号为一般为4位，少数为3位（如北京，上海等）；而电话号码一般为7位或8位。
//所以，拨打国内固定电话，一般为11位或12位，而拨打本地固定电话只需要输入7位或8位，因为不需要拨打区号。

// GetAttrByMobilePhonePrefix 根据手机号前7位获取属性
func GetAttrByMobilePhonePrefix(prefix string) (Attribute, error) {
	if len(prefix) < 7 {
		return Attribute{}, fmt.Errorf("prefix length too short")
	}
	prefix = prefix[:7]
	attr, find := prefixDict[prefix]
	if find {
		return *attr, nil
	}
	return Attribute{}, fmt.Errorf("invalid prefix %s", prefix)
}

// GetAttrByAreaCode 根据固话区号获取属性
func GetAttrByAreaCode(zcode int) (CityAttribute, error) {
	attr, find := zoneDict[zcode]
	if find {
		return *(attr.CityAttribute), nil
	}
	return CityAttribute{}, fmt.Errorf("invalid areaCode %d", zcode)
}

// GetAreaCode 根据号码获取区号，手机号暂时没有区号数据的返回0
func GetAreaCode(phone string) (int, error) {
	phoneProperty, err := GetPhoneProperty(phone)
	if err != nil {
		return 0, err
	}
	if phoneProperty.AreaCode > 0 {
		return phoneProperty.AreaCode, nil
	}
	//400和95的号码确实没有区号一说，返回区号0
	if phoneProperty.Type != MOBILEPHONE_NUMBER {
		return 0, nil
	}

	//手机号暂时没有区号数据的返回0
	attr, err := GetAttrByMobilePhonePrefix(phone)
	if err != nil {
		return 0, nil
	}
	return attr.ZoneCode, nil
}
