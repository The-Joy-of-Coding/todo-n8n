package cache

import "todo-n8n/module/types"

type Storage types.Todos

var cache Storage

func (s *Storage) writeLog(req string) {}

func (s *Storage) Get() {}

func (s *Storage) Set() {}

func (s *Storage) pending() {}
