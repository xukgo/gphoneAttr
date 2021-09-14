package gphoneAttr

import (
	"path"
	"testing"
)

var currentDir = "/mnt/diske/GitProject/Go/gphoneAttr"

func TestInitFromCsvFile(t *testing.T) {
	err := InitFromFile(path.Join(currentDir, "prefix2020.csv"))
	if err != nil {
		t.FailNow()
	}
	pp, err := GetPhoneProperty("15986400521")
	if err != nil {
		t.FailNow()
	}
	attr, err := GetAttrByMobilePhonePrefix(pp.FullNumber)
	if err != nil {
		t.FailNow()
	}
	if attr.City != "广州" || attr.IspName != "中国移动" || attr.ZoneCode != "20" || attr.Cid != "440100" {
		t.FailNow()
	}
}

func TestInitFromGzFile(t *testing.T) {
	err := InitFromFile(path.Join(currentDir, "prefix2020.csv.gz"))
	if err != nil {
		t.FailNow()
	}
	pp, err := GetPhoneProperty("15986400521")
	if err != nil {
		t.FailNow()
	}
	attr, err := GetAttrByMobilePhonePrefix(pp.FullNumber)
	if err != nil {
		t.FailNow()
	}
	if attr.City != "广州" || attr.IspName != "中国移动" || attr.ZoneCode != "20" || attr.Cid != "440100" {
		t.FailNow()
	}
}
