package testing

import (
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestFiberRequestBody(t *testing.T) {
	type Expense struct {
		CategoryID  int     `json:"category_id"`
		Title       string  `json:"title"`
		Amount      int     `json:"amount"`
		Description *string `json:"description,omitempty"`
	}

	type ResponseAPI struct {
		StatusCode int    `json:"status_code,omitempty"`
		Status     string `json:"status,omitempty"`
		Message    string `json:"message,omitempty"`
		Data       any    `json:"data,omitempty"`
	}

	scenario := []struct {
		Name             string
		Request          map[string]any
		ExpectStatusCode int
	}{
		{
			Name: "test with description",
			Request: map[string]any{
				"category_id": 1,
				"title":       "Hello with description",
				"amount":      1,
				"description": "hello world",
			},
			ExpectStatusCode: http.StatusOK,
		},
		{
			Name: "test without description",
			Request: map[string]any{
				"category_id": 2,
				"title":       "Hello without description",
				"amount":      2,
			},
			ExpectStatusCode: http.StatusOK,
		},
	}

	for _, testCase := range scenario {
		t.Run(testCase.Name, func(t *testing.T) {
			app := fiber.New()
			app.Post("/test", func(ctx *fiber.Ctx) error {
				var request Expense
				err := ctx.BodyParser(&request)
				if err != nil {
					statusCode := http.StatusBadRequest
					ctx.Status(statusCode)
					return ctx.JSON(&ResponseAPI{
						StatusCode: statusCode,
						Status:     "bad request",
						Message:    err.Error(),
					})
				}

				// ok
				statusCode := http.StatusOK
				ctx.Status(statusCode)
				return ctx.JSON(ResponseAPI{
					StatusCode: statusCode,
					Status:     "ok",
					Message:    "success",
					Data: map[string]any{
						"title":       request.Title,
						"description": request.Description,
					},
				})
			})

			reqJson, _ := json.Marshal(&testCase.Request)
			request := httptest.NewRequest(http.MethodPost, "/test", strings.NewReader(string(reqJson)))
			request.Header.Add("Content-Type", "application/json")

			response, err := app.Test(request)
			assert.Nil(t, err)
			assert.NotNil(t, response)

			assert.Equal(t, response.StatusCode, testCase.ExpectStatusCode)

			// get response body
			body, err := io.ReadAll(response.Body)
			assert.NotNil(t, body)
			assert.Nil(t, err)

			// log response body
			log.Println(string(body))
		})
	}
}
