package rbtree

import (
	"math/rand"
	"reflect"
	"testing"
)

func TestNewRBTree(t *testing.T) {
	tree := New[int, string]()
	if tree == nil {
		t.Fatal("New should return a non-nil tree")
	}
	if tree.Len() != 0 {
		t.Errorf("New tree should have length 0, got %d", tree.Len())
	}
}

func TestInsertAndFind(t *testing.T) {
	tree := New[int, string]()

	// Test inserting a single element
	node, inserted := tree.Insert(1, "one")
	if !inserted {
		t.Error("Insert should return true for a new key")
	}
	if node == nil {
		t.Fatal("Insert should return a non-nil node")
	}
	if node.Key() != 1 || node.Value() != "one" {
		t.Errorf("Inserted node has incorrect key/value, got %d/%s, want 1/one", node.Key(), node.Value())
	}
	if tree.Len() != 1 {
		t.Errorf("Tree length should be 1 after insertion, got %d", tree.Len())
	}

	// Test inserting a duplicate key
	node, inserted = tree.Insert(1, "uno")
	if inserted {
		t.Error("Insert should return false for an existing key")
	}
	if node.Value() != "uno" {
		t.Errorf("Insert should update the value for an existing key, got %s, want uno", node.Value())
	}
	if tree.Len() != 1 {
		t.Errorf("Tree length should still be 1 after inserting a duplicate, got %d", tree.Len())
	}

	// Test inserting multiple elements
	for i := 2; i <= 10; i++ {
		_, inserted := tree.Insert(i, "value")
		if !inserted {
			t.Errorf("Failed to insert key %d", i)
		}
	}
	if tree.Len() != 10 {
		t.Errorf("Tree length should be 10 after insertions, got %d", tree.Len())
	}

	// Test finding an existing key
	node = tree.Find(3)
	if node == nil {
		t.Fatal("Find should return a non-nil node for an existing key")
	}
	if node.Key() != 3 || node.Value() != "value" {
		t.Errorf("Find returned incorrect node, got %d/%s, want 3/value", node.Key(), node.Value())
	}

	// Test finding a non-existing key
	node = tree.Find(11)
	if node != nil {
		t.Error("Find should return nil for a non-existing key")
	}
}

func TestRemoveAndErase(t *testing.T) {
	tree := New[int, string]()
	for i := 1; i <= 10; i++ {
		tree.Insert(i, "value")
	}

	// Test removing an existing key
	removed := tree.Remove(3)
	if !removed {
		t.Error("Remove should return true for an existing key")
	}
	if tree.Len() != 9 {
		t.Errorf("Tree length should be 9 after removal, got %d", tree.Len())
	}
	if tree.Contains(3) {
		t.Error("Tree should not contain removed key")
	}

	// Test removing a non-existing key
	removed = tree.Remove(11)
	if removed {
		t.Error("Remove should return false for a non-existing key")
	}
	if tree.Len() != 9 {
		t.Errorf("Tree length should still be 9 after attempting to remove a non-existing key, got %d", tree.Len())
	}

	// Test removing the root
	removed = tree.Remove(1)
	if !removed {
		t.Error("Remove should return true when removing the root")
	}
	if tree.Len() != 8 {
		t.Errorf("Tree length should be 8 after removing the root, got %d", tree.Len())
	}
	if tree.Contains(1) {
		t.Error("Tree should not contain removed root key")
	}

	// Test Erase
	node := tree.Find(5)
	erased := tree.Erase(node)
	if !erased {
		t.Error("Erase should return true for an existing node")
	}
	if tree.Contains(5) {
		t.Error("Tree should not contain erased node")
	}

	// Test erasing a nil node
	erased = tree.Erase(nil)
	if erased {
		t.Error("Erase should return false for a nil node")
	}

	// Test erasing a node not in the tree
	anotherTree := New[int, string]()
	anotherTree.Insert(100, "hundred")
	node = anotherTree.Find(100)
	erased = tree.Erase(node)
	if erased {
		t.Error("Erase should return false for a node not in the tree")
	}
}

func TestClear(t *testing.T) {
	tree := New[int, string]()
	for i := 1; i <= 5; i++ {
		tree.Insert(i, "value")
	}

	tree.Clear()
	if tree.Len() != 0 {
		t.Errorf("Tree length should be 0 after Clear, got %d", tree.Len())
	}
	if tree.Root() != nil {
		t.Error("Tree root should be nil after Clear")
	}
}

func TestKeysAndValues(t *testing.T) {
	tree := New[int, string]()
	expectedKeys := []int{1, 2, 3, 4, 5}
	expectedValues := []string{"one", "two", "three", "four", "five"}

	for i, key := range expectedKeys {
		tree.Insert(key, expectedValues[i])
	}

	keys := tree.Keys()
	if !reflect.DeepEqual(keys, expectedKeys) {
		t.Errorf("Keys() returned incorrect keys, got %v, want %v", keys, expectedKeys)
	}

	values := tree.Values()
	if !reflect.DeepEqual(values, expectedValues) {
		t.Errorf("Values() returned incorrect values, got %v, want %v", values, expectedValues)
	}
}

func TestIterationOrder(t *testing.T) {
	tree := New[int, string]()
	for i := 1; i <= 10; i++ {
		tree.Insert(i, "value")
	}

	// Test forward iteration
	node := tree.First()
	for i := 1; i <= 10; i++ {
		if node == nil {
			t.Fatalf("Iteration ended prematurely at %d", i)
		}
		if node.Key() != i {
			t.Errorf("Incorrect key during forward iteration, got %d, want %d", node.Key(), i)
		}
		node = node.Next()
	}
	if node != nil {
		t.Error("Iteration should end with nil")
	}

	// Test backward iteration
	node = tree.Last()
	for i := 10; i >= 1; i-- {
		if node == nil {
			t.Fatalf("Backward iteration ended prematurely at %d", i)
		}
		if node.Key() != i {
			t.Errorf("Incorrect key during backward iteration, got %d, want %d", node.Key(), i)
		}
		node = node.Prev()
	}
	if node != nil {
		t.Error("Backward iteration should end with nil")
	}
}

func TestCustomComparator(t *testing.T) {
	tree := NewFunc[int, string](func(a, b int) bool { return a > b }) // Reverse order

	for i := 1; i <= 5; i++ {
		tree.Insert(i, "value")
	}

	node := tree.First()
	for i := 5; i >= 1; i-- {
		if node == nil {
			t.Fatalf("Iteration ended prematurely at %d", i)
		}
		if node.Key() != i {
			t.Errorf("Incorrect key with custom comparator, got %d, want %d", node.Key(), i)
		}
		node = node.Next()
	}
}

func TestRedBlackProperties(t *testing.T) {
	tree := New[int, string]()
	for i := 1; i <= 1000; i++ {
		tree.Insert(rand.Intn(10000), "value")
	}

	if !isValidRedBlackTree(tree.root) {
		t.Error("Tree does not satisfy Red-Black properties")
	}
}

func TestGet(t *testing.T) {
	tree := New[int, string]()
	tree.Insert(1, "one")
	tree.Insert(2, "two")

	value := tree.Get(1)
	if value != "one" {
		t.Errorf("Expected 'one', got '%s'", value)
	}

	value = tree.Get(3)
	if value != "" {
		t.Errorf("Expected empty string for non-existing key, got '%s'", value)
	}
}

func TestRotations(t *testing.T) {
	tree := New[int, string]()

	// Test left rotation
	tree.Insert(1, "one")
	tree.Insert(2, "two")
	tree.Insert(3, "three")
	if tree.root.key != 2 {
		t.Errorf("Expected root key 2 after left rotation, got %d", tree.root.key)
	}

	tree.Clear()

	// Test right rotation
	tree.Insert(3, "three")
	tree.Insert(2, "two")
	tree.Insert(1, "one")
	if tree.root.key != 2 {
		t.Errorf("Expected root key 2 after right rotation, got %d", tree.root.key)
	}
}

func TestNodeMethods(t *testing.T) {
	tree := New[int, string]()
	tree.Insert(2, "two")
	tree.Insert(1, "one")
	tree.Insert(3, "three")

	root := tree.Root()

	if root.Parent() != nil {
		t.Error("Root should not have a parent")
	}
	if root.left.Parent() != root {
		t.Error("Left child's parent should be root")
	}

	if root.NumChild() != 2 {
		t.Errorf("Expected 2 children for root, got %d", root.NumChild())
	}
	if root.left.NumChild() != 0 {
		t.Errorf("Expected 0 children for leaf, got %d", root.left.NumChild())
	}

	if root.GetChildByIndex(0) != root.left {
		t.Error("GetChildByIndex(0) should return left child")
	}
	if root.GetChildByIndex(1) != root.right {
		t.Error("GetChildByIndex(1) should return right child")
	}

	root.SetValue("new_two")
	if root.Value() != "new_two" {
		t.Errorf("SetValue failed, expected 'new_two', got '%s'", root.Value())
	}
}

func TestEmptyTreeOperations(t *testing.T) {
	tree := New[int, string]()

	if tree.Remove(1) {
		t.Error("Remove should return false on empty tree")
	}

	if tree.Erase(nil) {
		t.Error("Erase should return false on empty tree")
	}

	if tree.Find(1) != nil {
		t.Error("Find should return nil on empty tree")
	}

	if tree.First() != nil {
		t.Error("First should return nil on empty tree")
	}

	if tree.Last() != nil {
		t.Error("Last should return nil on empty tree")
	}
}

func TestTreeModificationDuringTraversal(t *testing.T) {
	tree := New[int, string]()
	for i := 1; i <= 10; i++ {
		tree.Insert(i, "value")
	}

	node := tree.First()
	for node != nil {
		key := node.Key()
		next := node.Next()
		tree.Remove(key)
		node = next
	}

	if tree.Len() != 0 {
		t.Errorf("Expected empty tree after removal during traversal, got length %d", tree.Len())
	}
}

func isValidRedBlackTree[K comparable, V any](node *Node[K, V]) bool {
	if node == nil || node.null() {
		return true
	}

	if node.color != red && node.color != black {
		return false
	}

	if node.parent == nil && node.color != black {
		return false
	}

	if node.color == red {
		if !node.left.null() && node.left.color != black {
			return false
		}
		if !node.right.null() && node.right.color != black {
			return false
		}
	}

	leftBlackCount := countBlackNodes(node.left)
	rightBlackCount := countBlackNodes(node.right)
	if leftBlackCount != rightBlackCount {
		return false
	}

	return isValidRedBlackTree(node.left) && isValidRedBlackTree(node.right)
}

func countBlackNodes[K comparable, V any](node *Node[K, V]) int {
	if node == nil || node.null() {
		return 1 // Null nodes are considered black
	}
	count := countBlackNodes(node.left)
	if node.color == black {
		count++
	}
	return count
}

func BenchmarkInsert(b *testing.B) {
	tree := New[int, int]()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tree.Insert(i, i)
	}
}

func BenchmarkFind(b *testing.B) {
	tree := New[int, int]()
	for i := 0; i < 1000000; i++ {
		tree.Insert(i, i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tree.Find(rand.Intn(1000000))
	}
}

func BenchmarkRemove(b *testing.B) {
	tree := New[int, int]()
	for i := 0; i < 1000000; i++ {
		tree.Insert(i, i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tree.Remove(rand.Intn(1000000))
	}
}

func TestRemoveRootWithChildren(t *testing.T) {
	tree := New[int, string]()
	tree.Insert(2, "two")
	tree.Insert(1, "one")
	tree.Insert(3, "three")

	removed := tree.Remove(2)
	if !removed {
		t.Error("Failed to remove root with children")
	}
	if tree.Len() != 2 {
		t.Errorf("Expected tree length 2 after root removal, got %d", tree.Len())
	}
	if tree.Contains(2) {
		t.Error("Tree should not contain removed root key")
	}
	if !isValidRedBlackTree(tree.root) {
		t.Error("Tree is not valid after removing root with children")
	}
}

func TestRemoveWithRebalancing(t *testing.T) {
	tree := New[int, string]()
	for i := 1; i <= 7; i++ {
		tree.Insert(i, "value")
	}

	// Remove nodes to force rebalancing
	tree.Remove(2)
	tree.Remove(4)
	tree.Remove(6)

	if !isValidRedBlackTree(tree.root) {
		t.Error("Tree is not valid after removals requiring rebalancing")
	}
}

func TestNodeNullMethods(t *testing.T) {
	tree := New[int, string]()
	tree.Insert(1, "one")

	nullNode := tree.root.left // Assuming left child is null

	if !nullNode.null() {
		t.Error("Null node should return true for null() method")
	}

	if nullNode.NumChild() != 0 {
		t.Error("Null node should have 0 children")
	}

	if nullNode.GetChildByIndex(0) != nil {
		t.Error("GetChildByIndex on null node should return nil")
	}
}

func TestPrevAndNextEdgeCases(t *testing.T) {
	tree := New[int, string]()
	tree.Insert(2, "two")
	tree.Insert(1, "one")
	tree.Insert(3, "three")

	firstNode := tree.First()
	if firstNode.Prev() != nil {
		t.Error("Prev of first node should be nil")
	}

	lastNode := tree.Last()
	if lastNode.Next() != nil {
		t.Error("Next of last node should be nil")
	}
}

func TestRemoveNonExistentNode(t *testing.T) {
	tree := New[int, string]()
	tree.Insert(1, "one")

	removed := tree.Remove(2)
	if removed {
		t.Error("Remove should return false for non-existent key")
	}

	if tree.Len() != 1 {
		t.Errorf("Tree length should remain 1, got %d", tree.Len())
	}
}

func TestGetChildByIndexPanic(t *testing.T) {
	tree := New[int, string]()
	tree.Insert(1, "one")
	root := tree.Root()

	defer func() {
		if r := recover(); r == nil {
			t.Error("GetChildByIndex should panic for invalid index")
		}
	}()

	root.GetChildByIndex(2) // This should panic
}

func TestRemoveWithColorChanges(t *testing.T) {
	tree := New[int, string]()
	for i := 1; i <= 5; i++ {
		tree.Insert(i, "value")
	}

	// Remove a node that will cause color changes
	tree.Remove(2)

	if !isValidRedBlackTree(tree.root) {
		t.Error("Tree is not valid after removal causing color changes")
	}
}

func TestRemoveComplexCases(t *testing.T) {
	tree := New[int, string]()

	// Case 1: Remove when right child is not null
	tree.Insert(2, "two")
	tree.Insert(1, "one")
	tree.Insert(3, "three")
	tree.Insert(4, "four")

	removed := tree.Remove(2)
	if !removed {
		t.Error("Failed to remove node with non-null right child")
	}
	if tree.Contains(2) {
		t.Error("Tree should not contain removed key 2")
	}
	if !isValidRedBlackTree(tree.root) {
		t.Error("Tree is not valid after removing node with non-null right child")
	}

	// Case 2: Remove root when it's the only node
	tree.Clear()
	tree.Insert(1, "one")
	removed = tree.Remove(1)
	if !removed {
		t.Error("Failed to remove root when it's the only node")
	}
	if tree.root != nil {
		t.Error("Tree root should be nil after removing the only node")
	}

	// Case 3: Remove causing red sibling case
	tree.Clear()
	tree.Insert(5, "five")
	tree.Insert(2, "two")
	tree.Insert(7, "seven")
	tree.Insert(1, "one")
	tree.Insert(3, "three")
	tree.Insert(6, "six")
	tree.Insert(8, "eight")

	removed = tree.Remove(1) // This should cause the red sibling case
	if !removed {
		t.Error("Failed to remove node causing red sibling case")
	}
	if !isValidRedBlackTree(tree.root) {
		t.Error("Tree is not valid after removal causing red sibling case")
	}
}

func TestDoRemoveComplexCases(t *testing.T) {
	tree := New[int, string]()

	// Setup a tree structure that will trigger complex removal cases
	tree.Insert(5, "five")
	tree.Insert(2, "two")
	tree.Insert(7, "seven")
	tree.Insert(1, "one")
	tree.Insert(3, "three")
	tree.Insert(6, "six")
	tree.Insert(8, "eight")

	// Case 1: Remove causing black sibling with red child
	removed := tree.Remove(1)
	if !removed {
		t.Error("Failed to remove node causing black sibling with red child case")
	}
	if !isValidRedBlackTree(tree.root) {
		t.Error("Tree is not valid after removal causing black sibling with red child case")
	}

	// Case 2: Remove causing black sibling with black children
	removed = tree.Remove(2)
	if !removed {
		t.Error("Failed to remove node causing black sibling with black children case")
	}
	if !isValidRedBlackTree(tree.root) {
		t.Error("Tree is not valid after removal causing black sibling with black children case")
	}

	// Case 3: Remove causing red parent with black sibling and black children
	tree.Clear()
	tree.Insert(5, "five")
	tree.Insert(3, "three")
	tree.Insert(7, "seven")
	tree.Insert(1, "one")
	tree.Insert(4, "four")
	tree.Insert(6, "six")
	tree.Insert(8, "eight")

	removed = tree.Remove(1)
	if !removed {
		t.Error("Failed to remove node causing red parent with black sibling and black children case")
	}
	if !isValidRedBlackTree(tree.root) {
		t.Error("Tree is not valid after removal causing red parent with black sibling and black children case")
	}
}

func TestEraseNilNode(t *testing.T) {
	tree := New[int, string]()
	tree.Insert(1, "one")

	// Attempt to erase a nil node
	erased := tree.Erase(nil)
	if erased {
		t.Error("Erase should return false for nil node")
	}
	if tree.Len() != 1 {
		t.Errorf("Tree length should remain 1, got %d", tree.Len())
	}
}

func TestEraseNodeFromDifferentTree(t *testing.T) {
	tree1 := New[int, string]()
	tree1.Insert(1, "one")

	tree2 := New[int, string]()
	tree2.Insert(2, "two")

	node := tree2.Find(2)

	// Attempt to erase a node from a different tree
	erased := tree1.Erase(node)
	if erased {
		t.Error("Erase should return false for node from different tree")
	}
	if tree1.Len() != 1 {
		t.Errorf("Tree1 length should remain 1, got %d", tree1.Len())
	}
	if tree2.Len() != 1 {
		t.Errorf("Tree2 length should remain 1, got %d", tree2.Len())
	}
}

func TestEmptyTreeKeysAndValues(t *testing.T) {
	tree := New[int, string]()

	keys := tree.Keys()
	if keys != nil {
		t.Errorf("Expected nil slice for Keys() on empty tree, got %v", keys)
	}

	values := tree.Values()
	if values != nil {
		t.Errorf("Expected nil slice for Values() on empty tree, got %v", values)
	}
}

func TestRemoveRootWithOneChild(t *testing.T) {
	tree := New[int, string]()
	tree.Insert(1, "one")
	tree.Insert(2, "two")

	removed := tree.Remove(1)
	if !removed {
		t.Error("Failed to remove root with one child")
	}
	if tree.root.key != 2 {
		t.Errorf("Expected new root key to be 2, got %d", tree.root.key)
	}
	if tree.root.color != black {
		t.Error("New root should be black")
	}
}

func TestDoRemoveComplexCases2(t *testing.T) {
	tree := New[int, string]()
	tree.Insert(5, "five")
	tree.Insert(2, "two")
	tree.Insert(7, "seven")
	tree.Insert(1, "one")
	tree.Insert(3, "three")
	tree.Insert(6, "six")
	tree.Insert(8, "eight")

	// Case: Red sibling
	tree.Remove(1)
	if !isValidRedBlackTree(tree.root) {
		t.Error("Tree is not valid after removal with red sibling case")
	}

	// Case: Black sibling with black children
	tree.Remove(2)
	if !isValidRedBlackTree(tree.root) {
		t.Error("Tree is not valid after removal with black sibling and black children case")
	}

	// Case: Black sibling with at least one red child
	tree.Insert(4, "four")
	tree.Remove(3)
	if !isValidRedBlackTree(tree.root) {
		t.Error("Tree is not valid after removal with black sibling and red child case")
	}
}

func TestNodeRelationships(t *testing.T) {
	tree := New[int, string]()
	tree.Insert(5, "five")
	tree.Insert(3, "three")
	tree.Insert(7, "seven")
	tree.Insert(2, "two")
	tree.Insert(4, "four")
	tree.Insert(6, "six")
	tree.Insert(8, "eight")

	root := tree.Root()
	leftChild := root.GetChildByIndex(0)
	rightChild := root.GetChildByIndex(1)

	// Test grandparent
	if leftChild.grandparent() != nil {
		t.Error("Grandparent of root's child should be nil")
	}
	if leftChild.left.grandparent() != root {
		t.Error("Incorrect grandparent relationship")
	}

	// Test sibling
	if leftChild.sibling() != rightChild {
		t.Error("Incorrect sibling relationship")
	}
	if root.sibling() != nil {
		t.Error("Root should not have a sibling")
	}

	// Test uncle
	if leftChild.left.uncle() != rightChild {
		t.Error("Incorrect uncle relationship")
	}
	if root.uncle() != nil {
		t.Error("Root should not have an uncle")
	}
}

func TestRemoveRootWithSingleLeftChild(t *testing.T) {
	rbtree := New[int, string]()

	// Insert root
	rbtree.Insert(10, "ten")

	// Insert left child
	rbtree.Insert(5, "five")

	// Verify initial structure
	if rbtree.root.key != 10 || rbtree.root.left.key != 5 || !rbtree.root.right.null() {
		t.Fatalf(
			"Initial tree structure is not as expected, root=%v",
			rbtree.root,
		)
	}

	// Remove the root
	removed := rbtree.Remove(10)

	if !removed {
		t.Error("Failed to remove root")
	}

	// Check new root
	if rbtree.root == nil {
		t.Fatal("Tree root should not be nil after removal")
	}

	if rbtree.root.key != 5 {
		t.Errorf("New root key should be 5, got %d", rbtree.root.key)
	}

	if rbtree.root.color != black {
		t.Error("New root should be black")
	}

	if rbtree.root.parent != nil {
		t.Error("New root's parent should be nil")
	}

	if rbtree.Len() != 1 {
		t.Errorf("Tree length should be 1 after removal, got %d", rbtree.Len())
	}

	// Verify tree properties
	if !isValidRedBlackTree(rbtree.root) {
		t.Error("Tree is not valid after removing root with single left child")
	}
}
