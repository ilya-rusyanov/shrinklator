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

	b := strings.Builder{}

	if _, err := io.Copy(&b, res.Body); err != nil {
		fmt.Println(err)
	}

	if err := res.Body.Close(); err != nil {
		fmt.Println(err)
	}

	fmt.Println(b.String())

	// Output:
	// http://localhost/664b8054bac1af66baafa7a01acd15ee
}
