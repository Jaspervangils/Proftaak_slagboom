package main

import (
	"database/sql"
	"io/ioutil"

	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/yaml.v2"
)

var (
	ceifers string = "0123456789"
	letters string = "DFGHJKLNPRSTXYZ"
)

type Config struct {
	Database struct {
		Username string `yaml:"user"`
		Password string `yaml:"pass"`
		server   string `yaml:"serv"`
	} `yaml:"database"`
}

func main() {
	//fmt.Println(kentekenGen())
	databaseFunc()
}

func databaseFunc() {
	file, err := ioutil.ReadFile("config.yaml")

	if err != nil {
		panic(err)
	}
	var ConfigStruct Config

	err = yaml.Unmarshal(file, &ConfigStruct)
	if err != nil {
		panic(err)
	}

	db, err := sql.Open("mysql", ConfigStruct.Database.Username+":"+ConfigStruct.Database.Password+"@tcp(capitaselectadb.mysql.database.azure.com:3306)/csdb?tls=true")
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()
	rows, err := db.Query("INSERT INTO `csdb`.`kenteken` (`KentekenID`) VALUES ('testets')") // select + from db
	if err != nil {
		fmt.Println(err)
	}
	defer rows.Close()
}

// func kentekenGen() string {
// 	rand.Seed(time.Now().Unix())

// 	var kenteken strings.Builder

// 	random := rand.Intn(len(letters))
// 	kenteken.WriteString(string(letters[random]))
// 	kenteken.WriteString("-")

// 	for i := 0; i < 3; i++ {
// 		randomC := rand.Intn(len(ceifers))
// 		kenteken.WriteString(string(ceifers[randomC]))
// 	}
// 	kenteken.WriteString("-")
// 	for i := 0; i < 2; i++ {
// 		randomC := rand.Intn(len(letters))
// 		kenteken.WriteString(string(letters[randomC]))
// 	}

// 	return kenteken.String()
// }
