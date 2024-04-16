package apigee

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

func GetSharedFlowList() ([]string, error) {
	return GetItemList(baseURL + "/sharedflows/")
}

func GetSharedflowDeployments(list []string, environment string) chan SharedflowDeployment {
	genericDeployments := GetDeployments(list, environment, "/sharedflows/")
	specificDeployments := make(chan SharedflowDeployment)

	go func() {
		for deployment := range genericDeployments {
			if sfd, ok := deployment.(*SharedflowDeployment); ok {
				specificDeployments <- *sfd
			}
		}
		close(specificDeployments)
	}()

	return specificDeployments
}

func DownloadSharedflowRevision(in chan SharedflowDeployment, environment string) {
	dirPath := filepath.Join(environment, "sharedflows")

	if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
		fmt.Printf("Error creating directory: %v\n", err)
		return
	}

	var wg sync.WaitGroup

	for sharedflow := range in {
		if sharedflow.Environment != environment {
			continue
		}

		for _, deployment := range sharedflow.Revision {
			if deployment.State != "deployed" {
				continue
			}
			wg.Add(1)

			go func(sharedflowName, deploymentName string) {
				url := baseURL + "/sharedflows/" + sharedflowName + "/revisions/" + deploymentName + "?format=bundle"

				folderName := fmt.Sprintf("%s-%s", sharedflowName, deploymentName)
				outputPath := filepath.Join(dirPath, folderName+".zip")

				// Download the file
				if err := DownloadBinaryContent(url, outputPath); err != nil {
					fmt.Printf("Error downloading %s: %v\n", url, err)
					return
				}

			}(sharedflow.Name, deployment.Name)

		}
	}

	wg.Done()

}
