package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/doctori/medcase/BDPM"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func check(e error) {
	if e != nil {
		log.Panic(e)
	}
}

// Application Configuration

// DBConfig support the database connection parameters
// within the main confguration
type DBConfig struct {
	Host     string `json:"host"`
	Username string `json:"username"`
	Password string `json:"password"`
	Port     int    `json:"port"`
	DBname   string `json:"dbname"`
}

// Config is the application configuration model
type Config struct {
	DB DBConfig `json:"db"`
}

var db = &gorm.DB{}

func main() {
	// Load Application configuration
	file, err := os.Open("conf/mainConfig.json")
	if err != nil {
		log.Fatal(err)
	}
	var config Config
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	log.Printf("%#v", config)
	if err != nil {
		log.Fatal(err)
	}

	initDB(config)
	defer db.Close()
	CISs, err := BDPM.LoadCIS(BDPM.CisDataFile)
	check(err)
	CIPs, err := BDPM.LoadCIP(BDPM.CipDataFile)
	check(err)
	fmt.Printf("loaded : %d entries\n", len(CISs)+len(CIPs))
	var associatedCIPs []BDPM.CIP
	for _, cis := range CISs {
		associatedCIPs = []BDPM.CIP{}
		// Try to find all the associated CIPs
		for _, cip := range CIPs {
			if cip.CIS == cis.CIS {
				//		fmt.Printf("Adding %d for medecine [%s][%d]\n", cip.CIS, cis.Name, cis.CIS)
				associatedCIPs = append(associatedCIPs, cip)
			}
		}
		// create the medecine definition
		var med Medecine
		err = med.loadMedecineFromBDPM(cis, associatedCIPs)
		if err != nil {
			panic(err)
		}
		//fmt.Printf("Loaded med : %#v\n", med)
		//med.Save()

	}
	//quick fetch test
	med := GetMedByPresShortID(3541759)
	fmt.Printf("fetched :%s\n", med.Name)
	fmt.Printf("\t with %d Different Presentations\n", len(med.Presentations))

}

func initDB(config Config) *gorm.DB {
	connectionString := fmt.Sprintf(
		"user=%s password='%s' host=%s dbname=%s",
		config.DB.Username,
		config.DB.Password,
		config.DB.Host,
		config.DB.DBname)
	log.Printf("Connecting to %s", connectionString)
	db, err := gorm.Open("postgres", connectionString)
	check(err)
	// Debug Mode
	db.LogMode(true)
	db.CreateTable(&Medecine{}, &Presentation{})
	//	db.Model(&Medecine{}).AddUniqueIndex("medecine_uniqueness", "national_pres_long_code")
	db.AutoMigrate(&Medecine{}, &Presentation{})
	check(err)

	return db
}
