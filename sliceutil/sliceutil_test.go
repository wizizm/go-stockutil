package sliceutil

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestContainsString(t *testing.T) {
	assert := require.New(t)

	input := []string{"one", "three", "five"}

	assert.True(ContainsString(input, "one"))
	assert.True(ContainsString(input, "three"))
	assert.True(ContainsString(input, "five"))
	assert.False(ContainsString(input, "One"))
	assert.False(ContainsString(input, "two"))
	assert.False(ContainsString(input, "Three"))
	assert.False(ContainsString(input, "four"))
	assert.False(ContainsString(input, "Five"))
	assert.False(ContainsString([]string{}, "one"))
	assert.False(ContainsString([]string{}, "two"))
	assert.False(ContainsString([]string{}, ""))
}

func TestContainsAnyString(t *testing.T) {
	assert := require.New(t)

	input := []string{"one", "three", "five"}
	any := []string{"one", "two", "four"}

	assert.True(ContainsAnyString(input, any...))
	assert.False(ContainsAnyString(input))
	assert.False(ContainsAnyString([]string{}, "one"))
	assert.False(ContainsAnyString([]string{}, "two"))
	assert.False(ContainsAnyString([]string{}, ""))
	assert.False(ContainsAnyString(input, []string{"six", "seven"}...))
}

func TestContainsAllStrings(t *testing.T) {
	assert := require.New(t)

	input := []string{"one", "three", "five"}

	assert.True(ContainsAllStrings(input, "one"))
	assert.True(ContainsAllStrings(input, "three"))
	assert.True(ContainsAllStrings(input, "five"))
	assert.True(ContainsAllStrings(input, "one", "three"))
	assert.True(ContainsAllStrings(input, "one", "five"))
	assert.True(ContainsAllStrings(input, "one", "three", "five"))
	assert.False(ContainsAllStrings(input, "one", "four"))
	assert.True(ContainsAllStrings(input))
}

func TestCompact(t *testing.T) {
	assert := require.New(t)

	assert.Nil(Compact(nil))

	assert.Equal([]interface{}{
		0, 1, 2, 3,
	}, Compact([]interface{}{
		0, 1, 2, 3,
	}))

	assert.Equal([]interface{}{
		1, 3, 5,
	}, Compact([]interface{}{
		nil, 1, nil, 3, nil, 5,
	}))

	assert.Equal([]interface{}{
		`one`, `three`, ` `, `five`,
	}, Compact([]interface{}{
		`one`, ``, `three`, ``, ` `, `five`,
	}))

	assert.Equal([]interface{}{
		[]int{1, 2, 3},
	}, Compact([]interface{}{
		nil, []string{}, []int{1, 2, 3}, map[string]bool{},
	}))
}

func TestCompactString(t *testing.T) {
	assert := require.New(t)

	assert.Nil(CompactString(nil))

	assert.Equal([]string{
		`one`, `three`, `five`,
	}, CompactString([]string{
		`one`, `three`, `five`,
	}))

	assert.Equal([]string{
		`one`, `three`, ` `, `five`,
	}, CompactString([]string{
		`one`, ``, `three`, ``, ` `, `five`,
	}))
}

func TestStringify(t *testing.T) {
	assert := require.New(t)

	assert.Nil(Stringify(nil))

	assert.Equal([]string{
		`0`, `1`, `2`,
	}, Stringify([]interface{}{
		0, 1, 2,
	}))

	assert.Equal([]string{
		`0.5`, `0.55`, `0.555`, `0.555001`,
	}, Stringify([]interface{}{
		0.5, 0.55, 0.55500, 0.555001,
	}))

	assert.Equal([]string{
		`true`, `true`, `false`,
	}, Stringify([]interface{}{
		true, true, false,
	}))
}

func TestOr(t *testing.T) {
	assert := require.New(t)

	assert.Nil(Or())
	assert.Nil(Or(nil))
	assert.Equal(1, Or(0, 1, 0, 2, 0, 3, 4, 5, 6))
	assert.Equal(true, Or(false, false, true))
	assert.Equal(`one`, Or(`one`))
	assert.Equal(4.0, Or(nil, ``, false, 0, 4.0))
	assert.Nil(Or(false, false, false))
	assert.Nil(Or(0, 0, 0))

	assert.Equal(`three`, Or(``, ``, `three`))

	type testStruct struct {
		name string
	}

	assert.Equal(testStruct{`three`}, Or(testStruct{}, testStruct{}, testStruct{`three`}))
}

func TestOrString(t *testing.T) {
	assert := require.New(t)

	assert.Equal(``, OrString())
	assert.Equal(``, OrString(``))

	assert.Equal(`one`, OrString(`one`))
	assert.Equal(`two`, OrString(``, `two`, ``, `three`))
}

func TestEach(t *testing.T) {
	assert := require.New(t)

	assert.Nil(Each(nil, nil))

	assert.Nil(Each([]string{`one`, `two`, `three`}, func(i int, v interface{}) error {
		return Stop
	}))

	count := 0
	assert.Nil(Each([]string{`one`, `two`, `three`}, func(i int, v interface{}) error {
		if v.(string) == `two` {
			return Stop
		} else {
			count += 1
			return nil
		}
	}))

	assert.Equal(1, count)
}
