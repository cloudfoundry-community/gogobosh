package api

import (
	"fmt"

	"github.com/cloudfoundry-community/gogobosh/models"
	"github.com/cloudfoundry-community/gogobosh/net"
)

// GetTaskStatuses returns a list of most recent task statuses
func (repo BoshDirectorRepository) GetTaskStatuses() (tasks []models.TaskStatus, apiResponse net.ApiResponse) {
	taskResponses := []taskStatusResponse{}

	path := fmt.Sprintf("/tasks")
	apiResponse = repo.gateway.GetResource(repo.config.TargetURL+path, repo.config.Username, repo.config.Password, &taskResponses)
	if apiResponse.IsNotSuccessful() {
		return
	}

	for _, resource := range taskResponses {
		tasks = append(tasks, resource.ToModel())
	}

	return
}

// GetTaskStatusesWithLimit returns a max of 'limit' task statuses
func (repo BoshDirectorRepository) GetTaskStatusesWithLimit(limit int) (tasks []models.TaskStatus, apiResponse net.ApiResponse) {
	taskResponses := []taskStatusResponse{}

	path := fmt.Sprintf("/tasks?limit=%d", limit)
	apiResponse = repo.gateway.GetResource(repo.config.TargetURL+path, repo.config.Username, repo.config.Password, &taskResponses)
	if apiResponse.IsNotSuccessful() {
		return
	}

	for _, resource := range taskResponses {
		tasks = append(tasks, resource.ToModel())
	}

	return
}

// GetTaskStatus returns details of a specific task to allow polling for state change
func (repo BoshDirectorRepository) GetTaskStatus(taskID int) (task models.TaskStatus, apiResponse net.ApiResponse) {
	taskResponse := taskStatusResponse{}

	path := fmt.Sprintf("/tasks/%d", taskID)
	apiResponse = repo.gateway.GetResource(repo.config.TargetURL+path, repo.config.Username, repo.config.Password, &taskResponse)
	if apiResponse.IsNotSuccessful() {
		return
	}

	task = taskResponse.ToModel()

	return
}

type taskStatusResponse struct {
	ID          int    `json:"id"`
	State       string `json:"state"`
	Description string `json:"description"`
	TimeStamp   int    `json:"timestamp"`
	Result      string `json:"result"`
	User        string `json:"user"`
}

func (resource taskStatusResponse) ToModel() (task models.TaskStatus) {
	task = models.TaskStatus{}

	task.ID = resource.ID
	task.State = resource.State
	task.Description = resource.Description
	task.TimeStamp = resource.TimeStamp
	task.Result = resource.Result
	task.User = resource.User

	return
}
