package api

import (
	"fmt"
	"net/url"
	"time"
	"github.com/cloudfoundry-community/gogobosh"
	"github.com/cloudfoundry-community/gogobosh/net"
)

func (repo BoshDirectorRepository) GetReleases() (releases []gogobosh.Release, apiResponse net.ApiResponse) {
	releasesResponse := []releaseResponse{}

	path := "/releases"
	apiResponse = repo.gateway.GetResource(repo.config.TargetURL+path, repo.config.Username, repo.config.Password, &releasesResponse)
	if apiResponse.IsNotSuccessful() {
		return
	}

	for _, resource := range releasesResponse {
		releases = append(releases, resource.ToModel())
	}

	return
}

func (repo BoshDirectorRepository) DeleteReleases(name string) (apiResponse net.ApiResponse) {
	path := fmt.Sprintf("/releases/%s?force=true", name)
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

func (repo BoshDirectorRepository) DeleteRelease(name string, version string) (apiResponse net.ApiResponse) {
	path := fmt.Sprintf("/releases/%s?force=true&version=%s", name, version)
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

type releaseResponse struct {
	Name string                       `json:"name"`
	Versions []releaseVersionResponse `json:"release_versions"`
}

type releaseVersionResponse struct {
	Version string          `json:"version"`
	CommitHash string       `json:"commit_hash"`
	UncommittedChanges bool `json:"uncommitted_changes"`
	CurrentlyDeployed bool  `json:"currently_deployed"`
}

func (resource releaseResponse) ToModel() (release gogobosh.Release) {
	release = gogobosh.Release{}
	release.Name = resource.Name
	for _, versionResponse := range resource.Versions {
		version := gogobosh.ReleaseVersion{}
		version.Version = versionResponse.Version
		version.CommitHash = versionResponse.CommitHash
		version.UncommittedChanges = versionResponse.UncommittedChanges
		version.CurrentlyDeployed = versionResponse.CurrentlyDeployed

		release.Versions = append(release.Versions, version)
	}
	return
}