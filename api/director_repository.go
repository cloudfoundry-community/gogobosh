package api

import (
	"github.com/cloudfoundry-community/gogobosh"
	"github.com/cloudfoundry-community/gogobosh/net"
)

// DirectorRepository is the interface for accessing a BOSH director
type DirectorRepository interface {
	GetInfo() (directorInfo gogobosh.DirectorInfo, apiResponse net.ApiResponse)

	GetStemcells() (stemcells []gogobosh.Stemcell, apiResponse net.ApiResponse)
	DeleteStemcell(name string, version string) (apiResponse net.ApiResponse)

	GetReleases() (releases []gogobosh.Release, apiResponse net.ApiResponse)
	DeleteReleases(name string) (apiResponse net.ApiResponse)
	DeleteRelease(name string, version string) (apiResponse net.ApiResponse)

	GetDeployments() (deployments []gogobosh.Deployment, apiResponse net.ApiResponse)
	GetDeploymentManifest(deploymentName string) (manifest *gogobosh.DeploymentManifest, apiResponse net.ApiResponse)
	DeleteDeployment(deploymentName string) (apiResponse net.ApiResponse)
	ListDeploymentVMs(deploymentName string) (deploymentVMs []gogobosh.DeploymentVM, apiResponse net.ApiResponse)
	FetchVMsStatus(deploymentName string) (vmsStatus []gogobosh.VMStatus, apiResponse net.ApiResponse)

	GetTaskStatuses() (task []gogobosh.TaskStatus, apiResponse net.ApiResponse)
	GetTaskStatus(taskID int) (task gogobosh.TaskStatus, apiResponse net.ApiResponse)
}

// BoshDirectorRepository represents a Director
type BoshDirectorRepository struct {
	config  *gogobosh.Director
	gateway net.Gateway
}

// NewBoshDirectorRepository is a constructor for a BoshDirectorRepository
func NewBoshDirectorRepository(config *gogobosh.Director, gateway net.Gateway) (repo BoshDirectorRepository) {
	repo.config = config
	repo.gateway = gateway
	return
}
