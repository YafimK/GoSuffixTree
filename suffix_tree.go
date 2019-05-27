package SuffixTree

import (
	"fmt"
	"strings"
)

type SuffixNode struct {
	Id          int
	Children    []*SuffixNode
	Value       byte
	Parent      *SuffixNode
	NodeCounter int
	IsWordEnd   bool
}

func (node *SuffixNode) SetId(id int) {
	node.Id = id
}

func (node *SuffixNode) AddChild(newNode *SuffixNode) {
	node.Children = append(node.Children, newNode)
	node.NodeCounter++
}

func (node *SuffixNode) GetNewNodeId() int {
	var newId int
	if node.Parent == nil { //this is the Root
		newId = node.NodeCounter
		node.NodeCounter++
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
	Root SuffixNode
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

func (tree *SuffixTree) LookupString(searchString []byte) (Match, bool) {
	node, startIndex, endIndex, isFound := LookupWord(&tree.Root, searchString)
	matches := Match{
		node, startIndex, endIndex,
	}
	return matches, isFound

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

func (tree *SuffixTree) InsertWord(word []byte) {
	node := &tree.Root
	for wordIndex := 0; wordIndex < len(word); {
		cursorWord := word[wordIndex:]
		insertionNode, index, isFound := FindInsertionBranch(node, cursorWord)
		if !isFound {
			//Now we have the tip of the insertion edge, and we know that the word doesn't exist
			cursorNode := insertionNode
			for _, char := range cursorWord[index:] {
				newNode := &SuffixNode{
					Id:        node.GetNewNodeId(),
					Value:     char,
					IsWordEnd: false,
					Parent:    cursorNode,
				}
				cursorNode.AddChild(newNode)
				cursorNode = newNode
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
		Root: newNode,
	}
}
