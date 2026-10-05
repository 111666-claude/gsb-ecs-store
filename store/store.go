// Package store 按组件索引查询实体。
package store

import "sort"

// Store 是实体与组件的存储。
type Store struct {
	entities map[string]map[string]bool
	index    map[string]map[string]bool
	order    []string
	disabled map[string]bool
	Work     int
}

func New() *Store {
	return &Store{
		entities: map[string]map[string]bool{},
		index:    map[string]map[string]bool{},
		disabled: map[string]bool{},
	}
}

// Set 给实体挂上一个组件。
func (s *Store) Set(entity, comp string) {
	delete(s.disabled, entity)
	if s.entities[entity] == nil {
		s.entities[entity] = map[string]bool{}
		s.order = append(s.order, entity)
	}
	s.entities[entity][comp] = true
}

// Disable 把实体排除在查询之外。
func (s *Store) Disable(entity string) {
	s.disabled[entity] = true
}

// Query 返回同时拥有全部给定组件的实体，按 ID 升序。
func (s *Store) Query(comps ...string) []string {
	var out []string
	for _, id := range s.order {
		s.Work++
		have := s.entities[id]
		ok := true
		for _, c := range comps {
			if !have[c] {
				ok = false
				break
			}
		}
		if ok {
			out = append(out, id)
		}
	}
	return out
}

// SortResults 是给调用方用的确定性排序助手。
func SortResults(ids []string) []string {
	out := append([]string(nil), ids...)
	sort.Strings(out)
	return out
}
