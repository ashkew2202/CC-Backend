package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"
)

type PokemonData struct {
	ID        int
	Type      []string
	Abilities []string
	Legendary bool
	Mythical  bool
}

func main() {

	filename := flag.String("file", "default.json", "Path to the file to store the pokemon data")
	flag.Parse()

	pokemonUrl := "https://pokeapi.co/api/v2/pokemon/"
	specieUrl := "https://pokeapi.co/api/v2/pokemon-species/"
	var pokipoki PokemonData
	FinalJsonResponse := make(map[string]PokemonData)

	data, err := os.ReadFile("pokemon.txt")
	if err != nil {
		panic(err)
	}

	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		name := strings.TrimSpace(line)
		if name == "" {
			continue
		}

		// Get abilities
		resp, err := http.Get(pokemonUrl + name)
		if err != nil {
			fmt.Printf("Error fetching %s: %v\n", name, err)
			continue
		}

		var pokeData struct {
			Name      string `json:"name"`
			Abilities []struct {
				Ability struct {
					Name string `json:"name"`
				} `json:"ability"`
			} `json:"abilities"`
			Types []struct {
				Type struct {
					Name string `json:"name"`
				} `json:"type"`
			} `json:"types"`
		}

		err = json.NewDecoder(resp.Body).Decode(&pokeData)
		resp.Body.Close()
		if err != nil {
			fmt.Printf("Error decoding %s: %v\n", name, err)
			continue
		}

		var abilities []string
		for _, ab := range pokeData.Abilities {
			abilities = append(abilities, ab.Ability.Name)
		}

		var types []string
		for _, t := range pokeData.Types {
			types = append(types, t.Type.Name)
		}

		speciesResp, err := http.Get(specieUrl + name)
		if err != nil {
			fmt.Printf("Error fetching species %s: %v\n", name, err)
			continue
		}

		var speciesData struct {
			IsLegendary bool `json:"is_legendary"`
			IsMythical  bool `json:"is_mythical"`
		}

		err = json.NewDecoder(speciesResp.Body).Decode(&speciesData)
		speciesResp.Body.Close()
		if err != nil {
			fmt.Printf("Error decoding species %s: %v\n", name, err)
			continue
		}

		pokipoki = PokemonData{
			ID:        len(FinalJsonResponse) + 1, // Self-increment ID based on current map length
			Type:      types,
			Abilities: abilities,
			Legendary: speciesData.IsLegendary,
			Mythical:  speciesData.IsMythical,
		}

		FinalJsonResponse[pokeData.Name] = pokipoki

	}

	file, err := os.Create(*filename)
	if err != nil {
		fmt.Printf("Error creating file: %v\n", err)
		return
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(FinalJsonResponse); err != nil {
		fmt.Printf("Error encoding JSON to file: %v\n", err)

	}
}
