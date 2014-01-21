package gogobosh

type DirectorRepository interface {
	GetInfo() (directorInfo DirectorInfo, apiResponse ApiResponse)

	GetStemcells() (stemcells []Stemcell, apiResponse ApiResponse)
	DeleteStemcell(name string, version string) (apiResponse ApiResponse)

	GetReleases() (releases []Release, apiResponse ApiResponse)
	DeleteReleases(name string) (apiResponse ApiResponse)
	DeleteRelease(name string, version string) (apiResponse ApiResponse)

	GetDeployments() (deployments []Deployment, apiResponse ApiResponse)
	DeleteDeployment(deploymentName string) (apiResponse ApiResponse)

	GetTaskStatus(taskID int) (task TaskStatus, apiResponse ApiResponse)
	ListDeploymentVMs(deploymentName string) (deploymentVMs []DeploymentVM, apiResponse ApiResponse)
	FetchVMsStatus(deploymentName string) (vmsStatus []VMStatus, apiResponse ApiResponse)
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

