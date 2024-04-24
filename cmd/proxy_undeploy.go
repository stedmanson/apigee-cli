/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"sync"

	"github.com/davecgh/go-spew/spew"
	"github.com/spf13/cobra"
	"github.com/stedmanson/apigee-cli/internal/apigee"
	"github.com/stedmanson/apigee-cli/internal/output"
)

// proxyUnDeployCmd represents the proxy:deploy command
var proxyUnDeployCmd = &cobra.Command{
	Use:   "undeploy",
	Short: "deploy an indivdual or collection of proxies to an environment",
	Long: `Provide either a proxy name and version for an individal or a headerless csv
	in the format "name, revision" to deploy proxies to the provided environment.`,
	Run: func(cmd *cobra.Command, args []string) {

		list, err := collateProxyList()
		if err != nil {
			fmt.Println("Error: error reading file: " + filename)
			os.Exit(1)
		}

		results := make(chan []string, len(list)) // Initialize the channel with a buffer

		var wg sync.WaitGroup

		for _, proxy := range list {
			wg.Add(1)
			go undeployProxy(proxy.name, proxy.environment, proxy.revision, results, &wg)
		}

		// Close the channel once all goroutines have finished
		go func() {
			wg.Wait()
			close(results)
		}()

		var data [][]string
		for result := range results {
			data = append(data, result)
		}

		headers := []string{"Proxy Name", "Revision", "Environment"}
		output.DisplayAsTable(headers, data)

	},
	PreRun: func(cmd *cobra.Command, args []string) {
		// Check if the environment flag was set by the user
		if environment == "" {
			fmt.Println("Error: --env flag is required")
			os.Exit(1)
		}

		if filename == "" {
			if proxyName == "" {
				fmt.Println("Error: --proxy flag is required")
				os.Exit(1)
			}

			if revision == "" {
				fmt.Println("Error: --rev flag is required")
				os.Exit(1)
			}
		}

		if filename == "" && proxyName == "" && revision == "" {
			fmt.Println("Error: --filename flag is required unless --proxy and --rev are provided")
			os.Exit(1)
		}
	},
}

func init() {
	proxyCmd.AddCommand(proxyUnDeployCmd)
}

func undeployProxy(name string, environment string, revision string, results chan []string, wg *sync.WaitGroup) {
	defer wg.Done()
	resp, err := apigee.UndeployProxyRevision(name, environment, revision)
	if err != nil {
		spew.Dump(err)
		os.Exit(1)
	}

	results <- []string{resp.APIProxy, resp.Revision, resp.Environment}
}
