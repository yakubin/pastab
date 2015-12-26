package main

import (
	"strings"
	"sync"
)

func indexesToString(keywords []string, indexes []int) string {
	perm := make([]string, len(indexes))

	for i, v := range indexes {
		perm[i] = keywords[v-1]
	}

	return strings.Join(perm, "")
}

func totalPerms(n int) int {
	var total int
	last := 1

	for i := 0; i < n; i++ {
		last *= n - i
		total += last
	}

	return total
}

func combine(n int) chan []int {
	results := make([][]int, 0, (1<<uint(n))-1)
	c := make(chan []int, (1<<uint(n))-1)

	go func() {
		for i := 1; i-1 < n; i++ {
			for _, r := range results {
				msg := append(r, i)
				c <- msg
				results = append(results, msg)
			}

			msg := []int{i}
			c <- msg
			results = append(results, msg)
		}
	}()

	return c
}

func redirectChan(out, in chan string, wg *sync.WaitGroup) {
	go func() {
		for v := range in {
			out <- v
		}

		wg.Done()
	}()
}

func genPermCombs(keywords []string) chan string {
	n := len(keywords)
	perms := make(chan string, 1<<20)
	combs := combine(n)

	var wg sync.WaitGroup
	wg.Add((1 << uint(n)) - 1)

	go func() {
		for c := range combs {
			redirectChan(perms, permute(c, keywords), &wg)
		}
	}()

	go func() {
		wg.Wait()
		close(perms)
	}()

	return perms
}

func permute(comb []int, keywords []string) chan string {
	n := len(keywords)
	c := make(chan string, 1000*n)

	go func() {
		for comb != nil {
			c <- indexesToString(keywords, comb)
			comb = nextPermutation(comb)
		}

		close(c)
	}()

	return c
}

// generates next (in lexicographical order) permutation
// returns nil if a is last permutation
func nextPermutation(a []int) []int {
	k := -1
	for i := range a[1:] {
		if a[i] < a[i+1] {
			k = i
		}
	}

	if k == -1 {
		return nil
	}

	l := k + 1
	for i, v := range a[k+2:] {
		if a[k] < v {
			l = i + k + 2
		}
	}

	a[k], a[l] = a[l], a[k]

	revlen := len(a) - (k + 1)
	for i := 0; 2*i < revlen; i++ {
		a[k+1+i], a[revlen+k-i] = a[revlen+k-i], a[k+1+i]
	}

	return a
}
