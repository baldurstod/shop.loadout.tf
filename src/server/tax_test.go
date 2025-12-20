package main_test

import (
	"path"
	"runtime"
	"testing"

	"shop.loadout.tf/src/server/tax"
)

func TestLoadUsTax(t *testing.T) {
	_, filename, _, _ := runtime.Caller(0)
	dir := path.Join(path.Dir(filename), "../../var")
	err := tax.LoadUSTax(path.Join(dir, "AS_zip4_basic_12_25.txt"))
	if err != nil {
		t.Error(err)
		return
	}

}
func TestLoadCanadaTax(t *testing.T) {
	_, filename, _, _ := runtime.Caller(0)
	dir := path.Join(path.Dir(filename), "../../var")
	err := tax.LoadCanadaTax(path.Join(dir, "canada_sales_tax.txt"))
	if err != nil {
		t.Error(err)
		return
	}
}

func TestLoadWorldTax(t *testing.T) {
	_, filename, _, _ := runtime.Caller(0)
	dir := path.Join(path.Dir(filename), "../../var")
	err := tax.LoadWorldSalesTax(path.Join(dir, "world_sales_tax.txt"))
	if err != nil {
		t.Error(err)
		return
	}
}
