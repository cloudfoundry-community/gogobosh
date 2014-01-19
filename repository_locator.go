package gogobosh

type RepositoryLocator struct {
	directorRepo  DirectorRepository
}

func NewRepositoryLocator(config *Director, gatewaysByName map[string]Gateway) (loc RepositoryLocator) {
	boshGateway := gatewaysByName["bosh"]

	loc.directorRepo = NewBoshDirectorRepository(config, boshGateway)

	return
}

func (locator RepositoryLocator) GetDirectorRepository() DirectorRepository {
	return locator.directorRepo
}
