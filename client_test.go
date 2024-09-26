package cambridgeapi

import (
	"encoding/json"
	"testing"
)

func TestRebuildDatabase(t *testing.T) {
	c := New()

	c.RebuildDatabase("./cambridge.db")
}

func TestQuery(t *testing.T) {
	c := New()
	w, err := c.Query("bridesmaid")
	if err != nil {
		t.Error(err)
	}

	data, _ := json.Marshal(w)
	c.log.Info().RawJSON("word", data).Msg("result")

	_, err = c.Query("crashx")
	if err == nil {
		t.Error("invalid word expected error")
	}
}
