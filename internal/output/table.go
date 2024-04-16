package output

import (
	"fmt"
	"os"

	"github.com/olekukonko/tablewriter"
)

func DisplayAsTable(headers []string, data [][]string) {
	if len(data) == 0 {
		fmt.Println("No data.")
		return
	}

	// Initialize the tablewriter with os.Stdout as the output
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(headers)
	table.SetAutoWrapText(false)
	table.SetAutoFormatHeaders(true)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetCenterSeparator("")
	table.SetColumnSeparator("")
	table.SetRowSeparator("")
	table.SetHeaderLine(false)
	table.SetBorder(false)
	table.SetTablePadding("\t")
	table.SetNoWhiteSpace(true)

	for _, row := range data {
		table.Append(row)
	}

	table.Render()
}
