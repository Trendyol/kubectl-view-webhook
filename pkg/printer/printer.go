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
	"fmt"
	"github.com/hako/durafmt"
	"github.com/olekukonko/tablewriter"
	"github.com/pterm/pterm"
	"io"
	"os"
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

//modifyNamespaces returns BulletListItem's for Namespaces with customizable fields in order to give custom string and styles
func modifyNamespaces(str string) (text string, textStyle *pterm.Style, bullet string, bulletStyle *pterm.Style) {
	return str, pterm.NewStyle(pterm.FgGreen), pterm.DefaultBulletList.Bullet, pterm.NewStyle(pterm.FgLightWhite)
}

type BulletItem struct {
	Modify func(str string) (text string, textStyle *pterm.Style, bullet string, bulletStyle *pterm.Style)
	Items  []string
}

//convertStringArrayToBulletListItem converts given string array to
//pterm's BulletListItem array and returns as []pterm.BulletListItem
func convertStringArrayToBulletListItem(s BulletItem) []pterm.BulletListItem {
	var bulletItems []pterm.BulletListItem

	if s.Items != nil {
		for _, t := range s.Items {
			if s.Modify != nil {
				t, tS, b, bS := s.Modify(t)
				bulletItems = append(bulletItems, pterm.BulletListItem{
					Level:       0,
					Text:        t,
					TextStyle:   tS,
					Bullet:      b,
					BulletStyle: bS})
			} else {
				bulletItems = append(bulletItems, pterm.BulletListItem{
					Level: 0,
					Text:  t,
				})
			}
		}
	} else {
		bulletItems = append(bulletItems, pterm.BulletListItem{
			Level:       0,
			Text:        "No Active Namespaces",
			Bullet:      "✖",
			TextStyle:   pterm.NewStyle(pterm.FgRed),
			BulletStyle: pterm.NewStyle(pterm.FgLightRed),
		})
	}
	return bulletItems
}

//Print reads given PrintModel and prints as
//table using tablewriter.
func (p *Printer) Print(model *PrintModel) {
	var data [][]string

	for _, item := range model.Items {
		namespacesData, _ := pterm.DefaultBulletList.WithItems(
			convertStringArrayToBulletListItem(BulletItem{Items: item.ActiveNamespaces, Modify: modifyNamespaces})).Srender()

		resourcesLeveledList := pterm.LeveledList{}
		for _, rm := range item.ResourceModels {
			for _, rs := range rm.Resources {
				resourcesLeveledList = append(resourcesLeveledList, pterm.LeveledListItem{Level: 0, Text: pterm.NewStyle(pterm.FgWhite).Sprint(rs)})
			}
			for _, op := range rm.Operations {
				switch strings.ToUpper(op) {
				case "CREATE":
					op = pterm.NewStyle(pterm.FgGreen).Sprint("+", op)
				case "UPDATE":
					op = pterm.NewStyle(pterm.FgBlue).Sprint("^", op)
				case "DELETE":
					op = pterm.NewStyle(pterm.FgRed).Sprint("-", op)
				}
				resourcesLeveledList = append(resourcesLeveledList, pterm.LeveledListItem{Level: 1, Text: op})
			}
		}

		remainingTime := func(t time.Duration) string {
			if t == 0 {
				return pterm.Red("No CABundle")
			}
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

		serviceLeveledList := pterm.LeveledList{}

		service := item.Webhook.Service

		if service.Found {
			serviceLeveledList = append(serviceLeveledList, pterm.LeveledListItem{Level: 0, Text: service.Name})
			serviceLeveledList = append(serviceLeveledList, pterm.LeveledListItem{Level: 1, Text: "NS  : " + service.Namespace})
			if service.Path != nil {
				serviceLeveledList = append(serviceLeveledList, pterm.LeveledListItem{Level: 1, Text: "Path: " + *service.Path})
			}
			serviceLeveledList = append(serviceLeveledList, pterm.LeveledListItem{Level: 1, Text: fmt.Sprintf("IP  : %s (%s)", service.ClusterIP, service.Type)})
			if service.Ports != nil {
				for _, p := range service.Ports {
					getPortInfo := func() string {
						if p.TargetPort == 0 {
							return fmt.Sprintf("%d/%s", p.Port, p.Protocol)
						}
						return fmt.Sprintf("%d::%d/%s", p.Port, p.TargetPort, p.Protocol)
					}
					serviceLeveledList = append(serviceLeveledList, pterm.LeveledListItem{Level: 2, Text: getPortInfo()})
				}
			}
		} else {
			serviceLeveledList = append(serviceLeveledList, pterm.LeveledListItem{Level: 0, Text: pterm.NewStyle(pterm.FgRed).Sprintf("✖ %s", service.Name)})
			serviceLeveledList = append(serviceLeveledList, pterm.LeveledListItem{Level: 1, Text: "NS  : " + service.Namespace})
		}

		if len(serviceLeveledList) == 0 {
			serviceLeveledList = append(serviceLeveledList, pterm.LeveledListItem{Level: 0, Text: pterm.NewStyle(pterm.FgRed).Sprint("✖ No Services")})
		}

		webhookTreeList := pterm.NewTreeFromLeveledList(serviceLeveledList)
		resourcesTreeList := pterm.NewTreeFromLeveledList(resourcesLeveledList)

		wt, _ := pterm.DefaultTree.WithRoot(webhookTreeList).Srender()
		rt, _ := pterm.DefaultTree.WithRoot(resourcesTreeList).Srender()

		data = append(data, []string{item.Kind, item.Name, item.Webhook.Name, strings.TrimSuffix(wt, "\n"), strings.TrimSuffix(rt, "\n"), remainingTime(item.ValidUntil), namespacesData})
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Kind", "Name", "Webhook", "Service", "Resources&Operations", "Remaining Day", "Active NS"})
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
