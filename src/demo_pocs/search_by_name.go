//******************************************************************
//
// Exploring how to present a search text, with API
//
// Prepare:
// cd examples
// ../src/N4L-db -u chinese.n4l
//
//******************************************************************

package main

import (
	"fmt"
        SST "SSTorytime"
)

//******************************************************************

const (
	host     = "localhost"
	port     = 5432
	user     = "sstoryline"
	password = "sst_1234"
	dbname   = "newdb"
)

//******************************************************************

func main() {

	load_arrows := false
	ctx := SST.Open(load_arrows)

	cntx := []string{ "yes", "thankyou", "rhyme"}
	chapter := "chinese"
	name := "g"
	nptrs := SST.GetDBNodePtrMatchingNCC(ctx,chapter,name,cntx)

	fmt.Println("RETURNED",nptrs)

	fmt.Println("\nExpanding..")

	for n := range nptrs {
		node := SST.GetDBNodeByNodePtr(ctx,nptrs[n])
		fmt.Println("Found:",node.S)
	}

	SST.Close(ctx)
}

