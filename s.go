package main

import "fmt"

func main() {
	s1 := []int{1, 2, 3, 4, 5}
	s2 := []int{1, 4, 7}
	s3 := []int{1, 2, 3, 4, 5, 6}
	//p1, p2, p3 := 0
	p1 :=0
	p2 :=0
	p3 := 0
	var res []int
	low := 1
	high := 6
	for i := 0; i < len(s3); i++ {
		if len(res) >= 3 {
			break
		}
		if s1[p1] > low && s1[p1] < high {
			if len(res) < 3 && s1[p1] <= s2[p2] && s1[p1] <= s3[p3] {
				res = append(res, s1[i])
				p1++
			}
		}
		if s2[p2] > low && s2[p2] < high {
			if len(res) < 3 && s2[p2] <= s1[p1] && s2[p2] <= s3[p3] {
				res = append(res, s2[i])
				p2++
			}
		}
		if s3[p3] > low && s3[p3] < high {
			if len(res) < 3 && s3[p3] <= s1[p1] && s3[p3] <= s1[p1] {
				res = append(res, s3[i])
				p3++
			}
		}
	}
	fmt.Println(res)
}
