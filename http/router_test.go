package http

import (
	"net/http"
	"testing"
	"fmt"
	"time"
	"encoding/json"
)


type In struct {
	Name     string `json:"name"`
	Birthday string `json:"birthday"`
}

type Out struct {
	Text string `json:"text"`
}

func HelloPost(in *In, out *Out) *ErrorResponse {
	out.Text = fmt.Sprintf(
		"Hello %s, today is %s, your birthday is %s.",
		in.Name,
		time.Now().Format(DefaultDateFM),
		in.Birthday,
	)
	return nil
}

func HelloPost2(in *In, out *Out, c *Context) *ErrorResponse {
	out.Text = fmt.Sprintf(
		"Hello %s, today is %s, your birthday is %s.",
		in.Name,
		time.Now().Format(DefaultDateFM),
		c.GetUrlParam("birthday"),
	)
	return nil
}

func HelloGet(c *Context) *ErrorResponse {
	text := fmt.Sprintf(
		"Hello %s, today is %s, your birthday is %s.",
		c.GetUrlParam("name"),
		time.Now().Format(DefaultDateFM),
		c.GetUrlParam("birthday"),
	)
	c.Json(map[string]interface{}{"text": text})
	return nil
}

func TestHttp(t *testing.T) {
	router := New()
	router.POST("/hello", HelloPost)
	router.POST("/hello2", HelloPost2)
	router.GET("/hello", HelloGet)
	go http.ListenAndServe(":8080", router)
}

func TestRouter_POST(t *testing.T) {
	url := "http://127.0.0.1:8080/hello"
	data := `{"name":"Lywane","birthday":"1994-06-25"}`

	response, err := Post(url, []byte(data))
	if err != nil {
		t.Fatal(err)
	}
	validResult(t, response)
}

func TestRouter_POST2(t *testing.T) {
	url := "http://127.0.0.1:8080/hello2?birthday=1994-06-25"
	data := `{"name":"Lywane"}`

	response, err := Post(url, []byte(data))
	if err != nil {
		t.Fatal(err)
	}
	validResult(t, response)
}

func TestRouter_GET(t *testing.T) {
	url := "http://127.0.0.1:8080/hello?name=Lywane&birthday=1994-06-25"

	response, err := Get(url)
	if err != nil {
		t.Fatal(err)
	}
	validResult(t, response)
}

func validResult(t *testing.T, response []byte) {
	res := struct {
		Status  int    `json:"status"`
		Message string `json:"message"`
		Data *struct {
			Text string `json:"text"`
		} `json:"data"`
	}{}
	err := json.Unmarshal(response, &res)
	if err != nil {
		t.Fatal(response,err)
	}
	if res.Status != 0 {
		t.Fatal(res.Message)
	}
	if res.Data.Text != fmt.Sprintf(
		"Hello %s, today is %s, your birthday is %s.",
		"Lywane",
		time.Now().Format(DefaultDateFM),
		"1994-06-25") {

		t.Fatal("hanler err", res.Data.Text)
	}
}
