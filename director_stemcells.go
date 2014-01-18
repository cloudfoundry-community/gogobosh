package gogobosh

type StemcellResponse struct {
	Name string    `json:"name"`
	Version string `json:"version"`
	Cid string     `json:"cid"`
}

func (resource StemcellResponse) ToModel() (stemcell Stemcell) {
	stemcell = Stemcell{}
	stemcell.Name = resource.Name
	stemcell.Version = resource.Version
	stemcell.Cid = resource.Cid

	return
}