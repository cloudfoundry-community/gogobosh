package gogobosh_test

import (
	gogobosh "github.com/cloudfoundry-community/gogobosh"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"net/http"
	"fmt"
)

var _ = Describe("Deployments", func() {
	It("GetDeployments() - list of deployments", func() {
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

	It("DeleteDeployment(name) forcefully", func() {
		request := gogobosh.NewDirectorTestRequest(gogobosh.TestRequest{
			Method: "DELETE",
			Path:   "/deployments/cf-warden?force=true",
			Response: gogobosh.TestResponse{
				Status: http.StatusFound,
				Header: http.Header{
					"Location":{"https://some.host/tasks/20"},
				},
			}})
		queuedTaskRequest := gogobosh.NewDirectorTestRequest(gogobosh.TestRequest{
			Method: "GET",
			Path:   "/tasks/20",
			Response: gogobosh.TestResponse{
				Status: http.StatusOK,
				Body: taskResponseJSONwithState("queued")}})
		processingTaskRequest := gogobosh.NewDirectorTestRequest(gogobosh.TestRequest{
			Method: "GET",
			Path:   "/tasks/20",
			Response: gogobosh.TestResponse{
				Status: http.StatusOK,
				Body: taskResponseJSONwithState("processing")}})
		doneTaskRequest := gogobosh.NewDirectorTestRequest(gogobosh.TestRequest{
			Method: "GET",
			Path:   "/tasks/20",
			Response: gogobosh.TestResponse{
				Status: http.StatusOK,
				Body: taskResponseJSONwithState("done")}})
		ts, handler, repo := createDirectorRepo(request, queuedTaskRequest, processingTaskRequest, doneTaskRequest)
		defer ts.Close()

		apiResponse := repo.DeleteDeployment("cf-warden")

		Expect(apiResponse.IsSuccessful()).To(Equal(true))
		Expect(handler.AllRequestsCalled()).To(Equal(true))
	})
})

func taskResponseJSONwithState(state string) (json string) {
	baseJSON := `{
	  "id": 20,
	  "state": "%s",
	  "description": "some task",
	  "timestamp": 1390174354,
	  "result": null,
	  "user": "admin"
	}`
	return fmt.Sprintf(baseJSON, state)
}
