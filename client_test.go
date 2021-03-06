package gsclient

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestClient_waitForRequestCompleted(t *testing.T) {
	requestTestCases := []string{"done", "failed", "pending"}
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	config := NewConfiguration(
		server.URL,
		"uuid",
		"token",
		true,
		true,
		1,
		100,
		5,
	)
	client := NewClient(config)
	defer server.Close()
	var isFailed bool
	var reqStatus string
	mux.HandleFunc(requestBase, func(w http.ResponseWriter, r *http.Request) {
		if isFailed {
			w.WriteHeader(400)
		} else {
			fmt.Fprint(w, fmt.Sprintf(`{"%s": {"status":"%s", "isFailed" : %v}}`, dummyUUID, reqStatus, isFailed))
		}
	})
	for _, reqStatusTest := range requestTestCases {
		reqStatus = reqStatusTest
		for _, serverTest := range commonSuccessFailTestCases {
			isFailed = serverTest.isFailed
			for _, testUUID := range uuidCommonTestCases {
				err := client.waitForRequestCompleted(emptyCtx, testUUID.testUUID)
				if isFailed || testUUID.isFailed || reqStatus != requestDoneStatus {
					assert.NotNil(t, err)
				} else {
					assert.Nil(t, err, "waitForRequestCompleted returned an error %v", err)
				}
			}
		}
	}
}
