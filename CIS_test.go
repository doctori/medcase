package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCIS_ArrayToCIS(t *testing.T) {
	type args struct {
		line []string
	}
	tests := []struct {
		name    string
		cis     CIS
		args    args
		wantErr bool
	}{
		{name: "Hein",
			cis: CIS{},
			args: args{
				line: []string{
					"je", "ne", "sais", "pas", "trop",
				},
			},
			wantErr: true,
		},
		{name: "good_cis",
			cis: CIS{},
			args: args{
				line: []string{
					"67320915",
					"ACIDE FOLIQUE ARROW 5 mg, comprimé",
					"comprimé",
					"orale",
					"Autorisation active",
					"Procédure nationale",
					"Commercialisée",
					"29/03/2011",
					"",
					"",
					" ARROW GENERIQUES",
					"Non",
				},
			},
			wantErr: false,
		},
		{name: "cis_with_bad_date",
			cis: CIS{},
			args: args{
				line: []string{
					"67320915",
					"ACIDE FOLIQUE ARROW 5 mg, comprimé",
					"comprimé",
					"orale",
					"Autorisation active",
					"Procédure nationale",
					"Commercialisée",
					"29/33/2011",
					"",
					"",
					" ARROW GENERIQUES",
					"Non",
				},
			},
			wantErr: true,
		},
		{name: "cis_with_bad_CIS",
			cis: CIS{},
			args: args{
				line: []string{
					"NOT A VALID NUMBER",
					"ACIDE FOLIQUE ARROW 5 mg, comprimé",
					"comprimé",
					"orale",
					"Autorisation active",
					"Procédure nationale",
					"Commercialisée",
					"29/33/2011",
					"",
					"",
					" ARROW GENERIQUES",
					"Non",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.cis.ArrayToCIS(tt.args.line); (err != nil) != tt.wantErr {
				t.Errorf("CIS.ArrayToCIS() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
func TestCallArrayToCIS(t *testing.T) {
	var cis CIS
	cis.ArrayToCIS([]string{
		"61104024",
		"ACT-HIB 10 microgrammes/0,5 ml, poudre et solvant pour solution injectable en seringue préremplie. Vaccin conjugué de l'Haemophilus type b",
		"poudre et  solvant pour solution injectable",
		"intramusculaire;sous-cutanée",
		"Autorisation active",
		"Procédure nationale",
		"Commercialisée",
		"06/02/1992",
		"",
		"",
		" SANOFI PASTEUR",
		"Non",
	})
	assert.Equal(t, 61104024, cis.CIS, "They Should be Equal")
	assert.Equal(t, false, cis.ExtremeMonitoring, "string to bool test")
}

func TestLoadCIS(t *testing.T) {
	// Init the data file with a few entries
	tmpDir, err := ioutil.TempDir(".", "testLoad")
	assert.Nil(t, err)
	defer os.RemoveAll(tmpDir) // clean up
	tmpFile, err := ioutil.TempFile(tmpDir, "testLoadCIS")
	assert.Nil(t, err)
	defer tmpFile.Close()
	lines := []string{
		"62872308	ACNETRAIT 40 mg, capsule molle	capsule molle	orale	Autorisation active	Procédure nationale	Commercialisée	25/11/2009			 ARROW GENERIQUES	Non",
		"64168039	ACNETRAIT 5 mg, capsule molle	capsule molle	orale	Autorisation active	Procédure nationale	Commercialisée	20/11/2008			 ARROW GENERIQUES	Non",
		"61736892	ACONITUM FEROX BOIRON, degré de dilution compris entre 2CH et 30CH ou entre 4DH et 60DH	 comprimé et solution(s) et granules et poudre et pommade	cutanée;orale;sublinguale	Autorisation active	Enreg homéo (Proc. Nat.)	Commercialisée	23/07/2010			 BOIRON	Non",
		"64692465	ACONITUM NAPELLUS BOIRON, degré de dilution compris entre 2CH et 30CH ou entre 4DH et 60DH	 comprimé et solution(s) et granules et poudre et pommade	cutanée;orale;sublinguale	Autorisation active	Enreg homéo (Proc. Nat.)	Commercialisée	05/11/2010			 BOIRON	Non",
		"69203615	ACONITUM NAPELLUS LEHNING, degré de dilution compris entre 2CH et 30CH ou entre 4DH et 60DH	 comprimé et solution(s) et granules et poudre et pommade	cutanée;orale;sublinguale	Autorisation active	Enreg homéo (Proc. Nat.)	Commercialisée	26/10/2010			 LEHNING	Non",
		"66506742	ACONITUM NAPELLUS WELEDA, degré de dilution compris entre 2CH et 30CH ou entre 4DH et 60DH	granules et  crème et  solution en gouttes en gouttes	cutanée;orale;sublinguale	Autorisation active	Enreg homéo (Proc. Nat.)	Commercialisée	14/02/2011			 WELEDA	Non",
		"69152302	ACORSPRAY 200 microgrammes/dose, solution pour inhalation en flacon pressurisé	solution pour inhalation	inhalée	Autorisation active	Procédure nationale	Commercialisée	04/10/2005			 CHIESI	Non",
		"69515643	RIMIFON 500 mg/5 ml, solution injectable/pour perfusion	solution injectable pour perfusion	intramusculaire;intraveineuse	Autorisation active	Procédure nationale	Commercialisée	29/11/2016			 GALIEN EUROPE (CHYPRE)	Non",
	}
	for _, line := range lines {
		tmpFile.WriteString(fmt.Sprintf("%s\r\n", line))
	}
	CISs, err := LoadCIS(tmpFile.Name())
	assert.Nil(t, err)
	assert.Equal(t, len(lines), len(CISs))
	assert.Equal(t, 62872308, CISs[0].CIS)
}
