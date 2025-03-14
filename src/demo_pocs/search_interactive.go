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
	//"os"
	//"bufio"
	"strings"

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

	//reader := bufio.NewReader(os.Stdin)

	for goes := 0; goes < 10; goes ++ {

		/*
		fmt.Println("\n\nEnter chapter text:")
		chaptext, _ := reader.ReadString('\n')
		fmt.Println("\n\nEnter context text:")
		context, _ := reader.ReadString('\n')
		fmt.Println("\n\nEnter search text:")
		searchtext, _ := reader.ReadString('\n')
		*/
		searchtext := "tiger"
		chaptext := "chinese"
		context := "poem"

		Search(ctx,chaptext,context,searchtext)
	}

	SST.Close(ctx)
}

//******************************************************************

func Search(ctx SST.PoSST, chaptext,context,searchtext string) {

	chaptext = strings.TrimSpace(chaptext)
	context = strings.TrimSpace(context)
	searchtext = strings.TrimSpace(searchtext)

	// **** Look for meaning in the arrows ***

	var ama map[SST.ArrowPtr][]SST.NodePtr

	ama = SST.GetMatroidArrayByArrow(ctx,context,chaptext)

	fmt.Println("--------------------------------------------------")
	fmt.Println("Looking for relevant arrows by",context,chaptext)
	fmt.Println("--------------------------------------------------")
	
	for arrowptr := range ama {
		arr_dir := SST.GetDBArrowByPtr(ctx,arrowptr)

		if strings.Contains(arr_dir.Long,context) {

			fmt.Println("\nArrow --(",arr_dir.Long,")--> points to a group of nodes with a similar role in the context of",context,"in the chapter",chaptext,"\n")
			
			for n := 0; n < len(ama[arrowptr]); n++ {
				node := SST.GetDBNodeByNodePtr(ctx,ama[arrowptr][n])
				SST.NewLine(n)
				fmt.Print("..  ",node.S,",")
				
			}
			fmt.Println()
			fmt.Println("............................................")
		}
	}

	fmt.Println("--------------------------------------------------")
	fmt.Println("Looking for relevant nodes by",searchtext)
	fmt.Println("--------------------------------------------------")

	const maxdepth = 5
	
	var start_set []SST.NodePtr
	
	search_items := strings.Split(searchtext," ")
	
	for w := range search_items {
		fmt.Print("Looking for nodes like ",search_items[w],"...")
		start_set = append(start_set,SST.GetDBNodePtrMatchingName(ctx,search_items[w])...)
	}

	fmt.Println("Found possible relevant nodes:",start_set)

	for start := range start_set {

		for sttype := -SST.EXPRESS; sttype <= SST.EXPRESS; sttype++ {

			fmt.Println("   ... Searching these",len(start_set),"nodes by",SST.STTypeName(sttype))
			
			name :=  SST.GetDBNodeByNodePtr(ctx,start_set[start])

			allnodes := SST.GetFwdConeAsNodes(ctx,start_set[start],sttype,maxdepth)
			
			if len(allnodes) > 1 {			
				fmt.Println()
				fmt.Println("    -------------------------------------------")
				fmt.Printf("     SEARCH TEXT MATCH #%d by %s : (%s -> %s)\n",start+1,SST.STTypeName(sttype),searchtext,name.S)
				fmt.Println("    -------------------------------------------")
				
				for l := range allnodes {
					fullnode := SST.GetDBNodeByNodePtr(ctx,allnodes[l])
					fmt.Println("   - Fwd ",SST.STTypeName(sttype)," cone item: ",fullnode.S,", found in",fullnode.Chap)
				}
			
				alt_paths,path_depth := SST.GetFwdPathsAsLinks(ctx,start_set[start],sttype,maxdepth)
				
				if alt_paths != nil {
					
					fmt.Printf("\n-- Forward",SST.STTypeName(sttype),"cone stories ----------------------------------\n")
					
					for p := 0; p < path_depth; p++ {
						SST.PrintLinkPath(ctx,alt_paths,p,"\nStory:")
					}
				}
				fmt.Printf("     (END %d)\n",start)
			}
		}
		
		// Now look at the arrow content
		
		matching_arrows := SST.GetDBArrowsMatchingArrowName(ctx,searchtext)
		
		relns := SST.GetDBNodeArrowNodeMatchingArrowPtrs(ctx,matching_arrows)
		
		for r := range relns {
			
			from := SST.GetDBNodeByNodePtr(ctx,relns[r].NFrom)
			to := SST.GetDBNodeByNodePtr(ctx,relns[r].NFrom)
			//st := relns[r].STType
			arr := SST.ARROW_DIRECTORY[relns[r].Arr].Long
			wgt := relns[r].Wgt
			actx := relns[r].Ctx
			fmt.Println("See also: ",from.S,"--(",arr,")->",to.S,"\n       (... wgt",wgt,"in the contexts",actx,")\n")
			
		}
	}
}









