package main

import (
	"bufio"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProcessOneVerticesFile(t *testing.T) {
	t.Parallel()
	buffer := strings.Builder{}
	for i := 0; i < 50000; i++ {
		buffer.WriteString(fmt.Sprintf("%d\t%d.example.com\n", i, i))
	}
	fileLength := buffer.Len()
	scanner := bufio.NewScanner(strings.NewReader(buffer.String()))

	result, err := processOneVerticesFile(scanner, "vertices.txt")
	require.NoError(t, err)
	assert.Len(t, result, 3)
	assert.Equal(t, "0.example.com\t0\t0\tvertices.txt", result[0].String())
	assert.Equal(t, "44617.example.com\t1048588\t44617\tvertices.txt", result[1].String())
	assert.Equal(t, "49999.example.com\t1177780\t49999\tvertices.txt", result[2].String())
	assert.Equal(t, fileLength, result[2].Offset())
}

func TestProcessOneVerticesFile_InvalidLine(t *testing.T) {
	t.Parallel()
	data := "domain1\tvalue1\ninvalid_line\ndomain3\tvalue3\n"
	scanner := bufio.NewScanner(strings.NewReader(data))

	_, err := processOneVerticesFile(scanner, "vertices.txt")
	require.Error(t, err)
}

func TestProcessOneVerticesFile_ScannerError(t *testing.T) {
	t.Parallel()
	data := "domain1\tvalue1\ndomain2\tvalue2\ndomain3\tvalue3\n"
	scanner := bufio.NewScanner(strings.NewReader(data))

	scanner.Split(func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		return 0, nil, fmt.Errorf("scanner error")
	})

	_, err := processOneVerticesFile(scanner, "vertices.txt")
	require.Error(t, err)
}
