package main

import "fmt"

type RollingCounter struct{
	TimeInMs int;
	NumberOfBuckets int;
	BucketSizeInMs int;

}

func main(){
	fmt.Println("RollingCounter")
}

