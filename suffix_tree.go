package SuffixTree

import (
    "fmt"
    "strings"
)

type SuffixNode struct {
    Id        int
    Children  []*SuffixNode
    Value     byte
    Parent    *SuffixNode
    IsWordEnd bool
}

func NewSuffixNode(id int, value byte, parent *SuffixNode) *SuffixNode {
    return &SuffixNode{Id: id, Value: value, Parent: parent}
}

func (node *SuffixNode) SetId(id int) {
    node.Id = id
}

func (node *SuffixNode) AddChild(newNode *SuffixNode) {
    node.Children = append(node.Children, newNode)
}

func (node SuffixNode) String() string {
    sb := strings.Builder{}
    for _, branch := range node.Children {
        sb.WriteString(fmt.Sprintf("#%v: {\n%v\n}\n", branch.Id, branch))
    }
    return sb.String()
}

type SuffixTree struct {
    Root SuffixNode
    Size int
}

func (tree *SuffixTree) GetSize() int {
    return tree.Size
}

func FindInsertionBranch(rootNode *SuffixNode, word []byte) (cursorNode *SuffixNode, index int, isFound bool) {
    cursorNode = rootNode
    for index = 0; index < len(word); index++ {
        cursorChar := word[index]
        isFound, cursorNode = queryChildren(cursorNode, cursorChar)
        if !isFound {
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

//LookupString Looks for a string that matches the given searchString in the tree event if searchString is a substring of one of the strings in the tree
//
func (tree *SuffixTree) LookupString(searchString string) (Match, bool) {
    byteSearchString := []byte(searchString)
    node, endIndex, isFound := LookupWord(&tree.Root, byteSearchString)
    matches := Match{
        node, 0, endIndex,
    }
    return matches, isFound

}

//LookupFullString Looks for the whole search string as one word, i.e. if the <search string> wasn't inserted as a whole
//word to the tree. It won't find any match.
func (tree *SuffixTree) LookupFullString(searchString string) (*Match, bool) {
    byteSearchString := []byte(searchString)
    node, endIndex, isFound := LookupWord(&tree.Root, byteSearchString)
    if isFound && node.IsWordEnd {
        matches := &Match{
            node, 0, endIndex,
        }
        return matches, isFound
    }
    return nil, false
}

//LookupSubString Looks for any substring of the given search string. i.e. if any suffix of <search string> appears in the tree it will return as a match.
func (tree *SuffixTree) LookupSubString(searchString string) (matches []Match, isFound bool) {
    byteSearchString := []byte(searchString)
    for startIndex := 0; startIndex < len(byteSearchString); startIndex++ {
        cursorWord := byteSearchString[startIndex:]
        node, endIndexOffset, isFound := LookupWord(&tree.Root, cursorWord)
        if isFound {
            match := Match{
                node, startIndex, startIndex + endIndexOffset,
            }
            matches = append(matches, match)
        }
    }

    return matches, len(matches) > 0
}

//LookupMaxContinuousSubStrings Looks for any sub string of the given search string, but returns only the max continuous substring found.
// i.e. if matching strings with
func (tree *SuffixTree) LookupMaxContinuousSubStrings(searchString string) ([]Match, bool) {
    byteSearchString := []byte(searchString)
    var matches []Match
    for startIndex := 0; startIndex < len(byteSearchString); startIndex++ {
        cursorWord := byteSearchString[startIndex:]
        cursorNode := &tree.Root
        index := 0
        isFound := false

        cursorNode, index, isFound = LookupWord(cursorNode, cursorWord)

        if isFound || index > 0 {
            match := Match{
                cursorNode, startIndex, startIndex + index,
            }
            matches = append(matches, match)
        }
    }
    if len(matches) == 0 {
        return nil, false
    }
    currentMaxLength := -1
    var maxMatches []Match
    for _, item := range matches {
        length := item.EndIndex - item.StartIndex
        if length > currentMaxLength {
            currentMaxLength = length
            maxMatches = []Match{}
            maxMatches = append(maxMatches, item)
        } else if length == currentMaxLength {
            maxMatches = append(maxMatches, item)
        }
    }
    return maxMatches, true
}

func LookupWord(rootNode *SuffixNode, word []byte) (*SuffixNode, int, bool) {
    node, maxMatchEndIndex, isFound := FindInsertionBranch(rootNode, word)
    if isFound {
        return node, maxMatchEndIndex, true
    }
    return node, maxMatchEndIndex, false
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

func (tree *SuffixTree) InsertWord(word []byte) *SuffixTree {
    node := &tree.Root
    for wordIndex := 0; wordIndex < len(word); wordIndex++ {
        cursorWord := word[wordIndex:]
        insertionNode, index, isFound := FindInsertionBranch(node, cursorWord)
        if isFound {
            continue
        }
        //Now we have the tip of the insertion edge, and we know that the word doesn't exist
        cursorNode := insertionNode
        for _, char := range cursorWord[index:] {
            newNode := NewSuffixNode(tree.Size, char, cursorNode)
            cursorNode.AddChild(newNode)
            cursorNode = newNode
            tree.Size++
        }
        if wordIndex == 0 {
            cursorNode.IsWordEnd = true
        }
    }
    return tree
}

//InsertFullWord Allows the suffix tree to become a trie tree containing only full words.
//We get less space usage then hash map by using similar prefixes
func (tree *SuffixTree) InsertFullWord(word []byte) *SuffixTree {
    node := &tree.Root
    insertionNode, index, isFound := FindInsertionBranch(node, word)
    if isFound {
        return tree
    }
    //Now we have the tip of the insertion edge, and we know that the word doesn't exist
    cursorNode := insertionNode
    for _, char := range word[index:] {
        newNode := NewSuffixNode(tree.Size, char, cursorNode)

        cursorNode.AddChild(newNode)
        cursorNode = newNode
        tree.Size++
    }
    cursorNode.IsWordEnd = true
    return tree
}

func NewSuffixTree() *SuffixTree {
    newNode := SuffixNode{}
    tree := &SuffixTree{
        Root: newNode,
    }
    tree.Size++
    return tree
}
