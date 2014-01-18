package gogobosh_test

import (
	gogobosh "github.com/cloudfoundry-community/gogobosh"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"encoding/json"
)

var _ = Describe("get basic status", func() {
	It("returns Director", func() {
		responseJSON := `[
		  {
		    "agent_id": "dc4f1174-4aaf-427f-ad45-4f73defef626",
		    "cid": "vm-2af28cfc-42df-43bb-86f6-cbfb9ba71a06",
		    "job": "postgres_z1",
		    "index": 0
		  },
		  {
		    "agent_id": "ac3c3277-946c-422a-8308-2f95e1365bb5",
		    "cid": "vm-96727db4-ccee-47f4-b6b6-aebc5a663620",
		    "job": "ha_proxy_z1",
		    "index": 1
		  }
		]`

		resources := []gogobosh.VMStatusResponse{}

		b := []byte(responseJSON)
		err := json.Unmarshal(b, &resources)
		Expect(err).NotTo(HaveOccurred())

		Expect(len(resources)).To(Equal(2))
		resource := resources[0]
		Expect(resource.JobName).To(Equal("postgres_z1"))
		Expect(resource.Index).To(Equal(0))


		resource = resources[1]
		Expect(resource.JobName).To(Equal("ha_proxy_z1"))
		Expect(resource.Index).To(Equal(1))

		var vm_status gogobosh.VMStatus
		vm_status = resource.ToModel()
		Expect(vm_status.JobName).To(Equal("ha_proxy_z1"))
		Expect(vm_status.Index).To(Equal(1))
	})
})

