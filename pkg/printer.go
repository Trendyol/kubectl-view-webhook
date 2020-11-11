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

package pkg

import (
	"github.com/pterm/pterm"
	"github.com/olekukonko/tablewriter"
	"io"
	"os"
)

// Printer formats and prints check results and warnings.
type Printer struct {
	out  io.Writer
}

// NewPrinter constructs a new Printer with the specified output io.Writer
// and output format.
func NewPrinter(out io.Writer) *Printer {
	return &Printer{
		out:  out,
	}
}

func (p *Printer) Print() {
	leveledList := pterm.LeveledList{
		pterm.LeveledListItem{Level: 0, Text: "C:"},
		pterm.LeveledListItem{Level: 1, Text: "Users"},
		pterm.LeveledListItem{Level: 1, Text: "Windows"},
		pterm.LeveledListItem{Level: 1, Text: "Programs"},
		pterm.LeveledListItem{Level: 1, Text: "Programs(x86)"},
		pterm.LeveledListItem{Level: 1, Text: "dev"},
		pterm.LeveledListItem{Level: 0, Text: "D:"},
		pterm.LeveledListItem{Level: 0, Text: "E:"},
		pterm.LeveledListItem{Level: 1, Text: "Movies"},
		pterm.LeveledListItem{Level: 1, Text: "Music"},
		pterm.LeveledListItem{Level: 2, Text: "LinkinPark"},
		pterm.LeveledListItem{Level: 1, Text: "Games"},
		pterm.LeveledListItem{Level: 2, Text: "Shooter"},
		pterm.LeveledListItem{Level: 3, Text: "CallOfDuty"},
		pterm.LeveledListItem{Level: 3, Text: "CS:GO"},
		pterm.LeveledListItem{Level: 3, Text: "Battlefield"},
		pterm.LeveledListItem{Level: 4, Text: "Battlefield 1"},
		pterm.LeveledListItem{Level: 4, Text: "Battlefield 2"},
		pterm.LeveledListItem{Level: 0, Text: "F:"},
		pterm.LeveledListItem{Level: 1, Text: "dev"},
		pterm.LeveledListItem{Level: 2, Text: "dops"},
		pterm.LeveledListItem{Level: 2, Text: "PTerm"},
	}

	// Generate tree from LeveledList.
	root := pterm.NewTreeFromLeveledList(leveledList)

	// Render TreePrinter
	s, _ := pterm.DefaultTree.WithRoot(root).Srender()

	pterm.NewRGB(178, 44, 199).Println("This text is printed with a custom RGB!")

	data := [][]string{
		[]string{"Service 1", s, "1234", "$10.98"},
		[]string{"Service 2", "January Hosting", "1234", "$10.98"},
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Date", "Description", "CV2", "Amount"})
	//table.SetFooter([]string{"", "", "Total", "$146.93"})
	table.SetAutoWrapText(false)
	table.SetRowLine(true)
	table.AppendBulk(data)
	table.Render()

}