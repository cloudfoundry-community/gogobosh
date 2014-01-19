package gogobosh

func (repo BoshDirectorRepository) GetReleases() (releases []Release, apiResponse ApiResponse) {
	releasesResponse := []ReleaseResponse{}

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

type ReleaseResponse struct {
	Name string                       `json:"name"`
	Versions []releaseVersionResponse `json:"release_versions"`
}

type releaseVersionResponse struct {
	Version string          `json:"version"`
	CommitHash string       `json:"commit_hash"`
	UncommittedChanges bool `json:"uncommitted_changes"`
	CurrentlyDeployed bool  `json:"currently_deployed"`
}

func (resource ReleaseResponse) ToModel() (stemcell Release) {
	stemcell = Release{}
	stemcell.Name = resource.Name
	for _, versionResponse := range resource.Versions {
		version := ReleaseVersion{}
		version.Version = versionResponse.Version
		version.CommitHash = versionResponse.CommitHash
		version.UncommittedChanges = versionResponse.UncommittedChanges
		version.CurrentlyDeployed = versionResponse.CurrentlyDeployed

		stemcell.Versions = append(stemcell.Versions, version)
	}
	return
}