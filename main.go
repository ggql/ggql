package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/go-git/go-git/v5"

	"github.com/ggql/ggql/ast"
	"github.com/ggql/ggql/cli"
	"github.com/ggql/ggql/engine"
	"github.com/ggql/ggql/parser"
)

func main() {
	args := os.Args
	command := cli.ParseArguments(args)

	if len(command.ReplMode.Repos) != 0 {
		launchGitqlRepl(command.ReplMode)
	}
	if command.QueryMode.Query != "" {
		reporter := cli.DiagnosticReporter{}
		gitReposResult := validateGitRepositories(command.QueryMode.Arguments.Repos)
		if gitReposResult.err != nil {
			reporter.ReportDiagnostic("", &parser.Diagnostic{Message: gitReposResult.err.Error()})
			return
		}

		repos := gitReposResult.ok
		env := ast.Environment{}
		executeGitqlQuery(command.QueryMode.Query, command.QueryMode.Arguments, repos, &env, &reporter)
	}
	if command.Help {
		cli.PrintHelpList()
	}
	if command.Version != "" {
		fmt.Printf("GitQL version %s\n", command.Version)
	}
	if command.Error != "" {
		fmt.Println(command.Error)
	}
}

func launchGitqlRepl(args cli.Arguments) {
	reporter := cli.DiagnosticReporter{}
	gitReposResult := validateGitRepositories(args.Repos)
	if gitReposResult.err != nil {
		reporter.ReportDiagnostic("", &parser.Diagnostic{Message: gitReposResult.err.Error()})
		return
	}
	globalEnv := ast.Environment{}
	gitRepositories := gitReposResult.ok

	scanner := bufio.NewScanner(os.Stdin)
	for {
		if isTerminal(os.Stdin) {
			fmt.Print("gql > ")
		}

		if !scanner.Scan() {
			break
		}

		input := scanner.Text()
		if input == "exit" {
			fmt.Println("Goodbye!")
			break
		}

		executeGitqlQuery(input, args, gitRepositories, &globalEnv, &reporter)
		globalEnv.ClearSession()
	}
}

// nolint:funlen,gocyclo
func executeGitqlQuery(
	query string,
	args cli.Arguments,
	repos []*git.Repository,
	env *ast.Environment,
	reporter *cli.DiagnosticReporter,
) {
	frontStart := time.Now()
	tokens, err := parser.Tokenize(query)
	if err != nil {
		reporter.ReportDiagnostic(query, err)
		return
	}

	if len(tokens) == 0 {
		return
	}

	queryNode, err1 := parser.ParserGql(tokens, env)
	if &err1 != nil {
		reporter.ReportDiagnostic(query, &err1)
		return
	}

	frontDuration := time.Since(frontStart)

	engineStart := time.Now()
	evaluationResult, err2 := engine.Evaluate(env, repos, queryNode)
	if err2 != nil {
		reporter.ReportDiagnostic(query, &parser.Diagnostic{Message: err2.Error()})
		return
	}

	if evaluationResult.SelectedGroups.Obj.Len() != 0 {
		groups := evaluationResult
		hiddenSelection := err2
		switch args.OutputFormat {
		case cli.Render:
			err := cli.RenderObjects(&groups.SelectedGroups.Obj, []string{hiddenSelection.Error()}, args.Pagination, args.PageSize)
			if err != nil {
				fmt.Println(err)
			}
		case cli.JSON:
			var indexes []int
			for index, title := range groups.SelectedGroups.Obj.Titles {
				if strings.Contains(hiddenSelection.Error(), title) {
					indexes = append([]int{index}, indexes...)
				}
			}

			if groups.SelectedGroups.Obj.Len() > 1 {
				groups.SelectedGroups.Obj.Flat()
			}

			for _, index := range indexes {
				groups.SelectedGroups.Obj.Titles = append(groups.SelectedGroups.Obj.Titles[:index], groups.SelectedGroups.Obj.Titles[index+1:]...)

				for _, row := range groups.SelectedGroups.Obj.Groups[0].Rows {
					row.Values = append(row.Values[:index], row.Values[index+1:]...)
				}
			}

			if json, err := groups.SelectedGroups.Obj.AsJson(); err == nil {
				fmt.Println(json)
			}
		case cli.CSV:
			indexes := []int{}
			for index, title := range groups.SelectedGroups.Obj.Titles {
				if strings.Contains(hiddenSelection.Error(), title) {
					indexes = append([]int{index}, indexes...)
				}
			}

			if len(groups.SelectedGroups.Obj.Groups) > 1 {
				groups.SelectedGroups.Obj.Flat()
			}

			for _, index := range indexes {
				groups.SelectedGroups.Obj.Titles = append(groups.SelectedGroups.Obj.Titles[:index], groups.SelectedGroups.Obj.Titles[index+1:]...)

				for _, row := range groups.SelectedGroups.Obj.Groups[0].Rows {
					row.Values = append(row.Values[:index], row.Values[index+1:]...)
				}
			}

			if csv, err := groups.SelectedGroups.Obj.AsCsv(); err == nil {
				fmt.Println(csv)
			}
		}
	}

	engineDuration := time.Since(engineStart)

	if args.Analysis {
		fmt.Println("Analysis:")
		fmt.Println("Frontend : ", frontDuration)
		fmt.Println("Engine   : ", engineDuration)
		fmt.Println("Total    : ", frontDuration+engineDuration)
	}
}

func validateGitRepositories(repositories []string) result {
	var gitRepositories []*git.Repository
	for _, repository := range repositories {
		gitRepository, err := git.PlainOpen(repository)
		if err != nil {
			return result{err: err}
		}
		gitRepositories = append(gitRepositories, gitRepository)
	}
	return result{ok: gitRepositories}
}

type result struct {
	ok  []*git.Repository
	err error
}

func isTerminal(f *os.File) bool {
	// Check file descriptor if used as terminal or not with os.Isatty()
	if stat, err := f.Stat(); err == nil {
		return (stat.Mode() & os.ModeCharDevice) != 0
	}
	return false
}
