package gogobosh

// Info struct
type Info struct {
	Name               string             `json:"name"`
	UUID               string             `json:"uuid"`
	Version            string             `json:"version"`
	User               string             `json:"user"`
	CPI                string             `json:"cpi"`
	UserAuthentication UserAuthentication `json:"user_authentication"`
}

// UserAuthentication struct
type UserAuthentication struct {
	Type    string `json:"type"`
	Options struct {
		URL string `json:"url"`
	} `json:"options"`
}

// Stemcell struct
type Stemcell struct {
	Name            string `json:"name"`
	OperatingSystem string `json:"operating_system"`
	Version         string `json:"version"`
	CID             string `json:"cid"`
	CPI             string `json:"cpi"`
	Deployments     []struct {
		Name string `json:"name"`
	} `json:"deployments"`
}

// Release struct
type Release struct {
	Name            string           `json:"name"`
	ReleaseVersions []ReleaseVersion `json:"release_versions"`
}

// ReleaseVersion struct
type ReleaseVersion struct {
	Version            string   `json:"version"`
	CommitHash         string   `json:"commit_hash"`
	UncommittedChanges bool     `json:"uncommitted_changes"`
	CurrentlyDeployed  bool     `json:"currently_deployed"`
	JobNames           []string `json:"job_names"`
}

// Deployment struct
type Deployment struct {
	Name        string     `json:"name"`
	CloudConfig string     `json:"cloud_config"`
	Releases    []Resource `json:"releases"`
	Stemcells   []Resource `json:"stemcells"`
}

// Resource struct
type Resource struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

// Manifest struct
type Manifest struct {
	Manifest string `json:"manifest"`
}

// VM struct
type VM struct {
	VMCID              string    `json:"vm_cid"`
	IPs                []string  `json:"ips"`
	DNS                []string  `json:"dns"`
	AgentID            string    `json:"agent_id"`
	JobName            string    `json:"job_name"`
	Index              int       `json:"index"`
	JobState           string    `json:"job_state"`
	State              string    `json:"state"`
	ResourcePool       string    `json:"resource_pool"`
	VMType             string    `json:"vm_type"`
	Vitals             Vitals    `json:"vitals"`
	Processes          []Process `json:"processes"`
	ResurrectionPaused bool      `json:"resurrection_paused"`
	AZ                 string    `json:"az"`
	ID                 string    `json:"id"`
	Bootstrap          bool      `json:"bootstrap"`
	Ignore             bool      `json:"ignore"`
}

// Vitals for a VM
type Vitals struct {
	Disk Disk     `json:"disk"`
	Load []string `json:"load"`
	Mem  Memory   `json:"mem"`
	Swap Memory   `json:"swap"`
	CPU  CPU      `json:"cpu"`
}

// Disk struct
type Disk struct {
	Ephemeral  DiskStats `json:"ephemeral"`
	System     DiskStats `json:"system"`
	Persistent DiskStats `json:"persistent"`
}

// CPU struct
type CPU struct {
	Sys  string `json:"sys"`
	User string `json:"user"`
	Wait string `json:"wait"`
}

// DiskStats struct
type DiskStats struct {
	Percent      string `json:"percent"`
	InodePercent string `json:"inode_percent"`
}

// Memory struct
type Memory struct {
	Percent string `json:"percent"`
	KB      string `json:"KB"`
}

// Process running on a VM
type Process struct {
	Name   string        `json:"name"`
	State  string        `json:"state"`
	Uptime Uptime        `json:"uptime"`
	Mem    ProcessMemory `json:"mem"`
	CPU    ProcessCPU    `json:"cpu"`
}

// Uptime struct
type Uptime struct {
	Secs int `json:"secs"`
}

// ProcessCPU struct
type ProcessCPU struct {
	Total float64 `json:"total"`
}

// ProcessMemory struct
type ProcessMemory struct {
	Percent float64 `json:"percent"`
	KB      int     `json:"KB"`
}

// Task struct
type Task struct {
	ID          int    `json:"id"`
	State       string `json:"state"`
	Description string `json:"description"`
	Timestamp   int    `json:"timestamp"`
	Result      string `json:"result"`
	User        string `json:"user"`
}

// Event struct
type Event struct {
	ID         string                 `json:"id"`
	ParentID   string                 `json:"parent_id"`
	Timestamp  int                    `json:"timestamp"`
	User       string                 `json:"user"`
	Action     string                 `json:"action"`
	ObjectType string                 `json:"object_type"`
	ObjectName string                 `json:"object_name"`
	Task       string                 `json:"task"`
	Deployment string                 `json:"deployment"`
	Error      string                 `json:"error"`
	Context    map[string]interface{} `json:"context"`
}

// TaskEvent struct
type TaskEvent struct {
	Time     int      `json:"time"`
	Stage    string   `json:"stage"`
	Tags     []string `json:"tags"`
	Total    int      `json:"total"`
	Task     string   `json:"task"`
	Index    int      `json:"index"`
	State    string   `json:"state"`
	Progress int      `json:"progress"`

	Error struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

// Cfg struct
type Cfg struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Type      string `json:"type"`
	Content   string `json:"content"`
	CreatedAt int    `json:"int"`
	Deleted   bool   `json:"deleted"`
}
