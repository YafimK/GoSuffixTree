package SuffixTree

import "fmt"

func printChildren(t *SuffixNode, pre string) {
	for _, cursorNode := range t.Children {
		children := cursorNode.Children
		if len(children) == 0 {
			fmt.Println("╴", string(cursorNode.Value))
			return
		}
		fmt.Println("┐", string(cursorNode.Value))
		last := len(children) - 1
		for i := 0; i < last; i++ {
			fmt.Print(pre, "├─")
			printChildren(children[i], pre+"  ")
		}
		for _, ch := range children[:last+1] {
			printChildren(ch, pre+"│ ")
		}
		fmt.Print(pre, "└─")
		printChildren(children[last], pre+"  ")
	}
}

func vis(t *SuffixNode) {
	if len(t.Children) == 0 {
		fmt.Println("<empty tree>")
		return
	}
	printChildren(t, "")
}

func main() {
	tree := NewSuffixTree()
	tree.InsertWord([]byte("cgi"))
	testStrings := []string{"cgi", "xcgi", "cgi-bin"}
	for _, testString := range testStrings {
		match, isFound := tree.LookupString([]byte(testString))
		if isFound {

			fmt.Printf("Found match of substring [%v], from string [%v]\n", testString[match.StartIndex:match.EndIndex+1], testString)

		} else {
			fmt.Printf("Didn't find match for string - [%v]\n", testString)
		}

	}
	//vis(&tree.Root)
}
