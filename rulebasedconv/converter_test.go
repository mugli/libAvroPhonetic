package rulebasedconv_test

import (
	"bufio"
	"os"
	"strings"
	"testing"

	"github.com/mugli/libavrophonetic/rulebasedconv"
	"github.com/stretchr/testify/assert"
)

type testCase struct {
	input string
	want  string
}

func buildTestCases() ([]testCase, error) {
	f, err := os.OpenFile("testdata.txt", os.O_RDONLY, os.ModePerm)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	retval := make([]testCase, 0)

	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())

		if line != "" && !strings.HasPrefix(line, "#") {
			parts := strings.Split(line, " ")

			tc := testCase{
				input: parts[0],
				want:  parts[1],
			}

			retval = append(retval, tc)
		}
	}

	if err := sc.Err(); err != nil {
		return nil, err
	}

	return retval, nil
}

func TestConvertWord(t *testing.T) {
	converter := rulebasedconv.NewConverter()
	testCases, err := buildTestCases()
	assert.NoError(t, err)

	for _, testCase := range testCases {
		t.Run(testCase.input, func(t *testing.T) {
			got := converter.ConvertWord(testCase.input)
			assert.Equal(t, testCase.want, got)
		})
	}
}

func BenchmarkConvertWord(b *testing.B) {
	converter := rulebasedconv.NewConverter()
	testCases, _ := buildTestCases()
	numTests := len(testCases) - 1

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		converter.ConvertWord(testCases[i%numTests].input)
	}
}
