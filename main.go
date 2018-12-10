package main

import "github.com/ShrewdSpirit/credman/cmd"

func main() {
	cmd.CheckDataDir()
	cmd.Execute()
}
