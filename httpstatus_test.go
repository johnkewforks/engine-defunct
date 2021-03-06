package engine

import (
	"fmt"
	"net/http"
	"testing"

	"golang.org/x/net/context"
)

func testCallException(code int, t *testing.T) {
	e, _ := New()
	e.Take("/test", "GET", func(c context.Context) { currentCtx(c).Status(code) })

	w := PerformRequest(e, "GET", "/test")

	if w.Code != code {
		t.Errorf("Status code should be %v, was %d", http.StatusNotFound, w.Code)
	}
}

func TestCallException(t *testing.T) {
	testCallException(404, t)
	testCustomException(418, t)
	testCallException(500, t)
}

func testCustomException(code int, t *testing.T) {
	expected := fmt.Sprintf("CUSTOM %d", code)
	e, _ := New()
	e.HttpStatuses[code].Update(func(c context.Context) { currentCtx(c).RW.Write([]byte(expected)) })

	e.Take("/test", "GET", func(c context.Context) { currentCtx(c).Status(code) })

	w := PerformRequest(e, "GET", "/test")

	if w.Body.String() != expected {
		t.Errorf("Body should be '%s', was but was '%s'.", expected, w.Body.String())
	}
	if w.Code != code {
		t.Errorf("Status code should be %d, was %d", code, w.Code)
	}
}

func TestCustomException(t *testing.T) {
	testCustomException(404, t)
	testCustomException(418, t)
	testCustomException(500, t)
}
