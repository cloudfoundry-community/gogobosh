package gogobosh

import (
	"fmt"
)

func (repo BoshDirectorRepository) GetTaskStatus(taskID int) (task TaskStatus, apiResponse ApiResponse) {
	taskResponse := TaskStatusResponse{}

	path := fmt.Sprintf("/tasks/%d", taskID)
	username := "admin"
	password := "admin"
	apiResponse = repo.gateway.GetResource(repo.config.TargetURL+path, username, password, &taskResponse)
	if apiResponse.IsNotSuccessful() {
		return
	}

	task = taskResponse.ToModel()

	return
}

type TaskStatusResponse struct {
	ID int             `json:"id"`
	State string       `json:"state"`
	Description string `json:"description"`
	TimeStamp int      `json:"timestamp"`
	Result string      `json:"result"`
	User string        `json:"user"`
}

func (resource TaskStatusResponse) ToModel() (task TaskStatus) {
	task = TaskStatus{}

	task.ID = resource.ID
	task.State = resource.State
	task.Description = resource.Description
	task.TimeStamp = resource.TimeStamp
	task.Result = resource.Result
	task.User = resource.User

	return
}