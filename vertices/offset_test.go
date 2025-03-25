package vertices_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/dharnitski/cc-hosts/access"
	"github.com/dharnitski/cc-hosts/vertices"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestVerticesOffset_String(t *testing.T) {
	t.Parallel()
	vo := vertices.NewOffset(123, "com.example", 42, "vertices.txt")
	expected := "com.example\t123\t42\tvertices.txt"
	assert.Equal(t, expected, vo.String())
}

func TestOffsets_SaveLoad(t *testing.T) {
	t.Parallel()
	items := []vertices.Offset{
		vertices.NewOffset(123, "com.example", 42, "vertices.txt"),
		vertices.NewOffset(456, "org.example", 84, "vertices.txt"),
	}
	offsets := vertices.Offsets{}
	offsets.Append(items)

	fileName := "test_offsets.txt"
	err := offsets.Save(fileName)
	require.NoError(t, err)
	defer os.Remove(fileName)

	// Verify the file content
	content, err := os.ReadFile(fileName)
	require.NoError(t, err)

	expected := "com.example\t123\t42\tvertices.txt\norg.example\t456\t84\tvertices.txt\n"
	assert.Equal(t, expected, string(content))

	actual := vertices.Offsets{}
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
	var offsets vertices.Offsets
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
	var offsets vertices.Offsets
	err = offsets.Load(fileName)
	if err == nil {
		t.Fatal("Expected an error, but got nil")
	}
}

func TestOffsets_Validate(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		offsets  []vertices.Offset
		expected string
	}{
		{
			name: "Valid offsets",
			offsets: []vertices.Offset{
				vertices.NewOffset(123, "com.example", 42, "vertices.txt"),
				vertices.NewOffset(456, "org.example", 84, "vertices.txt"),
			},
			expected: "",
		},
		{
			name:     "No offsets",
			offsets:  []vertices.Offset{},
			expected: "No offsets found",
		},
		{
			name: "Invalid offset",
			offsets: []vertices.Offset{
				vertices.NewOffset(-123, "com.example", 42, "vertices.txt"),
			},
			expected: "Invalid offset: -123",
		},
		{
			name: "Offset goes down",
			offsets: []vertices.Offset{
				vertices.NewOffset(456, "com.example", 42, "vertices.txt"),
				vertices.NewOffset(123, "org.example", 84, "vertices.txt"),
			},
			expected: "Offset goes down: 123, previous 456",
		},
		{
			name: "Empty domain",
			offsets: []vertices.Offset{
				vertices.NewOffset(123, "", 42, "vertices.txt"),
			},
			expected: "Empty domain",
		},
		{
			name: "Domain goes down",
			offsets: []vertices.Offset{
				vertices.NewOffset(123, "org.example", 42, "vertices.txt"),
				vertices.NewOffset(456, "com.example", 84, "vertices.txt"),
			},
			expected: "Domain goes down: com.example, previous org.example",
		},
		{
			name: "ID goes down",
			offsets: []vertices.Offset{
				vertices.NewOffset(123, "com.example", 84, "vertices.txt"),
				vertices.NewOffset(456, "org.example", 42, "vertices.txt"),
			},
			expected: "ID goes down: 42, previous 84",
		},
		{
			name: "Empty file",
			offsets: []vertices.Offset{
				vertices.NewOffset(123, "org.example", 42, ""),
			},
			expected: "Empty file",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			offsets := vertices.Offsets{}
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

func TestOffsetsFindForDomain(t *testing.T) {
	t.Parallel()

	offsets := vertices.Offsets{}
	err := offsets.Load(fmt.Sprintf("../data/%s", access.VerticesOffsetsFile))
	require.NoError(t, err)

	tests := []string{
		"aaa.11111",
		"ae.regards",
		"com.example",
		"org.example",
		"zw.zzs.th.ac.lpru.arounduniversity.ixiz.qoo",
		"asia.fjs.xr",
	}

	for _, domain := range tests {
		t.Run(domain, func(t *testing.T) {
			t.Parallel()
			start, finish := offsets.FindForDomain(domain)
			assert.LessOrEqual(t, start.Domain(), domain)
			assert.GreaterOrEqual(t, finish.Domain(), domain)
		})
	}
}
