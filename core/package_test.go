package core_test

import (
	"os"

	"github.com/gyf1214/docker-build-go/core"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("package info", func() {
	var (
		wd  string
		err error
	)

	BeforeEach(func() {
		wd, err = os.Getwd()
		Expect(err).NotTo(HaveOccurred())
	})

	It("returns the package info based on path", func() {
		pkg, err := core.GetPackageInfo(wd, "", "")
		Expect(err).NotTo(HaveOccurred())
		Expect(pkg.Short).To(Equal("core"))
		Expect(pkg.Full).To(Equal("github.com/gyf1214/docker-build-go/core"))
		Expect(pkg.Cmd).To(Equal("github.com/gyf1214/docker-build-go/core"))
		Expect(pkg.Deps).To(Equal(""))
		Expect(pkg.Path).To(Equal(wd))
		Expect(pkg.Build).To(Equal("__build"))
		Expect(pkg.Output).To(Equal("core"))
	})

	It("receive cmd as the cmd to build", func() {
		pkg, err := core.GetPackageInfo(wd, "cmd/foobar", "")
		Expect(err).NotTo(HaveOccurred())
		Expect(pkg.Cmd).To(Equal("github.com/gyf1214/docker-build-go/core/cmd/foobar"))
		Expect(pkg.Output).To(Equal("foobar"))
		pkg, err = core.GetPackageInfo(wd, ".", "")
		Expect(err).NotTo(HaveOccurred())
		Expect(pkg.Cmd).To(Equal("github.com/gyf1214/docker-build-go/core"))
		Expect(pkg.Output).To(Equal("core"))
	})

	It("receive deps as apt deps", func() {
		pkg, err := core.GetPackageInfo(wd, "", "a,b,c")
		Expect(err).NotTo(HaveOccurred())
		Expect(pkg.Deps).To(Equal("a b c"))
	})
})
