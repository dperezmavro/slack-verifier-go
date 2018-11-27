package middleware

import (
	"crypto/hmac"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/dperezmavro/slack-verifier-go/util"
)

type SlackVerifier func(http.Handler) http.Handler

const (
	slackAuthHeader      string = "X-Slack-Signature"
	slackTimestampHeader string = "X-Slack-Request-Timestamp"
)

func checkHeader(h string, r *http.Request) error {
	if len(r.Header[h]) != 1 || r.Header[h][0] == "" {
		return fmt.Errorf("missing %s header", h)
	}
	return nil
}

// VerifySlackMessage is an http middleware that
// will perform message verification on the incoming
// request to verify it is coming from slack
func VerifySlackMessage(slackSecret []byte) SlackVerifier {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				if err := checkHeader(slackAuthHeader, r); err != nil {
					http.Error(w, err.Error(), http.StatusBadRequest)
					return
				}

				if err := checkHeader(slackTimestampHeader, r); err != nil {
					http.Error(w, err.Error(), http.StatusBadRequest)
					return
				}

				if r.Method != http.MethodPost {
					http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
				}

				var body []byte
				body, err := ioutil.ReadAll(r.Body)
				if err != nil {
					http.Error(w, fmt.Sprintf("failed to ready body: %+v", err), http.StatusInternalServerError)
					return
				}

				slackSignature := r.Header[slackAuthHeader][0]
				slackTimeStamp := r.Header[slackTimestampHeader][0]

				messageDigest := util.BuildMessageForAuth(string(body), slackTimeStamp, slackSecret)
				if !hmac.Equal(messageDigest, []byte(slackSignature)) {
					http.Error(w, fmt.Sprintf("signature check failed"), http.StatusUnauthorized)
					return
				}
				h.ServeHTTP(w, r)
			},
		)
	}
}
