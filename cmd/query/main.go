// Command query 先挂组件建索引，再跑一组固定查询。
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"example.com/ecs-store/store"
)

func build(entities int, compsPer int) *store.Store {
	s := store.New()
	for i := 0; i < entities; i++ {
		id := fmt.Sprintf("e%04d", i)
		for c := 0; c < compsPer; c++ {
			s.Set(id, fmt.Sprintf("c%d", (i+c)%8))
		}
	}
	return s
}

func scenario(name string) (*store.Store, [][]string) {
	switch name {
	case "order":
		s := store.New()
		s.Set("b", "pos")
		s.Set("a", "pos")
		return s, [][]string{{"pos"}}
	case "disabled":
		s := store.New()
		s.Set("a", "pos")
		s.Set("b", "pos")
		s.Disable("a")
		return s, [][]string{{"pos"}}
	case "revive":
		s := store.New()
		s.Set("a", "pos")
		s.Disable("a")
		s.Set("a", "pos")
		return s, [][]string{{"pos"}}
	case "work":
		s := build(2000, 3)
		queries := make([][]string, 500)
		for i := range queries {
			queries[i] = []string{"c7", "c6"}
		}
		return s, queries
	}
	fmt.Fprintf(os.Stderr, "未知场景 %q\n", name)
	os.Exit(2)
	return nil, nil
}

func main() {
	sample := flag.String("sample", "order", "场景名")
	flag.Parse()

	s, queries := scenario(*sample)
	hits := 0
	for _, q := range queries {
		got := s.Query(q...)
		hits += len(got)
		if *sample != "work" {
			out, _ := json.Marshal(map[string]any{"query": q, "hits": got})
			fmt.Println(string(out))
		}
	}
	if *sample == "work" {
		out, _ := json.Marshal(map[string]int{"hits": hits, "queries": len(queries), "work": s.Work})
		fmt.Println(string(out))
	}
}
