package gogobosh

import (
	"fmt"
)

func (repo BoshDirectorRepository) GetTaskStatuses() (tasks []TaskStatus, apiResponse ApiResponse) {
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

func (repo BoshDirectorRepository) GetTaskStatus(taskID int) (task TaskStatus, apiResponse ApiResponse) {
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
	ID int             `json:"id"`
	State string       `json:"state"`
	Description string `json:"description"`
	TimeStamp int      `json:"timestamp"`
	Result string      `json:"result"`
	User string        `json:"user"`
}

func (resource taskStatusResponse) ToModel() (task TaskStatus) {
	task = TaskStatus{}

	task.ID = resource.ID
	task.State = resource.State
	task.Description = resource.Description
	task.TimeStamp = resource.TimeStamp
	task.Result = resource.Result
	task.User = resource.User

	return
}