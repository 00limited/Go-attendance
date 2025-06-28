package helper

import (
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"path"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

// Generate slug from string
func GenerateSlug(s string) string {
	return strings.ToLower(strings.ReplaceAll(s, " ", "-"))
}

// Generate code contains YYYYMMDDHHmmss
func GenerateCode() string {
	return time.Now().Format("20060102150405")
}

// Convert string 2021-01-01 to time
func ConvertStringToDate(s string) (t time.Time, err error) {
	t, err = time.Parse("2006-01-02", s)
	return
}

// check if valid from is before valid to
func IsValidFromBeforeValidTo(from, to time.Time) bool {
	return from.Before(to)
}

// get total page
func GetTotalPage(totalRecord, totalRecordPerPage int64) int64 {
	if totalRecord == 0 {
		return 0
	}
	return int64(math.Ceil(float64(totalRecord) / float64(totalRecordPerPage)))
}

// get today date
func GetTodayDate() time.Time {
	return time.Now()
}

// string to int
func StringToInt(s string) int {
	i, _ := strconv.Atoi(s)
	return i
}

// string to int64
func StringToInt64(s string) int64 {
	i, _ := strconv.ParseInt(s, 10, 64)
	return i
}

// uint to string
func UintToString(i uint) string {
	return strconv.Itoa(int(i))
}

// array int64 to string with separate
func ArrayInt64ToString(a []int64, separate string) string {
	var s string
	for _, v := range a {
		s += strconv.FormatInt(v, 10) + separate
	}
	return s[:len(s)-1]
}

// array string to string with separate : ["Territory","Salesman"] to 'Territory','Salesman'
func ArrayStringToString(a []string, separate string) string {
	var s string
	for _, v := range a {
		s += v + separate
	}

	return s[:len(s)-1]
}

// array uint to string with separate
func ArrayUintToString(a []uint, separate string) string {
	var s string
	for _, v := range a {
		s += strconv.Itoa(int(v)) + separate
	}
	return s[:len(s)-1]
}

// string to time
func StringToTime(s string) (t time.Time) {
	t, _ = time.Parse("2006-01-02 15:04:05", s)
	return
}

// time to string date
func TimeToStringDate(t time.Time) string {
	return t.Format("2006-01-02")
}

// check is bool
func IsBool(s int) bool {
	return s == 1
}

// string to bool
func StringToBool(s string) bool {
	return s == "1"
}

// string to uint
func StringToUint(s string) uint {
	i, _ := strconv.ParseUint(s, 10, 64)
	return uint(i)
}

// interface to uint
func InterfaceToUint(i interface{}) uint {
	return uint(i.(float64))
}

// array to string : "Territory,Salesman" to 'Territory','Salesman'
func StringWithSeparate(a string, separate string) string {
	var s string
	for _, v := range strings.Split(a, separate) {
		s += "'" + v + "',"
	}
	return s[:len(s)-1]
}

// check is not exist in slice
func IsNotInSliceUint(s uint, a []uint) bool {
	for _, v := range a {
		if s == v {
			return false
		}
	}
	return true
}

// check is exist in slice string
func IsInSliceString(s string, a []string) bool {
	for _, v := range a {
		if s == v {
			return true
		}
	}
	return false
}

// payment method for excel collection plan
func PaymentMethodForExcelCollectionPlan(s string) string {
	// tunai: COD,CBD,CDEP
	// tf: VA, TRF
	if s == "COD" || s == "CBD" || s == "CDEP" {
		return "tunai"
	} else if s == "VA" || s == "TRF" {
		return "transfer"
	} else if s == "GR" {
		return "giro"
	}

	return ""
}

// random value from slice string
func RandomValueFromSliceString(a []string) string {
	return a[rand.Intn(len(a))]
}

// random value from slice int
func RandomValueFromSliceInt(a []int) int {
	return a[rand.Intn(len(a))]
}

// random value by range uint
func RandomValueByRangeUint(min, max uint) uint {
	return uint(rand.Intn(int(max-min)) + int(min))
}

// random Date time by range
func RandomDateTimeByRange(min, max string) time.Time {
	minTime, _ := time.Parse(time.RFC3339, min)
	maxTime, _ := time.Parse(time.RFC3339, max)
	delta := maxTime.Sub(minTime)
	randomDuration := time.Duration(delta)
	return minTime.Add(randomDuration)
}

// random int by range
func RandomIntByRange(min, max int) int {
	return rand.Intn(max-min) + min
}

// random float by range
func RandomFloatByRange(min, max float64) float64 {
	return min + rand.Float64()*(max-min)
}

// generate new uuid
func GenerateUUID() uuid.UUID {
	return uuid.New()
}

// random date string
func RandomDateString() string {
	return time.Now().AddDate(0, 0, RandomIntByRange(-30, 30)).Format("2006-01-02")
}

// random date time by slice string
func RandomDateTimeBySliceString(a []string) time.Time {
	// example: 2021-01-01, 2021-01-02, 2021-01-03
	for _, v := range a {
		t, _ := time.Parse("2006-01-02", v)
		return t
	}

	return time.Now()
}

// get date range from string
func GetDateRangeFromString(s string) (from, to string) {
	// example: 2021-01-01, 2021-01-02, 2021-01-03
	split := strings.Split(s, ",")
	from = split[0]
	to = split[len(split)-1]

	return
}

// get random index by length of slice
func GetRandomIndexByLengthOfSlice(a int) int {
	return rand.Intn(a)
}

// slice string to string with separate
func SliceStringToStringWithSeparate(a string, separate string) []string {
	// [1,2,3] to 1,2,3
	a = strings.ReplaceAll(a, "[", "")
	a = strings.ReplaceAll(a, "]", "")
	a = strings.ReplaceAll(a, " ", "")
	return strings.Split(a, separate)
}

// 2023-02-01 00:00:00 to 2023-02-01
func TimeToStringDateWithoutTime(t time.Time) string {
	return t.Format("2006-01-02")
}

// second to time by type
func SecondToTimeByType(s int, t string) int {
	switch t {
	case "hour":
		return (s / 3600)
	case "minute":
		return (s % 3600) / 60
	case "second":
		return (s % 3600) % 60
	default:
		return 0
	}
}

// concat condition query, example: "a = 1" + " AND " + "b = 2", if only one condition, example: "a = 1"
func ConcatConditionQuery(a, b string) string {
	if a == "" {
		return b
	}

	return a + " AND " + b
}

// merge concat string with separate
func MergeConcatStringWithSeparate(a, b, separate string) string {
	if a == "" {
		return b
	}

	return a + separate + b
}

// uint array with separate to string
func UintArrayWithSeparateToString(a []uint, separate string) string {
	var s string
	for _, v := range a {
		s += strconv.FormatUint(uint64(v), 10) + separate
	}
	return s[:len(s)-1]
}

// get root project
func GetRootProject() string {
	_, b, _, _ := runtime.Caller(0)
	d := path.Join(path.Dir(b))
	return d
}

// concat string
func ConcatString(s ...string) string {
	var r string
	for _, v := range s {
		r += v
	}
	return r
}

// float64 to string
func Float64ToString(f float64) string {
	return strconv.FormatFloat(f, 'f', 0, 64)
}

// struct to []map[string]interface{}
func StructToMapStringInterface(s interface{}) []map[string]interface{} {
	var a []map[string]interface{}
	b, _ := json.Marshal(s)
	json.Unmarshal(b, &a)
	return a
}

// float to round up
func FloatToRoundUp(f float64) float64 {
	return math.Round(f*100) / 100
}

// array struct to array map
func ArrayStructToArrayMap(a interface{}) []map[string]interface{} {
	var b []map[string]interface{}
	c, _ := json.Marshal(a)
	json.Unmarshal(c, &b)
	return b
}

// unique uint array
func UniqueUintArray(a []uint) []uint {
	keys := make(map[uint]bool)
	list := []uint{}
	for _, entry := range a {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

// format time with timezone
func FormatTimeWithTimezone(t time.Time, timezone string) time.Time {
	loc, _ := time.LoadLocation(timezone)
	return t.In(loc)
}

// format time by layout
func FormatTimeByLayout(t time.Time, layout string) string {
	return t.Format(layout)
}

// interface to map[string]interface{}
func InterfaceToMapStringInterface(i interface{}) map[string]interface{} {
	var a map[string]interface{}
	b, _ := json.Marshal(i)
	json.Unmarshal(b, &a)
	return a
}

// format number 1000 to 1,000.00
func FormatNumberWithComma(n float64) string {
	return fmt.Sprintf("%0.2f", n)
}

// constaint uint in array uint
func ConstaintUintInArrayUint(a []uint, b uint) bool {
	for _, v := range a {
		if v == b {
			return true
		}
	}
	return false
}

// constaint string in array string
func ConstaintStringInArrayString(a []string, b string) bool {
	for _, v := range a {
		if v == b {
			return true
		}
	}
	return false
}

// get index by value in slice string
func GetIndexUintInArrayUint(a []uint, b uint) int {
	for i, v := range a {
		if v == b {
			return i
		}
	}
	return -1
}

// get index by value in slice string
func GetIndexStringInArrayString(a []string, b string) int {
	for i, v := range a {
		if v == b {
			return i
		}
	}
	return -1
}

// remove key in map
func RemoveKeyInMap(m map[string]interface{}, k string) map[string]interface{} {
	delete(m, k)
	return m
}

// replace ? if contains ? on string
func ReplaceStringIfContains(s string, contains string, replace string) string {
	if strings.Contains(s, contains) {
		return strings.ReplaceAll(s, contains, replace)
	}
	return s
}

type ModelDocNumber struct {
	CompanyID   uint
	WarehouseID uint
	TerritoryID uint
	DocType     string
}

// in array string
func InArrayString(s string, a []string) bool {
	for _, v := range a {
		if s == v {
			return true
		}
	}
	return false
}

// get extension from file
func GetExtensionFromFile(s string) string {
	split := strings.Split(s, ".")
	return split[len(split)-1]
}

// array string to string with separate : ["Territory","Salesman"] to 'Territory','Salesman'
func ArrayStringToStringWithSeparator(a []string, separate string) string {
	var s string

	if len(a) == 1 {
		return "'" + a[0] + "'"
	} else {
		for _, v := range a {
			s += "'" + v + "'" + separate
		}
	}

	return s[:len(s)-1]
}

func RemoveDuplicateValues(intSlice []string) []string {
	keys := make(map[string]bool)
	list := []string{}

	for _, entry := range intSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

// mapping array of T to array of T2 by mapper func
func ArrayMapping[T any, T2 any](arr []T, mapper func(T) T2) []T2 {
	res := []T2{}
	for _, d := range arr {
		res = append(res, mapper(d))
	}
	return res
}

// check if item in the array (go slice).
// type of item must same as []array type
func InArr[T comparable](item T, arr []T) bool {
	for _, a := range arr {
		if item == a {
			return true
		}
	}
	return false
}

func RoundFloat(val float64, precision uint) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Round(val*ratio) / ratio
}

func ConvertExcelDateString(dateStr string) (string, error) {
	// Daftar format tanggal umum yang sering digunakan di Excel
	formats := []string{
		"02/01/2006", "01/02/2006", // dd/mm/yyyy atau mm/dd/yyyy
		"02-01-2006", "01-02-2006", // dd-mm-yyyy atau mm-dd-yyyy
		"2006/01/02", "2006-01-02", // yyyy/mm/dd atau yyyy-mm-dd
		"1/2/06", "2-Jan-06", // Format pendek
		"02 January 2006", "January 02, 2006", // Format panjang
		"02-01-06", // dd-mm-yy (2 digit tahun)
		"02/01/06", // dd/mm/yy
	}

	// Coba parsing dengan berbagai format
	for _, layout := range formats {
		if t, err := time.Parse(layout, dateStr); err == nil {
			return t.Format("2006-01-02"), nil
		}
	}

	return "", fmt.Errorf("format tanggal tidak dikenali: %s", dateStr)
}
