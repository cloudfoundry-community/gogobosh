package gogobosh

type Director struct {
	targetURL string
	username  string
	password  string
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
	director.targetURL = targetURL
	director.username = username
	director.password = password
	
	return
}

func (director Director) GetInfo() (info DirectorInfo) {
	info = DirectorInfo{}
	info.Name = "hi"
	return
}