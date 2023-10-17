package utils

import (
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestCoreUtils_AddMapToMap_NilToEmpty(t *testing.T) {
	var src map[string]any
	var dst = map[string]any{}
	MergeMapToMap(src, dst)
	require.Empty(t, dst)
}

func TestCoreUtils_AddMapToMap_NilToNotEmpty(t *testing.T) {
	var src map[string]any
	var dst = map[string]any{"k": "v"}
	MergeMapToMap(src, dst)

	require.NotNil(t, dst)
	require.Equal(t, "v", dst["k"])
}

func TestCoreUtils_AddMapToMap_NotEmptyToEmpty(t *testing.T) {
	var src = map[string]any{"k": "v"}
	var dst = map[string]any{}
	MergeMapToMap(src, dst)

	require.NotNil(t, dst)
	require.Equal(t, "v", dst["k"])
}

func TestCoreUtils_AddMapToMap_NotEmptyToNotEmpty(t *testing.T) {
	var src = map[string]any{"k": "v"}
	var dst = map[string]any{"e": "f"}
	MergeMapToMap(src, dst)

	require.NotNil(t, dst)
	require.Equal(t, "v", dst["k"])
	require.Equal(t, "f", dst["e"])
}

func TestStringPtr(t *testing.T) {
	str := "test"
	ptr := StringPtr(str)

	require.Equal(t, str, *ptr)
}

func TestIntPtr(t *testing.T) {
	num := 42
	ptr := IntPtr(num)

	require.Equal(t, num, *ptr)
}

func TestRemoveDuplicates(t *testing.T) {
	input := []string{"a", "b", "a", "c", "b"}
	expected := []string{"a", "b", "c"}

	result := RemoveDuplicates(input)
	require.Equal(t, expected, result)
}

func TestReverseMap(t *testing.T) {
	input := map[int]string{
		1: "a",
		2: "b",
	}

	expected := map[string]int{
		"a": 1,
		"b": 2,
	}

	result := ReverseMap(input)
	require.Equal(t, expected, result)
}

func TestNodePtr(t *testing.T) {
	node := dbtype.Node{}
	ptr := NodePtr(node)

	require.Equal(t, &node, ptr)
}

func TestContains(t *testing.T) {
	list := []string{"a", "b", "c"}

	require.True(t, Contains(list, "b"))
	require.False(t, Contains(list, "x"))
}

func TestBoolPtr(t *testing.T) {
	b := true
	ptr := BoolPtr(b)

	require.Equal(t, b, *ptr)
}

func TestTimePtr(t *testing.T) {
	now := time.Now()
	ptr := TimePtr(now)

	require.Equal(t, now, *ptr)
}

func TestSurroundWithSpaces(t *testing.T) {
	str := "test"
	expected := " test "

	result := SurroundWithSpaces(str)
	require.Equal(t, expected, result)
}

func TestInt64Ptr(t *testing.T) {
	num := int64(42)
	ptr := Int64Ptr(num)

	require.Equal(t, num, *ptr)
}

func TestLowercaseStrings(t *testing.T) {
	input := []string{"A", "B", "C"}
	expected := []string{"a", "b", "c"}

	LowercaseStrings(input)
	require.Equal(t, expected, input)
}

func TestStringPtrNil(t *testing.T) {
	ptr := StringPtrNillable("")
	require.Nil(t, ptr)

	ptr = StringPtrNillable("test")
	require.NotNil(t, ptr)
}

func TestIfNotNilStringDefault(t *testing.T) {
	result := IfNotNilStringWithDefault(nil, "default")
	require.Equal(t, "default", result)

	str := "test"
	result = IfNotNilStringWithDefault(&str, "default")
	require.Equal(t, "test", result)
}

func TestRelationshipPtr(t *testing.T) {
	rel := dbtype.Relationship{}
	ptr := RelationshipPtr(rel)

	require.Equal(t, &rel, ptr)
}

func TestFloatToString(t *testing.T) {
	f := 1.5
	ptr := &f

	result := FloatToString(ptr)
	require.Equal(t, "1.500000", result)
}

func TestToPtr(t *testing.T) {
	str := "test"
	ptr := ToPtr(str)

	require.Equal(t, &str, ptr)
}

func TestStringFirstNonEmpty(t *testing.T) {
	result := StringFirstNonEmpty("a", "", "b")
	require.Equal(t, "a", result)

	result = StringFirstNonEmpty("", "", "")
	require.Equal(t, "", result)
}

func TestTimePtrFirstNonNil(t *testing.T) {
	var time1 *time.Time
	time2 := time.Now()

	result := TimePtrFirstNonNilNillableAsAny(time1, &time2)

	require.Equal(t, time2, result)
}

func TestInt64PtrToIntPtr(t *testing.T) {
	num := int64(42)
	ptr := Int64PtrToIntPtr(&num)

	require.Equal(t, int(num), *ptr)
}

func TestIntPtrToInt64Ptr(t *testing.T) {
	num := 42
	ptr := IntPtrToInt64Ptr(&num)

	require.Equal(t, int64(num), *ptr)
}

func TestMergeMapToMap(t *testing.T) {
	src := map[string]any{"a": "1"}
	dst := map[string]any{"b": "2"}

	MergeMapToMap(src, dst)

	require.Equal(t, map[string]any{"a": "1", "b": "2"}, dst)
}

func TestSurroundWith(t *testing.T) {
	str := "test"
	expected := "|test|"

	result := SurroundWith(str, "|")
	require.Equal(t, expected, result)
}

func TestIfNotNilInt64(t *testing.T) {
	num := int64(42)

	result := IfNotNilInt64(num)
	require.Equal(t, num, result)

	var pNum *int64
	result = IfNotNilInt64(pNum)
	require.Equal(t, int64(0), result)
}

func TestIfNotNilBool(t *testing.T) {
	b := true
	result := IfNotNilBool(b)
	require.True(t, result)

	var pBool *bool
	result = IfNotNilBool(pBool)
	require.False(t, result)
}

func TestIfNotNilTimeDefault(t *testing.T) {
	timeVal := time.Now()
	def := time.Time{}

	result := IfNotNilTimeWithDefault(timeVal, def)
	require.Equal(t, timeVal, result)

	result = IfNotNilTimeWithDefault(nil, def)
	require.Equal(t, def, result)
}

func TestExtractDomain(t *testing.T) {
	// Positive test case: URL with http:// scheme
	inputURL := "http://www.example.com"
	expectedDomain := "example.com"
	actualDomain := ExtractDomain(inputURL)
	if actualDomain != expectedDomain {
		t.Errorf("Expected domain: %s, but got: %s", expectedDomain, actualDomain)
	}

	// Positive test case: URL with https:// scheme
	inputURL = "https://www.example.com"
	expectedDomain = "example.com"
	actualDomain = ExtractDomain(inputURL)
	if actualDomain != expectedDomain {
		t.Errorf("Expected domain: %s, but got: %s", expectedDomain, actualDomain)
	}

	// Negative test case: Invalid URL
	inputURL = "invalidurl"
	expectedDomain = ""
	actualDomain = ExtractDomain(inputURL)
	if actualDomain != expectedDomain {
		t.Errorf("Expected domain: %s, but got: %s", expectedDomain, actualDomain)
	}

	// Edge case: URL without scheme
	inputURL = "example.com"
	expectedDomain = "example.com"
	actualDomain = ExtractDomain(inputURL)
	if actualDomain != expectedDomain {
		t.Errorf("Expected domain: %s, but got: %s", expectedDomain, actualDomain)
	}

	// Edge case: URL without scheme and with subdomains
	inputURL = "hu.example.com"
	expectedDomain = "example.com"
	actualDomain = ExtractDomain(inputURL)
	if actualDomain != expectedDomain {
		t.Errorf("Expected domain: %s, but got: %s", expectedDomain, actualDomain)
	}

	// Edge case: URL without scheme and with subdomains
	inputURL = "ro.example.co.uk"
	expectedDomain = "example.co.uk"
	actualDomain = ExtractDomain(inputURL)
	if actualDomain != expectedDomain {
		t.Errorf("Expected domain: %s, but got: %s", expectedDomain, actualDomain)
	}

	// Edge case: URL without scheme and with subdomains
	inputURL = "home.iasi.ro.example.co.uk"
	expectedDomain = "example.co.uk"
	actualDomain = ExtractDomain(inputURL)
	if actualDomain != expectedDomain {
		t.Errorf("Expected domain: %s, but got: %s", expectedDomain, actualDomain)
	}

	// Edge case: URL without scheme and with subdomains
	inputURL = "example.cop.ro"
	expectedDomain = "cop.ro"
	actualDomain = ExtractDomain(inputURL)
	if actualDomain != expectedDomain {
		t.Errorf("Expected domain: %s, but got: %s", expectedDomain, actualDomain)
	}

	// Edge case: URL without scheme and with subdomains
	inputURL = "example.co.ro"
	expectedDomain = "example.co.ro"
	actualDomain = ExtractDomain(inputURL)
	if actualDomain != expectedDomain {
		t.Errorf("Expected domain: %s, but got: %s", expectedDomain, actualDomain)
	}

	// Edge case: URL without scheme
	inputURL = "example.co.uk"
	expectedDomain = "example.co.uk"
	actualDomain = ExtractDomain(inputURL)
	if actualDomain != expectedDomain {
		t.Errorf("Expected domain: %s, but got: %s", expectedDomain, actualDomain)
	}

	// Edge case: URL with www. prefix
	inputURL = "http://www.example.com"
	expectedDomain = "example.com"
	actualDomain = ExtractDomain(inputURL)
	if actualDomain != expectedDomain {
		t.Errorf("Expected domain: %s, but got: %s", expectedDomain, actualDomain)
	}

	// Edge case: URL with wrong TLD
	inputURL = "http://www.example.stupidme"
	expectedDomain = ""
	actualDomain = ExtractDomain(inputURL)
	if actualDomain != expectedDomain {
		t.Errorf("Expected domain: %s, but got: %s", expectedDomain, actualDomain)
	}

	// Positive test case: URL with path
	inputURL = "http://www.example.com/main/final?param1=true&param2=1"
	expectedDomain = "example.com"
	actualDomain = ExtractDomain(inputURL)
	if actualDomain != expectedDomain {
		t.Errorf("Expected domain: %s, but got: %s", expectedDomain, actualDomain)
	}
}
