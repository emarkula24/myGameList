package test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"example.com/mygamelist/handler"
	"example.com/mygamelist/utils"
)

var dummyEnv = &utils.Env{}

func TestGameHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/games?query=metroid", nil)

	rr := httptest.NewRecorder()

	handlerFunc := handler.Search(dummyEnv)
	handlerFunc.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status 200 but got %d", rr.Code)
	}

}
