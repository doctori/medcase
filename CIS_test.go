package main

import (
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
