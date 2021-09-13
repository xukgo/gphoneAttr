package gphoneAttr

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// CheckIsPhoneNumber 判断是否是合法的电话，只判断是纯数字，并且长度不超过16
func CheckIsPhoneNumber(str string) bool {
	if len(str) > 16 {
		return false
	}

	matched, _ := regexp.MatchString(`^\d+`, str)
	return matched
}

const MOBILEPHONE_NUMBER = 1 //手机号
const LINEPHONE_NUMBER = 2   //固话，包括02010086,020114,02096xxx这些
const F400_NUMBER = 3        //400号码，长度为10
const F800_NUMBER = 4        //800号码，长度为10
const F95_NUMBER = 5         //95号码，长度为5或者6
//const OTHER_NUMBER = 9 //其余号码

type PhoneProperty struct {
	Type       int    //1手机 2固话 3其他
	AreaCode   int    //区号，仅固话有效
	FullNumber string //过滤后的号码，比如86xxx会把86删掉,区号前面保证只有一个0
}

func GetPhoneProperty(number string) (PhoneProperty, error) {
	res := PhoneProperty{}
	if !CheckIsPhoneNumber(number) {
		return res, fmt.Errorf("invalid length number %s", number)
	}

	number = strings.TrimLeft(number, "0")
	if number[:2] == "86" {
		number = number[2:]
		number = strings.TrimLeft(number, "0")
	}

	//20110=区号+短服务号
	if len(number) < 5 {
		return res, fmt.Errorf("invalid length number %s", number)
	}
	if len(number) > 11 {
		return res, fmt.Errorf("invalid length number %s", number)
	}

	//手机
	if len(number) == 11 && number[0] == '1' {
		return PhoneProperty{
			Type:       MOBILEPHONE_NUMBER,
			AreaCode:   0,
			FullNumber: number,
		}, nil
	}

	//3区号+8固话
	if len(number) == 11 && number[0] != '1' && number[0] != '2' {
		zcode, _ := strconv.ParseInt(number[:3], 10, 32)
		if zcode >= 300 && zcode <= 999 {
			return PhoneProperty{
				Type:       LINEPHONE_NUMBER,
				AreaCode:   int(zcode),
				FullNumber: fmt.Sprintf("0%s", number),
			}, nil
		}
		return res, fmt.Errorf("invalid number %s", number)
	}

	//11位的全部处理完了
	if len(number) == 11 {
		return res, fmt.Errorf("invalid number %s", number)
	}

	sub2 := number[:2]
	sub3 := number[:3]

	//400号码
	if len(number) == 10 && sub3 == "400" {
		return PhoneProperty{
			Type:       F400_NUMBER,
			AreaCode:   0,
			FullNumber: number,
		}, nil
	}

	//800号码
	if len(number) == 10 && sub3 == "800" {
		return PhoneProperty{
			Type:       F800_NUMBER,
			AreaCode:   0,
			FullNumber: number,
		}, nil
	}

	//95号码
	if (len(number) == 5 || len(number) == 6) && sub2 == "95" {
		return PhoneProperty{
			Type:       F95_NUMBER,
			AreaCode:   0,
			FullNumber: number,
		}, nil
	}

	//2区号+8固话||3区号+7固话||更短的服务号
	if len(number) <= 10 {
		zcode, _ := strconv.ParseInt(sub2, 10, 32)
		if zcode == 10 || (zcode >= 20 && zcode <= 29) {
			return PhoneProperty{
				Type:       LINEPHONE_NUMBER,
				AreaCode:   int(zcode),
				FullNumber: fmt.Sprintf("0%s", number),
			}, nil
		}

		zcode, _ = strconv.ParseInt(sub3, 10, 32)
		if zcode >= 300 && zcode <= 999 {
			return PhoneProperty{
				Type:       LINEPHONE_NUMBER,
				AreaCode:   int(zcode),
				FullNumber: fmt.Sprintf("0%s", number),
			}, nil
		}

		return res, fmt.Errorf("invalid number %s", number)
	}

	return res, fmt.Errorf("invalid number %s", number)
}
