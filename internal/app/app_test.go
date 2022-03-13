package app_test

import (
	"bytes"
	"errors"
	"github.com/rtakhautdinov/myhttp/internal/app"
	"io/ioutil"
	"net/http"
	"testing"
)

// MockClient is the mock client
type MockClient struct {
	DoFunc func(url string) (*http.Response, error)
}

var (
	// getFunc fetches the mock client's `Get` func
	getFunc func(url string) (*http.Response, error)
)

// Get is the mock client's `Get` func
func (m *MockClient) Get(url string) (resp *http.Response, err error) {
	return getFunc(url)
}

func TestUnitWorker_InvalidUrlError(t *testing.T) {
	mockClient := MockClient{}

	buf := make(chan string)
	results := make(chan app.WorkerResult)
	go app.Worker(&mockClient, buf, results)

	go func() {
		for _, u := range []string{"Aasdasd{{{}}}"} {
			buf <- u
		}
	}()

	assertError(t, <-results)
}

func TestUnitWorker_SingleUrlClientResponseError(t *testing.T) {
	mockClient := MockClient{}

	getFunc = func(string) (*http.Response, error) {
		return &http.Response{}, errors.New("an error occurred")
	}

	buf := make(chan string)
	results := make(chan app.WorkerResult)
	go app.Worker(&mockClient, buf, results)

	go func() {
		for _, u := range []string{"test.com"} {
			buf <- u
		}
	}()

	assertError(t, <-results)
}

func TestUnitWorker_ClientResponseSuccess(t *testing.T) {
	mockClient := MockClient{}

	getFunc = func(string) (*http.Response, error) {
		return &http.Response{
			Body: ioutil.NopCloser(bytes.NewBufferString("Hello World")),
		}, nil
	}

	buf := make(chan string)
	results := make(chan app.WorkerResult)
	go app.Worker(&mockClient, buf, results)

	go func() {
		for _, u := range []string{"test.com"} {
			buf <- u
		}
	}()

	result := <-results
	assertSuccess(t, result)
	assertEquals(t, result.Url, "https://test.com")
	assertEquals(t, result.Md5Data, "b10a8db164e0754105b7a99be72e3fe5")
}

func TestUnitWorker_ClientResponseTwoSuccess(t *testing.T) {
	mockClient := MockClient{}

	getFunc = func(string) (*http.Response, error) {
		return &http.Response{
			Body: ioutil.NopCloser(bytes.NewBufferString("Hello World")),
		}, nil
	}

	buf := make(chan string)
	results := make(chan app.WorkerResult)
	go app.Worker(&mockClient, buf, results)

	go func() {
		for _, u := range []string{"test.com", "test2.com"} {
			buf <- u
		}
	}()

	result1 := <-results
	assertSuccess(t, result1)
	assertEquals(t, result1.Url, "https://test.com")
	assertEquals(t, result1.Md5Data, "b10a8db164e0754105b7a99be72e3fe5")

	result2 := <-results
	assertSuccess(t, result2)
	assertEquals(t, result2.Url, "https://test2.com")
	assertEquals(t, result2.Md5Data, "b10a8db164e0754105b7a99be72e3fe5")
}

func assertError(t *testing.T, actual app.WorkerResult) {
	t.Helper()

	if actual.IsSuccess == true {
		t.Error("it is expected to have error but received success")
	}
}

func assertSuccess(t *testing.T, actual app.WorkerResult) {
	t.Helper()

	if actual.IsSuccess == false {
		t.Error("it is expected to have success but received fail")
	}
}

func assertEquals(t *testing.T, actual, expected string) {
	t.Helper()

	if expected != actual {
		t.Errorf("expected != actual: %s != %s\n", expected, actual)
	}
}
