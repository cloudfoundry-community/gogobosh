package gogobosh

type VMStatusResponse struct {
	JobName string `json:"job"`
	Index int      `json:"index"`
}

func (resource VMStatusResponse) ToModel() (status VMStatus) {
	status = VMStatus{}
	status.JobName = resource.JobName
	status.Index = resource.Index

	return
}