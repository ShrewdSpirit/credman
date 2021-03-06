package cmdutility

import (
	"fmt"
	"runtime"

	"github.com/fatih/color"
)

var (
	Red          = color.New(color.FgRed)
	BoldRed      = color.New(color.FgRed, color.Bold)
	BoldHiWhite  = color.New(color.FgHiWhite, color.Bold)
	BoldHiYellow = color.New(color.FgHiYellow, color.Bold)
	Green        = color.New(color.FgGreen)
	HiGreen      = color.New(color.FgHiGreen)
)

func LogError(message string, err error) {
	if runtime.GOOS == "windows" {
		fmt.Printf("[ERROR] %s: %s\n", message, err)
		return
	}

	basep := BoldHiWhite.SprintFunc()
	secondp := BoldRed.SprintFunc()
	fmt.Fprintf(color.Output, "%s %s: %s\n", secondp("[ERROR]"), basep(message), secondp(err))
}

func LogColor(c *color.Color, format string, params ...interface{}) {
	if runtime.GOOS == "windows" {
		fmt.Printf(format+"\n", params...)
		return
	}

	newParams := make([]interface{}, 0)
	cprint := c.SprintFunc()
	for _, p := range params {
		newParams = append(newParams, cprint(p))
	}
	fmt.Fprintf(color.Output, format+"\n", newParams...)
}
