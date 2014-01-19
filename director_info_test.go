package gogobosh_test

import (
	gogobosh "github.com/cloudfoundry-community/gogobosh"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"encoding/json"
	"net/http"
	"net/http/httptest"
)

var _ = Describe("get director info", func() {
	It("GET /info to return Director{}", func() {
		request := gogobosh.NewDirectorTestRequest(gogobosh.TestRequest{
			Method: "GET",
			Path:   "/info",
			Response: gogobosh.TestResponse{
				Status: http.StatusOK,
				Body: `{
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
				}`}})
		ts, handler, repo := createDirectorRepo(request)
		defer ts.Close()

		info, apiResponse := repo.GetInfo()
		
		Expect(info.Name                           ).To(Equal("Bosh Lite Director"))
		Expect(info.UUID                           ).To(Equal("bd462a15-213d-448c-aa5b-66624dad3f0e"))
		Expect(info.Version                        ).To(Equal("1.5.0.pre.1657 (14bc162c)"))
		Expect(info.User                           ).To(Equal("admin"))
		Expect(info.CPI                            ).To(Equal("warden"))
		Expect(info.DNSEnabled                     ).To(Equal(false))
		Expect(info.DNSDomainName                  ).To(Equal("bosh"))
		Expect(info.CompiledPackageCacheEnabled    ).To(Equal(true))
		Expect(info.CompiledPackageCacheProvider   ).To(Equal("local"))
		Expect(info.SnapshotsEnabled               ).To(Equal(false))

		Expect(apiResponse.IsSuccessful()).To(Equal(true))
		Expect(handler.AllRequestsCalled()).To(Equal(true))
	})

	/*
	 * To get the director info:
	 *   curl -v -k -u admin:admin https://192.168.50.4:25555/info
	 *
	 * This will give one of the responseJSON items per VM:
	*/
	It("unmarshalls JSON into Director{}", func() {
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
		resource := new(gogobosh.DirectorInfoResponse)
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

		info := resource.ToModel()
		Expect(info.Name                           ).To(Equal("Bosh Lite Director"))
		Expect(info.UUID                           ).To(Equal("bd462a15-213d-448c-aa5b-66624dad3f0e"))
		Expect(info.Version                        ).To(Equal("1.5.0.pre.1657 (14bc162c)"))
		Expect(info.User                           ).To(Equal("admin"))
		Expect(info.CPI                            ).To(Equal("warden"))
		Expect(info.DNSEnabled                     ).To(Equal(false))
		Expect(info.DNSDomainName                  ).To(Equal("bosh"))
		Expect(info.CompiledPackageCacheEnabled    ).To(Equal(true))
		Expect(info.CompiledPackageCacheProvider   ).To(Equal("local"))
		Expect(info.SnapshotsEnabled               ).To(Equal(false))
	})
})

func createDirectorRepo(reqs ...gogobosh.TestRequest) (ts *httptest.Server, handler *gogobosh.TestHandler, repo gogobosh.DirectorRepository) {
	ts, handler = gogobosh.NewTLSServer(reqs)
	config := &gogobosh.Director{
		TargetURL: ts.URL,
		Username:  "admin",
		Password:  "admin",
	}
	gateway := gogobosh.NewDirectorGateway()
	repo = gogobosh.NewBoshDirectorRepository(config, gateway)
	return
}

