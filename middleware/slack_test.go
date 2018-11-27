package middleware

import (
	"net/http"
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
