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
	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/tools/clientcmd/api"
	"strings"
)

var (
	errNoContext = errors.New(fmt.Sprintf("no context provided, use '%q' to select", "kubectl config use-context <context>"))
)

type ViewWebhookOptions struct {
	configFlags *genericclioptions.ConfigFlags

	resultingContext     *api.Context
	resultingContextName string

	userSpecifiedCluster   string
	userSpecifiedContext   string
	userSpecifiedAuthInfo  string
	userSpecifiedNamespace string

	rawConfig    api.Config
	args         []string

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

// Complete sets all information required for updating the current context
func (o *ViewWebhookOptions) Complete(cmd *cobra.Command, args []string) error {
	o.args = args

	var err error
	o.rawConfig, err = o.configFlags.ToRawKubeConfigLoader().RawConfig()
	if err != nil {
		return err
	}

	o.userSpecifiedNamespace, err = cmd.Flags().GetString("namespace")
	if err != nil {
		return err
	}
	if len(args) > 0 {
		if len(o.userSpecifiedNamespace) > 0 {
			return fmt.Errorf("cannot specify both a --namespace value and a new namespace argument")
		}

		o.userSpecifiedNamespace = args[0]
	}

	// if no namespace argument or flag value was specified, then there
	// is no need to generate a resulting context
	if len(o.userSpecifiedNamespace) == 0 {
		return nil
	}

	o.userSpecifiedContext, err = cmd.Flags().GetString("context")
	if err != nil {
		return err
	}

	o.userSpecifiedCluster, err = cmd.Flags().GetString("cluster")
	if err != nil {
		return err
	}

	o.userSpecifiedAuthInfo, err = cmd.Flags().GetString("user")
	if err != nil {
		return err
	}

	currentContext, exists := o.rawConfig.Contexts[o.rawConfig.CurrentContext]
	if !exists {
		return errNoContext
	}

	o.resultingContext = api.NewContext()
	o.resultingContext.Cluster = currentContext.Cluster
	o.resultingContext.AuthInfo = currentContext.AuthInfo

	// if a target context is explicitly provided by the user,
	// use that as our reference for the final, resulting context
	if len(o.userSpecifiedContext) > 0 {
		o.resultingContextName = o.userSpecifiedContext
		if userCtx, exists := o.rawConfig.Contexts[o.userSpecifiedContext]; exists {
			o.resultingContext = userCtx.DeepCopy()
		}
	}

	// override context info with user provided values
	o.resultingContext.Namespace = o.userSpecifiedNamespace

	if len(o.userSpecifiedCluster) > 0 {
		o.resultingContext.Cluster = o.userSpecifiedCluster
	}
	if len(o.userSpecifiedAuthInfo) > 0 {
		o.resultingContext.AuthInfo = o.userSpecifiedAuthInfo
	}

	// generate a unique context name based on its new values if
	// user did not explicitly request a context by name
	if len(o.userSpecifiedContext) == 0 {
		o.resultingContextName = generateContextName(o.resultingContext)
	}

	return nil
}

func generateContextName(fromContext *api.Context) string {
	name := fromContext.Namespace
	if len(fromContext.Cluster) > 0 {
		name = fmt.Sprintf("%s/%s", name, fromContext.Cluster)
	}
	if len(fromContext.AuthInfo) > 0 {
		cleanAuthInfo := strings.Split(fromContext.AuthInfo, "/")[0]
		name = fmt.Sprintf("%s/%s", name, cleanAuthInfo)
	}

	return name
}

// Validate ensures that all required args and flags are provided
func (o *ViewWebhookOptions) Validate() error {
	if len(o.rawConfig.CurrentContext) == 0 {
		return errNoContext
	}

	return nil
}

// Run lists all available webhooks on a user's KUBECONFIG or updates the
// current context based on a provided namespace.
func (o *ViewWebhookOptions) Run() error {

	return nil
}
