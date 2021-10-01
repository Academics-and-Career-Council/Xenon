package main

import (
	"log"

	"github.com/AnC-IITK/Xenon/internal/cmd"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func main() {
	cmd.Execute()
}
