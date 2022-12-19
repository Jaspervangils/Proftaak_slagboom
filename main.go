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
	app := fiber.New() //api server

	app.Get("/:id", func(c *fiber.Ctx) error {
		param := c.Params("id")
		out := databaseFunc(param)
		return c.SendString(out)
	})
	log.Fatal(app.Listen(":3000"))
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
			return "connectie err"
		} else {
			defer db.Close()
			rows, err := db.Query("SELECT kenteken FROM `fonteyn-internal-db`.`kenteken` WHERE kenteken = \"" + p + "\"")
			if err != nil { //de query die de database ontvangt met p als ingevoerd kenteken
				resp = "error query"
			}
			var result Result
			for rows.Next() { // zet resultaat van de query in een struct
				rows.Scan(&result.KentekenResult)
			}
			if p == result.KentekenResult {
				resp = "true"
			} else {
				resp = "false"
			}
		}
	}
	return resp // returned een string met het kenteken of foutcode
}
