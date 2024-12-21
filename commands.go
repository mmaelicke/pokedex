package main

type cliCommand struct {
	name        string
	description string
	callback    func(*config, []string) error
}

var commands = make(map[string]cliCommand)

func registerCommand(commands *map[string]cliCommand, command cliCommand) {
	(*commands)[command.name] = command
}

func populateCommands(commands *map[string]cliCommand) {
	// register the command
	registerCommand(commands, cliCommand{
		name:        "help",
		description: "Displays a help message",
		callback:    commandHelp,
	})
	registerCommand(commands, cliCommand{
		name:        "exit",
		description: "Close the Pokedex",
		callback:    commandExit,
	})
	registerCommand(commands, cliCommand{
		name:        "map",
		description: "navigate through the locations of the PokeAPI",
		callback:    mapCallback,
	})
	registerCommand(commands, cliCommand{
		name:        "mapb",
		description: "navigate backwards through the locations of the PokeAPI",
		callback:    mapbCallback,
	})
	registerCommand(commands, cliCommand{
		name:        "explore",
		description: "explore a location of the PokeAPI",
		callback:    exploreCallback,
	})
	registerCommand(commands, cliCommand{
		name:        "catch",
		description: "catch a Pokemon",
		callback:    catchCallback,
	})
	registerCommand(commands, cliCommand{
		name:        "inspect",
		description: "get infomation about a Pokemon in your Pokedex",
		callback:    inspectCallback,
	})
	registerCommand(commands, cliCommand{
		name:        "pokedex",
		description: "View your Pokedex",
		callback:    pokedexCallback,
	})
}
