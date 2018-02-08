package main

import (
	"fmt"
	"log"

	"github.com/doctori/medcase/BDPM"
)

func check(e error) {
	if e != nil {
		log.Panic(e)
	}
}

func main() {

	CISs, err := BDPM.LoadCIS(BDPM.CisDataFile)
	check(err)
	CIPs, err := BDPM.LoadCIP(BDPM.CipDataFile)
	check(err)
	fmt.Printf("loaded : %d entries\n", len(CISs)+len(CIPs))
}
