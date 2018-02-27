package dao

import (
	"github.com/clsung/grcode"
)

//解析条形码
func BarDecodeByImage(filePath string) ([]string, error) {
	results, err := grcode.GetDataFromFile(filePath)
	if err != nil {
		return results, err
	}
	return results, nil
}
