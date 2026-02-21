package cache

import (
	"bytes"
	"encoding/gob"
	"io"
	"log/slog"
	"os"
	"path/filepath"

	"todo-n8n/module/types"
)

type Storage struct {
	Request string
	Todo    types.Todos
}

var (
	cache []Storage
	file  string = "cache.gob"
	dir   string = "todo-n8n"
)

func init() {
	tmp := getPath()
	data, err := os.ReadFile(tmp)
	if err != nil {
		if !os.IsNotExist(err) {
			slog.Error(err.Error())
		}
		return
	}
	read(data)
}

func getPath() string {
	baseDir, err := os.UserCacheDir()
	if err != nil {
		baseDir = os.TempDir()
	}
	tmp := filepath.Join(baseDir, dir, file)
	if err := os.MkdirAll(filepath.Dir(tmp), 0o755); err != nil {
		slog.Error(err.Error())
	}
	return tmp
}

func read(data []byte) {
	if len(data) == 0 {
		return
	}
	buf := bytes.NewBuffer(data)
	denc := gob.NewDecoder(buf)
	if err := denc.Decode(&cache); err != nil && err != io.EOF {
		slog.Error(err.Error())
	}
}

func (s *Storage) write() {
	var buff bytes.Buffer
	tmp := getPath()
	enc := gob.NewEncoder(&buff)
	if err := enc.Encode(cache); err != nil {
		slog.Error(err.Error())
		return
	}
	if err := os.WriteFile(tmp, buff.Bytes(), 0o644); err != nil {
		slog.Error(err.Error())
	}
}

func (s *Storage) Save() {
	cache = append(cache, *s)
	s.write()
}

func Get() []Storage {
	return cache
}

func (s *Storage) Pending() {
	if len(cache) > 0 {
		cache = cache[1:]
		s.write()
	}
}
