package myhttp

import (
	"net/http"
	"testing"
	"fmt"
	"time"
	"io"
	"io/ioutil"
	"bytes"
	"encoding/json"
)

const DefaultDateFM string = "2006-01-02"

type In struct {
	Name string `json:"name"`
}

type Out struct {
	Text string `json:"text"`
}

func Hello(in *In, out *Out, urlParam UrlParam) *ErrorResponse {
	birthday := urlParam.Get("birthday")
	if birthday == "" {
		birthday = "unknown"
	}
	out.Text = fmt.Sprintf(
		"Hello %s, today is %s, your birthday is %s.",
		in.Name,
		time.Now().Format(DefaultDateFM),
		birthday)
	return nil
}

func TestHttp(t *testing.T) {
	router := New()
	router.POST("/hello", Hello)
	go http.ListenAndServe(":8080", router)

	url := "http://127.0.0.1:8080/hello?birthday=1994-06-25"
	data := `{"name":"Lywane"}`

	response, err := post(url, []byte(data))
	if err != nil {
		t.Fatal(err)
	}
	res := struct {
		Status  int    `json:"status"`
		Message string `json:"message"`
		Data *struct {
			Text string `json:"text"`
		} `json:"data"`
	}{}
	err = json.Unmarshal(response, &res)
	if err != nil {
		t.Fatal(err)
	}
	if res.Status != 0 {
		t.Fatal(res.Message)
	}
	if res.Data.Text != fmt.Sprintf(
		"Hello %s, today is %s, your birthday is %s.",
		"Lywane",
		time.Now().Format(DefaultDateFM),
		"1994-06-25") {
		t.Fatal("hanler err")
	}
}

func post(url string, body []byte) ([]byte, error) {
	data := bytes.NewReader(body)
	return request("POST", url, data)
}

func request(method, url string, body io.Reader) ([]byte, error) {
	request, err := http.NewRequest(method, url, body)
	if err != nil {
		return []byte(""), err
	}
	if method == "POST" {
		request.Header.Set("Content-Type", "application/json")
	}
	client := http.Client{Timeout: time.Duration(20 * time.Second)}
	resp, err := client.Do(request)
	if err != nil {
		return []byte(""), err
	}
	defer resp.Body.Close()
	ret, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte(""), err
	}
	return ret, nil
}
