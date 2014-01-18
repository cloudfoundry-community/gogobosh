package gogobosh

type VMStatusResponse struct {
	JobName string  `json:"job_name"`
	Index int       `json:"index"`
	JobState string `json:"job_state"`
	VMCid string    `json:"vm_cid"`
	AgentID string  `json:"agent_id"`
	IPs []string    `json:"ips"`
	DNSs []string   `json:"dns"`
	ResourcePool string     `json:"resource_pool"`
	ResurrectionPaused bool `json:"resurrection_paused"`
	Vitals vmStatusVitalsResponse `json:"vitals"`
}

type vmStatusVitalsResponse struct {
	Load []string
/*	CPU
	Memory
	Swap
	Disk
*/}

func (resource VMStatusResponse) ToModel() (status VMStatus) {
	status = VMStatus{}
	status.JobName  = resource.JobName
	status.Index    = resource.Index
	status.JobState = resource.JobState
    status.VMCid    = resource.VMCid
    status.AgentID  = resource.AgentID
    status.ResourcePool = resource.ResourcePool
    status.ResurrectionPaused = resource.ResurrectionPaused

    status.IPs = resource.IPs
    status.DNSs = resource.DNSs

	return
}
