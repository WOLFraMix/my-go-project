package main

import (
	"reflect"
	"testing"
)

func TestCompressNumbers(t *testing.T) {
	tests := []struct {
		name string
		nums []int
		want []string
	}{
		{"empty", []int{}, []string{}}, // явно пустой слайс
		{"single", []int{10}, []string{"10"}},
		{"range and single", []int{1, 2, 3, 5, 6, 8}, []string{"1-3", "5-6", "8"}},
		{"mixed", []int{1, 3, 4, 5, 7}, []string{"1", "3-5", "7"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CompressNumbers(tt.nums)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CompressNumbers(%v) = %v; want %v", tt.nums, got, tt.want)
			}
		})
	}
}
