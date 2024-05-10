package apigee

import (
	"encoding/json"
	"errors"
	"fmt"
)

func GetProxyList() ([]string, error) {
	return GetItemList(baseURL + "/apis/")
}

func GetProxyDeployments(list []string, environment string) chan ProxyDeployment {
	genericDeployments := GetDeployments(list, environment, "/apis/")
	specificDeployments := make(chan ProxyDeployment)

	go func() {
		for deployment := range genericDeployments {
			if pd, ok := deployment.(*ProxyDeployment); ok {
				specificDeployments <- *pd
			}
		}
		close(specificDeployments)
	}()

	return specificDeployments
}

func DeployExistingProxyRevision(name string, environment string, revision string) (*DeployedProxyResponse, error) {
	url := baseURL + "/environments/" + environment + "/apis/" + name + "/revisions/" + revision + "/deployments?action=deploy&override=true&delay=60"
	data, err := Post(url)
	if err != nil {
		if errors.Is(err, ErrBadRequest) {
			return nil, err
		} else {
			return nil, fmt.Errorf("error deploying revision %s of proxy %s to environment %s: %v", revision, name, environment, err)
		}
	}

	var response = new(DeployedProxyResponse)
	json.Unmarshal(data, response)

	return response, nil
}

func UndeployProxyRevision(name string, environment string, revision string) (*DeployedProxyResponse, error) {
	url := baseURL + "/apis/" + name + "/revisions/" + revision + "/deployments?action=undeploy&env=" + environment
	data, err := Post(url)
	if err != nil {
		if errors.Is(err, ErrBadRequest) {
			return nil, err
		} else {
			return nil, fmt.Errorf("error undeploying revision %s of proxy %s from environment %s: %v", revision, name, environment, err)
		}
	}

	var response = new(DeployedProxyResponse)
	json.Unmarshal(data, response)

	return response, nil
}
