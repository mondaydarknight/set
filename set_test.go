package set

import (
	"strings"
	"testing"
)

func TestAdd(t *testing.T) {
	set := New(1, 3, 5)
	tests := []struct {
		element interface{}
		result  bool
	}{
		{2, true},
		{4, true},
		{1, false},
		{3, false},
		{5, false},
	}
	for _, tt := range tests {
		if ok := set.Add(tt.element); ok != tt.result {
			t.Fatalf("set.Add(%d) - unexpected result: %t", tt.element, ok)
		}
	}
}

func TestClone(t *testing.T) {
	tests := []struct {
		source *Set
		target *Set
		result bool
	}{
		{New(1, 3, 5), New(1, 3, 5), true},
		{New(1, 2, 3), New(1, 2, 3, 4), false},
	}
	for _, tt := range tests {
		if c := tt.source.Clone(); c.Equal(tt.target) != tt.result {
			t.Fatalf("set.Clone() - unexpected result: %t", tt.result)
		}
	}
}

func TestDifference(t *testing.T) {
	set := New(1, 3, 5)
	tests := []struct {
		source *Set
		target *Set
	}{
		{New(1, 3, 5), New()},
		{New(2, 4, 6), New(1, 3, 5)},
		{New(1, 2, 6), New(3, 5)},
	}
	for _, tt := range tests {
		if diff := set.Difference(tt.source); !diff.Equal(tt.target) {
			t.Fatalf("set.Difference() - sets not equal: %s != %s", diff.String(), tt.target.String())
		}
	}
}

func TestEqual(t *testing.T) {
	tests := []struct {
		source *Set
		target *Set
		result bool
	}{
		{New(1, 3, 5), New(1, 3, 5), true},
		{New(1, 3, 5), New(1, 2, 3, 4, 5), false},
	}
	for _, tt := range tests {
		if ok := tt.source.Equal(tt.target); ok != tt.result {
			t.Fatalf("set.Equal() - unexpected result: %t", tt.result)
		}
	}
}

func TestFilter(t *testing.T) {
	tests := []struct {
		source *Set
		target *Set
		fn     func(interface{}) bool
	}{
		{
			New(1, 2, 3, 4, 5, 6),
			New(2, 4, 6),
			func(value interface{}) bool {
				return value.(int)%2 == 0
			},
		},
		{
			New(0, 1, false, true, "0", "1"),
			New(false, true),
			func(value interface{}) bool {
				_, ok := value.(bool)
				return ok
			},
		},
	}
	for _, tt := range tests {
		if set := tt.source.Filter(tt.fn); !set.Equal(tt.target) {
			t.Fatalf("set.Filter() - sets not equal: %s != %s", set.String(), tt.target.String())
		}
	}
}

func TestHas(t *testing.T) {
	set := New(1, 3, 5)
	tests := []struct {
		element interface{}
		result  bool
	}{
		{1, true},
		{2, false},
		{3, true},
		{4, false},
		{5, true},
	}
	for _, tt := range tests {
		if ok := set.Has(tt.element); ok != tt.result {
			t.Fatalf("set.Has(%d) - unexpected result: %t", tt.element, ok)
		}
	}
}

func TestIntersection(t *testing.T) {
	set := New(1, 3, 5)
	tests := []struct {
		source *Set
		target *Set
	}{
		{New(1, 3, 5), New(1, 3, 5)},
		{New(2, 4, 6), New()},
		{New(1, 2, 6), New(1)},
	}
	for _, tt := range tests {
		if set := set.Intersect(tt.source); !set.Equal(tt.target) {
			t.Fatalf("set.Intersect() - sets not equal: %s != %s", set.String(), tt.target.String())
		}
	}
}

func TestMap(t *testing.T) {
	tests := []struct {
		source *Set
		target *Set
		fn     func(interface{}) interface{}
	}{
		{
			New(1, 3, 5),
			New(1, 9, 25),
			func(value interface{}) interface{} {
				return value.(int) * value.(int)
			},
		},
		{
			New("hello", "world"),
			New("hell", "worl"),
			func(value interface{}) interface{} {
				return value.(string)[:len(value.(string))-1]
			},
		},
	}
	for _, tt := range tests {
		if set := tt.source.Map(tt.fn); !set.Equal(tt.target) {
			t.Fatalf("set.Map() - sets not equal: %s != %s", set.String(), tt.target.String())
		}
	}
}

func TestReduce(t *testing.T) {
	tests := []struct {
		source *Set
		value  interface{}
		result interface{}
		fn     func(interface{}, interface{}) interface{}
	}{
		{
			New(1, 3, 5),
			0,
			9,
			func(prev interface{}, value interface{}) interface{} {
				return prev.(int) + value.(int)
			},
		},
		{
			New("abc", "def", "ghi"),
			"",
			"ABCDEFGHI",
			func(prev interface{}, value interface{}) interface{} {
				return prev.(string) + strings.ToUpper(value.(string))
			},
		},
	}
	for _, tt := range tests {
		if out := tt.source.Reduce(tt.value, tt.fn); out != tt.result {
			t.Fatalf("set.Reduce() - unexpected result: %v", out)
		}
	}
}

func TestRemove(t *testing.T) {
	set := New(1, 3, 5)
	tests := []struct {
		element interface{}
		result  bool
	}{
		{1, true},
		{2, false},
		{3, true},
		{4, false},
		{5, true},
	}
	for _, tt := range tests {
		if ok := set.Remove(tt.element); ok != tt.result {
			t.Fatalf("set.Remove(%d) - unexpected result: %t", tt.element, ok)
		}
	}
}

func TestUnion(t *testing.T) {
	set := New(1, 3, 5)
	tests := []struct {
		source *Set
		target *Set
	}{
		{New(1, 3, 5), New(1, 3, 5)},
		{New(2, 4, 6), New(1, 2, 3, 4, 5, 6)},
		{New(1, 2, 6), New(1, 2, 3, 5, 6)},
	}
	for _, tt := range tests {
		if set := set.Union(tt.source); !set.Equal(tt.target) {
			t.Fatalf("set.Union() - sets not equal: %s != %s", set.String(), tt.target.String())
		}
	}
}
