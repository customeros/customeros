package utils

import (
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/types/known/timestamppb"
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
	if len(str) == 0 {
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

func StringFirstNonEmpty(strs ...string) string {
	for _, s := range strs {
		if len(s) > 0 {
			return s
		}
	}
	return ""
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

func Now() time.Time {
	return time.Now().UTC()
}

func Contains(list []string, str string) bool {
	for _, v := range list {
		if v == str {
			return true
		}
	}
	return false
}

func ContainsAll(sourceSlice, itemsToCheck []string) bool {
	for _, item := range itemsToCheck {
		found := false
		for _, sourceItem := range sourceSlice {
			if sourceItem == item {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

func AnySliceToStringSlice(input []any) ([]string, error) {
	result := []string{}
	for _, item := range input {
		str, ok := item.(string)
		if !ok {
			return nil, fmt.Errorf("could not convert item to string")
		}
		result = append(result, str)
	}
	return result, nil
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

func ConvertTimeToTimestampPtr(input *time.Time) *timestamppb.Timestamp {
	if input == nil {
		return nil
	}
	return timestamppb.New(*input)
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

func FilterEmpty(vals []string) []string {
	var result []string
	for _, val := range vals {
		if val != "" {
			result = append(result, val)
		}
	}
	return result
}

func LowercaseStrings(arr []string) {
	for i, s := range arr {
		arr[i] = strings.ToLower(s)
	}
}

func RemoveDuplicates(arr []string) []string {
	var result []string
	seen := map[string]bool{}
	for _, val := range arr {
		if !seen[val] {
			result = append(result, val)
			seen[val] = true
		}
	}
	return result
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
