package pick

import "testing"

func TestOneReturnsSelected(t *testing.T) {
	orig := find
	defer func() { find = orig }()
	find = func(n int, label func(i int) string, multi bool) ([]int, error) {
		return []int{1}, nil // pretend user picked index 1
	}
	got, ok, err := One([]string{"a", "b", "c"}, func(s string) string { return s })
	if err != nil || !ok || got != "b" {
		t.Fatalf("got=%q ok=%v err=%v", got, ok, err)
	}
}

func TestManyReturnsAllSelected(t *testing.T) {
	orig := find
	defer func() { find = orig }()
	find = func(n int, label func(i int) string, multi bool) ([]int, error) {
		return []int{0, 2}, nil
	}
	got, err := Many([]string{"a", "b", "c"}, func(s string) string { return s })
	if err != nil || len(got) != 2 || got[0] != "a" || got[1] != "c" {
		t.Fatalf("got=%v err=%v", got, err)
	}
}
