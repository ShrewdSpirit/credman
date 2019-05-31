package utility

import (
	"fmt"
)

func Kbmbgb(value int64) string {
	if value/1000000000 > 0 {
		return fmt.Sprintf("%dgb", value/1000000000)
	} else if value/1000000 > 0 {
		return fmt.Sprintf("%dmb", value/1000000)
	} else if value/1000 > 0 {
		return fmt.Sprintf("%dkb", value/1000)
	} else {
		return fmt.Sprintf("%db", value)
	}
}
