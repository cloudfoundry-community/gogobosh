package gogobosh_test

import (
	gogobosh "github.com/cloudfoundry-community/gogobosh"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"encoding/json"
)

var _ = Describe("get task status", func() {
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
