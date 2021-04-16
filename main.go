package main

import (
	"log"

	"github.com/AnC-IITK/Xenon/cmd"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

}

func main() {
	cmd.Execute()
}
