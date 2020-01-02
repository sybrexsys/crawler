package manager

import (
	"testing"
)

func TestGetURL(t *testing.T) {
	wrk := &siteManager{
		urls: []string{"test1", "test2"},
	}

	url := wrk.getURL()
	if url == nil {
		t.Fail()
	}
	if *url != "test1" {
		t.Fail()
	}
	if len(wrk.urls) != 1 {
		t.Fail()
	}
	url = wrk.getURL()
	if url == nil {
		t.Fail()
	}
	if *url != "test2" {
		t.Fail()
	}
	if len(wrk.urls) != 0 {
		t.Fail()
	}
	url = wrk.getURL()
	if url != nil {
		t.Fail()
	}
}
