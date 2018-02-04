package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"os"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {

	file, err := os.Open("./data/CIS_CIP_bdpm.txt")
	check(err)
	defer file.Close()
	reader := csv.NewReader(bufio.NewReader(file))
	reader.LazyQuotes = true
	reader.Comma = '\t'
	var CIPs []CIP
	var lines int
	for {
		line, err := reader.Read()
		if err == io.EOF {
			break
		}
		var cip = new(CIP)
		err = cip.ArrayToCIP(line)

		check(err)
		lines = lines + 1
		CIPs = append(CIPs, *cip)

	}
	fmt.Printf("Lines = %d\n", lines)
	fmt.Printf("CIPs : %#v\n", CIPs)

}
