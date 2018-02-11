package main

import (
	"fmt"

	"github.com/doctori/medcase/BDPM"
	"github.com/jinzhu/gorm"
)

// Medecine represent the summary of a medecine,
// Extracted from opendata databases ?
// Doesn't have an owner or expiration date
// Just a representation of the medcine and it's representations (boxes ?)
type Medecine struct {
	gorm.Model
	Name              string
	NationalShortCode int `gorm:"primary_key"`
	Presentations     []Presentation
}

// Presentation is the container of the medecine,
// It should hold specific identifiers
type Presentation struct {
	gorm.Model
	MedecineID            int
	Label                 string
	NationalPresShortCode int
	//NationalPresLongCode  uint64 `gorm:"primary_key"`
	Price float32
}

// We need to save this in a database !
// We need to manage database connection at an upper level !

func (med *Medecine) loadMedecineFromBDPM(cis *BDPM.CIS, cips []BDPM.CIP) (err error) {
	var presentations []Presentation
	if cis.IsNil() {
		err = fmt.Errorf("Could not Load Medecine because cis appears to be Nil : [%#v]", cis)
		return
	}
	for _, cip := range cips {
		if cis.CIS != cip.CIS {
			err = fmt.Errorf("the CIP and CIS are diferent drugs : %d and %d", cip.CIS, cis.CIS)
			return
		}
		presentations = append(presentations, Presentation{
			Label: cip.Label,
			NationalPresShortCode: cip.CIP7,
			//			NationalPresLongCode:  cip.CIP13,
			//Price:                 cip.Prices[0],
		})
	}

	med.Name = cis.Name
	med.NationalShortCode = cis.CIS
	med.Presentations = presentations
	//fmt.Printf("med name is %s\n", med.Name)

	return
}

// IsNil will return true if the Medecine Struct is considered Null
func (med *Medecine) IsNil() bool {
	if med.Name == "" || med.NationalShortCode == 0 {
		return true
	}
	return false
}

// Save Will save the medecine struct into the database
func (med *Medecine) Save() (Medecine, error) {
	if med.IsNil() {
		err := fmt.Errorf("Could not save Empty Medecine : %v", *med)
		return *med, err
	}
	if db.NewRecord(med) {
		oldMed := new(Medecine)
		db.Where("national_short_code = ?", med.NationalShortCode).First(&oldMed)
		if oldMed.Name == "" {
			//log.Println("Recording the New Medecine")
			db.Create(&med)
		} else {
			//log.Println("Updating The Existing  Medecine")
			db.Model(&oldMed).Updates(&med)
			med = oldMed
			//log.Printf("Saving Medecine %#v", med)
		}
	} else {
		//log.Printf("Creating Medecine %#v", med)
		db.Save(&med)
	}
	return *med, nil
}

//GetMedByPresShortID will fetch the medecine associated with
// this specific presentatin short ID
func GetMedByPresShortID(id int) (med Medecine) {
	// Init Presentation
	pres := Presentation{
		NationalPresShortCode: id,
	}
	//Get presenations ...
	db.Where(&pres).First(&pres)
	//fmt.Printf("Presentation found : %#v\n", pres)
	db.Preload("Presentations").First(&med, pres.MedecineID)
	return
}
