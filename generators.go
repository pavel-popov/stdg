package main

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"time"

	"github.com/icrowley/fake"
)

// Enum picks randomly one value from slice of strings.
func Enum(vals []string) string {
	return vals[rand.Intn(len(vals))]
}

var uniqInt32ByKeySeq = make(map[string]int32)

// UniqInt32ByKey returns unique int32 within key values.
func UniqInt32ByKey(key string) string {
	var i int32
	i, ok := uniqInt32ByKeySeq[key]
	if !ok {
		i = 0
	} else {
		i = i + 1
	}

	//Debug.Printf("UniqInt32ByKey generated generated: %d", i)

	uniqInt32ByKeySeq[key] = i
	return strconv.FormatInt(int64(i), 10)
}

// Date returns random date.
func Date(format string, lowBoundary time.Time) string {
	result := lowBoundary.Add(time.Duration(rand.Int63n(int64(time.Since(lowBoundary)))))
	return result.Format(format)
}

// UnixTimestamp returns time in unix timestamp format.
func UnixTimestamp(lowBoundary time.Time) string {
	result := lowBoundary.Add(time.Duration(rand.Int63n(int64(time.Since(lowBoundary)))))
	return fmt.Sprintf("%d", result.Unix())
}

// NormInt32 return Int32 having Normal distribution across mean with stddev.
func NormInt32(mean, stddev float64) string {
	return fmt.Sprintf("%.0f", math.Abs(rand.NormFloat64()*stddev+mean))
}

// NormMultiplierKey returns key value multiplied by factor distributed normally witn mean and stddev.
func NormMultiplierKey(key string, mean, stddev float64) string {
	var keyVal float64
	_, err := fmt.Sscanf(key, "%f", &keyVal)
	if err != nil {
		Error.Printf("Error when scanning key value: %s", err)
	}
	avgPrice := math.Abs(rand.NormFloat64()*stddev + mean)
	//Debug.Printf("Avg price: %f", avgPrice)
	return fmt.Sprintf("%.2f", keyVal*avgPrice)
}

var uniqEmail = make(map[string]struct{})

// UniqEmail generates unique email address.
func UniqEmail() string {
	email := fake.EmailAddress()
	if _, ok := uniqEmail[email]; ok {
		return UniqEmail()
	}
	uniqEmail[email] = struct{}{}
	return email
}

func init() {
	rand.Seed(time.Now().UnixNano())
}
