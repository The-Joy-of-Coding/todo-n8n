package cache

import (
	"bytes"
	"encoding/gob"
	"os"
	"testing"

	"todo-n8n/module/types"
)

func TestGetPath(t *testing.T) {
	path := getPath()
	if path == "" {
		t.Errorf("The path is empty!")
	}
}

func TestRead(t *testing.T) {
	oldCache := cache
	defer func() { cache = oldCache }()
	cache = []Storage{}
	testCache := []Storage{
		{
			Request: "test_request",
			Todo: types.Todos{
				Id:   1,
				Task: "This is test request",
			},
		},
	}
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(testCache); err != nil {
		t.Errorf("Error occured while encoding test cache, %s", err.Error())
	}
	read(buf.Bytes())
	if len(cache) != 1 {
		t.Errorf("Cache length is not same, got %v", len(cache))
	}
	if cache[0].Request != "test_request" {
		t.Errorf("Cache title is not same, %s", cache[0].Request)
	}
}

func TestSave(t *testing.T) {
	oldFile := file
	file = "test_cache_save.gob"
	defer func() {
		os.Remove(getPath())
		file = oldFile
	}()
	oldCache := cache
	defer func() { cache = oldCache }()
	cache = []Storage{}
	s := Storage{
		Request: "test_request",
		Todo: types.Todos{
			Id:   1,
			Task: "This is test request",
		},
	}
	s.Save()
	if len(cache) != 1 {
		t.Errorf("Cache length is not same, got %v", len(cache))
	}
	if _, err := os.Stat(getPath()); os.IsNotExist(err) {
		t.Errorf("Cache title is not same, %s", err.Error())
	}
}

func TestGet(t *testing.T) {
	oldCache := cache
	defer func() { cache = oldCache }()
	cache = []Storage{
		{Request: "Test 1"},
		{Request: "Test 2"},
	}
	got := Get()
	if len(got) != 2 {
		t.Errorf("Cache length is not same, got %v", len(got))
	}
}

func TestPending(t *testing.T) {
	oldFile := file
	file = "test_cache_pending.gob"
	defer func() {
		os.Remove(getPath())
		file = oldFile
	}()
	oldCache := cache
	defer func() { cache = oldCache }()
	cache = []Storage{
		{Request: "Test 1"},
		{Request: "Test 2"},
	}
	s := &Storage{}
	s.Pending()
	if len(cache) != 1 {
		t.Errorf("Cache length is not same, got %v", len(cache))
	}
	if cache[0].Request != "Test 2" {
		t.Errorf("Expected Test 2 but got %s", cache[0].Request)
	}
}
