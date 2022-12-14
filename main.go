package main

import (
	"database/sql"
	"io/ioutil"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Host     string `yaml:"host"`
	Database struct {
		Username string `yaml:"user"`
		Password string `yaml:"pass"`
	} `yaml:"database"`
}

type Result struct {
	KentekenResult string
}

func main() {
	input()
}

func databaseFunc(p string) string {
	file, err := ioutil.ReadFile("config.yaml") // laad configuratie bestand

	if err != nil {
		panic(err)
	}
	var configStruct Config

	err = yaml.Unmarshal(file, &configStruct)
	if err != nil {
		panic(err)
	}
	var resp string  //p is een parameter die is meegegeven in de url van de API req
	if len(p) != 8 { //filter alles wat niet dezelfde lengte heeft als een kenteken er uit
		resp = "geen kenteken"
	} else {
		db, err := sql.Open("mysql", configStruct.Database.Username+":"+configStruct.Database.Password+"@tcp("+configStruct.Host+")/fonteyn-internal-db?tls=true")
		if err != nil { //open SQL verbinding met ww, naam en server uit conf bestand
			resp = "connectie err"
		} else {
			defer db.Close()
			rows, err := db.Query("SELECT kenteken FROM `fonteyn-internal-db`.`kenteken` WHERE kenteken = \"" + p + "\"")
			if err != nil { //de query die de database ontvangt met p als ingevoerd kenteken
				resp = "fout database"
			} else {
				for rows.Next() { // zet resultaat van de query in een struct
					var result Result
					err := rows.Scan(&result.KentekenResult)
					if err != nil {
						resp = "false"
					} else {
						resp = result.KentekenResult
					}
				}
			}
		}
	}
	return resp // returned een string met het kenteken of foutcode
}

func input() {
	app := fiber.New() //api server

	app.Get("/:id", func(c *fiber.Ctx) error {
		param := c.Params("id")
		out := databaseFunc(param)
		return c.SendString(out)
	})
	log.Fatal(app.Listen(":3000"))
}

// func test2() {
// 	app := fiber.New() //api server

// 	app.Get("/api", func(c *fiber.Ctx) error {
// 		out := test()
// 		return c.SendString(out)
// 	})
// 	log.Fatal(app.Listen(":3000"))
// }

// func test() string {
// 	file, err := ioutil.ReadFile("config.yaml") // laad configuratie bestand

// 	if err != nil {
// 		panic(err)
// 	}
// 	var configStruct Config

// 	err = yaml.Unmarshal(file, &configStruct)
// 	if err != nil {
// 		panic(err)
// 	}

// 	var resp string

// 	db, err := sql.Open("mysql", configStruct.Database.Username+":"+configStruct.Database.Password+"@tcp("+configStruct.Host+")/csdb?tls=true")
// 	if err != nil { //open SQL verbinding met ww, naam en server uit conf bestand
// 		resp = "connectie err"
// 	}
// 	defer db.Close()
// 	rows, err := db.Query("SELECT kenteken FROM `fonteyn-internal-db`.`kenteken` WHERE kenteken = \"12345678\"")
// 	if err != nil { //de query die de database ontvangt met p als ingevoerd kenteken
// 		resp = "fout database"
// 		fmt.Println(err)
// 	}
// 	fmt.Println(rows)
// 	for rows.Next() { // zet resultaat van de query in een struct
// 		var result Result
// 		err := rows.Scan(&result.KentekenResult)
// 		if err != nil {
// 			resp = "false"
// 		} else {
// 			resp = result.KentekenResult
// 		}
// 	}
// 	return resp
// }
