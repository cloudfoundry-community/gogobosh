package gogobosh

// Info struct
type Info struct {
	Name              string            `json:"name"`
	UUID              string            `json:"uuid"`
	Version           string            `json:"version"`
	User              string            `json:"user"`
	CPI               string            `json:"cpi"`
	UserAuthenication UserAuthenication `json:"user_authentication"`
}

// UserAuthenication struct
type UserAuthenication struct {
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
	Deployments     []struct {
		Name string `json:"name"`
	} `json:"deployments"`
}

// Release struct
type Release struct {
	Name            string           `json:"name"`
	ReleaseVersions []ReleaseVersion `json:"release_versions"`
}

//ReleaseVersion struct
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
	AgentID           string   `json:"agent_id"`
	VMCID             string   `json:"vm_cid"`
	CID               string   `json:"cid"`
	JobName           string   `json:"job_name"`
	Index             int      `json:"index"`
	IPs               []string `json:"ips"`
	DNS               []string `json:"dns"`
	ResurectionPaused bool     `json:"resurrection_paused"`
	Vitals            Vitals   `json:"vitals"`
	ID                string   `json:"id"`
}

// VM Vitals struct
type Vitals struct {
	Disk Disk     `json:"disk"`
	Load []string `json:"load"`
	Mem  Memory   `json:"mem"`
	Swap Memory   `json:"swap"`
	CPU  CPU      `json:"cpu"`
}
type Disk struct {
	Ephemeral DiskStats `json:"ephemeral"`
	System    DiskStats `json:"system"`
}

type CPU struct {
	Sys  string `json:"sys"`
	User string `json:"user"`
	Wait string `json:"wait"`
}

// Disk struct
type DiskStats struct {
	Percent      string `json:"percent"`
	InodePercent string `json:"inode_percent"`
}

// Memory struct
type Memory struct {
	Percent string `json:"percent"`
	KB      string `json:"KB"`
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
