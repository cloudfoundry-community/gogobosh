package gogobosh

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