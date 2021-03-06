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

package main

import (
	"github.com/Trendyol/kubectl-view-webhook/cmd"
	"github.com/spf13/pflag"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"os"
)

var (
	version, commit, date string
)

func main() {
	flags := pflag.NewFlagSet("kubectl-view-webhook", pflag.ExitOnError)
	pflag.CommandLine = flags

	root := cmd.NewCmdViewWebhook(genericclioptions.IOStreams{In: os.Stdin, Out: os.Stdout, ErrOut: os.Stderr}, version, commit, date)
	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}
