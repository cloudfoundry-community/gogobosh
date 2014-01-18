package gogobosh

type GetStatusResponse struct {
	Name string          `json:"name"`
	UUID string          `json:"uuid"`
	Version string       `json:"version"`
	User string          `json:"user"`
	CPI string           `json:"cpi"`
	Features getStatusFeaturesResponse `json:"features"`
}

type getStatusFeaturesResponse struct {
	DNS getStatusFeaturesDNS                                   `json:"dns"`
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
	director = Director{}
	director.Name = resource.Name
	director.Version = resource.Version
	director.User = resource.User
	director.UUID = resource.UUID
	director.CPI = resource.CPI

	director.DNSEnabled = resource.Features.DNS.Status
	director.DNSDomainName = resource.Features.DNS.Extras.DomainName
	director.CompiledPackageCacheEnabled = resource.Features.CompiledPackageCache.Status
	director.CompiledPackageCacheProvider = resource.Features.CompiledPackageCache.Extras.Provider
	director.SnapshotsEnabled = resource.Features.Snapshots.Status

	return
}
