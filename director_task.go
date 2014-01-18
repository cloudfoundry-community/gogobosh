package gogobosh

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