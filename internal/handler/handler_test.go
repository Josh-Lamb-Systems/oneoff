package handler

import (
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRespondSuccessEncodesNilSliceAsEmptyArray(t *testing.T) {
	h := &Handler{}
	var items []string

	rec := httptest.NewRecorder()
	h.respondSuccess(rec, 200, items)
	if got := strings.TrimSpace(rec.Body.String()); got != `{"data":[]}` {
		t.Errorf("respondSuccess = %s, want {\"data\":[]}", got)
	}

	rec = httptest.NewRecorder()
	h.respondList(rec, items, 0, 50, 0)
	if !strings.Contains(rec.Body.String(), `"data":[]`) {
		t.Errorf("respondList = %s, want data to be []", rec.Body.String())
	}

	rec = httptest.NewRecorder()
	respondListChains(rec, items, 0, 50, 0)
	if !strings.Contains(rec.Body.String(), `"data":[]`) {
		t.Errorf("respondListChains = %s, want data to be []", rec.Body.String())
	}
}

func TestEmptyIfNilSliceLeavesOtherValuesAlone(t *testing.T) {
	if got := emptyIfNilSlice(nil); got != nil {
		t.Errorf("emptyIfNilSlice(nil) = %v, want nil", got)
	}
	m := map[string]string{"a": "b"}
	if got := emptyIfNilSlice(m).(map[string]string); got["a"] != "b" {
		t.Errorf("map was altered: %v", got)
	}
}
