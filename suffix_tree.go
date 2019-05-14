package main

// Impl' of Ukkonenes algo to construct suffix tree in linear time
import (
	"fmt"
	"strings"
)

type SuffixNode struct {
	Id          int
	Children    []*SuffixNode
	StartIndex  int
	EndIndex    int
	Value       []byte
	Parent      *SuffixNode
	nodeCounter int
	cursorIndex int
	activeChild *SuffixNode
}

func (node *SuffixNode) SetId(id int) {
	node.Id = id
}

func (node *SuffixNode) SetValue(value []byte) {
	node.Value = value
}

func (node *SuffixNode) AddChild(newNode *SuffixNode) {
	node.Children = append(node.Children, newNode)
}

// Only increase end index for leaves.
func (node *SuffixNode) IncrementEndIndex() {
	if len(node.Children) > 0 {
		for _, branch := range node.Children {
			branch.IncrementEndIndex()
		}
	} else {
		node.EndIndex++
	}
}
func (node *SuffixNode) InsertChar(char byte) {
	node.Value = append(node.Value, char)
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

func (node *SuffixNode) SplitNode(activeIndex, currentIndex int) {
	node.AddChild(&SuffixNode{
		Id:         node.Id,
		StartIndex: activeIndex + 1,
		EndIndex:   currentIndex + 1,
	})

	node.EndIndex = activeIndex
	node.SetId(node.GetNewNodeId())

	node.AddChild(&SuffixNode{
		Id:         node.GetNewNodeId(),
		StartIndex: currentIndex,
		EndIndex:   currentIndex + 1,
	})
}

func FindInsertionBranch(startNode *SuffixNode, currentIndex int, currentWord []byte) (*SuffixNode, bool) {
	for _, cursorNode := range startNode.Children {
		if currentWord[currentIndex] == currentWord[cursorNode.StartIndex] {
			return cursorNode, true
		}
	}
	return nil, false
}

func (tree *SuffixTree) InsertWord(word []byte) {
	node := &tree.root
	for i := 0; i < len(word); i++ {
		isSuitablePrefixNodeFound := false
		if tree.currentEdge != nil {
			if word[i] == word[tree.currentEdge.StartIndex+tree.activeLength] {
				isSuitablePrefixNodeFound = true
				tree.activeLength++
			}
			tree.reminder++
		} else {
			var cursorNode *SuffixNode
			cursorNode, isSuitablePrefixNodeFound = FindInsertionBranch(tree.activeNode, i, word)
			if isSuitablePrefixNodeFound {
				tree.currentEdge = cursorNode
				tree.currentSuffix = word[i]
				tree.activeLength++
				tree.reminder++
				continue
			}
			if !isSuitablePrefixNodeFound {
				activeNode := &SuffixNode{
					Id:         node.GetNewNodeId(),
					StartIndex: i,
					EndIndex:   i + 1,
				}
				node.AddChild(activeNode)
				fmt.Printf("new node #%v - \"%v\" \n", activeNode.Id, word[i:i+1])
			}
		}

		if tree.currentEdge != nil && !isSuitablePrefixNodeFound {
			for tree.activeLength > 0 {

				tree.activeNode.activeChild.SplitNode(tree.activeIndex, i+1)
				tree.reminder--
				tree.activeLength--
				tree.currentSuffix = word[i-tree.activeLength]
				tree.activeNode.activeChild = nil
			}
		}

		// Create new suffix branch if we don't have branch for this suffix

		for _, cursorNode := range root.Children {
			cursorNode.IncrementEndIndex()
		}

	}
}

func (tree SuffixTree) String() string {
	return fmt.Sprintln(tree.root)
}

func NewSuffixTree() *SuffixTree {
	newNode := SuffixNode{}
	newNode.SetId(newNode.GetNewNodeId())
	return &SuffixTree{
		root:          newNode,
		activeNode:    &newNode,
		reminder:      1,
		activeLength:  0,
		currentSuffix: byte(nil),
	}
}

func main() {
	tree := NewSuffixTree()
	tree.InsertWord([]byte("abcabxabcd"))
	fmt.Println(tree)
}
