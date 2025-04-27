package main

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

type Node struct {
	key   int
	value int
	left  *Node
	right *Node
}

type OrderedMap struct {
	root *Node
	size int
}

func NewOrderedMap() OrderedMap {
	return OrderedMap{
		root: nil,
		size: 0,
	}
}

func (m *OrderedMap) Insert(key, value int) {
	if m.root == nil {
		m.root = &Node{key: key, value: value}
		m.size++
		return
	}

	current := m.root
	for {
		if current.key == key {
			current.value = value
			return
		}
		if current.key > key {
			if current.left == nil {
				current.left = &Node{key: key, value: value}
				m.size++
				return
			}
			current = current.left
		} else {
			if current.right == nil {
				current.right = &Node{key: key, value: value}
				m.size++
				return
			}
			current = current.right
		}
	}
}

func (m *OrderedMap) Contains(key int) bool {
	current := m.root
	for current != nil {
		if current.key == key {
			return true
		}
		if current.key > key {
			current = current.left
		} else {
			current = current.right
		}
	}
	return false
}

func (m *OrderedMap) Size() int {
	return m.size
}

func (m *OrderedMap) ForEach(action func(int, int)) {
	m.inOrderTraversal(m.root, action)
}

func (m *OrderedMap) inOrderTraversal(node *Node, action func(int, int)) {
	if node == nil {
		return
	}
	m.inOrderTraversal(node.left, action)
	action(node.key, node.value)
	m.inOrderTraversal(node.right, action)
}

func (m *OrderedMap) Erase(key int) {
	parentSide := false
	parent := m.root
	current := m.root

	for {
		if current.key == key {

			m.size--

			if current.left == nil {
				if parentSide {
					parent.right = current.right
				} else {
					parent.left = current.right
				}

				return
			}
			if current.right == nil {
				if parentSide {
					parent.right = current.left
				} else {
					parent.left = current.left
				}

				return
			}

			// Находим минимальный элемент в правом поддереве
			minNode := current.right
			minNodeParent := current
			for minNode.left != nil {
				minNodeParent = minNode
				minNode = minNode.left
			}

			// Копируем данные минимального элемента
			current.key = minNode.key
			current.value = minNode.value

			// Удаляем минимальный элемент
			if minNodeParent == current {
				minNodeParent.right = nil
			} else {
				minNodeParent.left = nil
			}

			return
		}
		if key < current.key {
			if current.left == nil {
				return
			}
			parentSide = false
			parent = current
			current = current.left
		} else {
			if current.right == nil {
				return
			}
			parentSide = true
			parent = current
			current = current.right
		}
	}
}

func TestCircularQueue(t *testing.T) {
	data := NewOrderedMap()
	assert.Zero(t, data.Size())

	data.Insert(10, 10)
	data.Insert(5, 5)
	data.Insert(15, 15)
	data.Insert(2, 2)
	data.Insert(4, 4)
	data.Insert(12, 12)
	data.Insert(14, 14)

	assert.Equal(t, 7, data.Size())
	assert.True(t, data.Contains(4))
	assert.True(t, data.Contains(12))
	assert.False(t, data.Contains(3))
	assert.False(t, data.Contains(13))

	var keys []int
	expectedKeys := []int{2, 4, 5, 10, 12, 14, 15}
	data.ForEach(func(key, _ int) {
		keys = append(keys, key)
	})

	assert.True(t, reflect.DeepEqual(expectedKeys, keys))

	data.Erase(15)
	data.Erase(14)
	data.Erase(2)

	assert.Equal(t, 4, data.Size())
	assert.True(t, data.Contains(4))
	assert.True(t, data.Contains(12))
	assert.False(t, data.Contains(2))
	assert.False(t, data.Contains(14))

	keys = nil
	expectedKeys = []int{4, 5, 10, 12}
	data.ForEach(func(key, _ int) {
		keys = append(keys, key)
	})

	assert.True(t, reflect.DeepEqual(expectedKeys, keys))
}
