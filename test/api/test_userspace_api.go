package main

import (
  "net/http/httptest"
  "net/http"
  "testing"
  "devops.io/cloud/api"
)

type step struct {
  query, expect string
}

func TestConnectivity(t *testing.T) {
  queries := [1]step{
    step{`{
      ping
    }
    `,
    `{"code": 200, "data": "pong"}`,
    },
  }

  for _, query := range queries {
    srv := api.NewApiServer("test")
    w := httptest.NewRecorder()
    r := srv.GetMuxer()

    r.ServeHTTP(w, httptest.NewRequest("PUT", "/query", query))

    if w.Code != http.StatusOK {
      t.Error("Did not get expected HTTP status code, got", w.Code)
    }
  }
}
