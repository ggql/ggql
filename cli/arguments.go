package cli

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

const (
	// Render the output as table
	Render OutputFormat = iota
	// JSON Print the output in json format
	JSON
	// CSV Print the output in csv format
	CSV
)

var (
	Version string
)

// OutputFormat Represent the different type of available formats
type OutputFormat int

// Arguments for GitQL
type Arguments struct {
	Repos        []string
	Analysis     bool
	Pagination   bool
	PageSize     int
	OutputFormat OutputFormat
}

// Command represents the possible GitQL commands
type Command struct {
	ReplMode  Arguments
	QueryMode struct {
		Query     string
		Arguments Arguments
	}
	MutateMode struct {
		Mutate    string
		Arguments Arguments
	}
	Help    bool
	Version string
	Error   string
}

// NewArguments creates a new instance of Arguments with default settings
func NewArguments() Arguments {
	return Arguments{
		Repos:        []string{},
		Analysis:     false,
		Pagination:   false,
		PageSize:     10,
		OutputFormat: Render,
	}
}

// ParseArguments parses the command-line arguments and returns the corresponding Command
// nolint: funlen,gocyclo
func ParseArguments(args []string) Command {
	argsLen := len(args)

	if contains(args, "--help", "-h") {
		return Command{Help: true}
	}

	if contains(args, "--version", "-v") {
		return Command{Version: Version}
	}

	var optionalQuery string
	var optionalMutate string
	arguments := NewArguments()

	argIndex := 1
	for argIndex < argsLen {
		arg := args[argIndex]
		if arg[0] != '-' {
			return Command{Error: fmt.Sprintf("Unknown argument %s", arg)}
		}
		switch arg {
		case "--repos", "-r":
			argIndex++
			if argIndex >= argsLen {
				return Command{Error: fmt.Sprintf("Argument %s must be followed by one or more paths", arg)}
			}
			for argIndex < argsLen {
				repo := args[argIndex]
				if repo[0] != '-' {
					arguments.Repos = append(arguments.Repos, repo)
					argIndex++
					continue
				}
				break
			}
		case "--query", "-q":
			argIndex++
			if argIndex >= argsLen {
				return Command{Error: fmt.Sprintf("Argument %s must be followed by the query", arg)}
			}
			optionalQuery = args[argIndex]
			argIndex++
		case "--mutate", "-m":
			argIndex++
			if argIndex >= argsLen {
				return Command{Error: fmt.Sprintf("Argument %s must be followed by the mutate", arg)}
			}
			optionalMutate = args[argIndex]
			argIndex++
		case "--analysis", "-a":
			arguments.Analysis = true
			argIndex++
		case "--pagination", "-p":
			arguments.Pagination = true
			argIndex++
		case "--pagesize", "-ps":
			argIndex++
			if argIndex >= argsLen {
				return Command{Error: fmt.Sprintf("Argument %s must be followed by the page size", arg)}
			}
			pageSize, err := strconv.Atoi(args[argIndex])
			if err != nil {
				return Command{Error: "Invalid page size"}
			}
			arguments.PageSize = pageSize
			argIndex++
		case "--output", "-o":
			argIndex++
			if argIndex >= len(args) {
				return Command{Error: fmt.Sprintf("argument %s must be followed by output format", arg)}
			}
			switch strings.ToLower(args[argIndex]) {
			case "csv":
				arguments.OutputFormat = CSV
			case "json":
				arguments.OutputFormat = JSON
			case "render":
				arguments.OutputFormat = Render
			default:
				return Command{Error: "invalid output format"}
			}
			argIndex++
		default:
			return Command{Error: fmt.Sprintf("Unknown command %s", arg)}
		}
	}

	// Add the current directory if no repository is passed
	if len(arguments.Repos) == 0 {
		currentDir, err := os.Getwd()
		if err != nil {
			return Command{Error: "Missing repository paths"}
		}

		arguments.Repos = append(arguments.Repos, currentDir)
	}

	if optionalQuery != "" {
		return Command{
			QueryMode: struct {
				Query     string
				Arguments Arguments
			}{
				Query:     optionalQuery,
				Arguments: arguments,
			},
		}
	}

	if optionalMutate != "" {
		return Command{
			MutateMode: struct {
				Mutate    string
				Arguments Arguments
			}{
				Mutate:    optionalMutate,
				Arguments: arguments,
			},
		}
	}

	return Command{ReplMode: arguments}
}

// PrintHelpList prints the help message for GitQL
func PrintHelpList() {
	fmt.Println("ggql is a SQL like query and mutate language to run on local repositories")
	fmt.Println()
	fmt.Println("Usage: ggql [OPTIONS]")
	fmt.Println()
	fmt.Println("Options:")
	fmt.Println("-r,  --repos <REPOS>        Path for local repositories to run query on")
	fmt.Println("-q,  --query <GQL Query>    GitQL query to run on selected repositories")
	fmt.Println("-m,  --mutate <GQL Mutate>  GitQL mutate to run on selected repositories")
	fmt.Println("-p,  --pagination           Enable print result with pagination")
	fmt.Println("-ps, --pagesize             Set pagination page size [default: 10]")
	fmt.Println("-o,  --output               Set output format [render, json, csv]")
	fmt.Println("-a,  --analysis             Print Query analysis")
	fmt.Println("-h,  --help                 Print GitQL help")
	fmt.Println("-v,  --version              Print GitQL Current Version")
}

func contains(arr []string, items ...string) bool {
	for _, item := range items {
		for _, a := range arr {
			if item == a {
				return true
			}
		}
	}
	return false
}
