package main

import (
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

func evaluate(perm []int, t [][]int, d, c []int) int {
	m := len(t)
	n := len(perm)
	machineReady := make([]int, m)
	completion := make([]int, n)
	for _, job := range perm {
		prevFinish := 0
		for j := 0; j < m; j++ {
			start := machineReady[j]
			if prevFinish > start {
				start = prevFinish
			}
			finish := start + t[j][job]
			machineReady[j] = finish
			prevFinish = finish
		}
		completion[job] = prevFinish
	}
	penalty := 0
	for i := 0; i < n; i++ {
		if completion[i] > d[i] {
			penalty += c[i] * (completion[i] - d[i])
		}
	}
	return penalty
}
func permString(p []int) string {
	var sb strings.Builder
	sb.WriteString("(")
	for i, v := range p {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(strconv.Itoa(v + 1))
	}
	sb.WriteString(")")
	return sb.String()
}
func buildSchedule(perm []int, t [][]int) (x, y [][]int) {
	n := len(perm)
	m := len(t)
	x = make([][]int, n)
	y = make([][]int, n)
	for i := range x {
		x[i] = make([]int, m)
		y[i] = make([]int, m)
	}
	machineReady := make([]int, m)
	for _, job := range perm {
		prevFinish := 0
		for j := 0; j < m; j++ {
			start := machineReady[j]
			if prevFinish > start {
				start = prevFinish
			}
			x[job][j] = start
			y[job][j] = start + t[j][job]
			machineReady[j] = y[job][j]
			prevFinish = y[job][j]
		}
	}
	return x, y
}

func writeResults(filename string, perm []int, t [][]int, d, c []int, x, y [][]int) error {
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	fmt.Fprintln(f, "\n\n===================")
	n := len(perm)
	m := len(t)
	fmt.Fprintf(f, "Permutation: %v\n\n", permString(perm))
	fmt.Fprintf(f, "%-8v", "Job")
	for i := 0; i < m; i++ {
		fmt.Fprintf(f, "%-18v", fmt.Sprintf("Machine %v", i+1))
	}
	fmt.Fprintf(f, "%-14v%-12v%-8v\n", "Completion", "Due", "Fine")
	for job := 0; job < n; job++ {
		fmt.Fprintf(f, "%-8v", job+1)
		for j := 0; j < m; j++ {
			cell := fmt.Sprintf("%v-%v", x[job][j], y[job][j])
			fmt.Fprintf(f, "%-18v", cell)
		}
		finish := y[job][m-1]
		tard := finish - d[job]
		if tard < 0 {
			tard = 0
		}
		fmt.Fprintf(f, "%-14v%-12v%-8v\n", finish, d[job], c[job]*tard)
	}
	return nil
}
func greedyHeuristic(t [][]int, d, c []int, method int) ([]int, int) {
	n := len(d)
	type Job struct {
		id       int
		priority float64
	}
	jobs := make([]Job, n)
	for i := 0; i < n; i++ {
		if method == 1 {
			jobs[i] = Job{id: i, priority: float64(d[i])}
		} else {
			sumTime := 0
			for j := 0; j < len(t); j++ {
				sumTime += t[j][i]
			}
			jobs[i] = Job{id: i, priority: float64(d[i]) / float64(sumTime)}
		}
	}
	sort.Slice(jobs, func(i, j int) bool {
		return jobs[i].priority < jobs[j].priority
	})
	perm := make([]int, n)
	for i, j := range jobs {
		perm[i] = j.id
	}
	penalty := evaluate(perm, t, d, c)
	return perm, penalty
}
func main() {
	data1, err1 := os.ReadFile("param.txt")
	if err1 != nil {
		fmt.Printf("Error: %v\n", err1)
		return
	}
	input1 := strings.Fields(string(data1))
	if len(input1) != 2 {
		fmt.Println("param.txt: need two numbers (m n)")
		return
	}
	m, err2 := strconv.Atoi(input1[0])
	if err2 != nil {
		fmt.Printf("Error: %v\n", err2)
		return
	}
	n, err3 := strconv.Atoi(input1[1])
	if err3 != nil {
		fmt.Printf("Error: %v\n", err3)
		return
	}
	data2, err4 := os.ReadFile("time_w.txt")
	if err4 != nil {
		fmt.Printf("Error: %v\n", err4)
		return
	}
	lines1 := strings.Split(string(data2), "\n")
	var t [][]int
	for _, line := range lines1 {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.Fields(line)
		row := make([]int, n)
		for i, p := range parts {
			v, err := strconv.Atoi(p)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				return
			}
			row[i] = v
		}
		t = append(t, row)
	}
	if len(t) != m {
		fmt.Printf("Error: the file contains %v lines, but %v were expected\n", len(t), m)
		return
	}
	for i, row := range t {
		if len(row) != n {
			fmt.Printf("Error: in row %v, there are %v columns, but %v were expected\n", i, len(row), n)
			return
		}
	}
	data3, err5 := os.ReadFile("time_st.txt")
	if err5 != nil {
		fmt.Printf("Error: %v\n", err5)
		return
	}
	lines2 := strings.Split(string(data3), "\n")
	d := make([]int, n)
	c := make([]int, n)
	for i, line := range lines2 {
		if i == 1 || i == 3 {
			line = strings.TrimSpace(line)
			parts := strings.Fields(line)
			for j, p := range parts {
				v, err := strconv.Atoi(p)
				if err != nil {
					fmt.Printf("Error: %v\n", err)
					return
				}
				if i == 1 {
					d[j] = v
				} else {
					c[j] = v
				}
			}
		}
	}
	perm1, pen1 := greedyHeuristic(t, d, c, 1)
	fmt.Println("Option 1: Priority = mandatory deadline")
	fmt.Printf("Permutation: %v\n", perm1)
	fmt.Printf("Penalty: %v\n", pen1)
	perm2, pen2 := greedyHeuristic(t, d, c, 2)
	fmt.Println("Option 2: Priority = directive deadline / time amount")
	fmt.Printf("Permutation: %v\n", perm2)
	fmt.Printf("Penalty: %v\n", pen2)
	var finalPerm []int
	var finalPenalty int
	if pen1 <= pen2 {
		finalPerm = perm1
		finalPenalty = pen1
	} else {
		finalPerm = perm2
		finalPenalty = pen2
	}
	x, y := buildSchedule(finalPerm, t)
	if err := writeResults("result.txt", finalPerm, t, d, c, x, y); err != nil {
		fmt.Printf("Error: %v\n", err)
	}
	fmt.Printf("The best result: %v, penalty = %v\n", finalPerm, finalPenalty)
}
