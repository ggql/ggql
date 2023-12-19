package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/pkg/errors"
	"github.com/pterm/pterm"

	"github.com/ggql/ggql/ast"
)

const (
	nextPage = iota
	previousPage
	quit
)

const (
	pageNum = 2
)

// nolint: gocyclo
func RenderObjects(groups *[][]ast.GQLObject, hiddenSelections []string, pagination bool, pageSize int) error {
	contains := func(slice []string, item string) bool {
		for _, a := range slice {
			if a == item {
				return true
			}
		}
		return false
	}

	var currentPage = 1
	var err error
	var tableHeaders []string
	var titles []string

	if len(*groups) > 1 {
		if e := ast.FlatGQLGroups(groups); e != nil {
			return errors.Wrap(e, "failed to run FlatGQLGroups")
		}
	}

	if len(*groups) == 0 || len((*groups)[0]) == 0 {
		return nil
	}

	gqlGroup := (*groups)[0]
	gqlGroupLen := len(gqlGroup)

	// Setup titles
	for key, val := range (*groups)[0][0].Attributes {
		if !contains(hiddenSelections, key) {
			titles = append(titles, val.(string))
		}
	}

	// Setup table headers
	for _, item := range titles {
		tableHeaders = append(tableHeaders, pterm.Green(item))
	}

	// Print all data without pagination
	if !pagination || pageSize >= gqlGroupLen {
		return printGroupAsTable(titles, tableHeaders, gqlGroup)
	}

	// Setup the pagination mode
	var numberOfPages = gqlGroupLen / pageSize
	if gqlGroupLen%pageSize != 0 {
		numberOfPages++
	}

	for {
		startIndex := (currentPage - 1) * pageSize
		endIndex := min(startIndex+pageSize, gqlGroupLen)

		currentPageGroups := gqlGroup[startIndex:endIndex]

		fmt.Printf("Page %d/%d\n", currentPage, numberOfPages)
		err = printGroupAsTable(titles, tableHeaders, currentPageGroups)
		if err != nil {
			break
		}

		paginationInput := handlePaginationInput(currentPage, numberOfPages)
		switch paginationInput {
		case nextPage:
			currentPage++
		case previousPage:
			currentPage--
		case quit:
			break
		}
	}

	return err
}

func printGroupAsTable(titles, tableHeaders []string, group []ast.GQLObject) error {
	data := pterm.TableData{}
	data = append(data, tableHeaders)

	for _, object := range group {
		var row []string
		for _, title := range titles {
			row = append(row, object.Attributes[title].(string))
		}
		data = append(data, row)
	}

	if err := pterm.DefaultTable.WithHasHeader().WithRowSeparator("-").WithHeaderRowSeparator("-").WithData(data).Render(); err != nil {
		return errors.Wrap(err, "failed to run render")
	}

	return nil
}

func handlePaginationInput(currentPage, numberOfPages int) int {
	reader := bufio.NewReader(os.Stdin)

	for {
		if currentPage < pageNum {
			fmt.Println("Enter 'n' for next page, or 'q' to quit:")
		} else if currentPage == numberOfPages {
			fmt.Println("'p' for previous page, or 'q' to quit:")
		} else {
			fmt.Println("Enter 'n' for next page, 'p' for previous page, or 'q' to quit:")
		}

		input, _ := reader.ReadString('\n')
		input = strings.ToLower(strings.TrimSpace(input))

		if input == "q" || input == "n" || input == "p" {
			switch input {
			case "n":
				if currentPage < numberOfPages {
					return nextPage
				} else {
					fmt.Println("Already on the last page")
					continue
				}
			case "p":
				if currentPage > 1 {
					return previousPage
				} else {
					fmt.Println("Already on the first page")
					continue
				}
			case "q":
				return quit
			default:
				fmt.Println("Invalid input")
			}
		}

		fmt.Println("Invalid input")
	}
}
