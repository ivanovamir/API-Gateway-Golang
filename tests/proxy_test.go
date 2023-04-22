package tests

import (
	"fmt"
	"gateway/pkg/config"
	"gateway/pkg/proxy"
	"github.com/julienschmidt/httprouter"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRedirect(t *testing.T) {

	data := config.Data{
		Data: []config.Requests{
			{
				Path:      "/test_1",
				Url:       "https://jsonplaceholder.typicode.com/posts/1",
				Method:    "GET",
				MakeProxy: false,
			},
			{
				Path:      "/test_2",
				Url:       "https://jsonplaceholder.typicode.com/posts",
				Method:    "POST",
				MakeProxy: false,
			},
			{
				Path:      "/test_3",
				Url:       "https://jsonplaceholder.typicode.com/posts/1",
				Method:    "PUT",
				MakeProxy: false,
			},
			{
				Path:      "/test_4",
				Url:       "https://jsonplaceholder.typicode.com/posts/1",
				Method:    "PATCH",
				MakeProxy: false,
			},
			{
				Path:      "/test_5",
				Url:       "https://jsonplaceholder.typicode.com/posts/1",
				Method:    "DELETE",
				MakeProxy: false,
			},
		},
	}

	testTable := []struct {
		name                 string
		path                 string
		expectStatusCode     int
		expectedResponseBody string
	}{
		{
			name:                 "test_1",
			path:                 "/test_1",
			expectStatusCode:     http.StatusOK,
			expectedResponseBody: "{\n  \"userId\": 1,\n  \"id\": 1,\n  \"title\": \"sunt aut facere repellat provident occaecati excepturi optio reprehenderit\",\n  \"body\": \"quia et suscipit\\nsuscipit recusandae consequuntur expedita et cum\\nreprehenderit molestiae ut ut quas totam\\nnostrum rerum est autem sunt rem eveniet architecto\"\n}",
		},
		{
			name:                 "test_2",
			path:                 "/test_2",
			expectStatusCode:     http.StatusCreated,
			expectedResponseBody: "{\n  \"id\": 101\n}",
		},
		{
			name:                 "test_3",
			path:                 "/test_3",
			expectStatusCode:     http.StatusOK,
			expectedResponseBody: "{\n  \"id\": 1\n}",
		},
		{
			name:                 "test_4",
			path:                 "/test_4",
			expectStatusCode:     http.StatusOK,
			expectedResponseBody: "{\n  \"userId\": 1,\n  \"id\": 1,\n  \"title\": \"sunt aut facere repellat provident occaecati excepturi optio reprehenderit\",\n  \"body\": \"quia et suscipit\\nsuscipit recusandae consequuntur expedita et cum\\nreprehenderit molestiae ut ut quas totam\\nnostrum rerum est autem sunt rem eveniet architecto\"\n}",
		},
		{
			name:                 "test_5",
			path:                 "/test_5",
			expectStatusCode:     http.StatusOK,
			expectedResponseBody: "{}",
		},
	}

	for i, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			router := httprouter.New()

			router.Handle(data.Data[i].Method, data.Data[i].Path, proxy.NewProxy(
				proxy.WithProxy(data.Data[i].MakeProxy),
				proxy.WithProxyUrl(data.Data[i].ProxyUrl),
				proxy.WithRedirectUrl(data.Data[i].Url),
				proxy.WithLog(nil),
				proxy.WithExpectedStatusCodes(data.Data[i].ExpectedProxyStatusCodes),
				proxy.WithProxyMethod(data.Data[i].ProxyMethod),
			).Redirect())

			// Test request
			w := httptest.NewRecorder()
			req := httptest.NewRequest(data.Data[i].Method, testCase.path, nil)
			// Perform Request
			router.ServeHTTP(w, req)

			if w.Code != testCase.expectStatusCode {
				t.Errorf("Expected status code %d, got %d", testCase.expectStatusCode, w.Code)
			}

			body, err := io.ReadAll(w.Body)

			if err != nil {
				t.Errorf("Unexpected error while reading body: %v", err)
			}

			if string(body) != testCase.expectedResponseBody {
				fmt.Println()
				t.Errorf("Expected body to be empty")
			}
		})
	}
}

func TestProxy(t *testing.T) {
	data := config.Data{
		Data: []config.Requests{
			{
				Path:        "/test_1",
				Url:         "https://jsonplaceholder.typicode.com/posts/1",
				Method:      "GET",
				MakeProxy:   true,
				ProxyUrl:    "https://dummyjson.com/products/1",
				ProxyMethod: "GET",
				ExpectedProxyStatusCodes: []config.ExpectedStatusCodes{
					{
						StatusCode: "200",
					},
				},
			},
			{
				Path:        "/test_2",
				Url:         "https://dummyjson.com/products/1",
				Method:      "DELETE",
				MakeProxy:   true,
				ProxyUrl:    "https://dummyjson.com/products/123",
				ProxyMethod: "GET",
				ExpectedProxyStatusCodes: []config.ExpectedStatusCodes{
					{
						StatusCode: "200",
					},
				},
			},
			{
				Path:        "/test_3",
				Url:         "https://dummyjson.com/products/category/smartphones",
				Method:      "GET",
				MakeProxy:   true,
				ProxyUrl:    "https://dummyjson.com/products/add",
				ProxyMethod: "POST",
				ExpectedProxyStatusCodes: []config.ExpectedStatusCodes{
					{
						StatusCode: "200",
					},
				},
			},
			{
				Path:        "/test_4",
				Url:         "https://dummyjson.com/products/category/smartphones",
				Method:      "GET",
				MakeProxy:   true,
				ProxyUrl:    "https://dummyjson.com/products/add",
				ProxyMethod: "POST",
				ExpectedProxyStatusCodes: []config.ExpectedStatusCodes{
					{
						StatusCode: "200",
					},
				},
			},
		},
	}

	testTable := []struct {
		name                 string
		path                 string
		headerAuth           bool
		headerAuthValue      string
		expectStatusCode     int
		expectedResponseBody string
	}{
		{
			name:                 "test_1",
			path:                 "/test_1",
			expectStatusCode:     http.StatusOK,
			expectedResponseBody: "{\n  \"userId\": 1,\n  \"id\": 1,\n  \"title\": \"sunt aut facere repellat provident occaecati excepturi optio reprehenderit\",\n  \"body\": \"quia et suscipit\\nsuscipit recusandae consequuntur expedita et cum\\nreprehenderit molestiae ut ut quas totam\\nnostrum rerum est autem sunt rem eveniet architecto\"\n}",
		},
		{
			name:                 "test_2",
			path:                 "/test_2",
			expectStatusCode:     http.StatusNotFound,
			expectedResponseBody: "{\"message\":\"Product with id '123' not found\"}",
		},
		{
			name:                 "test_3",
			path:                 "/test_3",
			expectStatusCode:     http.StatusOK,
			expectedResponseBody: "{\"products\":[{\"id\":1,\"title\":\"iPhone 9\",\"description\":\"An apple mobile which is nothing like apple\",\"price\":549,\"discountPercentage\":12.96,\"rating\":4.69,\"stock\":94,\"brand\":\"Apple\",\"category\":\"smartphones\",\"thumbnail\":\"https://i.dummyjson.com/data/products/1/thumbnail.jpg\",\"images\":[\"https://i.dummyjson.com/data/products/1/1.jpg\",\"https://i.dummyjson.com/data/products/1/2.jpg\",\"https://i.dummyjson.com/data/products/1/3.jpg\",\"https://i.dummyjson.com/data/products/1/4.jpg\",\"https://i.dummyjson.com/data/products/1/thumbnail.jpg\"]},{\"id\":2,\"title\":\"iPhone X\",\"description\":\"SIM-Free, Model A19211 6.5-inch Super Retina HD display with OLED technology A12 Bionic chip with ...\",\"price\":899,\"discountPercentage\":17.94,\"rating\":4.44,\"stock\":34,\"brand\":\"Apple\",\"category\":\"smartphones\",\"thumbnail\":\"https://i.dummyjson.com/data/products/2/thumbnail.jpg\",\"images\":[\"https://i.dummyjson.com/data/products/2/1.jpg\",\"https://i.dummyjson.com/data/products/2/2.jpg\",\"https://i.dummyjson.com/data/products/2/3.jpg\",\"https://i.dummyjson.com/data/products/2/thumbnail.jpg\"]},{\"id\":3,\"title\":\"Samsung Universe 9\",\"description\":\"Samsung's new variant which goes beyond Galaxy to the Universe\",\"price\":1249,\"discountPercentage\":15.46,\"rating\":4.09,\"stock\":36,\"brand\":\"Samsung\",\"category\":\"smartphones\",\"thumbnail\":\"https://i.dummyjson.com/data/products/3/thumbnail.jpg\",\"images\":[\"https://i.dummyjson.com/data/products/3/1.jpg\"]},{\"id\":4,\"title\":\"OPPOF19\",\"description\":\"OPPO F19 is officially announced on April 2021.\",\"price\":280,\"discountPercentage\":17.91,\"rating\":4.3,\"stock\":123,\"brand\":\"OPPO\",\"category\":\"smartphones\",\"thumbnail\":\"https://i.dummyjson.com/data/products/4/thumbnail.jpg\",\"images\":[\"https://i.dummyjson.com/data/products/4/1.jpg\",\"https://i.dummyjson.com/data/products/4/2.jpg\",\"https://i.dummyjson.com/data/products/4/3.jpg\",\"https://i.dummyjson.com/data/products/4/4.jpg\",\"https://i.dummyjson.com/data/products/4/thumbnail.jpg\"]},{\"id\":5,\"title\":\"Huawei P30\",\"description\":\"Huawei’s re-badged P30 Pro New Edition was officially unveiled yesterday in Germany and now the device has made its way to the UK.\",\"price\":499,\"discountPercentage\":10.58,\"rating\":4.09,\"stock\":32,\"brand\":\"Huawei\",\"category\":\"smartphones\",\"thumbnail\":\"https://i.dummyjson.com/data/products/5/thumbnail.jpg\",\"images\":[\"https://i.dummyjson.com/data/products/5/1.jpg\",\"https://i.dummyjson.com/data/products/5/2.jpg\",\"https://i.dummyjson.com/data/products/5/3.jpg\"]}],\"total\":5,\"skip\":0,\"limit\":5}",
		},
		{
			name:                 "test_4_auth_header",
			path:                 "/test_4",
			headerAuth:           true,
			headerAuthValue:      "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6MTUsInVzZXJuYW1lIjoia21pbmNoZWxsZSIsImVtYWlsIjoia21pbmNoZWxsZUBxcS5jb20iLCJmaXJzdE5hbWUiOiJKZWFubmUiLCJsYXN0TmFtZSI6IkhhbHZvcnNvbiIsImdlbmRlciI6ImZlbWFsZSIsImltYWdlIjoiaHR0cHM6Ly9yb2JvaGFzaC5vcmcvYXV0cXVpYXV0LnBuZyIsImlhdCI6MTY4MjE3MjM5NiwiZXhwIjoxNjg0NzY0Mzk2fQ.G47X8LJUMMKPpIOTCNhsY6WXIEv3RdgHSTA3k8pnGww",
			expectStatusCode:     http.StatusOK,
			expectedResponseBody: "{\"products\":[{\"id\":1,\"title\":\"iPhone 9\",\"description\":\"An apple mobile which is nothing like apple\",\"price\":549,\"discountPercentage\":12.96,\"rating\":4.69,\"stock\":94,\"brand\":\"Apple\",\"category\":\"smartphones\",\"thumbnail\":\"https://i.dummyjson.com/data/products/1/thumbnail.jpg\",\"images\":[\"https://i.dummyjson.com/data/products/1/1.jpg\",\"https://i.dummyjson.com/data/products/1/2.jpg\",\"https://i.dummyjson.com/data/products/1/3.jpg\",\"https://i.dummyjson.com/data/products/1/4.jpg\",\"https://i.dummyjson.com/data/products/1/thumbnail.jpg\"]},{\"id\":2,\"title\":\"iPhone X\",\"description\":\"SIM-Free, Model A19211 6.5-inch Super Retina HD display with OLED technology A12 Bionic chip with ...\",\"price\":899,\"discountPercentage\":17.94,\"rating\":4.44,\"stock\":34,\"brand\":\"Apple\",\"category\":\"smartphones\",\"thumbnail\":\"https://i.dummyjson.com/data/products/2/thumbnail.jpg\",\"images\":[\"https://i.dummyjson.com/data/products/2/1.jpg\",\"https://i.dummyjson.com/data/products/2/2.jpg\",\"https://i.dummyjson.com/data/products/2/3.jpg\",\"https://i.dummyjson.com/data/products/2/thumbnail.jpg\"]},{\"id\":3,\"title\":\"Samsung Universe 9\",\"description\":\"Samsung's new variant which goes beyond Galaxy to the Universe\",\"price\":1249,\"discountPercentage\":15.46,\"rating\":4.09,\"stock\":36,\"brand\":\"Samsung\",\"category\":\"smartphones\",\"thumbnail\":\"https://i.dummyjson.com/data/products/3/thumbnail.jpg\",\"images\":[\"https://i.dummyjson.com/data/products/3/1.jpg\"]},{\"id\":4,\"title\":\"OPPOF19\",\"description\":\"OPPO F19 is officially announced on April 2021.\",\"price\":280,\"discountPercentage\":17.91,\"rating\":4.3,\"stock\":123,\"brand\":\"OPPO\",\"category\":\"smartphones\",\"thumbnail\":\"https://i.dummyjson.com/data/products/4/thumbnail.jpg\",\"images\":[\"https://i.dummyjson.com/data/products/4/1.jpg\",\"https://i.dummyjson.com/data/products/4/2.jpg\",\"https://i.dummyjson.com/data/products/4/3.jpg\",\"https://i.dummyjson.com/data/products/4/4.jpg\",\"https://i.dummyjson.com/data/products/4/thumbnail.jpg\"]},{\"id\":5,\"title\":\"Huawei P30\",\"description\":\"Huawei’s re-badged P30 Pro New Edition was officially unveiled yesterday in Germany and now the device has made its way to the UK.\",\"price\":499,\"discountPercentage\":10.58,\"rating\":4.09,\"stock\":32,\"brand\":\"Huawei\",\"category\":\"smartphones\",\"thumbnail\":\"https://i.dummyjson.com/data/products/5/thumbnail.jpg\",\"images\":[\"https://i.dummyjson.com/data/products/5/1.jpg\",\"https://i.dummyjson.com/data/products/5/2.jpg\",\"https://i.dummyjson.com/data/products/5/3.jpg\"]}],\"total\":5,\"skip\":0,\"limit\":5}",
		},
	}

	for i, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			router := httprouter.New()

			router.Handle(data.Data[i].Method, data.Data[i].Path, proxy.NewProxy(
				proxy.WithProxy(data.Data[i].MakeProxy),
				proxy.WithProxyUrl(data.Data[i].ProxyUrl),
				proxy.WithRedirectUrl(data.Data[i].Url),
				proxy.WithLog(nil),
				proxy.WithExpectedStatusCodes(data.Data[i].ExpectedProxyStatusCodes),
				proxy.WithProxyMethod(data.Data[i].ProxyMethod),
			).Redirect())

			// Test request
			w := httptest.NewRecorder()
			req := httptest.NewRequest(data.Data[i].Method, testCase.path, nil)

			if testCase.headerAuth {
				req.Header.Set("Authorization", testCase.headerAuthValue)
			}

			// Perform Request
			router.ServeHTTP(w, req)

			if w.Code != testCase.expectStatusCode {
				t.Errorf("Expected status code %d, got %d", testCase.expectStatusCode, w.Code)
			}

			body, err := io.ReadAll(w.Body)

			if err != nil {
				t.Errorf("Unexpected error while reading body: %v", err)
			}

			if string(body) != testCase.expectedResponseBody {
				fmt.Println()
				t.Errorf("Expected body to be empty")
			}
		})
	}
}
