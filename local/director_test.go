package local_test

import (
	"path/filepath"

	"github.com/cloudfoundry-community/gogobosh/local"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Local config", func() {
	It("Loads BOSH config", func() {
		configPath, err := filepath.Abs("../testhelpers/fixtures/bosh_config.yml")
		Expect(err).ShouldNot(HaveOccurred())

		config, err := local.LoadBoshConfig(configPath)
		Expect(err).ShouldNot(HaveOccurred())
		Expect(config).ToNot(BeNil())
		Expect(config.Name).To(Equal("Bosh Lite Director"))
		Expect(config.Authentication["https://192.168.50.4:25555"].Username).To(Equal("admin"))
	})

	It("CurrentBoshTarget", func() {
		configPath, err := filepath.Abs("../testhelpers/fixtures/bosh_config.yml")
		Expect(err).ShouldNot(HaveOccurred())
		config, err := local.LoadBoshConfig(configPath)
		Expect(err).ShouldNot(HaveOccurred())
		Expect(config).ToNot(BeNil())

		target, username, password, err := config.CurrentBoshTarget()
		Expect(err).ShouldNot(HaveOccurred())
		Expect(target).To(Equal("https://192.168.50.4:25555"))
		Expect(username).To(Equal("admin"))
		Expect(password).To(Equal("password"))
	})

	It("CurrentBoshDeployment", func() {
		configPath, err := filepath.Abs("../testhelpers/fixtures/bosh_config.yml")
		Expect(err).ShouldNot(HaveOccurred())
		config, err := local.LoadBoshConfig(configPath)
		Expect(err).ShouldNot(HaveOccurred())
		Expect(config).ToNot(BeNil())

		manifestPath := config.CurrentDeploymentManifest()
		Expect(manifestPath).To(Equal("path/to/manifest.yml"))
	})
})
