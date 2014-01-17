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

type GetStatusResponse struct {
	Name string          `json:"name"`
	UUID string          `json:"uuid"`
	Version string       `json:"version"`
	User string          `json:"user"`
	CPI string           `json:"cpi"`
	Features getStatusFeaturesResponse `json:"features"`
}

type getStatusFeaturesResponse struct {
	Dns getStatusFeaturesDNS                                   `json:"dns"`
	CompiledPackageCache getStatusFeaturesCompiledPackageCache `json:"compiled_package_cache"`
	Snapshots getStatusFeaturesSnapshots                       `json:"snapshots"`
}

type getStatusFeaturesDNS struct {
	Status bool                       `json:"status"`
	Extras getStatusFeaturesDNSExtras `json:"extras"`
}

type getStatusFeaturesDNSExtras struct {
	DomainName string `json:"domain_name"`
}

type getStatusFeaturesCompiledPackageCache struct {
	Status bool                                        `json:"status"`
	Extras getStatusFeaturesCompiledPackageCacheExtras `json:"extras"`
}

type getStatusFeaturesCompiledPackageCacheExtras struct {
	Provider string `json:"provider"`
}

type getStatusFeaturesSnapshots struct {
	Status bool `json:"status"`
}

func (resource GetStatusResponse) ToModel() (director Director) {
/*	director.ApplicationFields = resource.ToFields()
	routes := []cf.RouteSummary{}
	for _, route := range resource.Routes {
		routes = append(routes, route.ToModel())
	}
	app.RouteSummaries = routes

*/	return
}

func NewDirector(targetURL string, username string, password string) (director Director) {
	config := DirectorConfig{}
	config.targetURL = targetURL
	config.username = username
	config.password = password
	
	director = Director{}
	
	return
}
