package config_test

import (
	"testing"

	. "github.com/sjeandeaux/todo/pkg/config"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"os"
)

func TestConfig(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Config Suite")
}

var _ = Describe("Config", func() {
	BeforeEach(func() {
		os.Setenv("EXISTING_KEY_WITH_TEXT", "EXISTING_VALUE")
		os.Setenv("EXISTING_KEY_EMPTY_STRING", "")
	})

	AfterEach(func() {
		os.Unsetenv("EXISTING_KEY")
		os.Unsetenv("EXISTING_KEY_EMPTY_STRING")
	})

	Describe("Get the environmental variable", func() {
		Context("With an existing key with text", func() {
			It("should return the value with text", func() {
				Ω(LookupEnvOrString("EXISTING_KEY_WITH_TEXT", "DEFAULT_VALUE")).Should(Equal("EXISTING_VALUE"))
			})
		})

		Context("With an existing key with empty string", func() {
			It("should return an empty string", func() {
				Ω(LookupEnvOrString("EXISTING_KEY_EMPTY_STRING", "DEFAULT_VALUE")).To(Equal(""))
			})
		})

		Context("With non-existing key", func() {
			It("should return the default value", func() {
				Ω(LookupEnvOrString("NON_EXISTING_VALUE", "DEFAULT_VALUE")).Should(Equal("DEFAULT_VALUE"))
			})
		})
	})
})
