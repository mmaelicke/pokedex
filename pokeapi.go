package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
)

type LocationsResponse struct {
	Count    int        `json:"count"`
	Next     string     `json:"next"`
	Previous string     `json:"previous"`
	Results  []Location `json:"results"`
}

type LocationResponse struct {
	Location          Location            `json:"location"`
	PokemonEncounters []PokemonEncounters `json:"pokemon_encounters"`
}

// this is not fully mapped!
type PokemonResponse struct {
	Abilities      []Abilities `json:"abilities"`
	BaseExperience int         `json:"base_experience"`
	Name           string      `json:"name"`
	Height         int         `json:"height"`
	Weight         int         `json:"weight"`
	Stats          []Stats     `json:"stats"`
	Types          []Types     `json:"types"`
}
type Abilities struct {
	Ability  Ability `json:"ability"`
	IsHidden bool    `json:"is_hidden"`
	Slot     int     `json:"slot"`
}
type Ability struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}
type Types struct {
	Slot int  `json:"slot"`
	Type Type `json:"type"`
}
type Type struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}
type Stats struct {
	BaseStat int  `json:"base_stat"`
	Effort   int  `json:"effort"`
	Stat     Stat `json:"stat"`
}
type Stat struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type Location struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}
type Pokemon struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}
type PokemonEncounters struct {
	Pokemon Pokemon `json:"pokemon"`
}

func navigatePokeAPI(url string) (LocationsResponse, error) {
	// before we get the response, we check if there is a cached version
	var locations LocationsResponse

	cached, ok := Config.cache.Get(url)
	if !ok {
		// make the request
		res, err := http.Get(url)
		if err != nil {
			return LocationsResponse{}, fmt.Errorf("could not get the next page in PokeAPI: %v", err)
		}
		defer res.Body.Close()

		// get the body
		decoder := json.NewDecoder(res.Body)
		err = decoder.Decode(&locations)
		if err != nil {
			return LocationsResponse{}, fmt.Errorf("could not decode the response from PokeAPI: %v", err)
		}

		// cache the response
		body, err := json.Marshal(locations)
		if err != nil {
			return LocationsResponse{}, fmt.Errorf("could not store the response from PokeAPI into cache: %v", err)
		}
		Config.cache.Add(url, body)
	} else {
		// decode the cached response
		err := json.Unmarshal(cached, &locations)
		if err != nil {
			return LocationsResponse{}, fmt.Errorf("could not decode the cached response from PokeAPI: %v", err)
		}
	}

	return locations, nil
}

func mapCallback(Conf *config, _ []string) error {
	locations, err := navigatePokeAPI(Conf.nextUrl)
	if err != nil {
		return err
	}

	// print out all the results
	for _, location := range locations.Results {
		fmt.Println(location.Name)
	}

	// there were no errors, so we update the config
	Conf.nextUrl = locations.Next
	Conf.previousUrl = locations.Previous

	return nil
}

func mapbCallback(Conf *config, _ []string) error {
	if Conf.previousUrl == "" {
		fmt.Println("you're on the first page.")
		return nil
	}

	// get the last page
	locations, err := navigatePokeAPI(Conf.previousUrl)
	if err != nil {
		return err
	}

	// print out all the results
	for _, location := range locations.Results {
		fmt.Println(location.Name)
	}

	// there were no errors, so we update the config
	Conf.nextUrl = locations.Next
	Conf.previousUrl = locations.Previous

	return nil
}

func exploreCallback(Conf *config, args []string) error {
	// at first make sure, that there is exactly one argument
	if len(args) != 1 {
		return fmt.Errorf("explore needs exactly one argument: <area>")
	}

	// build the url
	url := Conf.baseUrl + "location-area/" + args[0]

	// get the data from the cache
	var location LocationResponse
	cached, ok := Conf.cache.Get(url)
	if !ok {
		// make the request
		res, err := http.Get(url)
		if err != nil {
			return fmt.Errorf("the location %v seems to be not valid: %v", args[0], err)
		}
		defer res.Body.Close()

		// get the body
		decoder := json.NewDecoder(res.Body)
		err = decoder.Decode(&location)
		if err != nil {
			return fmt.Errorf("could not decode response from PokeAPI: %v", err)
		}

		// cache the response
		body, err := json.Marshal(location)
		if err != nil {
			return fmt.Errorf("could not cache the response for %v: %v ", args[0], err)
		}
		Conf.cache.Add(url, body)
	} else {
		// we got a cached response and need to parse it
		err := json.Unmarshal(cached, &location)
		if err != nil {
			return fmt.Errorf("cached response for %v seems to be broken: %v", url, err)
		}
	}

	// print out the Pokemon
	fmt.Printf("Exploring %s...\n", args[0])
	fmt.Println("Found Pokemon:")
	for _, encounter := range location.PokemonEncounters {
		fmt.Printf(" - %s\n", encounter.Pokemon.Name)
	}
	return nil
}

func catchCallback(Conf *config, args []string) error {
	// make sure there is exactly one argument
	if len(args) != 1 {
		return fmt.Errorf("catch needs exactly one argument: <pokemon>")
	}

	// build the url
	url := Conf.baseUrl + "/pokemon/" + args[0]
	var pokemon PokemonResponse
	cached, ok := Conf.cache.Get(url)
	if !ok {
		// make the request
		res, err := http.Get(url)
		if err != nil {
			return fmt.Errorf("the pokemon %v seems to be invalid: %v", args[0], err)
		}
		defer res.Body.Close()

		// get the body
		decoder := json.NewDecoder(res.Body)
		err = decoder.Decode(&pokemon)
		if err != nil {
			return fmt.Errorf("could not decode the response from PokeAPI: %v", err)
		}

		// cache the response for subsequent usages
		body, err := json.Marshal(pokemon)
		if err != nil {
			return fmt.Errorf("could not cache the response for %v: %v", args[0], err)
		}
		Conf.cache.Add(url, body)
	} else {
		err := json.Unmarshal(cached, &pokemon)
		if err != nil {
			return fmt.Errorf("cached response for %v seems to be broken: %v", url, err)
		}
	}
	// throw the Pokeball
	fmt.Printf("Throwing a Pokeball at %s...\n", args[0])

	// generate a random number between 0 and 100
	chance := rand.Intn(150)
	if chance > pokemon.BaseExperience {
		fmt.Printf("You caught %s: %v > %v\n", args[0], chance, pokemon.BaseExperience)
		// add the pokemon to the pokedex
		Conf.pokedex[args[0]] = pokemon
	} else {
		fmt.Printf("You missed %s: %v < %v\n", args[0], chance, pokemon.BaseExperience)
	}

	return nil
}

func inspectCallback(Conf *config, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("inspect needs exactly one argument: <pokemon>")
	}

	// grab the Pokemon info out of the pokedex
	pokemon, ok := Conf.pokedex[args[0]]
	if !ok {
		fmt.Printf("you have not yet caught %v\n", args[0])
		return nil
	}

	// print out all necessary info
	fmt.Printf("Name: %v\n", pokemon.Name)
	fmt.Printf("Height: %v\n", pokemon.Height)
	fmt.Printf("Weight: %v\n", pokemon.Weight)

	fmt.Println("Stats:")
	for _, stat := range pokemon.Stats {
		fmt.Printf("  -%v: %v\n", stat.Stat.Name, stat.BaseStat)
	}

	fmt.Println("Types:")
	for _, typ := range pokemon.Types {
		fmt.Printf("  -%v\n", typ.Type.Name)
	}

	return nil
}

func pokedexCallback(Conf *config, args []string) error {
	if len(Conf.pokedex) == 0 {
		fmt.Println("You don't have any Pokemon so far")
		return nil
	}

	fmt.Println("Your Pokedex:")
	for _, pokemon := range Conf.pokedex {
		fmt.Printf("  - %v\n", pokemon.Name)
	}
	return nil
}
