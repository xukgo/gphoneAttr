package gphoneAttr

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

var prefixDict map[string]*Attribute = nil
var zoneDict map[int]*CityAttributeWithCounter = nil

// InitFromReader
//格式：3位前缀|7位前缀|省份|市|运营商|邮编|区号|Cid
func InitFromReader(srcReader io.Reader) error {
	var err error
	var line string
	//provinceNameDict := make(map[string]*string, 50)
	dict := make(map[string]*Attribute, 50*10000)
	zDict := make(map[string]*CityAttributeWithCounter, 500)

	var reader *bufio.Reader
	reader = bufio.NewReader(srcReader)
	for {
		line, err = reader.ReadString('\n')
		if err == io.EOF {
			break
		}

		if len(line) == 0 && err != nil {
			break
		}

		// Process the line here.
		line = strings.TrimSpace(line)
		if len(line) == 0 {
			continue
		}
		sarr := strings.Split(line, ",")
		if len(sarr) != 8 {
			fmt.Println("line format wrong:", line)
			continue
		}

		//每个城市对应的;城市名-邮编-区号-身份证前6位 都是一样的
		sarr[0] = strings.TrimSpace(sarr[0])   //手机前3位
		sarr[1] = strings.TrimSpace(sarr[1])   //手机前7位
		province := strings.TrimSpace(sarr[2]) //省份
		cityName := strings.TrimSpace(sarr[3]) //城市名
		sarr[4] = strings.TrimSpace(sarr[4])   //运营商
		zipCode := strings.TrimSpace(sarr[5])  //邮编
		sarr[6] = strings.TrimSpace(sarr[6])   //区号
		cid := strings.TrimSpace(sarr[7])      //身份证前6位

		if len(zipCode) == 0 || len(sarr[6]) == 0 || len(cid) == 0 {
			fmt.Printf("line wrong format:%s\n", line)
		}
		i64v, err := strconv.ParseInt(strings.TrimLeft(sarr[6], "0"), 10, 64)
		if err != nil {
			return err
		}
		areaCode := int(i64v)
		if areaCode <= 0 {
			return fmt.Errorf("解析区号出错:%s source:%s", sarr[6], line)
		}

		mainIspName, subIspName := extractIspName(sarr[4])
		if len(mainIspName) == 0 {
			fmt.Println("line data invalid ispName:", line)
			continue
		}
		isp := convertIspCode(mainIspName)
		if isp < 0 {
			fmt.Println("line data invalid ispCode:", line)
			continue
		}

		//_, find := provinceNameDict[province]
		//if !find {
		//	provinceNameDict[province] = &province
		//}

		cityBean := CityAttribute{
			Province: province,
			City:     cityName,
			ZipCode:  zipCode,
			ZoneCode: areaCode,
			Cid:      cid,
		}
		zDict = insertZdict(zDict, cityBean)
		dict[sarr[1]] = &Attribute{
			Isp:           isp,
			MainIspName:   &mainIspName,
			SubIspName:    subIspName,
			CityAttribute: zDict[formatZoneCityKey(areaCode, cityName)].CityAttribute,
		}
	}

	if err != io.EOF {
		return err
	}

	ndict := make(map[int]*CityAttributeWithCounter, len(zDict))
	for _, v1 := range zDict {
		zcode := v1.ZoneCode
		v2, find := ndict[zcode]
		if !find {
			ndict[zcode] = v1
			continue
		}
		if ndict[zcode].MatchCount >= v2.MatchCount {
			continue
		}
		ndict[zcode] = v1
	}
	zoneDict = ndict
	prefixDict = dict
	return nil
}

func formatZoneCityKey(zoneCode int, city string) string {
	return fmt.Sprintf("%d:%s", zoneCode, city)
}

func insertZdict(dict map[string]*CityAttributeWithCounter, bean CityAttribute) map[string]*CityAttributeWithCounter {
	key := formatZoneCityKey(bean.ZoneCode, bean.City)
	_, find := dict[key]
	if !find {
		dict[key] = &CityAttributeWithCounter{
			CityAttribute: &bean,
			MatchCount:    0,
		}
		return dict
	}

	if bean.Equals(dict[key].CityAttribute) {
		dict[key].MatchCount++
		return dict
	}

	//fmt.Printf("gphoneAttr zone info notmatch:%v #### %v\n", dict[key].CityAttribute, bean)
	dict[key].MatchCount--
	if dict[key].MatchCount < 0 {
		dict[key].CityAttribute = &bean
		dict[key].MatchCount = 0
	}
	return dict
}

func extractIspName(str string) (string, string) {
	var mainIsp = ""
	var subIsp = ""
	if strings.Index(str, "/") < 0 {
		mainIsp = extSureIspName(str)
	} else {
		arr := strings.Split(str, "/")
		for idx := range arr {
			arr[idx] = strings.TrimSpace(arr[idx])
			if len(arr[idx]) == 0 {
				continue
			}
			if len(mainIsp) == 0 {
				mainIsp = extSureIspName(arr[idx])
				continue
			}
			subIsp = arr[idx]
			//fmt.Printf("subIspName:%s\n", subIsp)
		}
	}
	return mainIsp, subIsp
}

func extSureIspName(mainIspName string) string {
	if strings.Index(mainIspName, "移动") >= 0 {
		return "中国移动"
	} else if strings.Index(mainIspName, "电信") >= 0 {
		return "中国电信"
	} else if strings.Index(mainIspName, "联通") >= 0 {
		return "中国联通"
	} else if strings.Index(mainIspName, "广电") >= 0 {
		return "中国广电"
	} else if strings.Index(mainIspName, "铁通") >= 0 {
		return "中国移动"
	} else {
		fmt.Printf("gphoneAttr main isp name invalid:%s\n", mainIspName)
		return mainIspName
	}
}

func convertIspCode(mainIspName string) int {
	var isp int
	if strings.Index(mainIspName, "移动") >= 0 {
		isp = 1
	} else if strings.Index(mainIspName, "电信") >= 0 {
		isp = 2
	} else if strings.Index(mainIspName, "联通") >= 0 {
		isp = 3
	} else if strings.Index(mainIspName, "广电") >= 0 {
		isp = 4
	} else {
		fmt.Printf("gphoneAttr main isp name invalid:%s\n", mainIspName)
		isp = -1
	}
	return isp
}
