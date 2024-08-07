/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>

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
	"fmt"

	"github.com/mc-admin-cli/mcc/src/common"
	"github.com/spf13/cobra"
)

// stopCmd represents the stop command
var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop M-CMP System",
	Long:  `Stop M-CMP System`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("\n[Stop M-CMP]")
		fmt.Println()

		if K8sFilePath == "" {
			fmt.Println("file is required")
		} else {
			// var cmdStr string

			cmdStr := fmt.Sprintf("helm uninstall --namespace %s %s", K8sNamespace, HelmReleaseName)
			common.SysCall(cmdStr)

		}

	},
}

func init() {
	k8sCmd.AddCommand(stopCmd)

	pf := stopCmd.PersistentFlags()
	pf.StringVarP(&K8sFilePath, "file", "f", DefaultKubernetesConfig, "User-defined configuration file")
	//	cobra.MarkFlagRequired(pf, "file")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// stopCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// stopCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
