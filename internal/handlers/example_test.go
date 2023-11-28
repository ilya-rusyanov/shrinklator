package handlers

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/ilya-rusyanov/shrinklator/internal/services"
	"github.com/ilya-rusyanov/shrinklator/internal/storage"
)

func ExampleShorten_Handler() {
	log := dummyLogger{}

	repository := storage.NewInMemory(&log)
	service := services.NewShortener(&log, repository, services.MD5Algo)

	handler := NewShorten(&log, service, "http://localhost")

	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader("http://yandex.ru"))
	rw := httptest.NewRecorder()

	handler.Handler(rw, req)
	res := rw.Result()
	defer res.Body.Close()
	b := strings.Builder{}
	io.Copy(&b, res.Body)
	fmt.Println(b.String())

	// Output:
	// http://localhost/664b8054bac1af66baafa7a01acd15ee
}
