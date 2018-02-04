package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

/* CIS : Code Identifiant de Spécialitées (Noms)
* CIS = Code Identifiant de Spécialité
 */
type CIS struct {
	CIS                   int
	Name                  string
	Form                  string
	Voies                 []string
	AMMState              string
	AMMProcedure          string
	Commercialisation     string
	AMMDate               time.Time
	BDMStatus             string
	EuropeanAuthorisation string
	Titulaires            []string
	ExtremeMonitoring     bool
}

/* ArrayToCIP
* Array to cis takes an array extrated from the CSV CIS_CIP_bdpm.txt
* and returns a cis struct
 */
func (cis *CIS) ArrayToCIS(line []string) (err error) {
	if len(line) < 11 {
		err = fmt.Errorf("Cannot convert into cis Struct because the line contains %d elements instead of a minimum of 11", len(line)+1)
		return
	}
	cis.CIS, err = strconv.Atoi(line[0])
	if err != nil {
		return
	}
	cis.Name = line[1]
	cis.Form = line[2]
	for _, voie := range strings.Split(line[3], ";") {
		cis.Voies = append(cis.Voies, voie)
	}
	cis.AMMState = line[4]
	cis.AMMProcedure = line[5]
	cis.Commercialisation = line[6]
	cis.AMMDate, err = time.Parse("02/01/2006", line[7])
	if err != nil {
		return
	}
	cis.BDMStatus = line[8]
	cis.EuropeanAuthorisation = line[9]
	for _, titulaire := range strings.Split(line[10], ";") {
		cis.Titulaires = append(cis.Titulaires, titulaire)
	}
	cis.ExtremeMonitoring = line[11] == "Oui"
	if err != nil {
		panic(err)
	}
	return
}
