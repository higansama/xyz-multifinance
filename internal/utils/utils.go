package utils

import (
	"fmt"
	"math"
	"math/rand"
	"path"
	"path/filepath"
	"reflect"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/golang-module/carbon/v2"
	"github.com/higansama/xyz-multi-finance/internal/app"
	"github.com/pkg/errors"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func GetFieldNameForNamespace(obj any, namespace string, finalField string) string {
	if !IsStructAlike(obj) {
		panic("Not a struct")
	}

	namespaceParts := strings.Split(namespace, ".")
	field := namespaceParts[0]

	m := GetStructValue(obj)
	for i := 0; i < m.NumField(); i++ {
		k := m.Type().Field(i)
		kf := m.Field(i)

		re := regexp.MustCompile(`^(?P<field>[\w_]+)\[(?P<index>[0-9]+)\]$`)
		match := re.FindStringSubmatch(field)

		if k.Name == field || (len(match) == 3 && k.Name == match[1]) {
			fld := GetFieldNameFromStructField(k)
			if k.Name != field {
				fld += "." + match[2]
			}
			sfld := []string{finalField, fld}
			nfld := strings.TrimLeft(strings.Join(sfld, "."), ".")
			nns := strings.Join(namespaceParts[1:], ".")

			if kf.Kind() == reflect.Struct {
				return GetFieldNameForNamespace(
					kf.Interface(),
					nns,
					nfld)
			} else if kf.Kind() == reflect.Slice {
				nest := reflect.MakeSlice(kf.Type(), 1, 1).Index(0)
				if nest.Kind() == reflect.Struct { // array of struct
					return GetFieldNameForNamespace(
						nest.Interface(),
						nns,
						nfld)
				}
			}
			return nfld
		}
	}

	return finalField
}

func GetFieldNameFromStructField(v reflect.StructField) string {
	val := v.Tag.Get("attr")
	if val == "" {
		val = v.Tag.Get("form")
		if val == "" {
			val = strings.SplitN(v.Tag.Get("json"), ",", 2)[0]
			if val == "" {
				val = v.Name
			}
		}
	}
	return val
}

func IsStructAlike(s any) bool {
	m := GetStructValue(s)
	return m.Kind() == reflect.Struct
}

func GetStructValue(s any) reflect.Value {
	m := reflect.ValueOf(s)
	if m.Kind() != reflect.Struct && m.Kind() == reflect.Ptr {
		m = reflect.ValueOf(reflect.Indirect(m).Interface())
	}
	return m
}

func ToSliceAny(v any) ([]any, error) {
	s := reflect.ValueOf(v)
	if s.Kind() != reflect.Slice {
		return nil, errors.New("value is not convertible to slice")
	}

	if s.IsNil() {
		return nil, nil
	}

	ret := make([]any, s.Len())
	for i := 0; i < s.Len(); i++ {
		ret[i] = s.Index(i).Interface()
	}

	return ret, nil
}

func TimeInLocal(val time.Time) time.Time {
	if val.IsZero() {
		return val
	}
	return val.Local()
}

func NewTimePtr(val time.Time) *time.Time {
	if val.IsZero() {
		return nil
	}
	return &val
}

func SliceToPipeRegex(values []string) string {
	regVal := ""
	for _, v := range values {
		if v != "" {
			regVal += regexp.QuoteMeta(v) + "|"
		}
	}
	return strings.TrimRight(regVal, "|")
}

func UpperFirst(val string) string {
	if len(val) == 0 {
		return ""
	}
	return strings.ToUpper(val[:1]) + val[1:]
}

func MapStringToLower(values []string) []string {
	for i, v := range values {
		values[i] = strings.ToLower(v)
	}
	return values
}

func MapStringToUpper(values []string) []string {
	for i, v := range values {
		values[i] = strings.ToUpper(v)
	}
	return values
}

func CompareIgnoreCase(v1 string, v2 string) bool {
	return strings.ToLower(v1) == strings.ToLower(v2)
}

func Float64OrZero(val *float64) float64 {
	if val == nil {
		return 0
	}
	return *val
}

func Int64OrZero(val *int64) int64 {
	if val == nil {
		return 0
	}
	return *val
}

func IntOrZero(val *int) int {
	if val == nil {
		return 0
	}
	return *val
}

func TimeParse(val string, format ...string) carbon.Carbon {
	ft := "Y-m-d H:i:s"
	if len(format) > 0 {
		ft = format[0]
	}

	return carbon.ParseByFormat(val, ft)
}

func AnyToInt64(val any) int64 {
	switch v := val.(type) {
	case string:
		r, _ := strconv.ParseInt(v, 10, 64)
		return r
	case float32:
		return int64(v)
	case float64:
		return int64(v)
	case int:
		return int64(v)
	case int8:
		return int64(v)
	case int16:
		return int64(v)
	case int32:
		return int64(v)
	case int64:
		return v
	case uint:
		return int64(v)
	case uint8:
		return int64(v)
	case uint16:
		return int64(v)
	case uint32:
		return int64(v)
	case uint64:
		return int64(v)
	}
	return 0
}

func StringToNullable(val string) *string {
	if val == "" {
		return nil
	}
	return &val
}

func FilterSliceString(data []string, predicate func(v string) bool) []string {
	if predicate == nil {
		predicate = func(v string) bool {
			return v != ""
		}
	}
	var res []string
	for _, v := range data {
		if predicate(v) {
			res = append(res, v)
		}
	}
	return res
}

func StrToSlice(v string, skipEmpty bool) []string {
	var res []string
	if v != "" || !skipEmpty {
		res = append(res, v)
	}
	return res
}

func StrOrDefault(v string, d string) string {
	if v == "" {
		return d
	}
	return v
}

func ChunkBy[T any](items []T, chunkSize int) (chunks [][]T) {
	for chunkSize < len(items) {
		items, chunks = items[chunkSize:], append(chunks, items[0:chunkSize:chunkSize])
	}
	return append(chunks, items)
}

func ArrayStringIncludes(data []string, val string, strict bool) bool {
	for _, v := range data {
		if v == val || (!strict && strings.ToLower(val) == strings.ToLower(v)) {
			return true
		}
	}
	return false
}

func MapStringIncludes(data map[string]string, val string, strict bool) bool {
	for k := range data {
		if k == val || (!strict && strings.ToLower(val) == strings.ToLower(k)) {
			return true
		}
	}
	return false
}

func UnwrapsError(err error) error {
	if err == nil {
		return nil
	}

	for {
		prev := err
		if err = errors.Unwrap(err); err == nil {
			return prev
		}
	}
}

func ExtractFirstAndLastNameFromName(name string) []string {
	parts := strings.Split(name, " ")
	last := ""
	first := ""
	if len(parts) > 1 {
		last = parts[len(parts)-1]
		first = strings.Join(parts[:len(parts)-1], " ")
	} else {
		first = parts[0]
	}

	return []string{first, last}
}

func JoinUrl(base string, paths ...string) string {
	url := strings.TrimRight(base, "/")

	for _, p := range paths {
		url += "/" + strings.TrimLeft(p, "/")
	}

	return url
}

func TitleCase(str string) string {
	title := cases.Title(language.English)
	return title.String(str)
}

func StartCase(str string) string {
	str = strings.ReplaceAll(str, "_", " ")
	str = strings.ReplaceAll(str, "-", " ")
	return TitleCase(str)
}

func KebabCase(str string) string {
	return strings.ReplaceAll(str, " ", "-")
}

func FormatNameForEnv(env app.Environment, name string) string {
	return fmt.Sprintf("%s_%s", name, env.Short())
}

func FormatQueueName(env app.Environment, appName string, name string) string {
	return fmt.Sprintf("%s_%s.%s", appName, env.Short(), name)
}

func PathFromRoot(p string) (string, error) {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return "", errors.New("unable to get the current filename")
	}

	root := filepath.Dir(filename)
	return path.Join(root, "../../", p), nil
}

func GetNextBillingDate(t time.Time) time.Time {
	return carbon.CreateFromStdTime(t).
		SetNanosecond(999999999).
		AddNanosecond().
		ToStdTime()
}

func ExcelTimeFormat(t time.Time) any {
	if t.IsZero() {
		v, _ := time.Parse("2006-01-02", "1970-01-01")
		return v
	}
	return t
}

func TimeOrDefault(t time.Time, d string) any {
	if t.IsZero() && d != "" {
		return d
	}
	return t
}

func NewIntPtr(val int) *int {
	return &val
}

func NewStringPtr(val string) *string {
	return &val
}

func NewBoolPtr(val bool) *bool {
	return &val
}

func TimePtrToTime(t *time.Time) time.Time {
	if t == nil {
		return time.Time{}
	}
	return *t
}

func AddToStringMap(m map[string]string, key string, value string) map[string]string {
	if m == nil {
		m = make(map[string]string)
	}
	m[key] = value
	return m
}

func GenRandomInt(min int, max int) int {
	return rand.Intn(max-min+1) + min
}

const randStrCharset = "abcdefghijklmnopqrstuvwxyz0123456789"

func Int64Range(from int64, to int64) []int64 {
	r := make([]int64, 0)

	if from < to {
		l := to - from + 1
		for i := int64(0); i < l; i++ {
			r = append(r, from+i)
		}
	} else {
		l := from - to + 1
		j := int64(0)
		for i := l; i > 0; i-- {
			r = append(r, from-j)
			j++
		}
	}

	return r
}

func GenRandomStr(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = randStrCharset[rand.Intn(len(randStrCharset))]
	}

	return string(b)
}

func EncNonAsciiToHtmlEntities(val string) string {
	x := []rune(val)

	s := ""
	for _, c := range x {
		v := string(c)
		if c > 127 { // ascii from 1 - 127
			v = "&#" + strconv.FormatInt(int64(c), 10) + ";"
		}
		s += v
	}

	return s
}

func IsUnique(items []any) bool {
	for i, vi := range items {
		for j, vj := range items {
			if i != j && vi == vj {
				return false
			}
		}
	}
	return true
}

func OnlyUniqueString(items []string) []string {
	mp := make(map[string]bool)
	res := make([]string, 0)
	for _, v := range items {
		if _, ok := mp[v]; !ok {
			mp[v] = true
			res = append(res, v)
		}
	}
	return res
}

func CarbonStartOfDayUTC(c carbon.Carbon) carbon.Carbon {
	return c.SetLocation(time.Local).StartOfDay().SetLocation(time.UTC)
}

func CarbonEndOfDayUTC(c carbon.Carbon) carbon.Carbon {
	return c.SetLocation(time.Local).EndOfDay().SetLocation(time.UTC)
}

func DisplayableFloat(v float64, format string) string {
	if math.Mod(v, 1) == 0 {
		return fmt.Sprintf("%d", int64(v))
	}
	return fmt.Sprintf(format, v)
}

func CreateAdminFee(hargaOtr int) int {
	return int(math.Ceil(float64(float64(hargaOtr) * 0.025)))
}

func GenerateKontrakCode() []string {
	// Seed the random number generator
	rand.Seed(time.Now().UnixNano())

	// Create a slice to hold the random integers
	var randomInts []string

	// Generate 5 random integers and convert to strings
	for i := 0; i < 5; i++ {
		num := rand.Intn(100)                              // Generate a random integer between 0 and 99
		randomInts = append(randomInts, strconv.Itoa(num)) // Convert to string and append
	}

	// Print the random integers as strings
	return randomInts
}

func TanggalJatuhTempo() string {
	date := time.Now()

	_, _, d := date.Date()

	if d > 30 {
		return strconv.Itoa(d)
	}

	return strconv.Itoa(d)
}

// FormatRupiah formats an integer to Rupiah format with thousand separators.
func FormatRupiah(amount int) string {
	// Convert integer to string
	strAmount := strconv.Itoa(amount)

	// Reverse the string for easier grouping
	reversed := reverseString(strAmount)

	// Group by three digits using a dot as a separator
	var grouped []string
	for i := 0; i < len(reversed); i += 3 {
		end := i + 3
		if end > len(reversed) {
			end = len(reversed)
		}
		grouped = append(grouped, reversed[i:end])
	}

	// Reverse the grouped string and join with dots
	result := reverseString(strings.Join(grouped, "."))

	return "Rp " + result
}

// Helper function to reverse a string
func reverseString(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}
