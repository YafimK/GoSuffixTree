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
			//We found the leaf node where the suffix should be inserted.
			break
		}
	}
	return cursorNode, index, isFound
}

func (tree *SuffixTree) LookupWord(word []byte) (*SuffixNode, int, int, bool) {
	return LookupWord(&tree.root, word)
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

func main() {
	tree := NewSuffixTree()
	tree.InsertString([]byte("cgi"))
	tree.InsertString([]byte("rtcgi"))
	tree.InsertString([]byte("cgigi"))
	testStrings := []string{"cgi", "xcgi", "cgi-bin"}
	for _, testString := range testStrings {
		word, startIndex, endIndex, isFound := tree.LookupWord([]byte(testString))
		if isFound {
			fmt.Printf("Found match of substring [%v], from string [%v]; in node [%v]\n", testString[startIndex:endIndex+1], testString, word.Id)
		} else {
			fmt.Printf("Didn't find match for string - [%v]\n", testString)
		}

	}
}
