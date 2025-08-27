package models

import (
	"gorm.io/gorm"
)

var err error

type Region struct {
	gorm.Model
	RegionId       int    `gorm:"uniqueIndex"`
	RegionName     string `gorm:"index"`
	PokemonSpecies []PokemonSpecies
}

type PokemonSpecies struct {
	gorm.Model
	RegionId    uint   `gorm:"uniqueIndex:idx_combined_unique"`
	PokemonId   int    `gorm:"uniqueIndex:idx_combined_unique"`
	PokemonName string `gorm:"index"`
}

func CreateRegion(db *gorm.DB, rid int, rname string) Region {
	if err != nil {
		panic("failed to connect database")
	} else {
		println("Database connection successful.")
	}

	var region Region
	result := db.Where(Region{RegionId: rid}).FirstOrCreate(&region, Region{RegionName: rname, RegionId: rid})

	if result.Error != nil {
		println("Error creating region:", result.Error.Error())
	} else {
		println("Region created successfully with ID:", region.RegionId)
	}

	return region
}

func CreatePokemonSpecies(db *gorm.DB, pid int, pname string, region Region) PokemonSpecies {
	if err != nil {
		panic("failed to connect database")
	} else {
		println("Database connection successful.")
	}
	var pokemon PokemonSpecies
	result := db.Where(PokemonSpecies{PokemonId: pid, RegionId: region.ID}).FirstOrCreate(&pokemon, PokemonSpecies{PokemonId: pid, PokemonName: pname, RegionId: region.ID})
	if result.Error != nil {
		println("Error creating pokemon species:", result.Error.Error())
	} else {
		println("Pokemon species created successfully with ID:", pokemon.PokemonId)
	}
	return pokemon
}
