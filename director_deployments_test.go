package gogobosh_test

import (
	gogobosh "github.com/cloudfoundry-community/gogobosh"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"encoding/json"
	"net/http"
)

var _ = Describe("get list of deployments", func() {
	It("GET /deployments to return []DirectorDeployment{}", func() {
		request := gogobosh.NewDirectorTestRequest(gogobosh.TestRequest{
			Method: "GET",
			Path:   "/deployments",
			Response: gogobosh.TestResponse{
				Status: http.StatusOK,
				Body: `[
				  {
				    "name": "cf-warden",
				    "deployments": [
				      {
				        "name": "cf",
				        "version": "153"
				      }
				    ],
				    "stemcells": [
				      {
				        "name": "bosh-stemcell",
				        "version": "993"
				      }
				    ]
				  }
				]`}})
		ts, handler, repo := createDirectorRepo(request)
		defer ts.Close()

		deployments, apiResponse := repo.GetDeployments()

		deployment := deployments[0]
		Expect(deployment.Name).To(Equal("cf-warden"))

		deployment_release := deployment.Releases[0]
		Expect(deployment_release.Name).To(Equal("cf"))
		Expect(deployment_release.Version).To(Equal("153"))

		deployment_stemcell := deployment.Stemcells[0]
		Expect(deployment_stemcell.Name).To(Equal("bosh-stemcell"))
		Expect(deployment_stemcell.Version).To(Equal("993"))

		Expect(apiResponse.IsSuccessful()).To(Equal(true))
		Expect(handler.AllRequestsCalled()).To(Equal(true))
	})

	/*
	 * To get the director deployments list:
	 *   curl -k -u admin:admin https://192.168.50.4:25555/deployments
	*/
	It("returns Deployment", func() {
		responseJSON := `[
		  {
		    "name": "cf-warden",
		    "deployments": [
		      {
		        "name": "cf",
		        "version": "153"
		      }
		    ],
		    "stemcells": [
		      {
		        "name": "bosh-stemcell",
		        "version": "993"
		      }
		    ]
		  }
		]`
		resources := []gogobosh.DeploymentResponse{}
		b := []byte(responseJSON)
		err := json.Unmarshal(b, &resources)
		Expect(err).NotTo(HaveOccurred())

		deployment := resources[0].ToModel()
		Expect(deployment.Name).To(Equal("cf-warden"))

		deployment_release := deployment.Releases[0]
		Expect(deployment_release.Name).To(Equal("cf"))
		Expect(deployment_release.Version).To(Equal("153"))

		deployment_stemcell := deployment.Stemcells[0]
		Expect(deployment_stemcell.Name).To(Equal("bosh-stemcell"))
		Expect(deployment_stemcell.Version).To(Equal("993"))
	})
})
