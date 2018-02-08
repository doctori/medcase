package BDPM

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"
)

// Standard Location for the CIP Datafile
var CipDataFile = "./data/CIS_CIP_bdpm.txt"

/* CIP : Code Identifiant de Presentation (boites)
* CIS = Code Identifiant de Spécialité
 */
type CIP struct {
	CIS                   int
	CIP7                  int
	Label                 string
	Status                string
	DeclaredState         string
	DeclaredDate          time.Time
	CIP13                 int16
	CollectivityAgreement bool //<== enum ?
	Remboursement         float32
	Prices                []float32
	RemboursementDetail   string
}

// Will read the CSV File on a given path
// And convert all entries to CIP Structs
func LoadCIP(source string) (CIPs []CIP, err error) {
	file, err := os.Open(source)
	if err != nil {
		err = fmt.Errorf("Could not open %s, because : %s", source, err.Error())
		return
	}
	defer file.Close()
	reader := csv.NewReader(bufio.NewReader(file))
	reader.LazyQuotes = true
	reader.Comma = '\t'
	for {
		line, err := reader.Read()
		if err == io.EOF {
			break
		}
		var cip = new(CIP)
		err = cip.ArrayToCIP(line)

		if err != nil {
			panic(err)
		}
		CIPs = append(CIPs, *cip)

	}
	return
}

/* ArrayToCIP
* Array to CIP takes an array extrated from the CSV CIS_CIP_bdpm.txt
* and returns a CIP struct
 */
func (cip *CIP) ArrayToCIP(line []string) (err error) {
	if len(line) < 10 {
		err = fmt.Errorf("Cannot convert into CIP Struct because the line contains %d elements instead of a minimum of 11", len(line)+1)
		return
	}
	cip.CIS, err = strconv.Atoi(line[0])
	if err != nil {
		return
	}
	cip.CIP7, err = strconv.Atoi(line[1])
	if err != nil {
		return
	}
	cip.Label = line[2]
	cip.Status = line[3]
	cip.DeclaredState = line[4]
	cip.DeclaredDate, err = time.Parse("02/01/2006", line[5])
	if err != nil {
		return
	}
	CIP13, err := strconv.ParseInt(line[6], 10, 16)
	cip.CIP13 = int16(CIP13)

	cip.CollectivityAgreement = (line[7] == "oui")
	remboursement, err := strconv.ParseFloat(strings.TrimSuffix(line[8], "%"), 32)
	if err != nil {
		cip.Remboursement = float32(0)
		err = nil
	} else {
		if remboursement < 0 || remboursement > 100 {
			err = fmt.Errorf("Remboursement should be between 0 and 100 percent is is : %f", remboursement)
			return
		}

		cip.Remboursement = float32(remboursement)
	}
	for i := 9; i < len(line)-2; i++ {
		if line[i] == "" {
			break
		}
		price, err := strconv.ParseFloat(line[i], 32)
		if err != nil {
			price = 0
			err = nil
		}
		cip.Prices = append(cip.Prices, float32(price))
	}
	if err != nil {
		panic(err)
	}
	cip.RemboursementDetail = line[len(line)-1]
	return
}
