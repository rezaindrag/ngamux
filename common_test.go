package ngamux

import (
	"context"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
)

func TestWithMiddlewares(t *testing.T) {
	result := WithMiddlewares()(func(rw http.ResponseWriter, r *http.Request) error {
		return nil
	})
	if result == nil {
		t.Errorf("TestWithMiddlewares need %v, but got %v", reflect.TypeOf(result), nil)
	}

	result = WithMiddlewares(nil)(func(rw http.ResponseWriter, r *http.Request) error {
		return nil
	})
	if result == nil {
		t.Errorf("TestWithMiddlewares need %v, but got %v", reflect.TypeOf(result), nil)
	}

	result = WithMiddlewares(nil)(nil)
	if result != nil {
		t.Errorf("TestWithMiddlewares need %v, but got %v", nil, reflect.TypeOf(result))
	}
}

func TestGetParam(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req = req.WithContext(context.WithValue(req.Context(), KeyContextParams, [][]string{{"id", "1"}}))
	result := GetParam(req, "id")

	if result != "1" {
		t.Errorf("TestGetParam need %v, but got %v", "1", result)
	}

	result = GetParam(req, "slug")
	if result != "" {
		t.Errorf("TestGetParam need %v, but got %v", "\"\"", result)
	}
}

func TestGetQuery(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/?id=1", nil)
	result := GetQuery(req, "id")

	if result != "1" {
		t.Errorf("TestGetQuery need %v, but got %v", "1", result)
	}

	result = GetQuery(req, "slug", "undefined")
	if result != "undefined" {
		t.Errorf("TestGetQuery need %v, but got %v", "undefined", result)
	}

	result = GetQuery(req, "slug")
	if result != "" {
		t.Errorf("TestGetQuery need %v, but got %v", "\"\"", result)
	}
}

func TestGetJSON(t *testing.T) {
	input := strings.NewReader(`{"id": 1}`)
	req := httptest.NewRequest(http.MethodGet, "/", input)

	var data map[string]interface{}
	err := GetJSON(req, &data)
	if err != nil {
		t.Errorf("TestGetJSON need %v, but got %v", "nil", err)
	}

	if data["id"] == nil {
		t.Errorf("TestGetJSON need %v, but got %v", "value", data["id"])
	}

	id, ok := data["id"].(float64)
	if !ok {
		t.Errorf("TestGetJSON need %v, but got %v", "true", ok)
	}

	if id != 1 {
		t.Errorf("TestGetJSON need %v, but got %v", 1, id)
	}
}

func TestSetContextValue(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req = SetContextValue(req, "id", 1)

	id := req.Context().Value("id")
	if id != 1 {
		t.Errorf("TestSetContextValue need %v, but got %v", 1, id)
	}

	slug := req.Context().Value("slug")
	if id != 1 {
		t.Errorf("TestSetContextValue need %v, but got %v", nil, slug)
	}
}

func TestGetContextValue(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req = req.WithContext(context.WithValue(req.Context(), "id", 1))

	id := GetContextValue(req, "id")
	if id != 1 {
		t.Errorf("TestGetContextValue need %v, but got %v", 1, id)
	}

	slug := GetContextValue(req, "slug")
	if id != 1 {
		t.Errorf("TestGetContextValue need %v, but got %v", nil, slug)
	}
}

func TestString(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	handler := http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		String(rw, "ok")
	})
	handler.ServeHTTP(rec, req)

	result := rec.Body.String()
	expected := "ok\n"
	if result != expected {
		t.Errorf("TestString need %v, but got %v", expected, result)
	}
}

func TestStringWithStatus(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	handler := http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		StringWithStatus(rw, http.StatusOK, "ok")
	})
	handler.ServeHTTP(rec, req)

	resultBody := rec.Body.String()
	expectedBody := "ok\n"
	if resultBody != expectedBody {
		t.Errorf("TestStringWithStatus need %v, but got %v", expectedBody, resultBody)
	}

	resultStatus := rec.Result().StatusCode
	expectedStatus := http.StatusOK
	if resultStatus != expectedStatus {
		t.Errorf("TestStringWithStatus need %v, but got %v", expectedStatus, resultStatus)
	}
}

func TestJSON(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	handler := http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		JSON(rw, Map{
			"id": 1,
		})
	})
	handler.ServeHTTP(rec, req)

	resultBody := strings.ReplaceAll(rec.Body.String(), "\n", "")
	expectedBody := `{"id":1}`
	if resultBody != expectedBody {
		t.Errorf("TestJSON need %v, but got %v", expectedBody, resultBody)
	}
}

func TestJSONWithStatus(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	handler := http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		JSONWithStatus(rw, http.StatusOK, Map{
			"id": 1,
		})
	})
	handler.ServeHTTP(rec, req)

	resultBody := strings.ReplaceAll(rec.Body.String(), "\n", "")
	expectedBody := `{"id":1}`
	if resultBody != expectedBody {
		t.Errorf("TestJSONWithStatus need %v, but got %v", expectedBody, resultBody)
	}

	resultStatus := rec.Result().StatusCode
	expectedStatus := http.StatusOK
	if resultStatus != expectedStatus {
		t.Errorf("TestJSONWithStatus need %v, but got %v", expectedStatus, resultStatus)
	}
}
