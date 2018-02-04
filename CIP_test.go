package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCIP_ArrayToCIP(t *testing.T) {
	type args struct {
		line []string
	}
	tests := []struct {
		name    string
		cip     CIP
		args    args
		wantErr bool
	}{
		{name: "Hein",
			cip: CIP{},
			args: args{
				line: []string{
					"je", "ne", "sais", "pas", "trop",
				},
			},
			wantErr: true,
		}, {
			name: "medoc_oui",
			cip:  CIP{},
			args: args{
				line: []string{
					"60026957",
					"3875331",
					"plaquette(s) thermoformée(s) PVC-Aluminium de 180 comprimé(s)",
					"Présentation active",
					"Déclaration de commercialisation",
					"25/11/2008",
					"3400938753318",
					"oui",
					"65%",
					"11,12",
					"12,14",
					"1,02",
				},
			},
			wantErr: false,
		}, {
			name: "medoc_oui_with_bad_date",
			cip:  CIP{},
			args: args{
				line: []string{
					"60026957",
					"3875331",
					"plaquette(s) thermoformée(s) PVC-Aluminium de 180 comprimé(s)",
					"Présentation active",
					"Déclaration de commercialisation",
					"25/11/3992008",
					"3400938753318",
					"oui",
					"65%",
					"11,12",
					"12,14",
					"1,02",
				},
			},
			wantErr: true,
		},
		{
			name: "medoc_oui_with_bad_date_percentage",
			cip:  CIP{},
			args: args{
				line: []string{
					"60026957",
					"3875331",
					"plaquette(s) thermoformée(s) PVC-Aluminium de 180 comprimé(s)",
					"Présentation active",
					"Déclaration de commercialisation",
					"25/11/2008",
					"3400938753318",
					"oui",
					"650%",
					"11,12",
					"12,14",
					"1,02",
				},
			},
			wantErr: true,
		},
		{
			name: "medoc_oui_with_no_remoursement",
			cip:  CIP{},
			args: args{
				line: []string{
					"60026494",
					"3000753",
					"plaquette(s) thermoformée(s) PVC PVDC aluminium de 21 comprimé(s)",
					"Présentation active",
					"Déclaration de commercialisation",
					"01/11/2017",
					"3400930007532",
					"non",
					"",
					"",
				},
			},
			wantErr: false,
		},
		{
			name: "medoc_oui_with_bad_CIP",
			cip:  CIP{},
			args: args{
				line: []string{
					"60026494",
					"300075BAD3",
					"plaquette(s) thermoformée(s) PVC PVDC aluminium de 21 comprimé(s)",
					"Présentation active",
					"Déclaration de commercialisation",
					"01/11/2017",
					"3400930007532",
					"non",
					"",
					"",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.cip.ArrayToCIP(tt.args.line); (err != nil) != tt.wantErr {
				t.Errorf("CIP.ArrayToCIP() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCallArrayToCIP(t *testing.T) {
	var cip CIP
	cip.ArrayToCIP([]string{
		"60026957",
		"3875331",
		"plaquette(s) thermoformée(s) PVC-Aluminium de 180 comprimé(s)",
		"Présentation active",
		"Déclaration de commercialisation",
		"25/11/2008",
		"3400938753318",
		"oui",
		"65%",
		"11,12",
		"12,14",
		"1,02",
	})
	assert.Equal(t, 60026957, cip.CIS, "They Should be Equal")
	assert.Equal(t, true, cip.CollectivityAgreement, "string to bool test")
}
