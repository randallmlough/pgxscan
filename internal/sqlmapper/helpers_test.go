package sqlmaper

import (
	"database/sql"
	"github.com/stretchr/testify/suite"
	"reflect"
	"testing"
)

var (
	uints = []interface{}{
		uint(10),
		uint8(10),
		uint16(10),
		uint32(10),
		uint64(10),
	}
	ints = []interface{}{
		int(10),
		int8(10),
		int16(10),
		int32(10),
		int64(10),
	}
	floats = []interface{}{
		float32(3.14),
		float64(3.14),
	}
	strs = []interface{}{
		"abc",
		"",
	}
	bools = []interface{}{
		true,
		false,
	}
	structs = []interface{}{
		sql.NullString{},
	}
	invalids = []interface{}{
		nil,
	}
	pointers = []interface{}{
		&sql.NullString{},
	}
)

type (
	TestInterface interface {
		A() string
	}
	TestInterfaceImpl struct {
		str string
	}
	TestStruct struct {
		arr  [0]string
		slc  []string
		mp   map[string]interface{}
		str  string
		bl   bool
		i    int
		i8   int8
		i16  int16
		i32  int32
		i64  int64
		ui   uint
		ui8  uint8
		ui16 uint16
		ui32 uint32
		ui64 uint64
		f32  float32
		f64  float64
		intr TestInterface
		ptr  *sql.NullString
	}
)

func (t TestInterfaceImpl) A() string {
	return t.str
}
func (rt *helperTest) SetupTest() {
	// reset the default annotation before each test
	NotatedByDefault(false)
}

type helperTest struct {
	suite.Suite
}

func (rt *helperTest) TestIsUint() {

	for _, v := range uints {
		rt.True(IsUint(reflect.ValueOf(v).Kind()))
	}

	for _, v := range ints {
		rt.False(IsUint(reflect.ValueOf(v).Kind()))
	}
	for _, v := range floats {
		rt.False(IsUint(reflect.ValueOf(v).Kind()))
	}
	for _, v := range strs {
		rt.False(IsUint(reflect.ValueOf(v).Kind()))
	}
	for _, v := range bools {
		rt.False(IsUint(reflect.ValueOf(v).Kind()))
	}
	for _, v := range structs {
		rt.False(IsUint(reflect.ValueOf(v).Kind()))
	}
	for _, v := range invalids {
		rt.False(IsUint(reflect.ValueOf(v).Kind()))
	}
	for _, v := range pointers {
		rt.False(IsUint(reflect.ValueOf(v).Kind()))
	}
}

func (rt *helperTest) TestIsInt() {
	for _, v := range ints {
		rt.True(IsInt(reflect.ValueOf(v).Kind()))
	}

	for _, v := range uints {
		rt.False(IsInt(reflect.ValueOf(v).Kind()))
	}
	for _, v := range floats {
		rt.False(IsInt(reflect.ValueOf(v).Kind()))
	}
	for _, v := range strs {
		rt.False(IsInt(reflect.ValueOf(v).Kind()))
	}
	for _, v := range bools {
		rt.False(IsInt(reflect.ValueOf(v).Kind()))
	}
	for _, v := range structs {
		rt.False(IsInt(reflect.ValueOf(v).Kind()))
	}
	for _, v := range invalids {
		rt.False(IsInt(reflect.ValueOf(v).Kind()))
	}
	for _, v := range pointers {
		rt.False(IsInt(reflect.ValueOf(v).Kind()))
	}
}

func (rt *helperTest) TestIsFloat() {
	for _, v := range floats {
		rt.True(IsFloat(reflect.ValueOf(v).Kind()))
	}

	for _, v := range uints {
		rt.False(IsFloat(reflect.ValueOf(v).Kind()))
	}
	for _, v := range ints {
		rt.False(IsFloat(reflect.ValueOf(v).Kind()))
	}
	for _, v := range strs {
		rt.False(IsFloat(reflect.ValueOf(v).Kind()))
	}
	for _, v := range bools {
		rt.False(IsFloat(reflect.ValueOf(v).Kind()))
	}
	for _, v := range structs {
		rt.False(IsFloat(reflect.ValueOf(v).Kind()))
	}
	for _, v := range invalids {
		rt.False(IsFloat(reflect.ValueOf(v).Kind()))
	}
	for _, v := range pointers {
		rt.False(IsFloat(reflect.ValueOf(v).Kind()))
	}
}

func (rt *helperTest) TestIsString() {
	for _, v := range strs {
		rt.True(IsString(reflect.ValueOf(v).Kind()))
	}

	for _, v := range uints {
		rt.False(IsString(reflect.ValueOf(v).Kind()))
	}
	for _, v := range ints {
		rt.False(IsString(reflect.ValueOf(v).Kind()))
	}
	for _, v := range floats {
		rt.False(IsString(reflect.ValueOf(v).Kind()))
	}
	for _, v := range bools {
		rt.False(IsString(reflect.ValueOf(v).Kind()))
	}
	for _, v := range structs {
		rt.False(IsString(reflect.ValueOf(v).Kind()))
	}
	for _, v := range invalids {
		rt.False(IsString(reflect.ValueOf(v).Kind()))
	}
	for _, v := range pointers {
		rt.False(IsString(reflect.ValueOf(v).Kind()))
	}
}

func (rt *helperTest) TestIsBool() {
	for _, v := range bools {
		rt.True(IsBool(reflect.ValueOf(v).Kind()))
	}

	for _, v := range uints {
		rt.False(IsBool(reflect.ValueOf(v).Kind()))
	}
	for _, v := range ints {
		rt.False(IsBool(reflect.ValueOf(v).Kind()))
	}
	for _, v := range floats {
		rt.False(IsBool(reflect.ValueOf(v).Kind()))
	}
	for _, v := range strs {
		rt.False(IsBool(reflect.ValueOf(v).Kind()))
	}
	for _, v := range structs {
		rt.False(IsBool(reflect.ValueOf(v).Kind()))
	}
	for _, v := range invalids {
		rt.False(IsBool(reflect.ValueOf(v).Kind()))
	}
	for _, v := range pointers {
		rt.False(IsBool(reflect.ValueOf(v).Kind()))
	}
}

func (rt *helperTest) TestIsStruct() {
	for _, v := range structs {
		rt.True(IsStruct(reflect.ValueOf(v).Kind()))
	}

	for _, v := range uints {
		rt.False(IsStruct(reflect.ValueOf(v).Kind()))
	}
	for _, v := range ints {
		rt.False(IsStruct(reflect.ValueOf(v).Kind()))
	}
	for _, v := range floats {
		rt.False(IsStruct(reflect.ValueOf(v).Kind()))
	}
	for _, v := range bools {
		rt.False(IsStruct(reflect.ValueOf(v).Kind()))
	}
	for _, v := range strs {
		rt.False(IsStruct(reflect.ValueOf(v).Kind()))
	}
	for _, v := range invalids {
		rt.False(IsStruct(reflect.ValueOf(v).Kind()))
	}
	for _, v := range pointers {
		rt.False(IsStruct(reflect.ValueOf(v).Kind()))
	}
}

func (rt *helperTest) TestIsSlice() {
	rt.True(IsSlice(reflect.ValueOf(uints).Kind()))
	rt.True(IsSlice(reflect.ValueOf(ints).Kind()))
	rt.True(IsSlice(reflect.ValueOf(floats).Kind()))
	rt.True(IsSlice(reflect.ValueOf(structs).Kind()))

	rt.False(IsSlice(reflect.ValueOf(structs[0]).Kind()))
}

func (rt *helperTest) TestIsInvalid() {
	for _, v := range invalids {
		rt.True(IsInvalid(reflect.ValueOf(v).Kind()))
	}

	for _, v := range uints {
		rt.False(IsInvalid(reflect.ValueOf(v).Kind()))
	}
	for _, v := range ints {
		rt.False(IsInvalid(reflect.ValueOf(v).Kind()))
	}
	for _, v := range floats {
		rt.False(IsInvalid(reflect.ValueOf(v).Kind()))
	}
	for _, v := range bools {
		rt.False(IsInvalid(reflect.ValueOf(v).Kind()))
	}
	for _, v := range strs {
		rt.False(IsInvalid(reflect.ValueOf(v).Kind()))
	}
	for _, v := range structs {
		rt.False(IsInvalid(reflect.ValueOf(v).Kind()))
	}
	for _, v := range pointers {
		rt.False(IsInvalid(reflect.ValueOf(v).Kind()))
	}
}

func (rt *helperTest) TestIsPointer() {
	for _, v := range pointers {
		rt.True(IsPointer(reflect.ValueOf(v).Kind()))
	}

	for _, v := range uints {
		rt.False(IsPointer(reflect.ValueOf(v).Kind()))
	}
	for _, v := range ints {
		rt.False(IsPointer(reflect.ValueOf(v).Kind()))
	}
	for _, v := range floats {
		rt.False(IsPointer(reflect.ValueOf(v).Kind()))
	}
	for _, v := range bools {
		rt.False(IsPointer(reflect.ValueOf(v).Kind()))
	}
	for _, v := range strs {
		rt.False(IsPointer(reflect.ValueOf(v).Kind()))
	}
	for _, v := range structs {
		rt.False(IsPointer(reflect.ValueOf(v).Kind()))
	}
	for _, v := range invalids {
		rt.False(IsPointer(reflect.ValueOf(v).Kind()))
	}
}

func (rt *helperTest) TestIsEmptyValue_emptyValues() {
	ts := TestStruct{}
	rt.True(IsEmptyValue(reflect.ValueOf(ts.arr)))
	rt.True(IsEmptyValue(reflect.ValueOf(ts.slc)))
	rt.True(IsEmptyValue(reflect.ValueOf(ts.mp)))
	rt.True(IsEmptyValue(reflect.ValueOf(ts.str)))
	rt.True(IsEmptyValue(reflect.ValueOf(ts.bl)))
	rt.True(IsEmptyValue(reflect.ValueOf(ts.i)))
	rt.True(IsEmptyValue(reflect.ValueOf(ts.i8)))
	rt.True(IsEmptyValue(reflect.ValueOf(ts.i16)))
	rt.True(IsEmptyValue(reflect.ValueOf(ts.i32)))
	rt.True(IsEmptyValue(reflect.ValueOf(ts.i64)))
	rt.True(IsEmptyValue(reflect.ValueOf(ts.ui)))
	rt.True(IsEmptyValue(reflect.ValueOf(ts.ui8)))
	rt.True(IsEmptyValue(reflect.ValueOf(ts.ui16)))
	rt.True(IsEmptyValue(reflect.ValueOf(ts.ui32)))
	rt.True(IsEmptyValue(reflect.ValueOf(ts.ui64)))
	rt.True(IsEmptyValue(reflect.ValueOf(ts.f32)))
	rt.True(IsEmptyValue(reflect.ValueOf(ts.f64)))
	rt.True(IsEmptyValue(reflect.ValueOf(ts.intr)))
	rt.True(IsEmptyValue(reflect.ValueOf(ts.ptr)))
}

func (rt *helperTest) TestIsEmptyValue_validValues() {
	ts := TestStruct{intr: TestInterfaceImpl{"hello"}}
	rt.False(IsEmptyValue(reflect.ValueOf([1]string{"a"})))
	rt.False(IsEmptyValue(reflect.ValueOf([]string{"a"})))
	rt.False(IsEmptyValue(reflect.ValueOf(map[string]interface{}{"a": true})))
	rt.False(IsEmptyValue(reflect.ValueOf("str")))
	rt.False(IsEmptyValue(reflect.ValueOf(true)))
	rt.False(IsEmptyValue(reflect.ValueOf(int(1))))
	rt.False(IsEmptyValue(reflect.ValueOf(int8(1))))
	rt.False(IsEmptyValue(reflect.ValueOf(int16(1))))
	rt.False(IsEmptyValue(reflect.ValueOf(int32(1))))
	rt.False(IsEmptyValue(reflect.ValueOf(int64(1))))
	rt.False(IsEmptyValue(reflect.ValueOf(uint(1))))
	rt.False(IsEmptyValue(reflect.ValueOf(uint8(1))))
	rt.False(IsEmptyValue(reflect.ValueOf(uint16(1))))
	rt.False(IsEmptyValue(reflect.ValueOf(uint32(1))))
	rt.False(IsEmptyValue(reflect.ValueOf(uint64(1))))
	rt.False(IsEmptyValue(reflect.ValueOf(float32(0.1))))
	rt.False(IsEmptyValue(reflect.ValueOf(float64(0.2))))
	rt.False(IsEmptyValue(reflect.ValueOf(ts.intr)))
	rt.False(IsEmptyValue(reflect.ValueOf(&TestStruct{str: "a"})))
}

func TestHelperSuite(t *testing.T) {
	suite.Run(t, new(helperTest))
}
