package main

import (
	"fmt"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {

	CISs, err := LoadCIS(cisDataFile)
	check(err)
	CIPs, err := LoadCIP(cipDataFile)
	check(err)
	fmt.Printf("loaded : %d entries", len(CISs)+len(CIPs))
}
