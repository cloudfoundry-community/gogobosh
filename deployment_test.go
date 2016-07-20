package gogobosh_test

import (
	. "github.com/cloudfoundry-community/gogobosh"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Deployment", func() {
	Describe("Test Deployment", func() {

		Describe("Test if release exist", func() {

			It("see if release exist", func() {
				deployment := &Deployment{
					Releases: []Resource{
						Resource{
							Name: "test",
						},
					},
				}
				Expect(deployment.HasRelease("test")).Should(Equal(true))
				Expect(deployment.HasRelease("bad")).Should(Equal(false))
			})
		})

	})
})
