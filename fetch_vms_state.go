package gogobosh

type VMStatusResponse struct {
	JobName string  `json:"job_name"`
	Index int   `json:"index"`
	JobState string `json:"job_state"`
	VMCid string    `json:"vm_cid"`
	AgentID string  `json:"agent_id"`
	IPs []string    `json:"ips"`
	DNSs []string   `json:"dns"`
	ResourcePool string     `json:"resource_pool"`
	ResurrectionPaused bool `json:"resurrection_paused"`
	Vitals vitalsResponse `json:"vitals"`
}

type vitalsResponse struct {
	Load []string            `json:"load"`
	CPU cpuResponse          `json:"cpu"`
	Memory percentKbResponse `json:"mem"`
	Swap percentKbResponse   `json:"swap"`
	Disk diskResponse        `json:"disk"`
}

type cpuResponse struct {
	User float64    `json:"user,string"`
	System float64  `json:"sys,string"`
	Wait float64    `json:"wait,string"`
}

type diskResponse struct {
	Persistent percentKbResponse `json:"persistent"`
}

type percentKbResponse struct {
	Percent float64 `json:"percent,string"`
	Kb int          `json:"kb,string"`
}

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

	status.CPUUser = resource.Vitals.CPU.User
	status.CPUSys = resource.Vitals.CPU.System
	status.CPUWait = resource.Vitals.CPU.Wait
	status.MemoryPercent = resource.Vitals.Memory.Percent
	status.MemoryKb = resource.Vitals.Memory.Kb
	status.SwapPercent = resource.Vitals.Swap.Percent
	status.SwapKb = resource.Vitals.Swap.Kb
	status.DiskPersistentPercent = resource.Vitals.Disk.Persistent.Percent

	return
}
