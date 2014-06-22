package main

import (
	"github.com/martini-contrib/binding"
	"net/http"
	"net/http/httptest"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("user_api", func() {
	BeforeEach(func() {
		InitializeEnvironment()
	})
	Describe("/login", func() {
		BeforeEach(func() {
			m.Post("/login", binding.Json(LoginCredentials{}), loginHandler)
		})
		Context("user should be able to log in", func() {
			It("should return 404 if there is no credentials", func() {
				request, _ := http.NewRequest("POST", "/login", CreateJSONBody(struct{}{}))
				response := httptest.NewRecorder()
				m.ServeHTTP(response, request)

				Expect(response.Code).To(Equal(http.StatusNotFound))
			})
		})
		Context("should return user session object if credentials match", func() {
			It("should return 404 if there is no credentials", func() {
				request, _ := http.NewRequest("POST", "/login", CreateJSONBody(struct{}{}))
				response := httptest.NewRecorder()
				m.ServeHTTP(response, request)

				Expect(response.Code).To(Equal(http.StatusNotFound))
			})
		})
	})
	Describe("/logout", func() {
		Context("user should be able to log out", func() {
			It("should return 404 if there is no credentials", func() {
				request, _ := http.NewRequest("POST", "/logout", nil)
				response := httptest.NewRecorder()
				m.ServeHTTP(response, request)

				Expect(response.Code).To(Equal(http.StatusNotFound))
			})
		})
	})
	Describe("/register", func() {
		Context("user should be able to register", func() {
			It("should return 404 if there is no form send", func() {
				request, _ := http.NewRequest("POST", "/register", nil)
				response := httptest.NewRecorder()
				m.ServeHTTP(response, request)

				Expect(response.Code).To(Equal(http.StatusNotFound))
			})
		})
	})
})
