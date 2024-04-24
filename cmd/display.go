package cmd

import "github.com/stedmanson/apigee-cli/internal/file"

// proxyName, revision, environment and file name currently sourced from global variables
// todo: clean up this function
func collateProxyList() ([]proxydeployment, error) {
	var list []proxydeployment

	if filename == "" {
		list = append(list, proxydeployment{
			name:        proxyName,
			environment: environment,
			revision:    revision,
		})
	} else {

		data, err := file.ReadCSV(filename)
		if err != nil {
			return nil, err
		}

		for _, d := range data {
			list = append(list, proxydeployment{
				name:        d.Name,
				environment: environment,
				revision:    d.Revision,
			})
		}
	}

	return list, nil
}
