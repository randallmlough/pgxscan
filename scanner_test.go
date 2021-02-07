package pgxscan

import (
	"testing"
	"time"

	"github.com/randallmlough/pgxscan/testdata"
)

func Test_scanStruct(t *testing.T) {
	dst := &testdata.TestStruct{}
	cols := testdata.TestCols

	fn := func(i ...interface{}) error {
		return nil
	}
	if err := ScanStruct(fn, dst, cols); err != nil {
		t.Errorf("scanStruct() failed to scan. Reason:  %v", err)
		return
	}

	if err := ScanStruct(fn, dst, []string{"id"}); err != nil {
		t.Errorf("scanStruct() failed to scan. Reason:  %v", err)
		return
	}

	if err := ScanStruct(fn, dst, nil); err == nil {
		t.Errorf("scanStruct() failed to scan. Reason:  %v", err)
		return
	}
}

func Test_isVariadic(t *testing.T) {
	type args []interface{}
	tests := []struct {
		name string
		test args
		want bool
	}{
		{
			name: "single string",
			test: args{
				"string",
			},
			want: true,
		},
		{
			name: "couple strings",
			test: args{
				"string", "string",
			},
			want: true,
		},
		{
			name: "single int",
			test: args{
				1,
			},
			want: true,
		},
		{
			name: "builtin slice",
			test: args{
				[]string{"one", "two"},
			},
			want: true,
		},
		{
			name: "struct",
			test: args{
				struct {
					Str string
				}{"one"},
			},
			want: false,
		},
		{
			name: "struct slice",
			test: args{
				[]struct {
					Str string
				}{{"one"}},
			},
			want: false,
		},
		{
			name: "interface slice",
			test: args{
				[]interface{}{"one", "two", 4},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isVariadic(tt.test...); got != tt.want {
				t.Errorf("isBuiltin() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_isBuiltin(t *testing.T) {
	tests := []struct {
		name string
		test interface{}
		want bool
	}{
		// things that are builtins
		{
			name: "string",
			test: "string",
			want: true,
		},
		{
			name: "uint",
			test: uint(1),
			want: true,
		},
		{
			name: "uint8",
			test: uint8(1),
			want: true,
		},
		{
			name: "uint16",
			test: uint16(1),
			want: true,
		},
		{
			name: "uint32",
			test: uint32(1),
			want: true,
		},
		{
			name: "uint64",
			test: uint64(1),
			want: true,
		},
		{
			name: "int",
			test: int(1),
			want: true,
		},
		{
			name: "int8",
			test: int8(1),
			want: true,
		},
		{
			name: "int16",
			test: int16(1),
			want: true,
		},
		{
			name: "int32",
			test: int32(1),
			want: true,
		},
		{
			name: "int64",
			test: int64(1),
			want: true,
		},
		{
			name: "float32",
			test: float32(1.56),
			want: true,
		},
		{
			name: "float64",
			test: float64(1.56),
			want: true,
		},
		{
			name: "complex64",
			test: complex64(1.56),
			want: true,
		},
		{
			name: "complex128",
			test: complex128(1.56),
			want: true,
		},
		{
			name: "bool",
			test: false,
			want: true,
		},
		{
			name: "time",
			test: time.Now(),
			want: true,
		},
		// things that are not. Should return false
		{
			name: "interface slice",
			test: []interface{}{1},
			want: false,
		},
		{
			name: "struct",
			test: struct {
				Str string
			}{
				Str: "hello",
			},
			want: false,
		},
		{
			name: "struct slice",
			test: []struct {
				Str string
			}{
				{
					Str: "hello",
				},
				{
					Str: "world",
				},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isBuiltin(tt.test); got != tt.want {
				t.Errorf("isBuiltin() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_isBuiltinPointers(t *testing.T) {
	tests := []struct {
		name string
		test interface{}
		want bool
	}{
		// things that are builtins
		{
			name: "string",
			test: stringToPtr("string"),
			want: true,
		},
		{
			name: "uint",
			test: uintToPtr(1),
			want: true,
		},
		{
			name: "uint8",
			test: uint8ToPtr(1),
			want: true,
		},
		{
			name: "uint16",
			test: uint16ToPtr(1),
			want: true,
		},
		{
			name: "uint32",
			test: uint32ToPtr(1),
			want: true,
		},
		{
			name: "uint64",
			test: uint64ToPtr(1),
			want: true,
		},
		{
			name: "int",
			test: intToPtr(int(1)),
			want: true,
		},
		{
			name: "int8",
			test: int8ToPtr(int8(1)),
			want: true,
		},
		{
			name: "int16",
			test: int16ToPtr(int16(1)),
			want: true,
		},
		{
			name: "int32",
			test: int32ToPtr(int32(1)),
			want: true,
		},
		{
			name: "int64",
			test: int64ToPtr(int64(1)),
			want: true,
		},
		{
			name: "byteToPointer",
			test: byteToPtr([]byte(`abc`)),
			want: true,
		},
		{
			name: "float32",
			test: float32ToPtr(float32(1.56)),
			want: true,
		},
		{
			name: "float64",
			test: float64ToPtr(float64(1.56)),
			want: true,
		},
		{
			name: "complex64",
			test: complex64ToPtr(complex64(1.56)),
			want: true,
		},
		{
			name: "complex128",
			test: complex128ToPtr(complex128(1.56)),
			want: true,
		},
		{
			name: "bool",
			test: boolToPtr(false),
			want: true,
		},
		{
			name: "time",
			test: timeToPtr(time.Now()),
			want: true,
		},
		// things that are not. Should return false
		{
			name: "interface slice",
			test: &[]interface{}{1},
			want: false,
		},
		{
			name: "struct",
			test: &struct {
				Str string
			}{
				Str: "hello",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isBuiltin(tt.test); got != tt.want {
				t.Errorf("isBuiltin() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_isBuiltinSlice(t *testing.T) {
	tests := []struct {
		name string
		test interface{}
		want bool
	}{
		// things that are builtins
		{
			name: "string slice",
			test: []string{"string"},
			want: true,
		},
		{
			name: "uint slice",
			test: []uint{1},
			want: true,
		},
		{
			name: "uint8 slice",
			test: []uint8{1},
			want: true,
		},
		{
			name: "uint16 slice",
			test: []uint16{1},
			want: true,
		},
		{
			name: "uint32 slice",
			test: []uint32{1},
			want: true,
		},
		{
			name: "uint64 slice",
			test: []uint64{1},
			want: true,
		},
		{
			name: "int slice",
			test: []int{1},
			want: true,
		},
		{
			name: "int8 slice",
			test: []int8{1},
			want: true,
		},
		{
			name: "int16 slice",
			test: []int16{1},
			want: true,
		},
		{
			name: "int32 slice",
			test: []int32{1},
			want: true,
		},
		{
			name: "int64 slice",
			test: []int64{1},
			want: true,
		},
		{
			name: "float32 slice",
			test: []float32{1.56},
			want: true,
		},
		{
			name: "float64 slice",
			test: []float64{1.56},
			want: true,
		},
		{
			name: "complex64 slice",
			test: []complex64{1.56},
			want: true,
		},
		{
			name: "complex128 slice",
			test: []complex128{1.56},
			want: true,
		},
		{
			name: "bool slice",
			test: []bool{false},
			want: true,
		},
		{
			name: "time slice",
			test: []time.Time{time.Now()},
			want: true,
		},
		// things that are not. Should return false
		{
			name: "interface slice",
			test: []interface{}{1},
			want: false,
		},
		{
			name: "struct slice",
			test: []struct {
				Str string
			}{
				{
					Str: "hello",
				},
				{
					Str: "world",
				},
			},

			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isBuiltin(tt.test); got != tt.want {
				t.Errorf("isBuiltin() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_isBuiltinSlicePointer(t *testing.T) {
	tests := []struct {
		name string
		test interface{}
		want bool
	}{
		// things that are builtins
		{
			name: "string slice pointer",
			test: &[]string{"string"},
			want: true,
		},
		{
			name: "uint slice pointer",
			test: &[]uint{1},
			want: true,
		},
		{
			name: "uint8 slice pointer",
			test: &[]uint8{1},
			want: true,
		},
		{
			name: "uint16 slice pointer",
			test: &[]uint16{1},
			want: true,
		},
		{
			name: "uint32 slice pointer",
			test: &[]uint32{1},
			want: true,
		},
		{
			name: "uint64 slice pointer",
			test: &[]uint64{1},
			want: true,
		},
		{
			name: "int slice pointer",
			test: &[]int{1},
			want: true,
		},
		{
			name: "int8 slice pointer",
			test: &[]int8{1},
			want: true,
		},
		{
			name: "int16 slice pointer",
			test: &[]int16{1},
			want: true,
		},
		{
			name: "int32 slice pointer",
			test: &[]int32{1},
			want: true,
		},
		{
			name: "int64 slice pointer",
			test: &[]int64{1},
			want: true,
		},
		{
			name: "float32 slice pointer",
			test: &[]float32{1.56},
			want: true,
		},
		{
			name: "float64 slice pointer",
			test: &[]float64{1.56},
			want: true,
		},
		{
			name: "complex64 slice pointer",
			test: &[]complex64{1.56},
			want: true,
		},
		{
			name: "complex128 slice pointer",
			test: &[]complex128{1.56},
			want: true,
		},
		{
			name: "bool slice pointer",
			test: &[]bool{false},
			want: true,
		},
		{
			name: "time slice pointer",
			test: &[]time.Time{time.Now()},
			want: true,
		},
		// things that are not. Should return false
		{
			name: "interface slice pointer",
			test: &[]interface{}{1},
			want: false,
		},
		{
			name: "struct slice pointer",
			test: &[]struct {
				Str string
			}{
				{
					Str: "hello",
				},
				{
					Str: "world",
				},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isBuiltin(tt.test); got != tt.want {
				t.Errorf("isBuiltin() = %v, want %v", got, tt.want)
			}
		})
	}
}
