package hclutil

import (
	"encoding/json"

	"github.com/zclconf/go-cty/cty"
)

// DumpCTYValue は cty.Value をJSON文字列に変換します。
//
//	これは、ログ出力等を行うときのデバッグ用途を想定しています。
func DumpCTYValue(v cty.Value) (string, error) {
	var raw json.RawMessage
	if err := UnmarshalCTYValue(v, &raw); err != nil {
		return "", err
	}
	return string(raw), nil
}

// MustDumpCtyValue は cty.Value をJSON文字列に変換します。
//
// DumpCTYValue と異なり、エラーが発生した場合は panic します
func MustDumpCtyValue(v cty.Value) string {
	s, err := DumpCTYValue(v)
	if err != nil {
		panic(err)
	}
	return s
}
