package BDPM

import (
	"fmt"
	"io/ioutil"
	"os"
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
			name: "medoc_oui_with_bad_CIS",
			cip:  CIP{},
			args: args{
				line: []string{
					"NOTANUMBER",
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
func TestLoadCIP(t *testing.T) {
	// Init the data file with a few entries
	tmpDir, err := ioutil.TempDir(".", "testLoad")
	assert.Nil(t, err)
	defer os.RemoveAll(tmpDir) // clean up
	tmpFile, err := ioutil.TempFile(tmpDir, "testLoadCIP")
	assert.Nil(t, err)
	defer tmpFile.Close()
	lines := []string{
		"65141739	5003049	1 pot(s) en verre de 6,5 g	Présentation active	Arrêt de commercialisation (le médicament n'a plus d'autorisation)	11/08/2014	3400950030497	non					",
		"65141760	3619243	pilulier(s) polypropylène de 60  comprimé(s)	Présentation active	Déclaration d'arrêt de commercialisation	09/09/2014	3400936192430	non	30%	7,75	8,77	1,02	",
		"65142848	3921516	plaquette(s) thermoformée(s) PVC polychlortrifluoroéthylène aluminium de 30 comprimé(s)	Présentation active	Déclaration d'arrêt de commercialisation	07/10/2017	3400939215167	non					",
		"65142848	5745377	56 plaquette(s) thermoformée(s) PVC polychlortrifluoroéthylène aluminium de 1 comprimé(s)	Présentation active	Déclaration d'arrêt de commercialisation	31/12/2014	3400957453770	non					",
		"65143010	3006566	plaquette(s) PVC aluminium de 30 comprimé(s)	Présentation active	Déclaration de commercialisation	19/04/2017	3400930065662	oui	65%	1,13	2,15	1,02	",
		"65143010	3006568	plaquette(s) PVC aluminium de 90 comprimé(s)	Présentation active	Déclaration de commercialisation	22/05/2017	3400930065686	oui	65%	2,89	3,91	1,02	",
		"65143010	3237197	plaquette(s) PVC aluminium de 50 comprimé(s) ( abrogée le 10/10/2017)	Présentation abrogée	Déclaration d'arrêt de commercialisation	24/03/2014	3400932371976	non					",
	}
	for _, line := range lines {
		tmpFile.WriteString(fmt.Sprintf("%s\r\n", line))
	}
	CIPs, err := LoadCIP(tmpFile.Name())
	assert.Nil(t, err)
	assert.Equal(t, len(lines), len(CIPs))
	assert.Equal(t, 65143010, CIPs[3006566].CIS)
	assert.Equal(t, uint64(3400930065662), CIPs[3006566].CIP13)
	// Load Bullshit file
	_, err = LoadCIP("Not_existing_file")
	assert.NotNil(t, err)
}
