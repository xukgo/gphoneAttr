package gphoneAttr

import (
	"bufio"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path"
	"strconv"
	"strings"
)

type Attribute struct {
	//Feature string `json:"feature"`
	Province string `json:"province"`
	City     string `json:"city"`

	//运营商，0未初始化1移动2电信3联通9虚拟运营商
	Isp      int    `json:"isp,omitempty"`
	IspName  string `json:"ispName,omitempty"`
	ZipCode  string `json:"zipCode"`
	ZoneCode string `json:"zoneCode"`
	Cid      string `json:"cid"`
}

var prefixDict map[string]*Attribute = nil
var zoneDict map[int]*Attribute = nil

//func main() {
//	fmt.Println("hello world")
//	InitFromFile("./prefix.csv")
//	printAttr("15986400521")
//	printAttr("008615986400521")
//	printAttr("0086020110")
//	printAttr("7815986400521")
//	printAttr("20114")
//	printAttr("008699010086")
//	printAttr("02010000")
//}
//
//func printAttr(phone string) {
//	attr, err := GetAttr(phone)
//	if err != nil {
//		fmt.Printf("%s:%s\n", phone, err)
//	} else {
//		fmt.Printf("%s:%s %s\n", phone, attr.IspName, attr.City)
//	}
//}

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
func GetAttrByAreaCode(zcode int) (Attribute, error) {
	attr, find := zoneDict[zcode]
	if find {
		return *attr, nil
	}
	return Attribute{}, fmt.Errorf("invalid areaCode %d", zcode)
}

// InitFromFile filePath支持csv和gz压缩包
//格式：3位前缀|7位前缀|省份|市|运营商|邮编|区号|Cid
func InitFromFile(filePath string) error {
	ext := strings.ToLower(path.Ext(filePath))
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer file.Close()

	var reader *bufio.Reader
	if ext == ".gzip" || (ext == ".gz" && !strings.Contains(filePath, ".tar")) {
		// 创建gzip.Reader
		gr, err := gzip.NewReader(file)
		if err != nil {
			return err
		}
		defer gr.Close()
		reader = bufio.NewReader(gr)
	} else {
		reader = bufio.NewReader(file)
	}

	var line string
	dict := make(map[string]*Attribute, 60*10000)
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
		ispName := sarr[4]
		isp := 0
		if strings.Index(ispName, "移动") >= 0 {
			isp = 1
		} else if strings.Index(ispName, "电信") >= 0 {
			isp = 2
		} else if strings.Index(ispName, "联通") >= 0 {
			isp = 3
		} else {
			isp = 9
		}
		if len(ispName) == 0 || len(sarr[5]) == 0 || len(sarr[6]) == 0 || len(sarr[7]) == 0 {
			fmt.Printf("line wrong format:%s\n", line)
		}
		dict[sarr[1]] = &Attribute{
			Province: sarr[2],
			City:     sarr[3],
			Isp:      isp,
			IspName:  sarr[4],
			ZipCode:  sarr[5],
			ZoneCode: sarr[6],
			Cid:      sarr[7],
		}

		if err != nil {
			break
		}
	}
	if err != io.EOF {
		return err
	}

	zdict := make(map[int]*Attribute, 256)
	for _, v := range dict {
		zcode, err := strconv.ParseInt(v.ZoneCode, 10, 64)
		if err != nil {
			fmt.Printf("zoneCode wrong format:%s\n", v.ZoneCode)
			continue
		}
		zdict[int(zcode)] = v
	}
	for k, v := range zdict {
		item := *v
		item.Isp = 0
		item.IspName = ""
		zdict[k] = &item
	}
	zoneDict = zdict
	prefixDict = dict
	return nil
}
