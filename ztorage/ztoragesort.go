package ztorage

import "sort"

type ZapMap struct {
	m map[string]int
	s []string
}

func (zm *ZapMap) Len() int {
	return len(zm.m)
}

func (zm *ZapMap) Less(i, j int) bool {
	return zm.m[zm.s[i]] > zm.m[zm.s[j]]
}

func (zm *ZapMap) Swap(i, j int) {
	zm.s[i], zm.s[j] = zm.s[j], zm.s[i]
}

func Top10(m map[string]int) []string {
	zm := new(ZapMap)
	zm.m = m
	zm.s = make([]string, len(m))
	i := 0
	for key, _ := range m {
		zm.s[i] = key
		i++
	}
	sort.Sort(zm)
	if len(zm.s) >= 10 {
		return zm.s[0:10]
	} else {
		return zm.s[0:]
	}
}
