package gogobosh_test

import (
	gogobosh "github.com/cloudfoundry-community/gogobosh"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"encoding/json"
)

var _ = Describe("get list of stemcells", func() {
	/*
	 * To get the director info:
	 *   curl -k -u admin:admin https://192.168.50.4:25555/stemcells
	 *
	 * This will give one of the responseJSON items per VM:
	*/
	It("returns Director", func() {
		responseJSON := `[
		  {
		    "name": "bosh-stemcell",
		    "version": "993",
		    "cid": "stemcell-6e6b9689-8b03-42cd-a6de-7784e3c421ec",
		    "deployments": [
		      "#<Bosh::Director::Models::Deployment:0x0000000474bdb0>"
		    ]
		  },
		  {
		    "name": "bosh-warden-boshlite-ubuntu",
		    "version": "24",
		    "cid": "stemcell-6936d497-b8cd-4e12-af0a-5f2151834a1a",
		    "deployments": [

		    ]
		  }
		]`
		resources := []gogobosh.StemcellResponse{}
		b := []byte(responseJSON)
		err := json.Unmarshal(b, &resources)
		Expect(err).NotTo(HaveOccurred())

		stemcell := resources[0].ToModel()
		Expect(stemcell.Name).To(Equal("bosh-stemcell"))
		Expect(stemcell.Version).To(Equal("993"))
		Expect(stemcell.Cid).To(Equal("stemcell-6e6b9689-8b03-42cd-a6de-7784e3c421ec"))

		/* TODO: deployments is returning an internal Ruby object string; not JSON. */
	})
})
