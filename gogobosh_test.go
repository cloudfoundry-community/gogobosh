package gogobosh_test

import (
	gogobosh "github.com/cloudfoundry-community/gogobosh"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"encoding/json"
)

var _ = Describe("GoGoBOSH", func() {
	It("parse response", func() {
		responseJSON := `{
		  "name": "Bosh Lite Director",
		  "uuid": "bd462a15-213d-448c-aa5b-66624dad3f0e",
		  "version": "1.5.0.pre.1657 (14bc162c)",
		  "user": "admin",
		  "cpi": "warden",
		  "features": {
		    "dns": {
		      "status": false,
		      "extras": {
		        "domain_name": "bosh"
		      }
		    },
		    "compiled_package_cache": {
		      "status": true,
		      "extras": {
		        "provider": "local"
		      }
		    },
		    "snapshots": {
		      "status": false
		    }
		  }
		}`
		resource := new(gogobosh.GetStatusResponse)
		b := []byte(responseJSON)
		err := json.Unmarshal(b, &resource)
		Expect(err).NotTo(HaveOccurred())
		Expect(resource.Name).To(Equal("Bosh Lite Director"))
		Expect(resource.UUID).To(Equal("bd462a15-213d-448c-aa5b-66624dad3f0e"))
		Expect(resource.Version).To(Equal("1.5.0.pre.1657 (14bc162c)"))
		Expect(resource.User).To(Equal("admin"))
		Expect(resource.CPI).To(Equal("warden"))
		Expect(resource.Features.DNS.Status).To(Equal(false))
		Expect(resource.Features.DNS.Extras.DomainName).To(Equal("bosh"))
		Expect(resource.Features.CompiledPackageCache.Status).To(Equal(true))
		Expect(resource.Features.CompiledPackageCache.Extras.Provider).To(Equal("local"))
		Expect(resource.Features.Snapshots.Status).To(Equal(false))
	})
})
