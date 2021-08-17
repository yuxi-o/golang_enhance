package main

import (
	"fmt"
	"time"
)

type RollingCounter struct {
	TimeInMs        int
	NumberOfBuckets int
	BucketSizeInMs  int
	BucketArray     []Bucket
	StartInMs       int64
	index           int
}

type Bucket struct {
	value int64
	tInMs int64
}

func NewRollingCounter(timeInMs int, numberOfBuckets int) *RollingCounter {
	bucketSizeInMs := timeInMs / numberOfBuckets
	buckets := make([]Bucket, bucketSizeInMs)

	nowInMs := time.Now().UnixNano() / 10e6
	for i := range buckets {
		buckets[i] = Bucket{0, nowInMs}
	}
	return &RollingCounter{
		TimeInMs:        timeInMs,
		NumberOfBuckets: numberOfBuckets,
		BucketSizeInMs:  bucketSizeInMs,
		BucketArray:     buckets,
		StartInMs:       nowInMs,
		index:           0,
	}
}

func (rc *RollingCounter) reset() {
	rc.index = 0
	rc.StartInMs = time.Now().UnixNano() / 10e6
	for i := range rc.BucketArray {
		rc.BucketArray[i].value = 0
		rc.BucketArray[i].tInMs = rc.StartInMs
	}
}

func (rc *RollingCounter) getCurrentBucket() *Bucket {
	nowInMs := time.Now().UnixNano() / 10e6
	if nowInMs < rc.BucketArray[rc.index].tInMs+int64(rc.BucketSizeInMs) {
		rc.BucketArray[rc.index].tInMs = nowInMs
		return &rc.BucketArray[rc.index]
	} else if nowInMs > rc.BucketArray[rc.index].tInMs+int64(rc.TimeInMs) {
		rc.reset()
		return &rc.BucketArray[0]
	}

	index := int((nowInMs-rc.BucketArray[rc.index].tInMs)/int64(rc.BucketSizeInMs) + 1)
	if rc.BucketArray[rc.index].value == 0 {
		rc.BucketArray[rc.index].tInMs = nowInMs
		return &rc.BucketArray[rc.index]
	}

	for i := 1; i < index; i++ {
		rc.BucketArray[(rc.index+i)%rc.NumberOfBuckets].value = rc.BucketArray[rc.index].value
	}

	rc.index = (rc.index + index) % rc.BucketSizeInMs
	rc.BucketArray[rc.index].tInMs = nowInMs
	return &rc.BucketArray[rc.index]
}

func (rc *RollingCounter) Add(value int64) {
	bucket := rc.getCurrentBucket()
	bucket.value = value
}

func (rc *RollingCounter) GetSum() int64 {
	var sum int64
	for i := range rc.BucketArray {
		sum += rc.BucketArray[i].value
	}
	return sum
}

func (rc *RollingCounter) GetMax() int64 {
	var max int64
	for i := range rc.BucketArray {
		if max < rc.BucketArray[i].value {
			max = rc.BucketArray[i].value
		}
	}
	return max
}

func main() {
	fmt.Println("RollingCounter")
}
