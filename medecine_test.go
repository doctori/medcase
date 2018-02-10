package main

import (
	"testing"

	"github.com/doctori/medcase/BDPM"
	"github.com/jinzhu/gorm"
)

func setup() *gorm.DB {
	config := Config{
		DB: DBConfig{
			Host:     "localhost",
			Username: "medcase_tests",
			Password: "medcase_tests",
			Port:     5432,
			DBname:   "medcase_tests",
		},
	}
	return initDB(config)
}
func tearDown(db *gorm.DB) {
	db.DropTableIfExists(&Medecine{}, &Presentation{})
}
func Test_medecine(t *testing.T) {
	// Setup the database
	db := setup()
	type args struct {
		cis  *BDPM.CIS
		cips []BDPM.CIP
	}
	tests := []struct {
		name    string
		args    args
		med     Medecine
		wantErr bool
	}{
		{
			name: "Load Fails becauseEmpty",
			med:  Medecine{},
			args: args{
				cis: &BDPM.CIS{
					CIS:  0,
					Name: "",
				},
				cips: []BDPM.CIP{
					{
						CIS:   0,
						Label: "",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "MinimalSucceed",
			med:  Medecine{},
			args: args{
				cis: &BDPM.CIS{
					CIS:  69811037,
					Name: "ACTI 5",
				},
				cips: []BDPM.CIP{
					{
						CIS:   69811037,
						CIP7:  3343597,
						CIP13: 3400933435974,
						Label: "30 ampoule(s) en verre brun de 5  ml",
					},
				},
			},
			wantErr: false,
		}, {
			name: "DfferentCIS",
			med:  Medecine{},
			args: args{
				cis: &BDPM.CIS{
					CIS:  69811037,
					Name: "ACTI 5",
				},
				cips: []BDPM.CIP{
					{
						CIS:   69811039,
						CIP7:  3343597,
						CIP13: 3400933435974,
						Label: "30 ampoule(s) en verre brun de 5  ml",
					},
				},
			},
			wantErr: true,
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.med.loadMedecineFromBDPM(tt.args.cis, tt.args.cips); (err != nil) != tt.wantErr {
				t.Errorf("CIP.ArrayToCIP() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
	tearDown(db)
}
