package gogobosh_test

import (
	gogobosh "github.com/cloudfoundry-community/gogobosh"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"encoding/json"
	"net/http"
)

var _ = Describe("parse full vms task output", func() {
	It("GET /deployments/cf-warden/vms?format=full to return Director{}", func() {
		vmsRequest := gogobosh.NewDirectorTestRequest(gogobosh.TestRequest{
			Method: "GET",
			Path:   "/deployments/cf-warden/vms?format=full",
			Response: gogobosh.TestResponse{
				Status: http.StatusFound,
				Header: http.Header{
					"Location":{"https://some.host/tasks/12"},
				},
			},
		})
		queuedTaskRequest := gogobosh.NewDirectorTestRequest(gogobosh.TestRequest{
			Method: "GET",
			Path:   "/tasks/12",
			Response: gogobosh.TestResponse{
				Status: http.StatusOK,
				Body: `{
				  "id": 12,
				  "state": "queued",
				  "description": "retrieve vm-stats",
				  "timestamp": 1390174354,
				  "result": null,
				  "user": "admin"
				}`}})
		doneTaskRequest := gogobosh.NewDirectorTestRequest(gogobosh.TestRequest{
			Method: "GET",
			Path:   "/tasks/12",
			Response: gogobosh.TestResponse{
				Status: http.StatusOK,
				Body: `{
				  "id": 12,
				  "state": "done",
				  "description": "retrieve vm-stats",
				  "timestamp": 1390174354,
				  "result": null,
				  "user": "admin"
				}`}})

		ts, handler, repo := createDirectorRepo(vmsRequest, queuedTaskRequest, doneTaskRequest)
		defer ts.Close()

		vm_statuses, apiResponse := repo.FetchVMsStatus("cf-warden")

		/* TODO convert vm_statuses to a chan and pluck first item from chan */
		Expect(len(vm_statuses)).To(Equal(2))
		vm_status := vm_statuses[0]
		Expect(vm_status.JobName).To(Equal("etcd_leader_z1"))
		Expect(vm_status.Index).To(Equal(0))
		Expect(vm_status.JobState).To(Equal("running"))
		Expect(vm_status.VMCid).To(Equal("vm-00b5c65f-d2f4-4289-ab8d-8ae413b4dc9b"))
		Expect(vm_status.AgentID).To(Equal("892d2de8-16aa-4567-b49a-45b1d99882b5"))
		Expect(vm_status.ResourcePool).To(Equal("medium_z1"))
		Expect(vm_status.ResurrectionPaused).To(Equal(false))

		Expect(len(vm_status.IPs)).To(Equal(1))
		Expect(vm_status.IPs[0]).To(Equal("10.244.0.38"))

		Expect(len(vm_status.DNSs)).To(Equal(1))
		Expect(vm_status.DNSs[0]).To(Equal("0.etcd-leader-z1.default.my-deployment.bosh"))

		Expect(vm_status.CPUUser).To(Equal(0.2))
		Expect(vm_status.CPUSys).To(Equal(1.5))
		Expect(vm_status.CPUWait).To(Equal(0.0))
		Expect(vm_status.MemoryPercent).To(Equal(43.1))
		Expect(vm_status.MemoryKb).To(Equal(2635712))
		Expect(vm_status.SwapPercent).To(Equal(0.0))
		Expect(vm_status.SwapKb).To(Equal(284))
		Expect(vm_status.DiskPersistentPercent).To(Equal(1.0))

		Expect(apiResponse.IsSuccessful()).To(Equal(true))
		Expect(handler.AllRequestsCalled()).To(Equal(true))
	})

	/*
	 * To get the full status of VMs, GET the following:
	 *   curl -v -k -L -u admin:admin https://192.168.50.4:25555/deployments/cf-warden/vms\?format\=full
	 * and with the Location redirect, extract the task_id, then run:
	 *   curl -k -u admin:admin https://192.168.50.4:25555/tasks/19/output\?type\=result | jazor
	 *
	 * This will give one of the responseJSON items per VM:
	*/
	It("returns VMState", func() {
		responseJSON := `{
          "job_name": "etcd_leader_z1",
          "index": 0,
          "job_state": "running",
          "vm_cid": "vm-00b5c65f-d2f4-4289-ab8d-8ae413b4dc9b",
          "agent_id": "892d2de8-16aa-4567-b49a-45b1d99882b5",
          "resource_pool": "medium_z1",
          "ips": [
            "10.244.0.38"
          ],
          "dns": [
            "0.etcd-leader-z1.default.my-deployment.bosh"
          ],
          "vitals": {
            "load": [
              "0.02",
              "0.60",
              "0.80"
            ],
            "cpu": {
              "user": "0.2",
              "sys": "1.5",
              "wait": "0.0"
            },
            "mem": {
              "percent": "43.1",
              "kb": "2635712"
            },
            "swap": {
              "percent": "0.0",
              "kb": "284"
            },
            "disk": {
              "persistent": {
                "percent": "1"
              }
            }
          },
          "resurrection_paused": false
        }`

		resource := gogobosh.VMStatusResponse{}

		b := []byte(responseJSON)
		err := json.Unmarshal(b, &resource)
		Expect(err).NotTo(HaveOccurred())

		vm_status := resource.ToModel()
		Expect(vm_status.JobName).To(Equal("etcd_leader_z1"))
		Expect(vm_status.Index).To(Equal(0))
		Expect(vm_status.JobState).To(Equal("running"))
		Expect(vm_status.VMCid).To(Equal("vm-00b5c65f-d2f4-4289-ab8d-8ae413b4dc9b"))
		Expect(vm_status.AgentID).To(Equal("892d2de8-16aa-4567-b49a-45b1d99882b5"))
		Expect(vm_status.ResourcePool).To(Equal("medium_z1"))
		Expect(vm_status.ResurrectionPaused).To(Equal(false))

		Expect(len(vm_status.IPs)).To(Equal(1))
		Expect(vm_status.IPs[0]).To(Equal("10.244.0.38"))

		Expect(len(vm_status.DNSs)).To(Equal(1))
		Expect(vm_status.DNSs[0]).To(Equal("0.etcd-leader-z1.default.my-deployment.bosh"))

		Expect(vm_status.CPUUser).To(Equal(0.2))
		Expect(vm_status.CPUSys).To(Equal(1.5))
		Expect(vm_status.CPUWait).To(Equal(0.0))
		Expect(vm_status.MemoryPercent).To(Equal(43.1))
		Expect(vm_status.MemoryKb).To(Equal(2635712))
		Expect(vm_status.SwapPercent).To(Equal(0.0))
		Expect(vm_status.SwapKb).To(Equal(284))
		Expect(vm_status.DiskPersistentPercent).To(Equal(1.0))
	})
})