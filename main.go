package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"

	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/yaml.v2"
)

var (
	ceifers string = "0123456789"
	letters string = "DFGHJKLNPRSTXYZ"
)

type Config struct {
	Host     string `yaml:"host"`
	Database struct {
		Username string `yaml:"user"`
		Password string `yaml:"pass"`
	} `yaml:"database"`
}

type result struct {
	KentekenResult string
}

func main() {
	databaseFunc()
}

func databaseFunc() {
	file, err := ioutil.ReadFile("config.yaml")

	if err != nil {
		panic(err)
	}
	var configStruct Config

	err = yaml.Unmarshal(file, &configStruct)
	if err != nil {
		panic(err)
	}

	db, err := sql.Open("mysql", configStruct.Database.Username+":"+configStruct.Database.Password+"@tcp("+configStruct.Host+")/csdb?tls=true")
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()
	rows, err := db.Query("SELECT * FROM `csdb`.`kenteken` ORDER BY `KentekenID`")
	if err != nil {
		fmt.Println(err)
	}
	defer rows.Close()

	for rows.Next() {
		var result result
		err := rows.Scan(&result.KentekenResult)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(result.KentekenResult[1])
	}
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
