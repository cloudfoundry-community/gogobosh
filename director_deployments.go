package gogobosh

type DeploymentResponse struct {
	Name string             `json:"name"`
	Releases []nameVersion  `json:"deployments"`
	Stemcells []nameVersion `json:"stemcells"`
}

type nameVersion struct {
	Name string    `json:"name"`
	Version string `json:"version"`
}

func (resource DeploymentResponse) ToModel() (deployment Deployment) {
	deployment = Deployment{}
	deployment.Name = resource.Name
	for _, releaseResponse := range resource.Releases {
		release := NameVersion{}
		release.Name = releaseResponse.Name
		release.Version = releaseResponse.Version

		deployment.Releases = append(deployment.Releases, release)
	}

	for _, stemcellResponse := range resource.Stemcells {
		stemcell := NameVersion{}
		stemcell.Name = stemcellResponse.Name
		stemcell.Version = stemcellResponse.Version

		deployment.Stemcells = append(deployment.Stemcells, stemcell)
	}
	return
}