package gogobosh_test

import (
	. "github.com/cloudfoundry-community/gogobosh"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"net/url"
)

var _ = Describe("Api", func() {
	Describe("Test API", func() {
		var client *Client

		Describe("Test get stemcells", func() {
			BeforeEach(func() {
				setup(MockRoute{"GET", "/stemcells", stemcells, ""}, "basic")
				config := &Config{
					BOSHAddress: server.URL,
					Username:    "admin",
					Password:    "admin",
				}

				client, _ = NewClient(config)
			})

			AfterEach(func() {
				teardown()
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
				setup(MockRoute{"GET", "/releases", releases, ""}, "basic")
				config := &Config{
					BOSHAddress: server.URL,
					Username:    "admin",
					Password:    "admin",
				}

				client, _ = NewClient(config)
			})

			AfterEach(func() {
				teardown()
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
					setup(MockRoute{"GET", "/deployments", deployments, ""}, "basic")
					config := &Config{
						BOSHAddress: server.URL,
						Username:    "admin",
						Password:    "admin",
					}

					client, _ = NewClient(config)
				})

				AfterEach(func() {
					teardown()
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
					setup(MockRoute{"POST", "/deployments", deploymentTask, ""}, "basic")
					config := &Config{
						BOSHAddress: server.URL,
						Username:    "admin",
						Password:    "admin",
					}

					client, _ = NewClient(config)
				})

				AfterEach(func() {
					teardown()
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
				setup(MockRoute{"GET", "/tasks", tasks, ""}, "basic")
				config := &Config{
					BOSHAddress: server.URL,
					Username:    "admin",
					Password:    "admin",
				}

				client, _ = NewClient(config)
			})

			AfterEach(func() {
				teardown()
			})

			It("can get tasks", func() {
				tasks, err := client.GetTasks()
				Expect(err).Should(BeNil())
				Expect(tasks[0].ID).Should(Equal(1180))
				Expect(tasks[0].State).Should(Equal("processing"))
				Expect(tasks[0].Description).Should(Equal("run errand acceptance_tests from deployment cf-warden"))
			})
		})

		Describe("Test tasks by query", func() {
			BeforeEach(func() {
				setup(MockRoute{"GET", "/tasks", tasks, ""}, "basic")
				config := &Config{
					BOSHAddress: server.URL,
					Username:    "admin",
					Password:    "admin",
				}

				client, _ = NewClient(config)
			})

			AfterEach(func() {
				teardown()
			})

			It("can get filtered tasks", func() {
				q := url.Values{}
				q.Set("state", "processing")
				tasks, err := client.GetTasksByQuery(q)
				Expect(err).Should(BeNil())
				Expect(tasks[0].ID).Should(Equal(1180))
				Expect(tasks[0].State).Should(Equal("processing"))
				Expect(tasks[0].Description).Should(Equal("run errand acceptance_tests from deployment cf-warden"))
			})
		})

		Describe("Test get deployment manifest", func() {
			BeforeEach(func() {
				setup(MockRoute{"GET", "/deployments/foo", manifest, ""}, "basic")
				config := &Config{
					BOSHAddress: server.URL,
					Username:    "admin",
					Password:    "admin",
				}

				client, _ = NewClient(config)
			})

			AfterEach(func() {
				teardown()
			})

			It("can get deployments manifest", func() {
				manifest, err := client.GetDeployment("foo")
				Expect(err).Should(BeNil())
				Expect(manifest.Manifest).Should(Equal("---\nfoo: bar\n"))
			})
		})

		Describe("Test get deployment vms", func() {
			BeforeEach(func() {
				setupMultiple([]MockRoute{
					{"GET", "/deployments/foo/vms", "", server.URL + "/tasks/2"},
					{"GET", "/tasks/2", task, ""},
					{"GET", "/tasks/2", task, ""},
					{"GET", "/tasks/2/output", vms, ""},
				}, "basic")

				config := &Config{
					BOSHAddress: server.URL,
					Username:    "admin",
					Password:    "admin",
				}

				client, _ = NewClient(config)
			})

			AfterEach(func() {
				teardown()
			})

			It("can get deployments vms", func() {
				vms, err := client.GetDeploymentVMs("foo")
				Expect(err).Should(BeNil())
				Expect(vms[0].VMCID).Should(Equal("ec974048-3352-4ba4-669d-beab87b16bcb"))
				Expect(vms[0].IPs[0]).Should(Equal("10.244.0.142"))
				Expect(vms[0].AgentID).Should(Equal("c5e7c705-459e-41c0-b640-db32d8dc6e71"))
				Expect(vms[0].JobName).Should(Equal("doppler_z1"))
				Expect(vms[0].Index).Should(Equal(0))
				Expect(vms[0].JobState).Should(Equal("running"))
				Expect(vms[0].State).Should(Equal("started"))
				Expect(vms[0].ResourcePool).Should(Equal("medium_z1"))
				Expect(vms[0].VMType).Should(Equal("default"))
				Expect(vms[0].Vitals.Disk.Ephemeral.Percent).Should(Equal("36"))
				Expect(vms[0].Vitals.Disk.Ephemeral.InodePercent).Should(Equal("11"))
				Expect(vms[0].Vitals.Disk.System.Percent).Should(Equal("36"))
				Expect(vms[0].Vitals.Disk.System.InodePercent).Should(Equal("11"))
				Expect(vms[0].Vitals.Load).Should(Equal([]string{"0.61", "0.74", "1.10"}))
				Expect(vms[0].Vitals.Mem.Percent).Should(Equal("41"))
				Expect(vms[0].Vitals.Mem.KB).Should(Equal("2520960"))
				Expect(vms[0].Vitals.Swap.Percent).Should(Equal("10"))
				Expect(vms[0].Vitals.Swap.KB).Should(Equal("102200"))
				Expect(vms[0].Vitals.CPU.Sys).Should(Equal("9.1"))
				Expect(vms[0].Vitals.CPU.User).Should(Equal("2.1"))
				Expect(vms[0].Vitals.CPU.Wait).Should(Equal("1.7"))
				Expect(vms[0].Processes[0].Name).Should(Equal("doppler"))
				Expect(vms[0].Processes[0].CPU.Total).Should(Equal(0.9))
				Expect(vms[0].Processes[0].Mem.KB).Should(Equal(2252))
				Expect(vms[0].Processes[0].Mem.Percent).Should(Equal(16.5))
				Expect(vms[0].Processes[0].State).Should(Equal("running"))
				Expect(vms[0].Processes[0].Uptime.Secs).Should(Equal(11794845))
				Expect(vms[0].ResurrectionPaused).Should(BeFalse())
				Expect(vms[0].AZ).Should(Equal("z1"))
				Expect(vms[0].ID).Should(Equal("4a9278c8-e93a-4d6a-b22c-13560208da9e"))
				Expect(vms[0].Bootstrap).Should(BeTrue())
				Expect(vms[0].Ignore).Should(BeFalse())
			})
		})

		Describe("Test stop instance", func() {
			BeforeEach(func() {
				setupMultiple([]MockRoute{
					{"PUT", "/deployments/deployment-foo/jobs/job-foo/id-foo", "", server.URL + "/tasks/3"},
					{"GET", "/tasks/3", task3, ""},
					{"GET", "/tasks/3", task3, ""},
				}, "basic")

				config := &Config{
					BOSHAddress: server.URL,
					Username:    "admin",
					Password:    "admin",
				}

				client, _ = NewClient(config)
			})

			AfterEach(func() {
				teardown()
			})

			It("can stop an instance", func() {
				task, err := client.Stop("deployment-foo", "job-foo", "id-foo")
				Expect(err).Should(BeNil())
				Expect(task.State).Should(Equal("done"))
			})
		})

		Describe("Test stop instance no converge", func() {
			BeforeEach(func() {
				setupMultiple([]MockRoute{
					{"PUT", "/deployments/deployment-foo/instance_groups/job-foo/id-foo/actions/stopped", "", server.URL + "/tasks/3"},
					{"GET", "/tasks/3", task3, ""},
					{"GET", "/tasks/3", task3, ""},
				}, "basic")

				config := &Config{
					BOSHAddress: server.URL,
					Username:    "admin",
					Password:    "admin",
				}

				client, _ = NewClient(config)
			})

			AfterEach(func() {
				teardown()
			})

			It("can stop an instance", func() {
				task, err := client.StopNoConverge("deployment-foo", "job-foo", "id-foo")
				Expect(err).Should(BeNil())
				Expect(task.State).Should(Equal("done"))
			})
		})

	})
})
