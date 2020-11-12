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
	"github.com/hako/durafmt"
	"github.com/olekukonko/tablewriter"
	"github.com/pterm/pterm"
	"io"
	"os"
	"strconv"
	"strings"
	"time"
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

//getItemValuesForString returns BulletListItem's customizable fields in order to give custom string and styles
func getItemValuesForString(str string) (text string, textStyle *pterm.Style, bullet string, bulletStyle *pterm.Style) {
	switch strings.ToUpper(str) {
	case "CREATE":
		return str, pterm.NewStyle(pterm.FgGreen), "+", pterm.NewStyle(pterm.FgLightGreen)
	case "UPDATE":
		return str, pterm.NewStyle(pterm.FgBlue), "^", pterm.NewStyle(pterm.FgLightBlue)
	case "DELETE":
		return str, pterm.NewStyle(pterm.FgRed), "-", pterm.NewStyle(pterm.FgLightRed)
	}

	return str, nil, pterm.DefaultBulletList.Bullet, nil
}

//convertStringArrayToBulletListItem converts given string array to
//pterm's BulletListItem array and returns as []pterm.BulletListItem
func convertStringArrayToBulletListItem(s []string) []pterm.BulletListItem {
	var bulletItems []pterm.BulletListItem
	for _, t := range s {
		t, tS, b, bS := getItemValuesForString(t)
		bulletItems = append(bulletItems, pterm.BulletListItem{
			Level:       0,
			Text:        t,
			TextStyle:   tS,
			Bullet:      b,
			BulletStyle: bS,
		})
	}
	return bulletItems
}

//Print reads given PrintModel and prints as
//table using tablewriter.
func (p *Printer) Print(model *PrintModel) {
	var data [][]string

	for _, item := range model.Items {
		operationsData, _ := pterm.DefaultBulletList.WithItems(convertStringArrayToBulletListItem(item.Operations)).Srender()
		resourcesData, _ := pterm.DefaultBulletList.WithItems(convertStringArrayToBulletListItem(item.Resources)).Srender()
		namespacesData, _ := pterm.DefaultBulletList.WithItems(convertStringArrayToBulletListItem(item.ActiveNamespaces)).Srender()

		remainingTime := func(t time.Duration) string {
			days := t.Hours() / 24

			N := func() int {
				if days < 2 {
					return 2
				} else {
					return 1
				}
			}

			str := durafmt.Parse(t).LimitFirstN(N()).String()

			if days < 4000 {
				return pterm.Red(str)
			} else if days < 60000 {
				return pterm.Yellow(str)
			} else {
				return pterm.Green(str)
			}
		}

		webhookTreeList := pterm.NewTreeFromLeveledList(pterm.LeveledList{
			pterm.LeveledListItem{Level: 0, Text: item.Webhook.ServiceName},
			pterm.LeveledListItem{Level: 1, Text: "NS  : " + item.Webhook.ServiceNamespace},
			pterm.LeveledListItem{Level: 1, Text: "Path: " + *item.Webhook.ServicePath},
			pterm.LeveledListItem{Level: 1, Text: "Port: " + strconv.Itoa(int(*item.Webhook.ServicePort))},
		})

		s, _ := pterm.DefaultTree.WithRoot(webhookTreeList).Srender()

		data = append(data, []string{item.Kind, item.Name, item.Webhook.Name, strings.TrimSuffix(s, "\n"), resourcesData, operationsData, remainingTime(item.ValidUntil), namespacesData})
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Kind", "Name", "Webhook", "Service", "Resources", "Operations", "Remaining Day", "Active NS"})
	table.SetRowLine(true)
	table.SetAutoMergeCells(true)
	table.SetHeaderLine(true)
	table.SetBorder(true)
	//table.SetReflowDuringAutoWrap(true)
	table.SetAutoWrapText(false)
	table.SetAutoFormatHeaders(true)
	//table.SetCenterSeparator("")
	//table.SetColumnSeparator("")
	//table.SetRowSeparator("")
	//table.SetTablePadding(" ")
	//table.SetNoWhiteSpace(true)
	table.AppendBulk(data)

	table.Render()
}
