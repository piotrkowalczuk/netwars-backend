package main

import (
	"net/http"
	"net/http/httptest"
	"github.com/martini-contrib/binding"

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
				Context("should return 404 if there is no credentials", func() {
						It("should be a novel", func() {
								request, _ := http.NewRequest("POST", "/login", CreateJSONBody(struct{}{}))
								response := httptest.NewRecorder()
								m.ServeHTTP(response, request)

								Expect(response.Code).To(Equal(http.StatusNotFound))
							})
					})
			})
		Describe("/logout", func() {
				Context("should return 404 if there is no credentials", func() {
						It("should be a novel", func() {
								request, _ := http.NewRequest("POST", "/logout", nil)
								response := httptest.NewRecorder()
								m.ServeHTTP(response, request)

								Expect(response.Code).To(Equal(http.StatusNotFound))
							})
					})
			})
		Describe("/register", func() {
				Context("should return 404 if there is no form send", func() {
						It("should be a novel", func() {
								request, _ := http.NewRequest("POST", "/register", nil)
								response := httptest.NewRecorder()
								m.ServeHTTP(response, request)

								Expect(response.Code).To(Equal(http.StatusNotFound))
							})
					})
			})
	})
