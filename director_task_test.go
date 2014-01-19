package gogobosh_test

import (
	gogobosh "github.com/cloudfoundry-community/gogobosh"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"encoding/json"
	"net/http"
)

var _ = Describe("get task status", func() {
	It("GET /tasks/1 to return TaskStatus{}", func() {
		request := gogobosh.NewDirectorTestRequest(gogobosh.TestRequest{
			Method: "GET",
			Path:   "/tasks/1",
			Response: gogobosh.TestResponse{
				Status: http.StatusOK,
				Body: `{
				  "id": 1,
				  "state": "done",
				  "description": "create release",
				  "timestamp": 1390068518,
				  "result": "Created release cf/153",
				  "user": "admin"
				}`}})
		ts, handler, repo := createDirectorRepo(request)
		defer ts.Close()

		task, apiResponse := repo.GetTaskStatus(1)
		
		Expect(task.ID).To(Equal(1))
		Expect(task.State).To(Equal("done"))
		Expect(task.Description).To(Equal("create release"))
		Expect(task.TimeStamp).To(Equal(1390068518))
		Expect(task.Result).To(Equal("Created release cf/153"))
		Expect(task.User).To(Equal("admin"))

		Expect(apiResponse.IsSuccessful()).To(Equal(true))
		Expect(handler.AllRequestsCalled()).To(Equal(true))
	})

	/*
	 * To get the director info:
	 *   curl -k -u admin:admin https://192.168.50.4:25555/tasks/1
	*/
	It("returns []TaskStatus", func() {
		responseJSON := `{
		  "id": 1,
		  "state": "done",
		  "description": "create release",
		  "timestamp": 1390068518,
		  "result": "Created release cf/153",
		  "user": "admin"
		}`

		resource := gogobosh.TaskStatusResponse{}
		b := []byte(responseJSON)
		err := json.Unmarshal(b, &resource)
		Expect(err).NotTo(HaveOccurred())

		task := resource.ToModel()
		Expect(task.ID).To(Equal(1))
		Expect(task.State).To(Equal("done"))
		Expect(task.Description).To(Equal("create release"))
		Expect(task.TimeStamp).To(Equal(1390068518))
		Expect(task.Result).To(Equal("Created release cf/153"))
		Expect(task.User).To(Equal("admin"))
	})
})
