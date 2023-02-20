package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"net/netip"
	"os"
	"sort"
	"strings"
)

func main() {
	f, err := os.Open("bad.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	r := bufio.NewReader(f)
	addrs := make(map[netip.Addr]bool)
	prefixs := make(map[netip.Prefix]bool)
NextLine:
	for {
		line, err := r.ReadString('\n')
		if err != nil {
			break
		}
		line = strings.TrimSpace(line)
		if prefix, err := netip.ParsePrefix(line); err == nil {
			prefixs[prefix] = true
			continue
		}
		if addr, err := netip.ParseAddr(line); err == nil {
			for _, v := range Keys(prefixs) {
				if v.Contains(addr) {
					continue NextLine
				}
			}
			addrs[addr] = true
		}
	}
	sortedAddrs := Keys(addrs)
	sort.Slice(sortedAddrs, func(i, j int) bool {
		return sortedAddrs[i].Less(sortedAddrs[j])
	})
	sortedPrefixs := Keys(prefixs)
	sort.Slice(sortedPrefixs, func(i, j int) bool {
		return sortedPrefixs[i].Addr().Less(sortedPrefixs[j].Addr())
	})
	b := new(bytes.Buffer)
	for _, prefix := range sortedPrefixs {
		fmt.Fprintln(b, prefix)
	}
	for _, addr := range sortedAddrs {
		fmt.Fprintln(b, addr)
	}
	err = os.WriteFile("bad.txt", b.Bytes(), 0644)
	if err != nil {
		log.Fatal(err)
	}

}

func Keys[T comparable, V any](m map[T]V) []T {
	keys := make([]T, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}
