package main

import (
	"github.com/cloudfoundry-community/gogobosh"
	"fmt"
)

func main() {
	director := gogobosh.NewDirector("https://192.168.50.4:25555", "admin", "admin")
	fmt.Println("Director")
	fmt.Printf("  Name       %s\n", director.Name)
	fmt.Printf("  URL        %s\n", director.URL)
	fmt.Printf("  Version    %s\n", director.Version)
	fmt.Printf("  User       %s\n", director.User)
	fmt.Printf("  UUID       %s\n", director.UUID)
	fmt.Printf("  CPI        %s\n", director.CPI)
	fmt.Printf("  dns        %s\n", director.DNSEnabled)
	fmt.Printf("  compiled_package_cache %#v (provider: %s)\n", director.CompiledPackageCacheEnabled, director.CompiledPackageCacheProvider)
	fmt.Printf("  snapshots  %#v\n", director.SnapshotsEnabled)
}
