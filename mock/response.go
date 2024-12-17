package mock

import (
	"net/http"

	"github.com/dxta-dev/app/mock/har"
)

func CreateSimpleMockHandler(entry *har.Entry) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		har.EntryToResponse(w, entry)
	}
}
