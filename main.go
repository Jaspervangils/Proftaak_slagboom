package main

import (
	"database/sql"
	"fmt"
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
	}

	db, err := sql.Open("mysql", configStruct.Database.Username+":"+configStruct.Database.Password+"@tcp("+configStruct.Host+")/csdb?tls=true")
	if err != nil { //open SQL verbinding met ww, naam en server uit conf bestand
		fmt.Println(err)
	}
	defer db.Close()
	rows, err := db.Query("SELECT * FROM `csdb`.`kenteken` WHERE KentekenID = \"" + p + "\"")
	if err != nil { //de query die de database ontvangt met p als ingevoerd kenteken
		resp = "code 0"
	} else {
		for rows.Next() { // zet resultaat van de query in een struct
			var result Result
			err := rows.Scan(&result.KentekenResult)
			if err != nil {
				resp = "code 1"
			}
			resp = result.KentekenResult
		}
	}
	return resp // returned een string met het kenteken of foutcode

}

func input() {
	app := fiber.New() //api server

	app.Get("api/:id", func(c *fiber.Ctx) error {
		param := c.Params("id")
		out := databaseFunc(param)
		return c.SendString(out)
	})
	log.Fatal(app.Listen(":3000"))
}

// var (
// 	ceifers string = "0123456789"
// 	letters string = "DFGHJKLNPRSTXYZ"
// )

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
