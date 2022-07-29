package gphoneAttr

import (
	"fmt"
	"strings"
	"testing"
)

func TestGetPhoneProperty(t *testing.T) {
	p, err := GetPhoneProperty("86015986400528")
	if err != nil || p.Type != MOBILEPHONE_NUMBER || p.FullNumber != "15986400528" {
		t.FailNow()
	}

	p, err = GetPhoneProperty("4001234567")
	if err != nil || p.Type != F400_NUMBER || p.FullNumber != "4001234567" {
		t.FailNow()
	}

	p, err = GetPhoneProperty("8001234567")
	if err != nil || p.Type != F800_NUMBER || p.FullNumber != "8001234567" {
		t.FailNow()
	}

	p, err = GetPhoneProperty("00095123")
	if err != nil || p.Type != F95_NUMBER || p.FullNumber != "95123" {
		t.FailNow()
	}

	p, err = GetPhoneProperty("000951234")
	if err != nil || p.Type != F95_NUMBER || p.FullNumber != "951234" {
		t.FailNow()
	}

	p, err = GetPhoneProperty("020110")
	if err != nil || p.Type != LINEPHONE_NUMBER || p.AreaCode != 20 || p.FullNumber != "020110" {
		t.FailNow()
	}
	p, err = GetPhoneProperty("0851114")
	if err != nil || p.Type != LINEPHONE_NUMBER || p.AreaCode != 851 || p.FullNumber != "0851114" {
		t.FailNow()
	}

	p, err = GetPhoneProperty("0851114")
	if err == nil && p.Type == LINEPHONE_NUMBER &&
		len(p.FullNumber) >= 5 && len(p.FullNumber) <= 7 && strings.HasSuffix(p.FullNumber, "114") {
		fmt.Printf("")
	}

	p, err = GetPhoneProperty("20961235")
	if err != nil || p.Type != LINEPHONE_NUMBER || p.AreaCode != 20 || p.FullNumber != "020961235" {
		t.FailNow()
	}

	p, err = GetPhoneProperty("0007321234567")
	if err != nil || p.Type != LINEPHONE_NUMBER || p.AreaCode != 732 || p.FullNumber != "07321234567" {
		t.FailNow()
	}

	p, err = GetPhoneProperty("00073212345678")
	if err != nil || p.Type != LINEPHONE_NUMBER || p.AreaCode != 732 || p.FullNumber != "073212345678" {
		t.FailNow()
	}
}
