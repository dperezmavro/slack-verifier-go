package util

import "testing"

func TestBuildMessage(t *testing.T) {
	tests := []struct {
		name      string
		want      string
		body      string
		timestamp string
		key       string
	}{
		{
			name:      "slack example",
			body:      "token=xyzz0WbapA4vBCDEFasx0q6G&team_id=T1DC2JH3J&team_domain=testteamnow&channel_id=G8PSS9T3V&channel_name=foobar&user_id=U2CERLKJA&user_name=roadrunner&command=%2Fwebhook-collect&text=&response_url=https%3A%2F%2Fhooks.slack.com%2Fcommands%2FT1DC2JH3J%2F397700885554%2F96rGlfmibIGlgcZRskXaIFfN&trigger_id=398738663015.47445629121.803a0bc887a14d10d2c447fce8b6703c",
			timestamp: "1531420618",
			key:       "8f742231b10e8888abcd99yyyzzz85a5",
			want:      "v0=a2114d57b48eac39b9ad189dd8316235a7b4a8d21a10bd27519666489c69b503",
		},
		{
			name:      "my test",
			body:      "user_name=roadrunner&text=hello-there",
			timestamp: "1531420618",
			key:       "8f742231b10e8888abcd99yyyzzz85a5",
			want:      "v0=cb29ed9065a48430d698a6c990c6cd994029bbe9467acc569c816a8a5eee80d8",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := BuildMessageForAuthString(tt.body, tt.timestamp, []byte(tt.key))
			if res != tt.want {
				t.Errorf("wanted %s, got %s", tt.want, res)
			}
		})
	}
}
