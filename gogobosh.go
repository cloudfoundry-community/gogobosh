package gogobosh

const (
	Version = "0.1.0"
)
type Director struct {
	TargetURL string
	Username  string
	Password  string
}

type DirectorInfo struct {
	Name string
	URL string
	Version string
	User string
	UUID string
	CPI string
	DNSEnabled bool
	DNSDomainName string
	CompiledPackageCacheEnabled bool
	CompiledPackageCacheProvider string
	SnapshotsEnabled bool
}

type Stemcell struct {
	Name string
	Version string
	Cid string
}

type Release struct {
	Name string
	Versions []ReleaseVersion
}

type ReleaseVersion struct {
	Version string
	CommitHash string
	UncommittedChanges bool
	CurrentlyDeployed bool
}

type Deployment struct {
	Name string
	Releases []NameVersion
	Stemcells []NameVersion
}

type DeploymentVM struct {
	JobName string
	Index int
	VMCid string
	AgentID string
}

type NameVersion struct {
	Name string
	Version string
}

type TaskStatus struct {
	ID int
	State string
	Description string
	TimeStamp int
	Result string
	User string
}

type VMStatus struct {
	JobName string
	Index int
	JobState string
	VMCid string
	AgentID string
	ResourcePool string
	ResurrectionPaused bool
	IPs []string
	DNSs []string
	CPUUser float64
	CPUSys float64
	CPUWait float64
	MemoryPercent float64
	MemoryKb int
	SwapPercent float64
	SwapKb int
	DiskPersistentPercent float64
}

func NewDirector(targetURL string, username string, password string) (director Director) {
	director = Director{}
	director.TargetURL = targetURL
	director.Username = username
	director.Password = password
	
	return
}
