package gogobosh_test

import (
	"net/http"

	. "github.com/cloudfoundry-community/gogobosh"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
)

var _ = Describe("Api", func() {
	Describe("Test API", func() {
		var server *ghttp.Server
		var client *Client

		BeforeEach(func() {
			server = ghttp.NewServer()
			config := &Config{
				BOSHAddress: server.URL(),
				Username:    "admin",
				Password:    "admin",
			}
			client, _ = NewClient(config)
		})

		AfterEach(func() {
			//shut down the server between tests
			server.Close()
		})

		Describe("Test get stemcells", func() {
			BeforeEach(func() {
				stemcells := []Stemcell{
					Stemcell{
						Name:            "bosh-warden-boshlite-ubuntu-trusty-go_agent",
						OperatingSystem: "ubuntu-trusty",
						Version:         "3126",
						CID:             "c3705a0d-0dd3-4b67-52b5-50533a432244",
					},
				}
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyBasicAuth("admin", "admin"),
						ghttp.VerifyRequest("GET", "/stemcells"),
						ghttp.RespondWithJSONEncoded(http.StatusOK, stemcells),
					),
				)
			})

			It("can get stemcells", func() {
				stemcells, err := client.GetStemcells()
				Expect(err).Should(BeNil())
				Expect(stemcells[0].Name).Should(Equal("bosh-warden-boshlite-ubuntu-trusty-go_agent"))
				Expect(stemcells[0].OperatingSystem).Should(Equal("ubuntu-trusty"))
				Expect(stemcells[0].Version).Should(Equal("3126"))
				Expect(stemcells[0].CID).Should(Equal("c3705a0d-0dd3-4b67-52b5-50533a432244"))
			})
		})

		Describe("Test get releases", func() {
			BeforeEach(func() {
				releases := []Release{
					Release{
						Name: "bosh-warden-cpi",
						ReleaseVersions: []ReleaseVersion{
							ReleaseVersion{
								Version:            "28",
								CommitHash:         "4c36884a",
								UncommittedChanges: false,
								CurrentlyDeployed:  true,
							},
						},
					},
				}
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyBasicAuth("admin", "admin"),
						ghttp.VerifyRequest("GET", "/releases"),
						ghttp.RespondWithJSONEncoded(http.StatusOK, releases),
					),
				)
			})

			It("can get releases", func() {
				releases, err := client.GetReleases()
				Expect(err).Should(BeNil())
				Expect(releases[0].Name).Should(Equal("bosh-warden-cpi"))
				Expect(releases[0].ReleaseVersions[0].Version).Should(Equal("28"))
				Expect(releases[0].ReleaseVersions[0].CommitHash).Should(Equal("4c36884a"))
				Expect(releases[0].ReleaseVersions[0].UncommittedChanges).Should(Equal(false))
				Expect(releases[0].ReleaseVersions[0].CurrentlyDeployed).Should(Equal(true))

			})
		})

		Describe("Test deployments", func() {
			Describe("get deployments", func() {
				BeforeEach(func() {
					deployments := []Deployment{
						Deployment{
							Name:        "cf-warden",
							CloudConfig: "none",
							Releases: []Resource{
								Resource{
									Name:    "cf",
									Version: "223",
								},
							},
							Stemcells: []Resource{
								Resource{
									Name:    "bosh-warden-boshlite-ubuntu-trusty-go_agent",
									Version: "3126",
								},
							},
						},
					}
					server.AppendHandlers(
						ghttp.CombineHandlers(
							ghttp.VerifyBasicAuth("admin", "admin"),
							ghttp.VerifyRequest("GET", "/deployments"),
							ghttp.RespondWithJSONEncoded(http.StatusOK, deployments),
						),
					)
				})

				It("can get deployments", func() {
					deployments, err := client.GetDeployments()
					Expect(err).Should(BeNil())
					Expect(deployments[0].Name).Should(Equal("cf-warden"))
					Expect(deployments[0].CloudConfig).Should(Equal("none"))
					Expect(deployments[0].Releases[0].Name).Should(Equal("cf"))
					Expect(deployments[0].Releases[0].Version).Should(Equal("223"))
					Expect(deployments[0].Stemcells[0].Name).Should(Equal("bosh-warden-boshlite-ubuntu-trusty-go_agent"))
					Expect(deployments[0].Stemcells[0].Version).Should(Equal("3126"))
				})
			})
			Describe("create deployments", func() {
				BeforeEach(func() {
					task := Task{
						ID:          2,
						State:       "processing",
						Description: "run errand acceptance_tests from deployment cf-warden",
					}
					server.AppendHandlers(
						ghttp.CombineHandlers(
							ghttp.VerifyBasicAuth("admin", "admin"),
							ghttp.VerifyContentType("text/yaml"),
							ghttp.VerifyRequest("POST", "/deployments"),
							ghttp.RespondWithJSONEncoded(http.StatusOK, task),
						),
					)
				})
				It("can create deployments", func() {
					task, err := client.CreateDeployment("---\nname: foo")
					Expect(err).Should(BeNil())
					Expect(task.ID).Should(Equal(2))
				})
			})
		})

		Describe("Test tasks", func() {
			BeforeEach(func() {
				tasks := []Task{
					Task{
						ID:          1180,
						State:       "processing",
						Description: "run errand acceptance_tests from deployment cf-warden",
					},
				}
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyBasicAuth("admin", "admin"),
						ghttp.VerifyRequest("GET", "/tasks"),
						ghttp.RespondWithJSONEncoded(http.StatusOK, tasks),
					),
				)
			})

			It("can get tasks", func() {
				tasks, err := client.GetTasks()
				Expect(err).Should(BeNil())
				Expect(tasks[0].ID).Should(Equal(1180))
				Expect(tasks[0].State).Should(Equal("processing"))
				Expect(tasks[0].Description).Should(Equal("run errand acceptance_tests from deployment cf-warden"))
			})
		})

		Describe("Test get deployment manifest", func() {
			BeforeEach(func() {
				manifest := Manifest{
					Manifest: "---\nfoo: bar\n",
				}
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyBasicAuth("admin", "admin"),
						ghttp.VerifyRequest("GET", "/deployments/foo"),
						ghttp.RespondWithJSONEncoded(http.StatusOK, manifest),
					),
				)
			})

			It("can get deployments manifest", func() {
				manifest, err := client.GetDeployment("foo")
				Expect(err).Should(BeNil())
				Expect(manifest.Manifest).Should(Equal("---\nfoo: bar\n"))
			})
		})

		Describe("Test get deployment vms", func() {
			BeforeEach(func() {
				task := Task{
					ID:          2,
					State:       "done",
					Description: "run errand acceptance_tests from deployment cf-warden",
				}
				vms := `{"vm_cid":"ec974048-3352-4ba4-669d-beab87b16bcb","disk_cid":null,"ips":["10.244.0.142"],"dns":[],"agent_id":"c5e7c705-459e-41c0-b640-db32d8dc6e71","job_name":"doppler_z1","index":0,"job_state":"running","resource_pool":"medium_z1","vitals":{"cpu":{"sys":"9.1","user":"2.1","wait":"1.7"},"disk":{"ephemeral":{"inode_percent":"11","percent":"36"},"system":{"inode_percent":"11","percent":"36"}},"load":["0.61","0.74","1.10"],"mem":{"kb":"2520960","percent":"41"},"swap":{"kb":"102200","percent":"10"}},"processes":[{"name":"doppler","state":"running"},{"name":"syslog_drain_binder","state":"running"},{"name":"metron_agent","state":"running"}],"resurrection_paused":false}`
				redirect := http.Header{}
				redirect.Add("Location", server.URL()+"/tasks/2")
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyBasicAuth("admin", "admin"),
						ghttp.VerifyRequest("GET", "/deployments/foo/vms"),
						ghttp.RespondWith(http.StatusMovedPermanently, nil, redirect),
					),
					ghttp.CombineHandlers(
						ghttp.VerifyBasicAuth("admin", "admin"),
						ghttp.VerifyRequest("GET", "/tasks/2"),
						ghttp.RespondWithJSONEncoded(http.StatusOK, task),
					),
					ghttp.CombineHandlers(
						ghttp.VerifyBasicAuth("admin", "admin"),
						ghttp.VerifyRequest("GET", "/tasks/2"),
						ghttp.RespondWithJSONEncoded(http.StatusOK, task),
					),
					ghttp.CombineHandlers(
						ghttp.VerifyBasicAuth("admin", "admin"),
						ghttp.VerifyRequest("GET", "/tasks/2/output", "type=result"),
						ghttp.RespondWith(http.StatusOK, vms, nil),
					),
				)
			})

			It("can get deployments manifest", func() {
				vms, err := client.GetDeploymentVMs("foo")
				Expect(err).Should(BeNil())
				Expect(vms[0].VMCID).Should(Equal("ec974048-3352-4ba4-669d-beab87b16bcb"))
				Expect(vms[0].IPs[0]).Should(Equal("10.244.0.142"))
				Expect(vms[0].AgentID).Should(Equal("c5e7c705-459e-41c0-b640-db32d8dc6e71"))
				Expect(vms[0].JobName).Should(Equal("doppler_z1"))
			})
		})

	})
})
