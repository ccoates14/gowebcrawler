package main

import "sync"

//later lets refactor this to use generics
type Queue struct {
    items []string
    lock sync.Mutex
}

var lock sync.Mutex    

func (q *Queue) Enqueue(item string) {
    q.lock.Lock()
    defer q.lock.Unlock()

    q.items = append(q.items, item)
}

func (q *Queue) Dequeue() string {
    q.lock.Lock()
    defer q.lock.Unlock()

    if len(q.items) == 0 {
        panic("queue is empty")
    }
    item := q.items[0]
    q.items = q.items[1:]
    return item
}

func (q *Queue) IsEmpty() bool {
    q.lock.Lock()
    defer q.lock.Unlock()

    return len(q.items) == 0
}