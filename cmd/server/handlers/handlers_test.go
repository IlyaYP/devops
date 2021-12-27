package handlers

import (
	"github.com/IlyaYP/devops/storage/inmemory"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUpdateHandler(t *testing.T) {
	st := inmemory.NewStorage()
	// определяем структуру теста
	type want struct {
		code        int
		response    string
		contentType string
	}
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name     string
		want     want
		args     args
		endpoint string
	}{
		{
			name:     "positive test #1",           // TODO: Add test cases.
			endpoint: "/update/gauge/Alloc/201456", //"http://localhost:8080/update/gauage/Alloc/201456",
			want: want{
				code:        200,
				response:    `OK`,
				contentType: "text/plain",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			//UpdateHandler(tt.args.w, tt.args.r)

			request := httptest.NewRequest(http.MethodPost, tt.endpoint, nil)
			// создаём новый Recorder
			w := httptest.NewRecorder()

			// определяем хендлер
			h := http.HandlerFunc(UpdateHandler(st))

			// запускаем сервер
			h.ServeHTTP(w, request)
			res := w.Result()

			// проверяем код ответа
			if res.StatusCode != tt.want.code {
				t.Errorf("Expected status code %d, got %d", tt.want.code, w.Code)
			}

			// получаем и проверяем тело запроса
			//defer res.Body.Close()
			//resBody, err := io.ReadAll(res.Body)
			//if err != nil {
			//	t.Fatal(err)
			//}
			//if string(resBody) != tt.want.response {
			//	t.Errorf("Expected body %s, got %s", tt.want.response, w.Body.String())
			//}

			// заголовок ответа
			//if res.Header.Get("Content-Type") != tt.want.contentType {
			//	t.Errorf("Expected Content-Type %s, got %s", tt.want.contentType, res.Header.Get("Content-Type"))
			//}

		})
	}
}
