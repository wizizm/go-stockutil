package stringutil

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestToBytes(t *testing.T) {
	expected := map[string]map[string]float64{
		//  numeric passthrough (no suffix)
		``: map[string]float64{
			`-1`: -1,
			`0`:  0,
			`1`:  1,
			`4611686018427387903`:  4611686018427387903,
			`4611686018427387904`:  4611686018427387904,
			`4611686018427387905`:  4611686018427387905,
			`9223372036854775807`:  9223372036854775807, // beyond this overflows the positive int64 bound
			`-4611686018427387903`: -4611686018427387903,
			`-4611686018427387904`: -4611686018427387904,
			`-4611686018427387905`: -4611686018427387905,
			`-9223372036854775807`: -9223372036854775807,
			`-9223372036854775808`: -9223372036854775808, // beyond this overflows the negative int64 bound
		},

		//  suffix: b/B
		`b`: map[string]float64{
			`-1`: -1,
			`0`:  0,
			`1`:  1,
			`4611686018427387903`:  4611686018427387903,
			`4611686018427387904`:  4611686018427387904,
			`4611686018427387905`:  4611686018427387905,
			`9223372036854775807`:  9223372036854775807,
			`-4611686018427387903`: -4611686018427387903,
			`-4611686018427387904`: -4611686018427387904,
			`-4611686018427387905`: -4611686018427387905,
			`-9223372036854775807`: -9223372036854775807,
			`-9223372036854775808`: -9223372036854775808,
		},

		//  suffix: k/K
		`k`: map[string]float64{
			`-1`:               -1024,
			`0`:                0,
			`1`:                1024,
			`0.5`:              512,
			`2`:                2048,
			`9007199254740992`: 9223372036854775808,
		},

		//  suffix: m/M
		`m`: map[string]float64{
			`-1`:            -1048576,
			`0`:             0,
			`1`:             1048576,
			`0.5`:           524288,
			`8796093022208`: 9223372036854775808,
		},

		//  suffix: g/G
		`g`: map[string]float64{
			`-1`:         -1073741824,
			`0`:          0,
			`1`:          1073741824,
			`0.5`:        536870912,
			`8589934592`: 9223372036854775808,
		},

		//  suffix: t/T
		`t`: map[string]float64{
			`-1`:      -1099511627776,
			`0`:       0,
			`1`:       1099511627776,
			`0.5`:     549755813888,
			`8388608`: 9223372036854775808,
		},

		//  suffix: p/P
		`p`: map[string]float64{
			`-1`:   -1125899906842624,
			`0`:    0,
			`1`:    1125899906842624,
			`0.5`:  562949953421312,
			`8192`: 9223372036854775808,
		},

		//  suffix: e/E
		`e`: map[string]float64{
			`-1`:  -1152921504606846976,
			`0`:   0,
			`1`:   1152921504606846976,
			`0.5`: 576460752303423488,
			`8`:   9223372036854775808,
		},

		//  suffix: z/Z
		`z`: map[string]float64{
			`-1`:  -1180591620717411303424,
			`0`:   0,
			`1`:   1180591620717411303424,
			`0.5`: 590295810358705651712,
		},

		//  suffix: y/Y
		`y`: map[string]float64{
			`-1`:  -1208925819614629174706176,
			`0`:   0,
			`1`:   1208925819614629174706176,
			`0.5`: 604462909807314587353088,
		},
	}

	testExpectations := func(expectedValues map[string]float64, appendToInput string) {
		for in, out := range expectedValues {
			in = in + appendToInput

			if v, err := ToBytes(in); err == nil {
				if v != out {
					t.Errorf("Conversion error on '%s': expected %f, got %f", in, out, v)
				}
			} else {
				t.Errorf("Got error converting '%s' to bytes: %v", in, err)
			}
		}
	}

	for suffix, expectations := range expected {
		testExpectations(expectations, suffix)

		//  only unleash testing hell on higher-order conversions
		if suffix != `` && suffix != `b` {
			testExpectations(expectations, suffix+`b`)
			testExpectations(expectations, suffix+`B`)
			testExpectations(expectations, suffix+`ib`)
			testExpectations(expectations, suffix+`iB`)
		}
	}

	if v, err := ToBytes(`potato`); err == nil {
		t.Errorf("Value 'potato' inexplicably returned a value (%v)", v)
	}

	if v, err := ToBytes(`potatoG`); err == nil {
		t.Errorf("Value 'potatoG' inexplicably returned a value (%v)", v)
	}

	if v, err := ToBytes(`123X`); err == nil {
		t.Errorf("Invalid SI suffix 'X' did not error, got: %v", v)
	}
}

func TestCamelize(t *testing.T) {
	assert := require.New(t)

	tests := map[string]string{
		`Test`:        `Test`,
		`test`:        `Test`,
		`test_value`:  `TestValue`,
		`test-Value`:  `TestValue`,
		`test-Val_ue`: `TestValUe`,
		`test value`:  `TestValue`,
		`TestValue`:   `TestValue`,
		`testValue`:   `TestValue`,
		`TeSt VaLue`:  `TeStVaLue`,
	}

	for have, want := range tests {
		assert.Equal(want, Camelize(have))
	}
}

func TestUnderscore(t *testing.T) {
	assert := require.New(t)

	tests := map[string]string{
		`Test`:       `test`,
		`test`:       `test`,
		`test_value`: `test_value`,
		`test-Value`: `test_value`,
		`test value`: `test_value`,
		`TestValue`:  `test_value`,
		`testValue`:  `test_value`,
		`TeSt VaLue`: `te_st_va_lue`,
	}

	for have, want := range tests {
		assert.Equal(want, Underscore(have))
	}
}

func TestIsMixedCase(t *testing.T) {
	assert := require.New(t)

	assert.False(IsMixedCase(``))
	assert.False(IsMixedCase(`0123456789`))
	assert.False(IsMixedCase(`abcdefghijklmnopqrstuvwxyz`))
	assert.False(IsMixedCase(`abcdefghijklm0123456789nopqrstuvwxyz`))
	assert.False(IsMixedCase(`ABCDEFGHIJKLMNOPQRSTUVWXYZ`))
	assert.False(IsMixedCase(`ABCDEFGHIJKLM0123456789NOPQRSTUVWXYZ`))
	assert.False(IsMixedCase(` ABCDEFGHIJKLM 0123456789 NOPQRSTUVWXYZ `))
	assert.False(IsMixedCase(`сою́з`))
	assert.False(IsMixedCase(`СОЮ́З`))

	assert.True(IsMixedCase(`AbCdEfGhIjKlMnOpQrStUvWxYz`))
	assert.True(IsMixedCase(`ABCDEFGHIJKLM0123456789nopqrstuvwxyz`))
	assert.True(IsMixedCase(`Сою́з`))
}

func TestIsHexadecimal(t *testing.T) {
	assert := require.New(t)

	for i := 0; i < 16; i++ {
		assert.True(IsHexadecimal(fmt.Sprintf("%x", i), -1))
		assert.True(IsHexadecimal(fmt.Sprintf("%X", i), -1))
	}

	for i := 10; i < 16; i++ {
		assert.False(IsHexadecimal(fmt.Sprintf("%x%X", i, i), -1))
		assert.False(IsHexadecimal(fmt.Sprintf("%X%x", i, i), -1))
		assert.False(IsHexadecimal(fmt.Sprintf("%x", i), 2))
		assert.False(IsHexadecimal(fmt.Sprintf("%X", i), 2))
	}

	assert.True(IsHexadecimal(`abc123`, -1))
	assert.True(IsHexadecimal(`ABC123`, -1))
	assert.True(IsHexadecimal(`abc123`, 6))
	assert.True(IsHexadecimal(`ABC123`, 6))
	assert.True(IsHexadecimal(`b26252862a11dd3221427bdbae6025604b1760e4`, 40))

	assert.False(IsHexadecimal(`aBc123`, -1))
	assert.False(IsHexadecimal(`abc123`, 32))
	assert.False(IsHexadecimal(`ABC123`, 32))

}
