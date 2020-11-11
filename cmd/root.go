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

package cmd

import (
	"context"
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd/api"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	errNoKubeConfig = errors.New(fmt.Sprintf("no kubeconfig provided"))
)

type ViewWebhookOptions struct {
	configFlags *genericclioptions.ConfigFlags

	restConfig *rest.Config
	rawConfig  api.Config
	args       []string

	genericclioptions.IOStreams
}

// NewViewWebhookOptions provides an instance of ViewWebhookOptions with default values
func NewViewWebhookOptions(streams genericclioptions.IOStreams) *ViewWebhookOptions {
	return &ViewWebhookOptions{
		configFlags: genericclioptions.NewConfigFlags(true),
		IOStreams:   streams,
	}
}

// NewCmdViewWebhook provides a cobra command wrapping NamespaceOptions
func NewCmdViewWebhook(streams genericclioptions.IOStreams) *cobra.Command {
	o := NewViewWebhookOptions(streams)

	cmd := &cobra.Command{
		Use:   "kubectl view-webhook [flags]", //TODO: KIND? (Validating, Admission, etc.)
		Short: "Visualize your webhook configurations of the Kubernetes resource",
		Long:  `Visualize your webhook configurations of the Kubernetes resource`,
		Example: fmt.Sprintf(`
%[1]s view-webhook
`, "kubectl"),
		SilenceErrors: false,
		SilenceUsage:  false,
		RunE: func(c *cobra.Command, args []string) error {
			if err := o.Complete(c, args); err != nil {
				return err
			}

			if err := o.Validate(); err != nil {
				return err
			}

			if err := o.Run(); err != nil {
				return err
			}

			return nil
		},
	}

	o.configFlags.AddFlags(cmd.Flags())

	return cmd
}

// Complete sets all information required for viewing webhook
func (o *ViewWebhookOptions) Complete(cmd *cobra.Command, args []string) error {
	o.args = args

	var err error
	o.restConfig, err = o.configFlags.ToRESTConfig()
	if err != nil {
		return err
	}

	o.rawConfig, err = o.configFlags.ToRawKubeConfigLoader().RawConfig()
	if err != nil {
		return err
	}

	return nil
}

// Validate ensures that all required args and flags are provided
func (o *ViewWebhookOptions) Validate() error {
	if len(o.rawConfig.CurrentContext) == 0 {
		return errNoKubeConfig
	}

	return nil
}

// Run lists all available webhooks on a user's KUBECONFIG or updates the
// current context based on a provided namespace.
func (o *ViewWebhookOptions) Run() error {
	ctx := context.TODO()
	// create the ClientSet
	clientset, err := kubernetes.NewForConfig(o.restConfig)
	if err != nil {
		return err
	}

	mutatingWebhookClient := clientset.AdmissionregistrationV1beta1().MutatingWebhookConfigurations()
	//validatingWebhookClient := clientset.AdmissionregistrationV1beta1().ValidatingWebhookConfigurations()

	mutatingWebhookConfigurationList, err := mutatingWebhookClient.List(ctx, metaV1.ListOptions{})
	if err != nil {
		return err
	}

	for _, mwc := range mutatingWebhookConfigurationList.Items {
		fmt.Printf("Mutating webhook name: %s \n", mwc.Name)
	}

	return nil
}
