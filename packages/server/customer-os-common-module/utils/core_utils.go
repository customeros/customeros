package utils

import (
	"encoding/json"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/publicsuffix"
	"net/url"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"time"
)

type Pair[T, U any] struct {
	First  T
	Second U
}

func ToPtr[T any](obj T) *T {
	return &obj
}

func StringPtr(str string) *string {
	return &str
}

func StringPtrNillable(str string) *string {
	if str == "" {
		return nil
	}
	return &str
}

func StringPtrFirstNonEmptyNillable(strs ...string) *string {
	for _, s := range strs {
		if len(s) > 0 {
			return &s
		}
	}
	return nil
}

func BoolPtr(b bool) *bool {
	return &b
}

func TimePtr(t time.Time) *time.Time {
	return &t
}

func TimePtrFirstNonNilNillableAsAny(times ...*time.Time) interface{} {
	for _, t := range times {
		if t != nil {
			return *t
		}
	}
	return nil
}

func NodePtr(node dbtype.Node) *dbtype.Node {
	return &node
}

func RelationshipPtr(relationship dbtype.Relationship) *dbtype.Relationship {
	return &relationship
}

func IntPtr(i int) *int {
	return &i
}

func Int64Ptr(i int64) *int64 {
	return &i
}

func Float64Ptr(f float64) *float64 {
	return &f
}

func Int64PtrToIntPtr(v *int64) *int {
	if v == nil {
		return nil
	}
	var output = int(*v)
	return &output
}

func IntPtrToInt64Ptr(v *int) *int64 {
	if v == nil {
		return nil
	}
	var output = int64(*v)
	return &output
}

func MergeMapToMap(src, dst map[string]any) {
	if dst == nil {
		logrus.Error("expecting not nil map")
	} else if src != nil {
		for k, v := range src {
			dst[k] = v
		}
	}
}

func SurroundWithSpaces(src string) string {
	return SurroundWith(src, " ")
}

func SurroundWithRoundParentheses(src string) string {
	return "(" + src + ")"
}

func SurroundWith(src, surround string) string {
	return surround + src + surround
}

func IfNotNilString(check any, valueExtractor ...func() string) string {
	if reflect.ValueOf(check).Kind() == reflect.String {
		return check.(string)
	}
	if reflect.ValueOf(check).Kind() == reflect.Pointer && reflect.ValueOf(check).IsNil() {
		return ""
	}
	if len(valueExtractor) > 0 {
		return valueExtractor[0]()
	}
	out := check.(*string)
	return *out
}

func IfNotNilStringWithDefault(check any, defaultValue string) string {
	if reflect.ValueOf(check).Kind() == reflect.String {
		return check.(string)
	}
	if (reflect.ValueOf(check).Kind() == reflect.Pointer && reflect.ValueOf(check).IsNil()) || check == nil {
		return defaultValue
	}
	out := check.(*string)
	return *out
}

func IfNotNilInt64(check any, valueExtractor ...func() int64) int64 {
	if reflect.ValueOf(check).Kind() == reflect.Int64 {
		return check.(int64)
	}
	if reflect.ValueOf(check).Kind() == reflect.Pointer && reflect.ValueOf(check).IsNil() {
		return 0
	}
	if len(valueExtractor) > 0 {
		return valueExtractor[0]()
	}
	out := check.(*int64)
	return *out
}

func IfNotNilFloat64(check any, valueExtractor ...func() float64) float64 {
	if reflect.ValueOf(check).Kind() == reflect.Int64 {
		return check.(float64)
	}
	if reflect.ValueOf(check).Kind() == reflect.Pointer && reflect.ValueOf(check).IsNil() {
		return 0
	}
	if len(valueExtractor) > 0 {
		return valueExtractor[0]()
	}
	out := check.(*float64)
	return *out
}

func IfNotNilBool(check any, valueExtractor ...func() bool) bool {
	if reflect.ValueOf(check).Kind() == reflect.Bool {
		return check.(bool)
	}
	if reflect.ValueOf(check).Kind() == reflect.Pointer && reflect.ValueOf(check).IsNil() {
		return false
	}
	if len(valueExtractor) > 0 {
		return valueExtractor[0]()
	}
	out := check.(*bool)
	return *out
}

func IfNotNilTimeWithDefault(check any, defaultValue time.Time) time.Time {
	if check == nil {
		return defaultValue
	}
	if reflect.ValueOf(check).Kind() != reflect.Pointer {
		return check.(time.Time)
	}
	if reflect.ValueOf(check).Kind() == reflect.Pointer && reflect.ValueOf(check).IsNil() {
		return defaultValue
	}
	out := check.(*time.Time)
	return *out
}

func ReverseMap[K comparable, V comparable](in map[K]V) map[V]K {
	out := make(map[V]K, len(in))
	for k, v := range in {
		out[v] = k
	}
	return out
}

func GetFunctionName() string {
	pc, _, _, _ := runtime.Caller(2)
	fullName := runtime.FuncForPC(pc).Name()
	lastSlash := strings.LastIndex(fullName, "/")
	if lastSlash >= 0 {
		fullName = fullName[lastSlash+1:]
	}
	return fullName
}

func ParseStringToFloat(input string) *float64 {
	if input == "" {
		return nil
	}

	parsedFloat, err := strconv.ParseFloat(input, 64)
	if err != nil {
		fmt.Printf("Error parsing string to float: %v\n", err)
		return nil
	}
	return &parsedFloat
}

func FloatToString(num *float64) string {
	if num == nil {
		return ""
	}
	return fmt.Sprintf("%f", *num)
}

func FirstNotEmpty(input ...string) *string {
	for _, item := range input {
		if item != "" {
			return &item
		}
	}
	return nil
}

func ExtractJsonFromString(str string) (string, error) {
	start := strings.IndexByte(str, '{')
	if start == -1 {
		return "", errors.New("could not find start of json")
	}

	end := strings.LastIndexByte(str, '}')
	if end == -1 {
		return "", errors.New("could not find end of json")
	}

	return str[start : end+1], nil
}

func ExtractAfterColon(s string) string {
	// Find first index of colon
	idx := strings.Index(s, ":")
	if idx == -1 {
		// No colon found, return original string
		return s
	}
	// Return substring after colon
	return s[idx+1:]
}

func Float64PtrEquals(a, b *float64) bool {
	if a == nil && b == nil {
		return true
	}
	if a != nil && b != nil {
		return *a == *b
	}
	return false
}

func FormatCurrencyAmount(value float64) string {
	// Format to 2 decimals
	str := fmt.Sprintf("%.2f", value)

	// Split into parts
	parts := strings.Split(str, ".")
	intPart := parts[0]
	fracPart := parts[1]

	// Trim trailing zeros
	fracPart = strings.TrimRight(fracPart, "0")

	// Add commas to int
	intPartFormatted := addThousandSeparators(intPart)

	// Join parts
	if fracPart != "" {
		return intPartFormatted + "." + fracPart
	} else {
		return intPartFormatted
	}
}

// Helper to add commas to an integer string
func addThousandSeparators(value string) string {
	var newParts []string
	// Get length of string
	strlen := len(value)

	for i, char := range value {

		// Insert comma every 3 digits from right
		if i > 0 && (strlen-i)%3 == 0 {
			newParts = append(newParts, ",")
		}

		newParts = append(newParts, string(char))
	}

	return strings.Join(newParts, "")
}

func ToJson(obj any) (string, error) {
	outputJson, err := json.Marshal(obj)
	if err != nil {
		return "", err
	}
	return string(outputJson), nil
}

func ExtractDomain(input string) string {
	if !strings.Contains(input, ".") {
		return ""
	}

	hostname := extractHostname(strings.TrimSpace(strings.ToLower(input)))

	domain, err := publicsuffix.EffectiveTLDPlusOne(hostname)
	if err != nil {
		return ""
	}

	if IsValidTLD(domain) {
		return domain
	}
	return ""
}

func extractHostname(inputURL string) string {
	// Prepend "http://" if the URL doesn't start with a scheme
	if !strings.HasPrefix(inputURL, "http://") && !strings.HasPrefix(inputURL, "https://") {
		inputURL = "http://" + inputURL
	}

	// Parse the URL
	u, err := url.Parse(inputURL)
	if err != nil {
		return ""
	}

	// Extract and return the hostname (domain)
	hostname := u.Hostname()

	// Remove "www." if it exists
	if strings.HasPrefix(hostname, "www.") {
		hostname = hostname[4:] // Remove the first 4 characters ("www.")
	}

	return strings.ToLower(hostname)
}

func IsValidTLD(input string) bool {
	etld, im := publicsuffix.PublicSuffix(input)
	var validtld = false
	if im { // ICANN managed
		validtld = true
	} else if strings.IndexByte(etld, '.') >= 0 { // privately managed
		validtld = true
	}
	return validtld
}

func IsEmptyString(s *string) bool {
	return s == nil || *s == ""
}
