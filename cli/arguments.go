package cli

import (
	"fmt"
	"os"
	"strconv"
)

// Arguments for GitQL
type Arguments struct {
	Repos      []string
	Analysis   bool
	Pagination bool
	PageSize   int
}

// Command represents the possible GitQL commands
type Command struct {
	ReplMode  Arguments
	QueryMode struct {
		Query     string
		Arguments Arguments
	}
	Help    bool
	Version string
	Error   string
}

// NewArguments creates a new instance of Arguments with default settings
func NewArguments() Arguments {
	return Arguments{
		Repos:      []string{},
		Analysis:   false,
		Pagination: false,
		PageSize:   10,
	}
}

// nolint: funlen,gocyclo
// ParseArguments parses the command-line arguments and returns the corresponding Command
func ParseArguments(args []string) Command {
	argsLen := len(args)

	if contains(args, "--help", "-h") {
		return Command{Help: true}
	}

	if contains(args, "--version", "-v") {
		return Command{Version: "1.0.0"}
	}

	var optionalQuery string
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

	return Command{ReplMode: arguments}
}

// PrintHelpList prints the help message for GitQL
func PrintHelpList() {
	fmt.Println("GitQL is a SQL like query language to run on local repositories")
	fmt.Println()
	fmt.Println("Usage: gitql.exe [OPTIONS]")
	fmt.Println()
	fmt.Println("Options:")
	fmt.Println("-r,  --repos <REPOS>        Path for local repositories to run query on")
	fmt.Println("-q,  --query <GQL Query>    GitQL query to run on selected repositories")
	fmt.Println("-p,  --pagination           Enable print result with pagination")
	fmt.Println("-ps, --pagesize             Set pagination page size [default: 10]")
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
