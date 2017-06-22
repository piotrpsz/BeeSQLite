# BeeSQLite
With this project you can use SQLite library direct from your's Go code

# Example: How to use
```Go
package main

import (
	"BeeSQLite"
	"fmt"
	"log"
	"os"
	"os/user"
)

var sqlm BeeSQLite.SQLite

type City struct {
	ID          int
	Name        string
	CountryName string
}

var cmd = `CREATE TABLE city
(
	id           INTEGER  NOT NULL PRIMARY KEY AUTOINCREMENT,
	name         TEXT     NOT NULL COLLATE NOCASE,
	country_name TEXT     NOT NULL COLLATE NOCASE,

	CONSTRAINT city_name         CHECK  (name!=''),
	CONSTRAINT city_country_name CHECK  (country_name!=''),
	CONSTRAINT city_unique       UNIQUE (country_name, name)
)`

func (c *City) Save() bool {
	var fields []BeeSQLite.Field

	if c.ID > 0 {
		id := BeeSQLite.Field{Name: "id"}
		id.SetValue(c.ID)
		fields = append(fields, id)
	}
	name := BeeSQLite.Field{Name: "name"}
	name.SetValue(c.Name)
	fields = append(fields, name)

	country := BeeSQLite.Field{Name: "country_name"}
	country.SetValue(c.CountryName)
	fields = append(fields, country)

	if c.ID == 0 {
		id, ok := sqlm.Insert("city", fields)
		if ok {
			c.ID = id
		}
		return ok
	}
	return sqlm.Update("city", fields)
}

func (c *City) Init(data BeeSQLite.Row) *City {
	id, ok := data["id"]
	if ok {
		c.ID = id.Int()
	}

	name, ok := data["name"]
	if ok {
		c.Name = name.String()
	}

	countryName, ok := data["country_name"]
	if ok {
		c.CountryName = countryName.String()
	}
	return c
}

func (c *City) String() string {
	return fmt.Sprintf("(ID: %d, name: %s, country: %s)", c.ID, c.Name, c.CountryName)
}

func main() {
	usr, _ := user.Current()
	fpath := usr.HomeDir + "/test.sqlite"
	os.Remove(fpath)

	ok := sqlm.Create(fpath)
	if !ok {
		log.Fatal()
	}
	fmt.Printf("Database created: %s\n", fpath)

	ok = sqlm.ExecQuery(cmd)
	if !ok {
		log.Fatal()
	}
	fmt.Println("Succes: table city created")

	//-----------------------------------------------

	var poznan = City{Name: "Poznań", CountryName: "Poland"}
	fmt.Printf("Record before save: %s\n", poznan.String())
	ok = poznan.Save()
	if !ok {
		log.Fatal()
	}
	fmt.Println("Succes: record inserted")
	fmt.Printf("Record after insert: %s\n", poznan.String())

	poznan.CountryName = "Polska"
	ok = poznan.Save()
	if !ok {
		log.Fatal()
	}
	fmt.Println("Succes: record updated")
	fmt.Printf("Record after update: %s\n", poznan.String())

	//-------------------------------------------------------

	var przemysl = City{Name: "Przemyśl", CountryName: "Poland"}
	ok = przemysl.Save()
	if !ok {
		log.Fatal()
	}
	fmt.Printf("Record after insert: %s\n", przemysl.String())

	//-------------------------------------------------

	var gdansk = City{Name: "Gdańsk", CountryName: "Poland"}
	ok = gdansk.Save()
	if !ok {
		log.Fatal()
	}
	fmt.Printf("Record after insert: %s\n", gdansk.String())

	//-------------------------------------------------

	query := fmt.Sprintf("SELECT * FROM City WHERE id=%d", 1)
	result, ok := sqlm.Select(query)

	if ok {
		if row := result.First(); row != nil {
			city := (&City{}).Init(row)
			fmt.Printf("Record after select: %s\n", city.String())
		}
	}

	fmt.Println("------- SELECT ALL ----------------")

	query = fmt.Sprintf("SELECT * FROM City")
	result, ok = sqlm.Select(query)
	if ok {
		if result.IsNotEmpty() {
			for index, row := range result {
				city := (&City{}).Init(row)
				fmt.Printf("Record %d: %s\n", index, city.String())
			}
		}
	}

}
```