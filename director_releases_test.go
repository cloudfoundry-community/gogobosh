package gogobosh_test

import (
	gogobosh "github.com/cloudfoundry-community/gogobosh"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"encoding/json"
)

var _ = Describe("get list of releases", func() {
	/*
	 * To get the director info:
	 *   curl -k -u admin:admin https://192.168.50.4:25555/releases
	 *
	 * This will give one of the responseJSON items per VM:
	*/
	It("returns Director", func() {
		responseJSON := `[
		  {
		    "name": "cf",
		    "release_versions": [
		      {
		        "version": "153",
		        "commit_hash": "009fdd9a",
		        "uncommitted_changes": true,
		        "currently_deployed": true,
		        "job_names": [
		          "cloud_controller_ng",
		          "nats",
		          "dea_next",
		          "login",
		          "health_manager_next",
		          "uaa",
		          "debian_nfs_server",
		          "loggregator",
		          "postgres",
		          "dea_logging_agent",
		          "syslog_aggregator",
		          "narc",
		          "haproxy",
		          "hm9000",
		          "saml_login",
		          "nats_stream_forwarder",
		          "collector",
		          "pivotal_login",
		          "loggregator_trafficcontroller",
		          "etcd",
		          "gorouter"
		        ]
		      }
		    ]
		  }
		]`
		resources := []gogobosh.ReleaseResponse{}
		b := []byte(responseJSON)
		err := json.Unmarshal(b, &resources)
		Expect(err).NotTo(HaveOccurred())

		release := resources[0].ToModel()
		Expect(release.Name).To(Equal("cf"))
		
		release_version := release.Versions[0]
		Expect(release_version.Version).To(Equal("153"))
		Expect(release_version.CommitHash).To(Equal("009fdd9a"))
		Expect(release_version.UncommittedChanges).To(Equal(true))
		Expect(release_version.CurrentlyDeployed).To(Equal(true))
	})
})
