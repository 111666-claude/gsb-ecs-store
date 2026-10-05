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

// Set 给实体挂上一个组件。重复挂同一组件是幂等的，且不会解除禁用。
func (s *Store) Set(entity, comp string) {
	comps := s.entities[entity]
	if comps == nil {
		comps = map[string]bool{}
		s.entities[entity] = comps
	}
	if comps[comp] {
		return
	}
	comps[comp] = true
	bucket := s.index[comp]
	if bucket == nil {
		bucket = map[string]bool{}
		s.index[comp] = bucket
	}
	bucket[entity] = true
	s.invalidate()
}

// Disable 把实体排除在查询之外，直到显式 Enable。
func (s *Store) Disable(entity string) {
	if s.disabled[entity] {
		return
	}
	s.disabled[entity] = true
	s.invalidate()
}

// Enable 解除实体的禁用状态。
func (s *Store) Enable(entity string) {
	if !s.disabled[entity] {
		return
	}
	delete(s.disabled, entity)
	s.invalidate()
}

func (s *Store) invalidate() {
	s.cache = map[string][]string{}
}

// Query 返回同时拥有全部给定组件的实体，按 ID 升序。
// 组件列表为空时返回全部实体。
func (s *Store) Query(comps ...string) []string {
	key := cacheKey(comps)
	if hit, ok := s.cache[key]; ok {
		return append([]string{}, hit...)
	}
	out := s.runQuery(comps)
	s.cache[key] = out
	return append([]string{}, out...)
}

// runQuery 由组件倒排表驱动，从最小的那个集合开始筛。
func (s *Store) runQuery(comps []string) []string {
	out := []string{}
	if len(comps) == 0 {
		s.Work += len(s.entities)
		for id := range s.entities {
			if !s.disabled[id] {
				out = append(out, id)
			}
		}
		sort.Strings(out)
		return out
	}
	smallest := comps[0]
	for _, c := range comps[1:] {
		if len(s.index[c]) < len(s.index[smallest]) {
			smallest = c
		}
	}
	s.Work += len(s.index[smallest])
	for id := range s.index[smallest] {
		if s.disabled[id] {
			continue
		}
		ok := true
		for _, c := range comps {
			if c != smallest && !s.index[c][id] {
				ok = false
				break
			}
		}
		if ok {
			out = append(out, id)
		}
	}
	sort.Strings(out)
	return out
}

// cacheKey 归一化组件列表，保证同一组组件命中同一条缓存。
func cacheKey(comps []string) string {
	seen := map[string]bool{}
	norm := make([]string, 0, len(comps))
	for _, c := range comps {
		if !seen[c] {
			seen[c] = true
			norm = append(norm, c)
		}
	}
	sort.Strings(norm)
	return strings.Join(norm, "\x1f")
}

// SortResults 是给调用方用的确定性排序助手。
func SortResults(ids []string) []string {
	out := append([]string(nil), ids...)
	sort.Strings(out)
	return out
}
