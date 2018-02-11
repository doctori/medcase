package main

import (
	"testing"

	"github.com/doctori/medcase/BDPM"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type MedecineTestSuite struct {
	suite.Suite
}

func setup() {
	config := Config{
		DB: DBConfig{
			Host:     "localhost",
			Username: "medcase_tests",
			Password: "medcase_tests",
			Port:     5432,
			DBname:   "medcase_tests",
		},
	}
	initDB(config)

}
func tearDown() {
	db.DropTableIfExists(&Medecine{}, &Presentation{})
	db.Close()
}
func (suite *MedecineTestSuite) SetupTest() {
	setup()
}

func (suite *MedecineTestSuite) TearDownTest() {
	tearDown()
}

func TestMedecine_IsNil(t *testing.T) {
	tests := []struct {
		name   string
		med    Medecine
		rValue bool
	}{
		{
			name:   "simpleNilCIS",
			med:    Medecine{},
			rValue: true,
		}, {
			name: "simpleNonNilSerial",
			med: Medecine{
				Name:              "I'm a Medecine",
				NationalShortCode: 33403495,
			},
			rValue: false,
		}, {
			name: "MixedWithNoSerial",
			med: Medecine{
				Name: "SuperMedoc2000",
			},
			rValue: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if ret := tt.med.IsNil(); ret != tt.rValue {
				t.Errorf("Medecine.IsNil wanted [%v] got : [%v]", tt.rValue, ret)
			}
		})
	}
}

func (suite *MedecineTestSuite) TestBullshit() {
	med := Medecine{
		Name:              "Medecine 3000",
		NationalShortCode: 60002746,
	}
	med, err := med.Save()
	assert.Nil(suite.T(), err)
}
func TestMedecineTestSuite(t *testing.T) {
	suite.Run(t, new(MedecineTestSuite))
}
func TestMedecine_Save(t *testing.T) {
	setup()
	tests := []struct {
		name    string
		med     Medecine
		rValue  Medecine
		wantErr bool
	}{
		// Are we in parrallel ?
		{
			name:    "bullshit",
			med:     Medecine{},
			wantErr: true,
		},
		{
			name: "Simple Medecine",
			med: Medecine{
				Name:              "Medecine 3000",
				NationalShortCode: 60002746,
			},
			wantErr: false,
			rValue: Medecine{
				Name:              "Medecine 3000",
				NationalShortCode: 60002746,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ret, err := tt.med.Save()
			if (err != nil) != tt.wantErr {
				t.Errorf("Medecine.Save() err = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil {
				if ret.NationalShortCode != tt.rValue.NationalShortCode {
					t.Errorf("Result not Expected : Wanted : %v, got : %v", tt.rValue, ret)
				}
			}
		})

	}
	tearDown()
}
func TestMedecine_loadMedecineFromBDPM(t *testing.T) {
	// Setup the database
	setup()
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
		// TODO: Add test moars cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.med.loadMedecineFromBDPM(tt.args.cis, tt.args.cips); (err != nil) != tt.wantErr {
				t.Errorf("Medecine.loadMedecineFromBDPM() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
	tearDown()
}
