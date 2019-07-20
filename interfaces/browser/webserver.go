package browser

import (
	"github.com/ShrewdSpirit/credman/gui"
)

func Run() {
	gui.Root.Get("favicon").Bytes()
}
