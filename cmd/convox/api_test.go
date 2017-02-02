package main

import (
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/convox/rack/client"
	"github.com/convox/rack/test"
	"github.com/stretchr/testify/assert"
)

func TestApiGet(t *testing.T) {
	ts := testServer(t,
		test.Http{Method: "GET", Path: "/foo", Code: 200, Response: "bar"},
	)
	defer ts.Close()

	test.Runs(t,
		test.ExecRun{
			Command: "convox api get /foo",
			Exit:    0,
			Stdout:  "\"bar\"\n",
		},
	)
}

func TestApiGetApps(t *testing.T) {
	ts := testServer(t,
		test.Http{Method: "GET", Path: "/apps", Code: 200, Response: client.Apps{
			client.App{Name: "sinatra", Status: "running"},
		}},
	)

	defer ts.Close()

	test.Runs(t,
		test.ExecRun{
			Command: "convox api get /apps",
			Exit:    0,
			Stdout:  "[\n  {\n    \"name\": \"sinatra\",\n    \"release\": \"\",\n    \"status\": \"running\"\n  }\n]\n",
		},
	)
}

func TestApiGet404(t *testing.T) {
	c := &http.Client{}
	req, err := http.NewRequest("GET", "https://console.convox.com/nonexistent", nil)
	assert.NoError(t, err)

	req.Header.Add("Accept", "application/json")

	resp, err := c.Do(req)
	assert.NoError(t, err)
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)

	body := string(b)

	ts := testServer(t,
		test.Http{
			Method:   "GET",
			Path:     "/nonexistent",
			Code:     404,
			Response: client.Error{Error: body},
		},
	)
	defer ts.Close()

	test.Runs(t,
		test.ExecRun{
			Command: "convox api get /nonexistent",
			Exit:    1,
			Stderr:  "ERROR: \"Not found\"",
		},
	)
}
