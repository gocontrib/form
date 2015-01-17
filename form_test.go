package form

import (
	"net/http"
	"strings"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type LoginForm struct {
	User     string `required:"true"`
	Password string `required:"true"`
}

func TestForm(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Form")
}

var _ = Describe("Form", func() {
	validLogin := func(v interface{}, err error) {
		Expect(err).To(BeNil())
		form, ok := v.(*LoginForm)
		Expect(ok).To(BeTrue())
		Expect(form).NotTo(BeNil())
		Expect(form.User).To(Equal("bob"))
		Expect(form.Password).To(Equal("b0b"))
	}

	It("should be decoded from map", func() {
		var input = map[string]interface{}{
			"user":     "bob",
			"password": "b0b",
		}
		var form = &LoginForm{}
		var err = Decode(form, input)
		validLogin(form, err)
	})

	It("should be decoded from GET request", func() {
		var input, _ = http.NewRequest("GET", "/?user=bob&password=b0b", nil)
		var form = &LoginForm{}
		var err = Decode(form, input)
		validLogin(form, err)
	})

	It("should be decoded from JSON post", func() {
		var body = strings.NewReader(`{"user":"bob", "password": "b0b"}`)
		var input, _ = http.NewRequest("POST", "/", body)
		input.Header.Set("Content-Type", "application/json")
		var form = &LoginForm{}
		var err = Decode(form, input)
		validLogin(form, err)
	})

	XIt("should be decoded from XML post", func() {
		var body = strings.NewReader(`<form><user>bob</user><password>b0b</password></form>`)
		var input, _ = http.NewRequest("POST", "/", body)
		input.Header.Set("Content-Type", "application/xml")
		var form = &LoginForm{}
		var err = Decode(form, input)
		validLogin(form, err)
	})

	It("should be decoded from FORM post", func() {
		var input = NewRequest("/", map[string]interface{}{
			"user":     "bob",
			"password": "b0b",
		})

		var form = &LoginForm{}
		var err = Decode(form, input)
		validLogin(form, err)
	})

	It("error when required field is not set", func() {
		var input = map[string]interface{}{
			"user": "bob",
		}
		var form = &LoginForm{}
		var err = Decode(form, input)
		Expect(err).NotTo(BeNil())
	})
})
