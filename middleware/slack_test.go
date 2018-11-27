package middleware

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCheckHeader(t *testing.T) {
	tests := []struct {
		name string
		fail bool
		r    *http.Request
		h    string
	}{
		{
			name: "empty",
			fail: true,
			r:    &http.Request{},
		},
		{
			name: "missing header",
			r: &http.Request{
				Header: http.Header{
					"A": []string{"B"},
					"C": []string{"D"},
				},
			},
			fail: true,
			h:    "myHeader",
		},
		{
			name: "missing value",
			r: &http.Request{
				Header: http.Header{
					"A":        []string{"B"},
					"myHeader": []string{},
				},
			},
			fail: true,
			h:    "myHeader",
		},
		{
			name: "empty value",
			r: &http.Request{
				Header: http.Header{
					"A":        []string{"B"},
					"myHeader": []string{""},
				},
			},
			fail: true,
			h:    "myHeader",
		},
		{
			name: "multiple values",
			fail: true,
			r: &http.Request{
				Header: http.Header{
					"myHeader": []string{"no", "pass"},
				},
			},
			h: "myHeader",
		},
		{
			name: "pass",
			r: &http.Request{
				Header: http.Header{
					"myHeader": []string{"pass"},
				},
			},
			h: "myHeader",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := checkHeader(tt.h, tt.r)
			if tt.fail {
				if err == nil {
					t.Errorf("didn't fail")
				}
			} else {
				if err != nil {
					t.Errorf("failed %+v", err)
				}
			}
		})
	}
}

func TestVerifySlackMessage(t *testing.T) {
	tests := []struct {
		name      string
		fail      bool
		r         *http.Request
		key, body []byte
	}{
		{
			name: "empty",
			fail: true,
			r:    &http.Request{},
		},
		{
			name: "missing signature",
			fail: true,
			r: &http.Request{
				Method: "POST",
				Header: http.Header{
					slackTimestampHeader: []string{"1531420618"},
				},
				Body: ioutil.NopCloser(bytes.NewBuffer([]byte("dummy"))),
			},

			key: []byte("8f742231b10e8888abcd99yyyzzz85a5"),
		},
		{
			name: "missing timestam[",
			fail: true,
			r: &http.Request{
				Method: "POST",
				Header: http.Header{
					slackAuthHeader: []string{"v0=dummy"},
				},
				Body: ioutil.NopCloser(bytes.NewBuffer([]byte("dummy"))),
			},

			key: []byte("8f742231b10e8888abcd99yyyzzz85a5"),
		},
		{
			name: "bad method",
			fail: true,
			r: &http.Request{
				Method: "GET",
				Header: http.Header{
					slackAuthHeader:      []string{"v0=a2114d57b48eac39b9ad189dd8316235a7b4a8d21a10bd27519666489c69b503"},
					slackTimestampHeader: []string{"1531420618"},
				},
				Body: ioutil.NopCloser(bytes.NewBuffer([]byte("token=xyzz0WbapA4vBCDEFasx0q6G&team_id=T1DC2JH3J&team_domain=testteamnow&channel_id=G8PSS9T3V&channel_name=foobar&user_id=U2CERLKJA&user_name=roadrunner&command=%2Fwebhook-collect&text=&response_url=https%3A%2F%2Fhooks.slack.com%2Fcommands%2FT1DC2JH3J%2F397700885554%2F96rGlfmibIGlgcZRskXaIFfN&trigger_id=398738663015.47445629121.803a0bc887a14d10d2c447fce8b6703c"))),
			},

			key: []byte("8f742231b10e8888abcd99yyyzzz85a5"),
		},
		{
			name: "works",
			r: &http.Request{
				Method: "POST",
				Header: http.Header{
					slackAuthHeader:      []string{"v0=a2114d57b48eac39b9ad189dd8316235a7b4a8d21a10bd27519666489c69b503"},
					slackTimestampHeader: []string{"1531420618"},
				},
				Body: ioutil.NopCloser(bytes.NewBuffer([]byte("token=xyzz0WbapA4vBCDEFasx0q6G&team_id=T1DC2JH3J&team_domain=testteamnow&channel_id=G8PSS9T3V&channel_name=foobar&user_id=U2CERLKJA&user_name=roadrunner&command=%2Fwebhook-collect&text=&response_url=https%3A%2F%2Fhooks.slack.com%2Fcommands%2FT1DC2JH3J%2F397700885554%2F96rGlfmibIGlgcZRskXaIFfN&trigger_id=398738663015.47445629121.803a0bc887a14d10d2c447fce8b6703c"))),
			},

			key: []byte("8f742231b10e8888abcd99yyyzzz85a5"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			handler := VerifySlackMessage(tt.key)(http.HandlerFunc(HealthCheckHandler))
			handler.ServeHTTP(rr, tt.r)

			// Check the status code is what we expect.
			if tt.fail {
				if rr.Code == http.StatusOK {
					t.Errorf("did not fail")
				}
			} else {
				if rr.Code != http.StatusOK {
					t.Errorf("failed with response code %d and message: %s", rr.Code, rr.Body.String())
				}
			}
		})
	}
}
