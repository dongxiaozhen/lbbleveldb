package main

import (
	"flag"
)

var index = flag.Int64("index", 1, "1:table,2:log,3:manifest")
var fileName = flag.String("file", "000033.ldb", "*ldb,MANIFEST*,*.log")

func main() {
	flag.Parse()

	switch *index {
	case 1:
		sstable(*fileName)
	case 2:
		log(*fileName)
	case 3:
		manifest(*fileName)
	default:
	}
}

const (
	fullChunkType   = 1
	firstChunkType  = 2
	middleChunkType = 3
	lastChunkType   = 4
)
