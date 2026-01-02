package handlers

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/as-ifn-at/REST/config"
	"github.com/as-ifn-at/REST/models"
	"github.com/gin-gonic/gin"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/rs/zerolog"
)

func TestHandlers(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Handlers Suite")
}

var _ = Describe("class test cases", func() {
	var (
		logger zerolog.Logger
		r      *gin.Engine
	)
	BeforeEach(func() {
		logger = zerolog.New(os.Stdout).With().Timestamp().Logger()
		logger = logger.Output(zerolog.ConsoleWriter{Out: os.Stderr})
		r = gin.Default()
		testClassHandler := NewClassHandler(config.Config{}, logger)
		routerG := r.Group("/classes/v1")
		routerG.POST("/create", testClassHandler.Save)
		routerG.GET("/:id", testClassHandler.Get)
	})
	When("save is called", func() {
		It("should give success response with created status", func() {
			payload := []byte(`{
			"class_name":"Pilates",
			"start_date":"01-12-2025",
			"end_date":"20-12-2025",
			"capacity": 20
			}`)
			bodyReader := bytes.NewReader(payload)

			w := httptest.NewRecorder()
			req, err := http.NewRequest(http.MethodPost, "/classes/v1/create", bodyReader)
			Expect(err).To(BeNil())

			r.ServeHTTP(w, req)
			Expect(w.Code).To(Equal(http.StatusCreated))
		})
	})
	When("payload is invalid", func() {
		It("should give error response", func() {
			invalid := []byte(`{
				"class_name":123,
				"start_date":"01-12-2025",
				"end_date":"20-12-2025",
				"capacity": "20"
			}`)
			bodyReader := bytes.NewReader(invalid)

			w := httptest.NewRecorder()
			req, err := http.NewRequest(http.MethodPost, "/classes/v1/create", bodyReader)
			Expect(err).To(BeNil())

			r.ServeHTTP(w, req)
			Expect(w.Code).To(Equal(http.StatusBadRequest))
		})
	})
	When("start-date and end-date is invalid", func() {
		It("should give error response", func() {
			invalid := []byte(`{
				"class_name":"Pilates",
				"start_date":"01-12-2025",
				"end_date":"20-11-2025",
				"capacity": 20
			}`)
			bodyReader := bytes.NewReader(invalid)

			w := httptest.NewRecorder()
			req, err := http.NewRequest(http.MethodPost, "/classes/v1/create", bodyReader)
			Expect(err).To(BeNil())

			r.ServeHTTP(w, req)
			Expect(w.Code).To(Equal(http.StatusInternalServerError))
		})
	})

	When("Get is called", func() {
		It("should return the class", func() {
			Classes["Pilates"] = models.Class{
				ClassName: "Pilates",
				StartDate: "01-12-2025",
				EndDate:   "20-12-2025",
				Capacity:  20,
			}

			w := httptest.NewRecorder()
			req, err := http.NewRequest(http.MethodGet, "/classes/v1/Pilates", nil)
			Expect(err).To(BeNil())

			r.ServeHTTP(w, req)
			data, _ := io.ReadAll(w.Body)

			fmt.Printf("data: %s\n", (data))

			// Expect(data).To(BeEquivalentTo(expectedData))
			Expect(w.Code).To(Equal(http.StatusOK))
		})
	})

	It("should return error if class is absent", func() {
		w := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodGet, "/classes/v1/test", nil)
		Expect(err).To(BeNil())

		r.ServeHTTP(w, req)
		Expect(w.Code).To(Equal(http.StatusNotFound))
	})
})

var _ = Describe("booking test cases", func() {
	var (
		logger zerolog.Logger
		r      *gin.Engine
	)

	BeforeEach(func() {
		Classes["Pilates"] = models.Class{
			ClassName: "Pilates",
			StartDate: "01-12-2025",
			EndDate:   "20-12-2025",
			Capacity:  20,
		}
		logger = zerolog.New(os.Stdout).With().Timestamp().Logger()
		logger = logger.Output(zerolog.ConsoleWriter{Out: os.Stderr})
		r = gin.Default()
		bookingHandler := NewBookingHandler(config.Config{}, logger)
		routerG := r.Group("/bookings/v1")
		routerG.POST("/book", bookingHandler.Save)
		routerG.GET("/:id", bookingHandler.Get)
	})
	When("Save is called", func() {
		It("should give success response", func() {
			payload := []byte(`{
    					"name":"Asif",
    					"date":"02-12-2025",
    					"class_name":"Pilates"
					}`)
			bodyReader := bytes.NewReader(payload)

			w := httptest.NewRecorder()
			req, err := http.NewRequest(http.MethodPost, "/bookings/v1/book", bodyReader)
			Expect(err).To(BeNil())

			r.ServeHTTP(w, req)
			Expect(w.Code).To(Equal(http.StatusCreated))
		})
	})
	When("payload is invalid", func() {
		It("should give error", func() {
			payload := []byte(`{
				"name":12,
				"date":"02-12-2025",
				"class_name":"test-class"
			}`)
			bodyReader := bytes.NewReader(payload)

			w := httptest.NewRecorder()
			req, err := http.NewRequest(http.MethodPost, "/bookings/v1/book", bodyReader)
			Expect(err).To(BeNil())

			r.ServeHTTP(w, req)
			Expect(w.Code).To(Equal(http.StatusBadRequest))
		})
	})

	When("class does not exist", func() {
		It("should give error", func() {
			payload := []byte(`{
				"name":"Asif",
				"date":"02-12-2025",
				"class_name":"test-class"
			}`)
			bodyReader := bytes.NewReader(payload)

			w := httptest.NewRecorder()
			req, err := http.NewRequest(http.MethodPost, "/bookings/v1/book", bodyReader)
			Expect(err).To(BeNil())

			r.ServeHTTP(w, req)
			Expect(w.Code).To(Equal(http.StatusInternalServerError))
		})
	})

	When("booking date is invalid", func() {
		It("should give error", func() {
			payload := []byte(`{
				"name":"Asif",
				"date":"21-12-2025",
				"class_name":"Pilates"
			}`)
			bodyReader := bytes.NewReader(payload)

			w := httptest.NewRecorder()
			req, err := http.NewRequest(http.MethodPost, "/bookings/v1/book", bodyReader)
			Expect(err).To(BeNil())

			r.ServeHTTP(w, req)
			Expect(w.Code).To(Equal(http.StatusInternalServerError))
		})
	})

	When("get is called", func() {
		It("should return the booking", func() {
			payload := []byte(`{
				"name":"Asif",
				"date":"19-12-2025",
				"class_name":"Pilates"
			}`)
			bodyReader := bytes.NewReader(payload)

			w := httptest.NewRecorder()
			req, err := http.NewRequest(http.MethodPost, "/bookings/v1/book", bodyReader)
			Expect(err).To(BeNil())
			r.ServeHTTP(w, req)
			Expect(w.Code).To(Equal(http.StatusCreated))

			w = httptest.NewRecorder()
			req, err = http.NewRequest(http.MethodGet, "/bookings/v1/Asif", nil)
			Expect(err).To(BeNil())
			r.ServeHTTP(w, req)

			Expect(w.Code).To(Equal(http.StatusOK))
		})
	})

	When("member is not present", func() {
		It("return error", func() {
			// payload := []byte(`{
			// 	"name":"test-name",
			// 	"date":"19-12-2025",
			// 	"class_name":"Pilates"
			// }`)
			// bodyReader := bytes.NewReader(payload)

			w := httptest.NewRecorder()
			req, err := http.NewRequest(http.MethodGet, "/bookings/v1/testUser", nil)
			Expect(err).To(BeNil())
			r.ServeHTTP(w, req)
			fmt.Printf("w.Code: %v\n", w.Code)
			Expect(w.Code).To(Equal(http.StatusNotFound))
		})
	})
})
