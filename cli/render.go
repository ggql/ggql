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
func RenderObjects(groups *ast.GitQLObject, hiddenSelections []string, pagination bool, pageSize int) error {
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

	if groups.Len() > 1 {
		groups.Flat()
	}

	if groups.IsEmpty() || groups.Groups[0].IsEmpty() {
		return nil
	}

	gqlGroup := groups.Groups[0]
	gqlGroupLen := gqlGroup.Len()

	// Setup titles
	for _, title := range groups.Titles {
		if !contains(hiddenSelections, title) {
			titles = append(titles, title)
		}
	}

	// Setup table headers
	for _, item := range titles {
		tableHeaders = append(tableHeaders, pterm.Green(item))
	}

	// Print all data without pagination
	if !pagination || pageSize >= gqlGroupLen {
		return printGroupAsTable(titles, tableHeaders, gqlGroup.Rows)
	}

	// Setup the pagination mode
	var numberOfPages = gqlGroupLen / pageSize
	if gqlGroupLen%pageSize != 0 {
		numberOfPages++
	}

	for {
		startIndex := (currentPage - 1) * pageSize
		endIndex := min(startIndex+pageSize, gqlGroupLen)

		currentPageGroups := gqlGroup.Rows[startIndex:endIndex]

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

func printGroupAsTable(titles, tableHeaders []string, rows []ast.Row) error {
	data := pterm.TableData{}
	data = append(data, tableHeaders)

	for _, row := range rows {
		var buf []string
		for index, _ := range titles {
			buf = append(buf, row.Values[index].AsText())
		}
		data = append(data, buf)
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
