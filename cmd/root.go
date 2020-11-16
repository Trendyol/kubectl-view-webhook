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
	"errors"
	"fmt"
	"github.com/Trendyol/kubectl-view-webhook/pkg/k8s"
	"github.com/Trendyol/kubectl-view-webhook/pkg/printer"
	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"os"
	"path/filepath"
)

type ViewWebhookOptions struct {
	configFlags *genericclioptions.ConfigFlags

	restConfig *rest.Config
	kubeconfig string
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
func NewCmdViewWebhook(streams genericclioptions.IOStreams, version, commit, date string) *cobra.Command {
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

			if err := o.Run(); err != nil {
				return err
			}

			return nil
		},
		Version: func(version, commit, date string) string {
			return fmt.Sprintf(`
Version: %s
Commit: %s
Date: %s`, version, commit, date)
		}(version, commit, date),
	}

	o.configFlags.AddFlags(cmd.Flags())

	return cmd
}

// Complete sets all information required for viewing webhook
func (o *ViewWebhookOptions) Complete(cmd *cobra.Command, args []string) error {
	o.args = args

	kubeconfig, err := cmd.Flags().GetString("kubeconfig")
	if err != nil {
		return err
	}

	if kubeconfig == "" {
		// fallback to kubeconfig
		kubeconfig = filepath.Join(homedir.HomeDir(), ".kube", "config")
		if envvar := os.Getenv("KUBECONFIG"); len(envvar) > 0 {
			kubeconfig = envvar
		}
	}

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		fmt.Printf("The kubeconfig cannot be loaded: %v\n", err)
		os.Exit(1)
	}

	o.restConfig = config
	o.kubeconfig = kubeconfig
	return nil
}

// Validate ensures that all required args and flags are provided
func (o *ViewWebhookOptions) Validate() error {
	if len(o.args) > 2 {
		return errors.New("more than one argument supplied , you can only give one argument for the webhook name")
	}
	return nil
}

// Run lists all available webhooks on a user's KUBECONFIG or updates the
// current context based on a provided namespace.
func (o *ViewWebhookOptions) Run() error {
	p := printer.NewPrinter(o.Out)

	// create the ClientSet from restConfig
	clientSet, err := kubernetes.NewForConfig(o.restConfig)
	if err != nil {
		return err
	}

	//validatingWebhookClient := clientset.AdmissionregistrationV1beta1().ValidatingWebhookConfigurations()

	mw := k8s.NewWebHookClient(clientSet)
	model, err := mw.Run(o.args)

	if err != nil {
		return err
	}

	p.Print(model)

	return nil
}
