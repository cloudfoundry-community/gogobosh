package gogobosh

type VMStatusResponse struct {
	JobName string  `json:"job_name"`
	Index int       `json:"index"`
	JobState string `json:"job_state"`
	VMCid string    `json:"vm_cid"`
	AgentID string  `json:"agent_id"`
	ResourcePool string `json:"resource_pool"`
	IPs []string    `json:"ips"`
	DNSs []string    `json:"dns"`
}

func (resource VMStatusResponse) ToModel() (status VMStatus) {
	status = VMStatus{}
	status.JobName  = resource.JobName
	status.Index    = resource.Index
	status.JobState = resource.JobState
    status.VMCid    = resource.VMCid
    status.AgentID  = resource.AgentID
    status.ResourcePool = resource.ResourcePool

    status.IPs = resource.IPs
    status.DNSs = resource.DNSs

	return
}
