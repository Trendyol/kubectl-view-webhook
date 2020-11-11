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

package k8s

import (
	"context"
	"github.com/Trendyol/kubectl-view-webhook/pkg/printer"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type MutatingWebHookClient struct {
	client  *kubernetes.Clientset
	context context.Context
}

// NewMutatingWebHookClient constructs a new MutatingWebHookClient with the specified output
// of *kubernetes.Clientset
func NewMutatingWebHookClient(client *kubernetes.Clientset) *MutatingWebHookClient {
	return &MutatingWebHookClient{
		client:  client,
		context: context.Background(),
	}
}

func (w *MutatingWebHookClient) Run() (*printer.PrintModel, error) {
	c := w.client.AdmissionregistrationV1().MutatingWebhookConfigurations()

	mutatingWebhookConfigurationList, err := c.List(w.context, metaV1.ListOptions{
		TypeMeta: metaV1.TypeMeta{
			Kind:       "MutatingWebhookConfiguration",
			APIVersion: "admissionregistration.k8s.io/v1",
		},
	})
	if err != nil {
		return nil, err
	}

	var items []printer.PrintItem

	for _, mwc := range mutatingWebhookConfigurationList.Items {
		items = append(items, printer.PrintItem{
			Kind: "MutatingWebhookConfiguration", //TODO: typeMeta nil
			Name: mwc.Name,
		})
	}

	return &printer.PrintModel{
		Items: items,
	}, nil
}
