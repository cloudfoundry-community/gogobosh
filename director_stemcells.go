package gogobosh

import (
	"fmt"
	"net/url"
	"time"
)

func (repo BoshDirectorRepository) GetStemcells() (stemcells []Stemcell, apiResponse ApiResponse) {
	stemcellsResponse := []stemcellResponse{}

	path := "/stemcells"
	apiResponse = repo.gateway.GetResource(repo.config.TargetURL+path, repo.config.Username, repo.config.Password, &stemcellsResponse)
	if apiResponse.IsNotSuccessful() {
		return
	}

	for _, resource := range stemcellsResponse {
		stemcells = append(stemcells, resource.ToModel())
	}

	return
}

func (repo BoshDirectorRepository) DeleteStemcell(name string, version string) (apiResponse ApiResponse) {
	path := fmt.Sprintf("/stemcells/%s/%s?force=true", name, version)
	apiResponse = repo.gateway.DeleteResource(repo.config.TargetURL+path, repo.config.Username, repo.config.Password)
	if apiResponse.IsNotSuccessful() {
		return
	}
	if !apiResponse.IsRedirection() {
		return
	}

	var taskStatus TaskStatus
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

type stemcellResponse struct {
	Name string    `json:"name"`
	Version string `json:"version"`
	Cid string     `json:"cid"`
}

func (resource stemcellResponse) ToModel() (stemcell Stemcell) {
	stemcell = Stemcell{}
	stemcell.Name = resource.Name
	stemcell.Version = resource.Version
	stemcell.Cid = resource.Cid

	return
}