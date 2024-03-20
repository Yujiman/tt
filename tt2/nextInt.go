package main

import "log"

func main() {
	index, exist := binarysearch([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, 10)
	if !exist {
		println("not exist")
	}
	println(index)
}
func binarysearch(array []int, target int) (int, bool) {

	low := 0
	high := len(array) - 1
	for low <= high {

		midllindex := (high + low) / 2
		if array[midllindex] == target {
			return midllindex, true
		}
		if array[midllindex] < target {
			low = midllindex + 1
		}
		if array[midllindex] > target {
			high = midllindex - 1
		}
		log.Printf("low: %d, high: %d, midllindex: %d", low, high, midllindex)
	}
	return 0, false
}
