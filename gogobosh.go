package gogobosh

type DirectorConfig struct {
	targetURL string
	username  string
	password  string
}

type Director struct {
	Name string
	URL string
	Version string
	User string
	UUID string
	CPI string
	DNSEnabled bool
	CompiledPackageCacheEnabled bool
	CompiledPackageCacheProvider string
	SnapshotsEnabled bool
}

func NewDirector(targetURL string, username string, password string) (director Director) {
	config := DirectorConfig{}
	config.targetURL = targetURL
	config.username = username
	config.password = password
	
	director = Director{}
	
	return
}
