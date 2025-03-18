package har

import (
	"compress/gzip"
	"net/http"

	"golang.org/x/net/http/httpguts"
)

func EntryToResponse(w http.ResponseWriter, entry *Entry) {
	for _, h := range entry.Response.Headers {
		if httpguts.ValidHeaderFieldName(h.Name) && httpguts.ValidHeaderFieldValue(h.Value) {
			w.Header().Set(h.Name, h.Value)
		}
	}

	for _, c := range entry.Response.Cookies {
		cookie := &http.Cookie{Name: c.Name, Value: c.Value, HttpOnly: c.HTTPOnly, Domain: c.Domain}
		http.SetCookie(w, cookie)
	}

	body := entry.Response.Content.Text

	gz := gzip.NewWriter(w)

	if _, err := gz.Write([]byte(body)); err != nil {
		http.Error(w, "Failed to write response body", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(entry.Response.Status)
}
