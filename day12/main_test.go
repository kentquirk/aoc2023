package main

import (
	"fmt"
	"reflect"
	"testing"
)

func Test_groups(t *testing.T) {
	tests := []struct {
		s    string
		want []int
	}{
		{"#", []int{1}},
		{"##", []int{2}},
		{"#.#", []int{1, 1}},
		{".###.##.#.", []int{3, 2, 1}},
		{"###.##.#.", []int{3, 2, 1}},
		{".###.##.#", []int{3, 2, 1}},
	}
	for _, tt := range tests {
		t.Run(tt.s, func(t *testing.T) {
			if got := groups([]byte(tt.s)); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("groups() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_all(t *testing.T) {
	tests := []struct {
		s     string
		count int
		want  int
	}{
		{"# 1", 1, 1},
		{"## 2", 1, 1},
		{"#.# 1,1", 1, 1},
		{".###.##.#. 3,2,1", 1, 1},
		{"###.##.#. 3,2,1", 1, 1},
		{".###.##.# 3,2,1", 1, 1},
		{"???.### 1,1,3", 1, 1},
		{".??..??...?##. 1,1,3", 1, 4},
		{"?#?#?#?#?#?#?#? 1,3,1,6", 1, 1},
		{"????.#...#... 4,1,1", 1, 1},
		{"????.######..#####. 1,6,5", 1, 4},
		{"?###???????? 3,2,1", 1, 10},
		{"???.### 1,1,3", 5, 1},
		{".??..??...?##. 1,1,3", 5, 16384},
		{"?#?#?#?#?#?#?#? 1,3,1,6", 5, 1},
		{"????.#...#... 4,1,1", 5, 16},
		{"????.######..#####. 1,6,5", 5, 2500},
		{"?###???????? 3,2,1", 5, 506250},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%s/%d", tt.s, tt.count), func(t *testing.T) {
			r := NewRow(tt.s, tt.count)
			// fmt.Println(r)
			if got := r.caa(r.groups, 0); got != tt.want {
				t.Errorf("n = %v, want %v", got, tt.want)
			}
		})
	}
}
