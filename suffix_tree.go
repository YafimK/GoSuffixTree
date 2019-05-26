package main

import (
	"fmt"
	"regexp"
	"strings"
)

var WordSegmentationRegex = `[\pL\p{Mc}\p{Mn}']+`

type SuffixNode struct {
	Id          int
	Children    []*SuffixNode
	Value       byte
	Parent      *SuffixNode
	nodeCounter int
	cursorIndex int
	activeChild *SuffixNode
	IsWordEnd   bool
}

func (node *SuffixNode) SetId(id int) {
	node.Id = id
}

func (node *SuffixNode) SetValue(value byte) {
	node.Value = value
}

func (node *SuffixNode) AddChild(newNode *SuffixNode) {
	node.Children = append(node.Children, newNode)
	node.nodeCounter++
}

func (node *SuffixNode) GetNewNodeId() int {
	var newId int
	if node.Parent == nil { //this is the root
		newId = node.nodeCounter
		node.nodeCounter++
	} else {
		newId = node.Parent.GetNewNodeId()
	}
	return newId
}

//
func (node SuffixNode) String() string {
	sb := strings.Builder{}
	for _, branch := range node.Children {
		sb.WriteString(fmt.Sprintf("#%v: {\n%v\n}\n", branch.Id, branch))
	}
	return sb.String()
}

type SuffixTree struct {
	root          SuffixNode
	activeNode    *SuffixNode
	activeChild   int
	activeIndex   int
	activeLength  int
	reminder      int
	currentSuffix byte
	currentEdge   *SuffixNode
}

func (tree *SuffixTree) WalkTree() *SuffixNode {

	return nil
}

func FindInsertionBranch(rootNode *SuffixNode, word []byte) (cursorNode *SuffixNode, index int, isFound bool) {
	var char byte
	cursorNode = rootNode
	for index, char = range word {
		isFound, cursorNode = queryChildren(cursorNode, char)
		if isFound == false {
			//We found the leaf Node where the suffix should be inserted.
			break
		}
	}
	return cursorNode, index, isFound
}

type Match struct {
	Node       *SuffixNode
	StartIndex int
	EndIndex   int
}

func (tree *SuffixTree) LookupString(searchString []byte) (matches []Match, isFound bool) {
	for _, word := range SegmentString(searchString) {
		node, startIndex, endIndex, isFound := LookupWord(&tree.root, word)
		if isFound {
			matches = append(matches, Match{
				node, startIndex, endIndex,
			})
		}
	}
	return matches, len(matches) > 0

}

func LookupWord(rootNode *SuffixNode, word []byte) (*SuffixNode, int, int, bool) {

	for startIndex := 0; startIndex < len(word); startIndex++ {
		cursorWord := word[startIndex:]
		node, index, isFound := FindInsertionBranch(rootNode, cursorWord)
		if isFound && node.IsWordEnd {
			return node, startIndex, startIndex + index, true
		}
	}
	return nil, -1, -1, false
}

func queryChildren(cursorNode *SuffixNode, char byte) (bool, *SuffixNode) {
	for _, node := range cursorNode.Children {
		if node.Value == char {
			cursorNode = node
			return true, cursorNode
		}
	}
	return false, cursorNode
}

func SegmentString(word []byte) [][]byte {
	reg := regexp.MustCompile(WordSegmentationRegex)
	return reg.FindAll(word, -1)
}

//Inserts several words - if common separators are found
func (tree *SuffixTree) InsertString(word []byte) {
	for _, segmentedWord := range SegmentString(word) {
		tree.InsertWord(segmentedWord)
	}
}

func (tree *SuffixTree) InsertWord(word []byte) {
	node := &tree.root
	for wordIndex := 0; wordIndex < len(word); {
		cursorWord := word[wordIndex:]
		insertionNode, index, isFound := FindInsertionBranch(node, cursorWord)
		if !isFound {
			//Now we have the tip of the insertion edge, and we know that the word doesn't exist
			cursorNode := insertionNode
			for _, char := range cursorWord[index:] {
				activeNode := &SuffixNode{
					Id:        node.GetNewNodeId(),
					Value:     char,
					IsWordEnd: false,
					Parent:    cursorNode,
				}
				cursorNode.AddChild(activeNode)
				cursorNode = activeNode
			}
			cursorNode.IsWordEnd = true
			wordIndex += index
		} else {
			wordIndex++
		}
	}

}

func NewSuffixTree() *SuffixTree {
	newNode := SuffixNode{}
	newNode.SetId(newNode.GetNewNodeId())
	return &SuffixTree{
		root:       newNode,
		activeNode: &newNode,
	}
}

func printChildren(t *SuffixNode, pre string) {
	for _, cursorNode := range t.Children {
		children := cursorNode.Children
		if len(children) == 0 {
			fmt.Println("╴", string(cursorNode.Value))
			return
		}
		fmt.Println("┐", string(cursorNode.Value))
		last := len(children) - 1
		for _, ch := range children[:last] {
			fmt.Print(pre, "├─")
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
	tree.InsertString([]byte("boris tries to take over the world"))
	tree.InsertString([]byte("yafim kazak tries too"))
	tree.InsertString([]byte("Kobi has already succedd"))
	testStrings := []string{"cgi", "xcgi", "cgi-bin"}
	for _, testString := range testStrings {
		matches, isFound := tree.LookupString([]byte(testString))
		if isFound {
			for _, match := range matches {
				fmt.Printf("Found match of substring [%v], from string [%v]; in Node [%v]\n", testString[match.StartIndex:match.EndIndex+1], testString, match.Node.Id)
			}
		} else {
			fmt.Printf("Didn't find match for string - [%v]\n", testString)
		}

	}
	vis(&tree.root)
}
