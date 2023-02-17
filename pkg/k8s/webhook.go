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
	"log"
	"time"

	"github.com/Trendyol/kubectl-view-webhook/pkg/printer"
	v1 "k8s.io/api/admissionregistration/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	typedV1 "k8s.io/client-go/kubernetes/typed/admissionregistration/v1"
	typedCoreV1 "k8s.io/client-go/kubernetes/typed/core/v1"
)

type WebHookClient struct {
	client  *kubernetes.Clientset
	wClient typedV1.MutatingWebhookConfigurationInterface
	vClient typedV1.ValidatingWebhookConfigurationInterface
	nClient typedCoreV1.NamespaceInterface
	context context.Context
}

// NewWebHookClient constructs a new WebHookClient with the specified output
// of *kubernetes.Clientset
func NewWebHookClient(client *kubernetes.Clientset) *WebHookClient {
	return &WebHookClient{
		client:  client,
		wClient: client.AdmissionregistrationV1().MutatingWebhookConfigurations(),
		vClient: client.AdmissionregistrationV1().ValidatingWebhookConfigurations(),
		nClient: client.CoreV1().Namespaces(),
		context: context.Background(),
	}
}

type Resource struct {
	Name       string
	Operations []string
}

// Run
// args[0]: self executable
// args[1]: may be 'webhookname' or '--kubeconfig'
func (w *WebHookClient) Run(args []string) (*printer.PrintModel, error) {
	var items []printer.PrintItem

	if len(args) == 0 {
		mutatingWebhookConfigurationList, _ := w.wClient.List(w.context, metaV1.ListOptions{})

		validatingWebhookConfigurationList, _ := w.vClient.List(w.context, metaV1.ListOptions{})

		for _, mwc := range mutatingWebhookConfigurationList.Items {
			w.fillMutatingWebhookConfigurations(mwc, &items)
		}
		for _, mwc := range validatingWebhookConfigurationList.Items {
			w.fillValidatingWebhookConfigurations(mwc, &items)
		}
	} else {
		mutatingWebhookConfiguration, _ := w.wClient.Get(w.context, args[0], metaV1.GetOptions{})

		validatingWebhookConfiguration, _ := w.vClient.Get(w.context, args[0], metaV1.GetOptions{})

		w.fillMutatingWebhookConfigurations(*mutatingWebhookConfiguration, &items)
		w.fillValidatingWebhookConfigurations(*validatingWebhookConfiguration, &items)
	}

	return &printer.PrintModel{
		Items: items,
	}, nil
}

func (w *WebHookClient) fillMutatingWebhookConfigurations(mwc v1.MutatingWebhookConfiguration, items *[]printer.PrintItem) {
	item := printer.PrintItem{
		Kind: "Mutating",
		Name: mwc.Name, //TODO: typeMeta nil
	}

	for _, webhook := range mwc.Webhooks {
		var activeNamespaces []string
		w.fillActiveNamespacesForMutating(webhook, &activeNamespaces)

		webhookItem := printer.PrintWebhookItem{
			Name: webhook.Name,
		}

		if webhook.ClientConfig.Service != nil {
			ss := w.GenerateServiceItem(webhook.ClientConfig.Service.Namespace, webhook.ClientConfig.Service.Name, webhook.ClientConfig.Service.Path, webhook.ClientConfig.Service.Port)
			webhookItem.Service = ss
		}

		item.Webhook = webhookItem
		resources := w.fillRulesForMutating(webhook)

		item.ResourceModels = resources
		item.ValidUntil = retrieveValidDateCount(webhook.ClientConfig.CABundle)
		item.ActiveNamespaces = activeNamespaces
		*items = append(*items, item)
	}
}
func (w *WebHookClient) fillValidatingWebhookConfigurations(mwc v1.ValidatingWebhookConfiguration, items *[]printer.PrintItem) {
	item := printer.PrintItem{
		Kind: "Validating",
		Name: mwc.Name, //TODO: typeMeta nil
	}

	for _, webhook := range mwc.Webhooks {
		var activeNamespaces []string
		w.fillActiveNamespacesForValidating(webhook, &activeNamespaces)

		webhookItem := printer.PrintWebhookItem{
			Name: webhook.Name,
		}

		if webhook.ClientConfig.Service != nil {
			ss := w.GenerateServiceItem(webhook.ClientConfig.Service.Namespace, webhook.ClientConfig.Service.Name, webhook.ClientConfig.Service.Path, webhook.ClientConfig.Service.Port)
			webhookItem.Service = ss
		}

		item.Webhook = webhookItem
		resources := w.fillRulesForValidating(webhook)

		item.ResourceModels = resources
		item.ValidUntil = retrieveValidDateCount(webhook.ClientConfig.CABundle)
		item.ActiveNamespaces = activeNamespaces
		*items = append(*items, item)
	}
}
func (w *WebHookClient) fillRulesForMutating(webhook v1.MutatingWebhook) []printer.ResourceModel {
	var resources []printer.ResourceModel

	for _, rule := range webhook.Rules {
		var ops, rs []string

		for _, op := range rule.Operations {
			ops = append(ops, string(op))
		}

		rs = append(rs, rule.Resources...)

		resources = append(resources, printer.ResourceModel{
			Operations: ops,
			Resources:  rs,
		})
	}
	return resources
}
func (w *WebHookClient) fillRulesForValidating(webhook v1.ValidatingWebhook) []printer.ResourceModel {
	var resources []printer.ResourceModel
	var ops, rs []string
	for _, rule := range webhook.Rules {

		for _, op := range rule.Operations {
			ops = append(rs, string(op))
		}

		rs = append(rs, rule.Resources...)

		resources = append(resources, printer.ResourceModel{
			Operations: ops,
			Resources:  rs,
		})
	}
	return resources
}
func (w *WebHookClient) fillActiveNamespacesForMutating(webhook v1.MutatingWebhook, activeNamespaces *[]string) {
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
func (w *WebHookClient) fillActiveNamespacesForValidating(webhook v1.ValidatingWebhook, activeNamespaces *[]string) {
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

func (w *WebHookClient) GenerateServiceItem(ns, name string, path *string, port *int32) printer.PrintServiceItem {
	service := w.client.CoreV1().Services(ns)

	ss, err := service.Get(w.context, name, metaV1.GetOptions{})

	result := printer.PrintServiceItem{
		Name:      name,
		Namespace: ns,
		Path:      path,
		Ports:     nil,
	}

	if err != nil {
		result.Found = false
		return result
	}

	for _, p := range ss.Spec.Ports {
		result.Ports = append(result.Ports, printer.PrintServicePortItem{
			Port:       p.Port,
			TargetPort: p.TargetPort.IntVal,
			Protocol:   string(p.Protocol),
		})
	}

	result.Found = true
	result.ClusterIP = ss.Spec.ClusterIP
	result.Type = string(ss.Spec.Type)

	return result
}

// retrieveValidDateCount returns remaining time of the given
// webhook's CABundle certificate.
func retrieveValidDateCount(certificate []byte) time.Duration {
	if certificate == nil {
		return 0
	}
	block, _ := pem.Decode(certificate)
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		log.Fatalf("x509.ParseCertificate - error occurred, detail: %v", err)
	}
	return time.Until(cert.NotAfter)
}
