package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/mmaelicke/pokedex/internal"
)

type config struct {
	baseUrl     string
	nextUrl     string
	previousUrl string
	cache       internal.Cache
	pokedex     map[string]PokemonResponse
}

var Config = config{
	baseUrl:     "https://pokeapi.co/api/v2/",
	nextUrl:     "https://pokeapi.co/api/v2/location-area/",
	previousUrl: "",
	cache:       internal.NewCache(time.Minute * 10),
	pokedex:     map[string]PokemonResponse{},
}

func cleanInput(text string) []string {
	names := []string{}

	for _, word := range strings.Fields(text) {
		lower := strings.TrimSuffix(strings.ToLower(word), ",")
		if len(lower) > 0 {
			names = append(names, lower)
		}
	}
	return names
}

func commandExit(Config *config, _ []string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(Config *config, _ []string) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println("")
	for _, cmd := range commands {
		fmt.Printf("  %s - %s\n", cmd.name, cmd.description)
	}
	return nil
}

func main() {
	// init the commands
	populateCommands(&commands)
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		prompt := scanner.Text()

		// get the commands
		args := cleanInput(prompt)
		if len(args) == 0 {
			continue
		}

		// check if the command exists
		command, ok := commands[args[0]]
		if !ok {
			fmt.Println("Unkown command")
			continue
		}

		err := command.callback(&Config, args[1:])
		if err != nil {
			fmt.Println(err)
			continue
		}
	}
}
