package main

import (
	"fmt"
	"testing"
)

func Test_generateAllGaps(t *testing.T) {
	tests := []struct {
		ngaps  int
		nextra int
		wantN  int
	}{
		{1, 1, 1},
		{2, 1, 2},
		{2, 3, 4},
		{2, 7, 8},
		{5, 0, 1},
		{5, 1, 5},
		{5, 2, 15},
		{5, 3, 35},
		{4, 4, 35},
	}
	for _, tt := range tests {
		name := fmt.Sprintf("g%d-x%d", tt.ngaps, tt.nextra)
		t.Run(name, func(t *testing.T) {
			gaps := make([]int, tt.ngaps+2)
			for i := 1; i < len(gaps)-1; i++ {
				gaps[i] = 1
			}
			got := generateAllGaps(gaps, tt.nextra)
			// fmt.Println(got)
			if len(got) != tt.wantN {
				t.Errorf("generateAllGaps() = %v items, want %v", len(got), tt.wantN)
				t.Error(got)
			}
		})
	}
}

func TestA(t *testing.T) {
	r := NewRow(".??..??...?##. 1,1,3")
	fmt.Println(r)
	all := r.arrangements()
	fmt.Println(len(all))
	for _, a := range all {
		fmt.Println(bstr(a, r.length))
	}
	t.Fail()
}

func Test_bstr(t *testing.T) {
	tests := []struct {
		name string
		i    int
		l    int
		want string
	}{
		{"0", 0, 1, "."},
		{"1", 1, 1, "#"},
		{"C", 12, 4, "##.."},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := bstr(tt.i, tt.l); got != tt.want {
				t.Errorf("bstr() = %v, want %v", got, tt.want)
			}
		})
	}
}
