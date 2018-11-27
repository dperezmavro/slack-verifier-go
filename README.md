# About

This middleware can be used in in an HTTP handler in order to verify that a call to the application is originalting from slack (perhaps as part of a Slack-application's callback).

# Usage

Top use the slack verifier, you need to have a key for all the crypto-operations, and initialise the middleware with that key.

You can use it directly, like this:

```http.Handle("/path", middleware.VerifySlackMessage(slackKey)(http.HandlerFunc(healthCheckHandlerSimple)))```


Or as part of a collection of other middleware, like this:
```
http.Handle("/path2",
    addMiddleware(
        healthCheckHandler(),
        middleware.VerifySlackMessage(slackKey),
    ),
)
```

`main.go`  has some examples on how to use it.