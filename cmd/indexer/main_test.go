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
	for i := 0; i < 2000; i++ {
		buffer.WriteString(fmt.Sprintf("%d\t%d.example.com\n", i, i))
	}
	fileLength := buffer.Len()
	scanner := bufio.NewScanner(strings.NewReader(buffer.String()))

	result, err := processOneVerticesFile(scanner, "vertices.txt")
	require.NoError(t, err)
	assert.Len(t, result, 3)
	assert.Equal(t, "0.example.com\t0\t0\tvertices.txt", result[0].String())
	assert.Equal(t, "1591.example.com\t32782\t1591\tvertices.txt", result[1].String())
	assert.Equal(t, "1999.example.com\t41780\t1999\tvertices.txt", result[2].String())
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

func TestProcessOneEdgesFile(t *testing.T) {
	t.Parallel()
	buffer := strings.Builder{}
	for i := 0; i < 20000; i++ {
		buffer.WriteString(fmt.Sprintf("%d\t%d\n", i, i))
	}
	scanner := bufio.NewScanner(strings.NewReader(buffer.String()))
	result, err := processOneEdgesFile(scanner, "edges.txt")

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 3, len(result))

	assert.Equal(t, "0\t0\tedges.txt", result[0].String())
	assert.Equal(t, "12775\t131080\tedges.txt", result[1].String())
	assert.Equal(t, "19999\t217780\tedges.txt", result[2].String())
}

func TestProcessOneEdgesFile_InvalidLine(t *testing.T) {
	t.Parallel()
	data := "bad_data\n"
	scanner := bufio.NewScanner(strings.NewReader(data))

	_, err := processOneEdgesFile(scanner, "edges.txt")
	require.Error(t, err)
}

func TestProcessOneEdgesFile_ScannerError(t *testing.T) {
	t.Parallel()
	data := "domain1\n"
	scanner := bufio.NewScanner(strings.NewReader(data))

	scanner.Split(func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		return 0, nil, fmt.Errorf("scanner error")
	})

	_, err := processOneEdgesFile(scanner, "vertices.txt")
	require.Error(t, err)
}
