package api

import (
	"fmt"
	"net/url"
	"time"
	"github.com/cloudfoundry-community/gogobosh"
	"github.com/cloudfoundry-community/gogobosh/net"
)

func (repo BoshDirectorRepository) GetDeployments() (deployments []gogobosh.Deployment, apiResponse net.ApiResponse) {
	deploymentsResponse := []deploymentResponse{}

	path := "/deployments"
	apiResponse = repo.gateway.GetResource(repo.config.TargetURL+path, repo.config.Username, repo.config.Password, &deploymentsResponse)
	if apiResponse.IsNotSuccessful() {
		return
	}

	for _, resource := range deploymentsResponse {
		deployments = append(deployments, resource.ToModel())
	}

	return
}

func (repo BoshDirectorRepository) DeleteDeployment(deploymentName string) (apiResponse net.ApiResponse) {
	path := fmt.Sprintf("/deployments/%s?force=true", deploymentName)
	apiResponse = repo.gateway.DeleteResource(repo.config.TargetURL+path, repo.config.Username, repo.config.Password)
	if apiResponse.IsNotSuccessful() {
		return
	}
	if !apiResponse.IsRedirection() {
		return
	}

	var taskStatus gogobosh.TaskStatus
	taskUrl, err := url.Parse(apiResponse.RedirectLocation)
	if err != nil {
		return
	}

	apiResponse = repo.gateway.GetResource(repo.config.TargetURL+taskUrl.Path, repo.config.Username, repo.config.Password, &taskStatus)
	if apiResponse.IsNotSuccessful() {
		return
	}

	/* Progression should be: queued, progressing, done */
	/* TODO task might fail; end states: done, error, cancelled */
	for taskStatus.State != "done" {
		time.Sleep(1)
		taskStatus, apiResponse = repo.GetTaskStatus(taskStatus.ID)
		if apiResponse.IsNotSuccessful() {
			return
		}
	}

	return
}

type deploymentResponse struct {
	Name string             `json:"name"`
	Releases []nameVersion  `json:"deployments"`
	Stemcells []nameVersion `json:"stemcells"`
}

type nameVersion struct {
	Name string    `json:"name"`
	Version string `json:"version"`
}

func (resource deploymentResponse) ToModel() (deployment gogobosh.Deployment) {
	deployment = gogobosh.Deployment{}
	deployment.Name = resource.Name
	for _, releaseResponse := range resource.Releases {
		release := gogobosh.NameVersion{}
		release.Name = releaseResponse.Name
		release.Version = releaseResponse.Version

		deployment.Releases = append(deployment.Releases, release)
	}

	for _, stemcellResponse := range resource.Stemcells {
		stemcell := gogobosh.NameVersion{}
		stemcell.Name = stemcellResponse.Name
		stemcell.Version = stemcellResponse.Version

		deployment.Stemcells = append(deployment.Stemcells, stemcell)
	}
	return
}