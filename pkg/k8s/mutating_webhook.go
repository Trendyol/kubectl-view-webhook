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
	"crypto/x509"
	"encoding/pem"
	"github.com/Trendyol/kubectl-view-webhook/pkg/printer"
	"k8s.io/api/admissionregistration/v1beta1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	typedV1beta1 "k8s.io/client-go/kubernetes/typed/admissionregistration/v1beta1"
	typedCoreV1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"log"
	"time"
)

type WebHookClient struct {
	wClient typedV1beta1.MutatingWebhookConfigurationInterface
	vClient typedV1beta1.ValidatingWebhookConfigurationInterface
	nClient typedCoreV1.NamespaceInterface
	context context.Context
}

// NewWebHookClient constructs a new WebHookClient with the specified output
// of *kubernetes.Clientset
func NewWebHookClient(client *kubernetes.Clientset) *WebHookClient {
	return &WebHookClient{
		wClient: client.AdmissionregistrationV1beta1().MutatingWebhookConfigurations(),
		vClient: client.AdmissionregistrationV1beta1().ValidatingWebhookConfigurations(),
		nClient: client.CoreV1().Namespaces(),
		context: context.Background(),
	}
}

// Run
func (w *WebHookClient) Run(args []string) (*printer.PrintModel, error) {
	var items []printer.PrintItem

	if len(args) == 1 {
		mutatingWebhookConfigurationList, _ := w.wClient.List(w.context, metaV1.ListOptions{})

		validatingWebhookConfigurationList, _ := w.vClient.List(w.context, metaV1.ListOptions{})

		for _, mwc := range mutatingWebhookConfigurationList.Items {
			w.fillMutatingWebhookConfigurations(mwc, &items)
		}
		for _, mwc := range validatingWebhookConfigurationList.Items {
			w.fillValidatingWebhookConfigurations(mwc, &items)
		}
	} else {
		mutatingWebhookConfiguration, _ := w.wClient.Get(w.context, args[1], metaV1.GetOptions{})

		validatingWebhookConfiguration, _ := w.vClient.Get(w.context, args[1], metaV1.GetOptions{})

		w.fillMutatingWebhookConfigurations(*mutatingWebhookConfiguration, &items)
		w.fillValidatingWebhookConfigurations(*validatingWebhookConfiguration, &items)
	}

	return &printer.PrintModel{
		Items: items,
	}, nil
}

func (w *WebHookClient) fillMutatingWebhookConfigurations(mwc v1beta1.MutatingWebhookConfiguration, items *[]printer.PrintItem) {
	item := printer.PrintItem{
		Kind: "Mutating",
		Name: mwc.Name, //TODO: typeMeta nil
	}

	for _, webhook := range mwc.Webhooks {
		var operations, resources, activeNamespaces []string
		w.fillActiveNamespacesForMutating(webhook, &activeNamespaces)

		item.Webhook = printer.PrintWebhookItem{
			Name:             webhook.Name,
			ServiceName:      webhook.ClientConfig.Service.Name,
			ServiceNamespace: webhook.ClientConfig.Service.Namespace,
			ServicePath:      webhook.ClientConfig.Service.Path,
			ServicePort:      webhook.ClientConfig.Service.Port,
		}

		w.fillRulesForMutating(webhook, &operations, &resources)

		item.Operations = operations
		item.Resources = resources
		item.ValidUntil = retrieveValidDateCount(webhook.ClientConfig.CABundle)
		item.ActiveNamespaces = activeNamespaces
		*items = append(*items, item)
	}
}
func (w *WebHookClient) fillValidatingWebhookConfigurations(mwc v1beta1.ValidatingWebhookConfiguration, items *[]printer.PrintItem) {
	item := printer.PrintItem{
		Kind: "Validating",
		Name: mwc.Name, //TODO: typeMeta nil
	}

	for _, webhook := range mwc.Webhooks {
		var operations, resources, activeNamespaces []string
		w.fillActiveNamespacesForValidating(webhook, &activeNamespaces)

		item.Webhook = printer.PrintWebhookItem{
			Name:             webhook.Name,
			ServiceName:      webhook.ClientConfig.Service.Name,
			ServiceNamespace: webhook.ClientConfig.Service.Namespace,
			ServicePath:      webhook.ClientConfig.Service.Path,
			ServicePort:      webhook.ClientConfig.Service.Port,
		}

		w.fillRulesForValidating(webhook, &operations, &resources)

		item.Operations = operations
		item.Resources = resources
		item.ValidUntil = retrieveValidDateCount(webhook.ClientConfig.CABundle)
		item.ActiveNamespaces = activeNamespaces
		*items = append(*items, item)
	}
}
func (w *WebHookClient) fillRulesForMutating(webhook v1beta1.MutatingWebhook, operations *[]string, resources *[]string) {
	for _, rule := range webhook.Rules {

		for _, op := range rule.Operations {
			*operations = append(*operations, string(op))
		}

		*resources = append(*resources, rule.Resources...)
	}
}
func (w *WebHookClient) fillRulesForValidating(webhook v1beta1.ValidatingWebhook, operations *[]string, resources *[]string) {
	for _, rule := range webhook.Rules {

		for _, op := range rule.Operations {
			*operations = append(*operations, string(op))
		}

		*resources = append(*resources, rule.Resources...)
	}
}
func (w *WebHookClient) fillActiveNamespacesForMutating(webhook v1beta1.MutatingWebhook, activeNamespaces *[]string) {
	if webhook.NamespaceSelector != nil {
		ncList, _ := w.nClient.List(w.context, metaV1.ListOptions{})

		if ncList != nil {
			for _, ns := range ncList.Items {
				available := false
				for k, v := range webhook.NamespaceSelector.MatchLabels {
					if ns.Labels[k] == v {
						available = true
					}
				}
				if available {
					*activeNamespaces = append(*activeNamespaces, ns.Name)
				}
			}
		}
	}
}
func (w *WebHookClient) fillActiveNamespacesForValidating(webhook v1beta1.ValidatingWebhook, activeNamespaces *[]string) {
	if webhook.NamespaceSelector != nil {
		ncList, _ := w.nClient.List(w.context, metaV1.ListOptions{})

		if ncList != nil {
			for _, ns := range ncList.Items {
				available := false
				for k, v := range webhook.NamespaceSelector.MatchLabels {
					if ns.Labels[k] == v {
						available = true
					}
				}
				if available {
					*activeNamespaces = append(*activeNamespaces, ns.Name)
				}
			}
		}
	}
}

//retrieveValidDateCount returns remaining time of the given
//webhook's CABundle certificate.
func retrieveValidDateCount(certificate []byte) time.Duration {
	block, _ := pem.Decode(certificate)
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		log.Fatalf("x509.ParseCertificate - error occurred, detail: %v", err)
	}
	return time.Until(cert.NotAfter)
}
