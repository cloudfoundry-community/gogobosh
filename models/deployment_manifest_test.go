package models_test

import (
	"github.com/cloudfoundry-community/gogobosh/models"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("DeploymentManifest", func() {
	It("FindJobTemplates", func() {
		manifest := &models.DeploymentManifest{
			Jobs: []*models.ManifestJob{
				{Name: "job1", JobTemplates: []*models.ManifestJobTemplate{{Name: "common"}}},
				{Name: "job2", JobTemplates: []*models.ManifestJobTemplate{{Name: "common"}}},
				{Name: "other", JobTemplates: []*models.ManifestJobTemplate{{Name: "other"}}},
			},
		}
		jobs := manifest.FindJobTemplates("common")
		Expect(len(jobs)).To(Equal(2))
	})
})
