package main

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

type CircularQueue struct {
	values []int
	back   int
	front  int
}

func NewCircularQueue(size int) CircularQueue {
	return CircularQueue{
		values: make([]int, size),
		back:   -1,
		front:  -1,
	}
}

func (q *CircularQueue) Push(value int) bool {
	if q.Full() {
		return false
	}

	if q.back == -1 {
		q.back = 0
		q.front = 0
		q.values[q.back] = value
	} else {
		q.back++
		q.back %= len(q.values)
		q.values[q.back] = value
	}

	return true
}

func (q *CircularQueue) Pop() bool {
	if q.Empty() {
		return false
	}

	if q.back == q.front {
		q.back = -1
		q.front = -1
	} else {
		q.front++
		q.front %= len(q.values)
	}

	return true
}

func (q *CircularQueue) Front() int {
	if q.front == -1 {
		return -1
	}

	return q.values[q.front]
}

func (q *CircularQueue) Back() int {
	if q.back == -1 {
		return -1
	}

	return q.values[q.back]
}

func (q *CircularQueue) Empty() bool {
	return q.back == -1
}

func (q *CircularQueue) Full() bool {
	return (q.back+1)%len(q.values) == q.front
}

func TestCircularQueue(t *testing.T) {
	const queueSize = 3
	queue := NewCircularQueue(queueSize)

	assert.True(t, queue.Empty())
	assert.False(t, queue.Full())

	assert.Equal(t, -1, queue.Front())
	assert.Equal(t, -1, queue.Back())
	assert.False(t, queue.Pop())

	assert.True(t, queue.Push(1))
	assert.True(t, queue.Push(2))
	assert.True(t, queue.Push(3))
	assert.False(t, queue.Push(4))

	assert.True(t, reflect.DeepEqual([]int{1, 2, 3}, queue.values))

	assert.False(t, queue.Empty())
	assert.True(t, queue.Full())

	assert.Equal(t, 1, queue.Front())
	assert.Equal(t, 3, queue.Back())

	assert.True(t, queue.Pop())
	assert.False(t, queue.Empty())
	assert.False(t, queue.Full())
	assert.True(t, queue.Push(4))

	assert.True(t, reflect.DeepEqual([]int{4, 2, 3}, queue.values))

	assert.Equal(t, 2, queue.Front())
	assert.Equal(t, 4, queue.Back())

	assert.True(t, queue.Pop())
	assert.True(t, queue.Pop())
	assert.True(t, queue.Pop())
	assert.False(t, queue.Pop())

	assert.True(t, queue.Empty())
	assert.False(t, queue.Full())
}
