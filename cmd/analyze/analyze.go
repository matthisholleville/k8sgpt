/*
Copyright 2023 The K8sGPT Authors.
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

package analyze

import (
	"fmt"
	"os"
	"time"

	"github.com/briandowns/spinner"
	"github.com/fatih/color"
	"github.com/k8sgpt-ai/k8sgpt/pkg/analysis"
	"github.com/spf13/cobra"
)

var (
	explain        bool
	backend        string
	output         string
	filters        []string
	language       string
	nocache        bool
	namespace      string
	anonymize      bool
	alert          string
	maxConcurrency int
	withDoc        bool
)

// AnalyzeCmd represents the problems command
var AnalyzeCmd = &cobra.Command{
	Use:     "analyze",
	Aliases: []string{"analyse"},
	Short:   "This command will find problems within your Kubernetes cluster",
	Long: `This command will find problems within your Kubernetes cluster and
	provide you with a list of issues that need to be resolved`,
	Run: func(cmd *cobra.Command, args []string) {

		if alert != "" {
			explain = true
		}

		// Create analysis configuration first.
		config, err := analysis.NewAnalysis(
			alert,
			backend,
			explain,
			filters,
			language,
			maxConcurrency,
			namespace,
			nocache,
			withDoc,
		)
		if err != nil {
			color.Red("Error: %v", err)
			os.Exit(1)
		}

		defer config.Close()

		s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
		done := make(chan bool)

		s.Suffix = " Analysis in progress..."
		s.Start()
		go func() {
			config.RunAnalysis()
			done <- true
		}()

		<-done
		s.Stop()

		if alert != "" {
			s.Suffix = " Evaluation in progress..."
			s.Start()
			go func() {
				if err := config.EvaluateResult(); err != nil {
					color.Red("Error: %v", err)
					os.Exit(1)
				}
				done <- true
			}()
			<-done
			s.Stop()
		}

		if explain {
			if err := config.GetAIResults(output, anonymize); err != nil {
				color.Red("Error: %v", err)
				os.Exit(1)
			}
		}

		// print results
		output, err := config.PrintOutput(output)
		if err != nil {
			color.Red("Error: %v", err)
			os.Exit(1)
		}
		fmt.Println(string(output))
	},
}

func init() {

	// namespace flag
	AnalyzeCmd.Flags().StringVarP(&namespace, "namespace", "n", "", "Namespace to analyze")
	// no cache flag
	AnalyzeCmd.Flags().BoolVarP(&nocache, "no-cache", "c", false, "Do not use cached data")
	// anonymize flag
	AnalyzeCmd.Flags().BoolVarP(&anonymize, "anonymize", "a", false, "Anonymize data before sending it to the AI backend. This flag masks sensitive data, such as Kubernetes object names and labels, by replacing it with a key. However, please note that this flag does not currently apply to events.")
	// alert flag
	AnalyzeCmd.Flags().StringVarP(&alert, "alert", "l", "", "Alert to be analyzed")
	// array of strings flag
	AnalyzeCmd.Flags().StringSliceVarP(&filters, "filter", "f", []string{}, "Filter for these analyzers (e.g. Pod, PersistentVolumeClaim, Service, ReplicaSet)")
	// explain flag
	AnalyzeCmd.Flags().BoolVarP(&explain, "explain", "e", false, "Explain the problem to me")
	// add flag for backend
	AnalyzeCmd.Flags().StringVarP(&backend, "backend", "b", "openai", "Backend AI provider")
	// output as json
	AnalyzeCmd.Flags().StringVarP(&output, "output", "o", "text", "Output format (text, json)")
	// add language options for output
	AnalyzeCmd.Flags().StringVarP(&language, "language", "l", "english", "Languages to use for AI (e.g. 'English', 'Spanish', 'French', 'German', 'Italian', 'Portuguese', 'Dutch', 'Russian', 'Chinese', 'Japanese', 'Korean')")
	// add max concurrency
	AnalyzeCmd.Flags().IntVarP(&maxConcurrency, "max-concurrency", "m", 10, "Maximum number of concurrent requests to the Kubernetes API server")
	// kubernetes doc flag
	AnalyzeCmd.Flags().BoolVarP(&withDoc, "with-doc", "d", false, "Give me the official documentation of the involved field")
}
