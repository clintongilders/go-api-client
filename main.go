package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"

	"github.com/clintongilders/go-api-client/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var err error
var DB *gorm.DB

type RegionResponse struct {
	Count  int      `json:"count"`
	Region []Region `json:"results"`
}

type Region struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

// A Response struct to map the Entire Response
type PokeResponse struct {
	Name    string    `json:"name"`
	Pokemon []Pokemon `json:"pokemon_entries"`
}

// A Pokemon Struct to map every Pokemon to.
type Pokemon struct {
	EntryNo int            `json:"entry_number"`
	Species PokemonSpecies `json:"pokemon_species"`
}

// A struct to map our Pokemon's Species, which includes its name
type PokemonSpecies struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

func InitDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("/tmp/test.db?_foreign_keys=on"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	} else {
		println("Database connection successful.")
	}
	// Migrate the schema
	db.AutoMigrate(&models.Region{}, &models.PokemonSpecies{})
	println("Database Migrated")
	return db
}

func main() {
	db := InitDB()
	// Get the region List
	response, err := http.Get("https://pokeapi.co/api/v2/pokedex/?limit=50")
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	var responseObject RegionResponse
	json.Unmarshal(responseData, &responseObject)

	fmt.Println("Total Records: ", responseObject.Count)
	re := regexp.MustCompile(`\/(\d{1,5})\/$`)
	for i := 0; i < len(responseObject.Region); i++ {
		//for i := 0; i < 1; i++ {

		fmt.Println(responseObject.Region[i].Name, " ", responseObject.Region[i].URL)
		// Extract the ID from the URL using regex
		matches := re.FindStringSubmatch(responseObject.Region[i].URL)

		if len(matches) > 0 {
			// Create Region in DB
			num, err := strconv.Atoi(matches[1])
			if err != nil {
				log.Fatal(err)
			}

			region := models.CreateRegion(db, num, responseObject.Region[i].Name)

			// Get the Pokemon List for this Region
			pokeUrl := fmt.Sprintf("https://pokeapi.co/api/v2/pokedex/%d/", region.RegionId)
			pokeResponse, err := http.Get(pokeUrl)

			if err != nil {
				fmt.Print(err.Error())
				os.Exit(1)
			}
			pokeResponseData, err := io.ReadAll(pokeResponse.Body)
			if err != nil {
				log.Fatal(err)
			}
			var pokeResponseObject PokeResponse
			json.Unmarshal(pokeResponseData, &pokeResponseObject)
			for i := 0; i < len(pokeResponseObject.Pokemon); i++ {
				// Extract the ID from the URL using regex
				pokeMatches := re.FindStringSubmatch(pokeResponseObject.Pokemon[i].Species.URL)

				if len(pokeMatches) > 0 {
					pokeNum, err := strconv.Atoi(pokeMatches[1])
					fmt.Println("Full Match:", pokeMatches[0])
					if err != nil {
						log.Fatal(err)
					}
					// Create Pokemon Species in DB
					pokeResult := models.CreatePokemonSpecies(db, pokeNum, pokeResponseObject.Pokemon[i].Species.Name, region)
					fmt.Printf("Created Pokemon Species: ID %d, Name: %s, Region: %s\n", pokeResult.PokemonId, pokeResult.PokemonName, region.RegionName)
				}
			}

		}
	}
}
