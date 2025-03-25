package edges_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/dharnitski/cc-hosts/access"
	"github.com/dharnitski/cc-hosts/edges"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOffset(t *testing.T) {
	t.Parallel()
	offset := edges.NewOffset(123, "id123", "edges.txt")
	expected := "id123\t123\tedges.txt"
	assert.Equal(t, expected, offset.String())
}

func TestOffsets_SaveLoad(t *testing.T) {
	t.Parallel()
	items := []edges.Offset{
		edges.NewOffset(123, "42", "vertices.txt"),
		edges.NewOffset(456, "84", "vertices.txt"),
	}
	offsets := edges.Offsets{}
	offsets.Append(items)

	fileName := "test_offsets.txt"
	err := offsets.Save(fileName)
	require.NoError(t, err)
	defer os.Remove(fileName)

	// Verify the file content
	content, err := os.ReadFile(fileName)
	require.NoError(t, err)

	expected := "42\t123\tvertices.txt\n84\t456\tvertices.txt\n"
	assert.Equal(t, expected, string(content))

	actual := edges.Offsets{}
	err = actual.Load(fileName)
	require.NoError(t, err)
	assert.Equal(t, len(items), actual.Len())
	for i, offset := range actual.Items() {
		assert.Equal(t, items[i], offset)
	}
}

func TestOffsets_Load_InvalidFile(t *testing.T) {
	t.Parallel()
	// Try to load offsets from a non-existent file
	var offsets edges.Offsets
	err := offsets.Load("non_existent_file.txt")
	if err == nil {
		t.Fatal("Expected an error, but got nil")
	}
}

func TestOffsets_Load_InvalidLine(t *testing.T) {
	t.Parallel()
	// Create a temporary file with invalid test data
	fileName := "test_invalid_offsets.txt"
	content := "invalid_line\norg.example\t456\tvertices.txt\n"
	err := os.WriteFile(fileName, []byte(content), 0o644)
	if err != nil {
		t.Fatalf("Expected no error, but got %v", err)
	}
	defer os.Remove(fileName)

	// Try to load the offsets from the file
	var offsets edges.Offsets
	err = offsets.Load(fileName)
	if err == nil {
		t.Fatal("Expected an error, but got nil")
	}
}

func TestOffsets_Validate(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		offsets  []edges.Offset
		expected string
	}{
		{
			name: "Valid offsets",
			offsets: []edges.Offset{
				edges.NewOffset(123, "42", "vertices.txt"),
				edges.NewOffset(456, "84", "vertices.txt"),
			},
			expected: "",
		},
		{
			name:     "No offsets",
			offsets:  []edges.Offset{},
			expected: "No offsets found",
		},
		{
			name: "Invalid offset",
			offsets: []edges.Offset{
				edges.NewOffset(-123, "42", "vertices.txt"),
			},
			expected: "Invalid offset: -123",
		},
		{
			name: "Offset goes down",
			offsets: []edges.Offset{
				edges.NewOffset(456, "42", "vertices.txt"),
				edges.NewOffset(123, "84", "vertices.txt"),
			},
			expected: "Offset goes down: 123, previous 456",
		},
		{
			name: "Empty ID",
			offsets: []edges.Offset{
				edges.NewOffset(123, "", "vertices.txt"),
			},
			expected: "Empty id",
		},
		{
			name: "ID not number",
			offsets: []edges.Offset{
				edges.NewOffset(123, "aaa", "vertices.txt"),
			},
			expected: "Error converting id to integer: strconv.Atoi: parsing \"aaa\": invalid syntax",
		},
		{
			name: "ID goes down",
			offsets: []edges.Offset{
				edges.NewOffset(123, "84", "vertices.txt"),
				edges.NewOffset(456, "42", "vertices.txt"),
			},
			expected: "ID goes down: 42, previous 84",
		},
		{
			name: "Empty file",
			offsets: []edges.Offset{
				edges.NewOffset(123, "42", ""),
			},
			expected: "Empty file",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			offsets := edges.Offsets{}
			offsets.Append(tt.offsets)
			err := offsets.Validate()
			if tt.expected == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tt.expected)
			}
		})
	}
}

func TestOffsetsFindForFromID(t *testing.T) {
	t.Parallel()

	offsets := edges.Offsets{}
	err := offsets.Load(fmt.Sprintf("../data/%s", access.EdgesOffsetsFile))
	require.NoError(t, err)

	tests := []struct {
		id   string
		from int
		to   int
	}{
		// 0 is not on file
		{id: "0", from: 0, to: 0},
		// 74 is not of file
		{id: "74", from: 0, to: 0},
		// first line
		{id: "75", from: 0, to: 1048590},
		{id: "96032", from: 0, to: 1048590},
		// second line
		{id: "96033", from: 0, to: 2097167},
		{id: "96034", from: 1048590, to: 2097167},
		// last line
		{id: "283704001", from: 3688922830, to: 3689816010},
	}
	for _, tt := range tests {
		t.Run(tt.id, func(t *testing.T) {
			t.Parallel()
			allOffsets := offsets.FindForFromID(tt.id)
			offset, ok := allOffsets["part-00000-02106921-c60f-49b6-912c-b03ea5690455-c000.txt"]
			assert.True(t, ok)
			assert.Equal(t, tt.from, offset.From.Offset())
			assert.Equal(t, tt.to, offset.To.Offset())
		})
	}
}
