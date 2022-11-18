package utils

import (
	"reflect"
	"testing"

	"github.com/elliotchance/pie/pie"
)

func TestDiffSliceV2(t *testing.T) {
	type args struct {
		a interface{}
		b interface{}
	}
	tests := []struct {
		name  string
		args  args
		want  interface{}
		want1 interface{}
	}{
		{
			name: "diff slice type",
			args: args{
				a: []int{1, 2},
				b: pie.Ints{2, 3},
			},
			want:  []int{1},
			want1: pie.Ints{3},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := DiffSliceV2(tt.args.a, tt.args.b)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DiffSliceV2() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("DiffSliceV2() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestRemoveSlice(t *testing.T) {
	type args struct {
		src interface{}
		rm  interface{}
	}
	tests := []struct {
		name string
		args args
		want interface{}
	}{
		{
			name: "diff slice type",
			args: args{
				src: []int{1, 2},
				rm:  pie.Ints{2, 3},
			},
			want: []int{1},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RemoveSlice(tt.args.src, tt.args.rm); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RemoveSlice() = %v, want %v", got, tt.want)
			}
		})
	}
}
