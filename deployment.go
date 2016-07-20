package gogobosh

// HasRelease if deployment has release
func (d *Deployment) HasRelease(name string) bool {
	for _, release := range d.Releases {
		if release.Name == name {
			return true
		}
	}
	return false
}
