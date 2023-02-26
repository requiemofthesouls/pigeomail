package main

import (
	"log"

	"github.com/requiemofthesouls/pigeomail/cmd"
)

func init() {
	log.SetFlags(log.Ldate | log.Lmicroseconds | log.Lshortfile | log.LUTC)
}

func main() {
	var err error
	if err = cmd.Execute(); err != nil {
		log.Fatal(err.Error())
	}
}
