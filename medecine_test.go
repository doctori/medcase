package main

import (
	"testing"

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
	tests := []struct {
		name string
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			main()
		})
	}
	tearDown(db)
}
