package store

import "testing"

func TestQueryReturnsMatchingEntity(t *testing.T) {
	s := New()
	s.Set("e1", "pos")
	s.Set("e1", "hp")
	s.Set("e2", "pos")
	got := s.Query("pos", "hp")
	if len(got) != 1 || got[0] != "e1" {
		t.Fatalf("应只剩 e1：%v", got)
	}
}

func TestQueryWithNoMatchIsEmpty(t *testing.T) {
	s := New()
	s.Set("e1", "pos")
	if got := s.Query("hp"); len(got) != 0 {
		t.Fatalf("没有匹配应为空：%v", got)
	}
}

func TestSingleComponentQuery(t *testing.T) {
	s := New()
	s.Set("e1", "pos")
	s.Set("e2", "pos")
	if got := s.Query("pos"); len(got) != 2 {
		t.Fatalf("两个实体都有 pos：%v", got)
	}
}

func TestSortResultsIsAscending(t *testing.T) {
	got := SortResults([]string{"b", "a"})
	if got[0] != "a" || got[1] != "b" {
		t.Fatalf("应升序：%v", got)
	}
}
