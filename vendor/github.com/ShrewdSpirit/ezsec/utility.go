package ezsec

import (
	"bytes"

	"github.com/gopherjs/gopherjs/js"
)

func Uint8ArrayToString(value []byte) string {
	return string(value)
}

func StringToUint8Array(value string) []byte {
	return []byte(value)
}

func artobar(arr []interface{}) []byte {
	buf := bytes.Buffer{}
	for _, v := range arr {
		switch v.(type) {
		case int:
			buf.WriteByte(byte(v.(int)))
		case float64:
			buf.WriteByte(byte(v.(float64)))
		}
	}
	return buf.Bytes()
}

func initUtility() {
	js.Module.Get("exports").Set("uint8ArrayToString", Uint8ArrayToString)
	js.Module.Get("exports").Set("stringToUint8Array", StringToUint8Array)
	js.Module.Get("exports").Set("artobar", artobar)
}
