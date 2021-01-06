package collection

import (
	"reflect"
	"testing"
)

func TestAdditions(t *testing.T) {
	type args struct {
		l []interface{}
		r []interface{}
	}
	tests := []struct {
		name          string
		args          args
		wantAdditions []interface{}
	}{
		{
			name: "Additions anywhere should be found",
			args: args{
				l: []interface{}{"a", "b"},
				r: []interface{}{"c", "a", "d", "b", "e"},
			},
			wantAdditions: []interface{}{"c", "d", "e"},
		},
		{
			name: "Removed items should not be found",
			args: args{
				l: []interface{}{"a", "b"},
				r: []interface{}{"c"},
			},
			wantAdditions: []interface{}{"c"},
		},
		{
			// Note: it is not the responsibility of Additions to remove duplicates from the right collection
			name: "Duplicate items should be found twice",
			args: args{
				l: []interface{}{"a", "b"},
				r: []interface{}{"c", "c", "d"},
			},
			wantAdditions: []interface{}{"c", "c", "d"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotAdditions := Additions(tt.args.l, tt.args.r); !reflect.DeepEqual(gotAdditions, tt.wantAdditions) {
				t.Errorf("Additions() = %v, want %v", gotAdditions, tt.wantAdditions)
			}
		})
	}
}

func TestContains(t *testing.T) {
	type args struct {
		needle   interface{}
		haystack []interface{}
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Finding the item should return true",
			args: args{
				needle:   "b",
				haystack: []interface{}{"a", "b"},
			},
			want: true,
		},
		{
			name: "Not finding the item should return false",
			args: args{
				needle:   "c",
				haystack: []interface{}{"a", "b"},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Contains(tt.args.needle, tt.args.haystack); got != tt.want {
				t.Errorf("Contains() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDiff(t *testing.T) {
	type args struct {
		l []interface{}
		r []interface{}
	}
	tests := []struct {
		name          string
		args          args
		wantAdditions []interface{}
		wantRemovals  []interface{}
	}{
		{
			name: "Removals and additions should be returned",
			args: args{
				l: []interface{}{"a", "b", "c"},
				r: []interface{}{"b", "d", "e", "f"},
			},
			wantAdditions: []interface{}{"d", "e", "f"},
			wantRemovals:  []interface{}{"a", "c"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotAdditions, gotRemovals := Diff(tt.args.l, tt.args.r)
			if !reflect.DeepEqual(gotAdditions, tt.wantAdditions) {
				t.Errorf("Diff() gotAdditions = %v, want %v", gotAdditions, tt.wantAdditions)
			}
			if !reflect.DeepEqual(gotRemovals, tt.wantRemovals) {
				t.Errorf("Diff() gotRemovals = %v, want %v", gotRemovals, tt.wantRemovals)
			}
		})
	}
}

func TestOrdered(t *testing.T) {
	type args struct {
		items []interface{}
		like  []interface{}
	}
	tests := []struct {
		name        string
		args        args
		wantOrdered []interface{}
	}{
		{
			name: "Ordering should be preserved",
			args: args{
				items: []interface{}{"a", "b", "c", "d"},
				like:  []interface{}{"c", "b"},
			},
			wantOrdered: []interface{}{"c", "b", "a", "d"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotOrdered := Ordered(tt.args.items, tt.args.like); !reflect.DeepEqual(gotOrdered, tt.wantOrdered) {
				t.Errorf("Ordered() = %v, want %v", gotOrdered, tt.wantOrdered)
			}
		})
	}
}

func TestToInterfaceSlice(t *testing.T) {
	type args struct {
		s []string
	}
	tests := []struct {
		name string
		args args
		want []interface{}
	}{
		{
			name: "Items should be retained",
			args: args{s: []string{"a", "b"}},
			want: []interface{}{"a", "b"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToInterfaceSlice(tt.args.s); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToInterfaceSlice() = %v, want %v", got, tt.want)
			}
		})
	}
}
