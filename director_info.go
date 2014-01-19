package gogobosh

func (repo BoshDirectorRepository) GetInfo() (directorInfo DirectorInfo, apiResponse ApiResponse) {
	infoResource := DirectorInfoResponse{}

	path := "/info"
	apiResponse = repo.gateway.GetResource(repo.config.TargetURL+path, repo.config.Username, repo.config.Password, &infoResource)
	if apiResponse.IsNotSuccessful() {
		return
	}

	directorInfo = infoResource.ToModel()

	return
}

type DirectorInfoResponse struct {
	Name string          `json:"name"`
	UUID string          `json:"uuid"`
	Version string       `json:"version"`
	User string          `json:"user"`
	CPI string           `json:"cpi"`
	Features directorInfoFeaturesResponse `json:"features"`
}

type directorInfoFeaturesResponse struct {
	DNS directorInfoFeaturesDNS                                   `json:"dns"`
	CompiledPackageCache directorInfoFeaturesCompiledPackageCache `json:"compiled_package_cache"`
	Snapshots directorInfoFeaturesSnapshots                       `json:"snapshots"`
}

type directorInfoFeaturesDNS struct {
	Status bool                       `json:"status"`
	Extras directorInfoFeaturesDNSExtras `json:"extras"`
}

type directorInfoFeaturesDNSExtras struct {
	DomainName string `json:"domain_name"`
}

type directorInfoFeaturesCompiledPackageCache struct {
	Status bool                                        `json:"status"`
	Extras directorInfoFeaturesCompiledPackageCacheExtras `json:"extras"`
}

type directorInfoFeaturesCompiledPackageCacheExtras struct {
	Provider string `json:"provider"`
}

type directorInfoFeaturesSnapshots struct {
	Status bool `json:"status"`
}

func (resource DirectorInfoResponse) ToModel() (director DirectorInfo) {
	director = DirectorInfo{}
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
