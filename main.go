package main

import (
	"fmt"
	"go-torrent/pkg/bencode"
)

func PrettyPrint(dict map[string]bencode.DictionaryElement){
	// put this in bencode
	// print in json ideally
}



func main() {
	_, _, err := bencode.ParseDictionary("d3:cowsl3:bob4:maryee")
	fmt.Printf("%v\n", err)
}