// Package store 按组件索引查询实体。
package store

import (
	"sort"
	"strings"
)

// Store 是实体与组件的存储。
type Store struct {
	entities map[string]map[string]bool
	index    map[string]map[string]bool
	disabled map[string]bool
	cache    map[string][]string
	Work     int
}

func New() *Store {
	return &Store{
		entities: map[string]map[string]bool{},
		index:    map[string]map[string]bool{},
		disabled: map[string]bool{},
		cache:    map[string][]string{},
	}
}

// Set 给实体挂上一个组件。重复挂同一个组件是幂等 no-op，
// 且不影响实体的禁用状态。
func (s *Store) Set(entity, comp string) {
	if s.entities[entity] == nil {
		s.entities[entity] = map[string]bool{}
	}
	if s.entities[entity][comp] {
		return
	}
	s.entities[entity][comp] = true
	if s.index[comp] == nil {
		s.index[comp] = map[string]bool{}
	}
	s.index[comp][entity] = true
	s.cache = map[string][]string{}
}

// Disable 把实体排除在查询之外，只能由显式启用解除。
func (s *Store) Disable(entity string) {
	if s.disabled[entity] {
		return
	}
	s.disabled[entity] = true
	s.cache = map[string][]string{}
}

// Query 返回同时拥有全部给定组件的实体，按 ID 升序。
// 由组件倒排表驱动，从最小的集合开始筛；组件列表为空时返回全部实体。
func (s *Store) Query(comps ...string) []string {
	key := queryKey(comps)
	if hit, ok := s.cache[key]; ok {
		return append([]string{}, hit...)
	}

	var candidates map[string]bool
	if len(comps) == 0 {
		candidates = make(map[string]bool, len(s.entities))
		for id := range s.entities {
			candidates[id] = true
		}
	} else {
		smallest := comps[0]
		for _, c := range comps[1:] {
			if len(s.index[c]) < len(s.index[smallest]) {
				smallest = c
			}
		}
		candidates = s.index[smallest]
	}

	out := []string{}
	for id := range candidates {
		s.Work++
		if s.disabled[id] {
			continue
		}
		ok := true
		for _, c := range comps {
			if !s.entities[id][c] {
				ok = false
				break
			}
		}
		if ok {
			out = append(out, id)
		}
	}
	sort.Strings(out)

	s.cache[key] = out
	return append([]string{}, out...)
}

func queryKey(comps []string) string {
	sorted := append([]string(nil), comps...)
	sort.Strings(sorted)
	return strings.Join(sorted, "\x00")
}

// SortResults 是给调用方用的确定性排序助手。
func SortResults(ids []string) []string {
	out := append([]string(nil), ids...)
	sort.Strings(out)
	return out
}
