/*
Copyright © 2020 Trendyol Tech

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package printer

import (
	"github.com/olekukonko/tablewriter"
	"github.com/pterm/pterm"
	"io"
	"os"
)

// Printer formats and prints check results and warnings.
type Printer struct {
	out io.Writer
}

// NewPrinter constructs a new Printer with the specified output io.Writer
// and output format.
func NewPrinter(out io.Writer) *Printer {
	return &Printer{
		out: out,
	}
}

func (p *Printer) Print(model *PrintModel) {
	leveledList := pterm.LeveledList{
		pterm.LeveledListItem{Level: 0, Text: "Foo"},
		pterm.LeveledListItem{Level: 1, Text: "Bar"},
		pterm.LeveledListItem{Level: 1, Text: "Baz"},
	}

	// Generate tree from LeveledList.
	root := pterm.NewTreeFromLeveledList(leveledList)

	// Render TreePrinter
	s, _ := pterm.DefaultTree.WithRoot(root).Srender()

	//pterm.NewRGB(178, 44, 199).Println("This text is printed with a custom RGB!")

	var data [][]string

	for _, item := range model.Items {
		data = append(data, []string{item.Kind, item.Name, s, "$10.98"})
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Kind", "Name", "CV2", "Amount"})
	//table.SetFooter([]string{"", "", "Total", "$146.93"})
	table.SetAutoWrapText(false)
	table.SetRowLine(true)
	table.AppendBulk(data)
	table.Render()

}
