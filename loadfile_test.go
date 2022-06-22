package gphoneAttr

import (
	"compress/gzip"
	"os"
	"path"
	"testing"
	"time"
)

var currentDir = "/home/hermes/work"

func BenchmarkInitFromReader(b *testing.B) {
	for i := 0; i < b.N; i++ {
		file, err := os.Open(path.Join(currentDir, "prefix.csv"))
		if err != nil {
			b.FailNow()
		}
		defer file.Close()

		err = InitFromReader(file)
		if err != nil {
			b.FailNow()
		}
	}
}

func TestInitFromCsvFile(t *testing.T) {
	file, err := os.Open(path.Join(currentDir, "prefix.csv"))
	if err != nil {
		t.FailNow()
	}
	defer file.Close()

	err = InitFromReader(file)
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
	if attr.City != "广州" || *attr.MainIspName != "中国移动" || attr.ZoneCode != 20 || attr.Cid != "440100" {
		t.FailNow()
	}
	//time.Sleep(time.Minute * 10)
}

func TestInitFromGzFile(t *testing.T) {
	file, err := os.Open(path.Join(currentDir, "prefix.csv.gz"))
	if err != nil {
		t.FailNow()
	}
	defer file.Close()
	gr, err := gzip.NewReader(file)
	if err != nil {
		t.FailNow()
	}
	defer gr.Close()

	err = InitFromReader(gr)
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
	if attr.City != "广州" || *attr.MainIspName != "中国移动" || attr.ZoneCode != 20 || attr.Cid != "440100" {
		t.FailNow()
	}
	time.Sleep(time.Hour)
}
