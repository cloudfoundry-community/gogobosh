package gogobosh

type DirectorRepository interface {
	GetInfo() (directorInfo DirectorInfo, apiResponse ApiResponse)
	GetStemcells() (stemcells []Stemcell, apiResponse ApiResponse)
	GetReleases() (releases []Release, apiResponse ApiResponse)
	GetDeployments() (deployments []Deployment, apiResponse ApiResponse)
}

type BoshDirectorRepository struct {
	config  *Director
	gateway Gateway
}


func NewBoshDirectorRepository(config *Director, gateway Gateway) (repo BoshDirectorRepository) {
	repo.config = config
	repo.gateway = gateway
	return
}

