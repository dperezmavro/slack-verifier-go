package middlware

import (
	"crypto/hmac"
	"fmt"
	"io"
	"net/http"

	"github.com/dperezmavro/slack-verifier-go/util"
)

type SlackVerifier func(http.Handler) http.Handler

const (
	SLACK_AUTH_HEADER      string = "X-Slack-Signature"
	SLACK_TIMESTAMP_HEADER string = "X-Slack-Request-Timestamp"
	SLACK_AUTH_ENV         string = "SLACKAUTHENV"
)

func checkHeader(h string, r *http.Request) error {
	if len(r.Header[h]) != 1 || r.Header[h][0] == "" {
		return fmt.Errorf("missing %s header", h)
	}
	return nil
}

func Verify(slackSecret []byte) SlackVerifier {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				if err := checkHeader(SLACK_AUTH_HEADER, r); err != nil {
					http.Error(w, err.Error(), http.StatusBadRequest)
					return
				}

				if err := checkHeader(SLACK_TIMESTAMP_HEADER, r); err != nil {
					http.Error(w, err.Error(), http.StatusBadRequest)
					return
				}

				var body []byte
				_, err := io.ReadFull(r.Body, body)
				if err != nil {
					http.Error(w, fmt.Sprintf("failed to ready body: %+v", err), http.StatusInternalServerError)
					return
				}

				slackSignature := r.Header[SLACK_AUTH_HEADER][0]
				slackTimeStamp := r.Header[SLACK_TIMESTAMP_HEADER][0]

				messageDigest := util.BuildMessageForAuth(string(body), slackTimeStamp, slackSecret)
				if !hmac.Equal(messageDigest, []byte(slackSignature)) {
					http.Error(w, fmt.Sprintf(""), http.StatusUnauthorized)
					return
				}
				h.ServeHTTP(w, r)
			},
		)
	}
}
