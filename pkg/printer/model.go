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

import "time"

type PrintModel struct {
	Items []PrintItem
}

type ResourceModel struct {
	Operations []string
	Resources  []string
}

type PrintItem struct {
	Name             string
	Webhook          PrintWebhookItem
	Kind             string
	ResourceModels   []ResourceModel
	ValidUntil       time.Duration
	ActiveNamespaces []string
}

type PrintWebhookItem struct {
	Name    string
	Service PrintServiceItem
}

type PrintServiceItem struct {
	Found     bool
	Name      string
	Namespace string
	Path      *string
	Ports     []PrintServicePortItem
	ClusterIP string
	Type      string
}

type PrintServicePortItem struct {
	Port       int32
	TargetPort int32
	Protocol   string
}
