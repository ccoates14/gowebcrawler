package main


//later lets refactor this to use generics
type Queue []string

func (q *Queue) Enqueue(item string) {
    *q = append(*q, item)
}

func (q *Queue) Dequeue() string {
    if len(*q) == 0 {
        panic("queue is empty")
    }
    item := (*q)[0]
    *q = (*q)[1:]
    return item
}

func (q *Queue) IsEmpty() bool {
    return len(*q) == 0
}