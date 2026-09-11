package slru

type Element struct {
	Key   string
	Value string
	prev  *Element
	next  *Element
}

type DoublyLinkedList struct {
	head *Element
	tail *Element
	len  int
}

func NewList() *DoublyLinkedList {
	head := &Element{}
	tail := &Element{}
	head.next = tail
	tail.prev = head
	return &DoublyLinkedList{head: head, tail: tail, len: 0}
}

func (l *DoublyLinkedList) PushFront(e *Element) {
	e.next = l.head.next
	e.prev = l.head
	l.head.next.prev = e
	l.head.next = e
	l.len++
}

func (l *DoublyLinkedList) MoveToFront(e *Element) {
	l.remove(e)
	l.PushFront(e)
}

func (l *DoublyLinkedList) Remove(e *Element) {
	l.remove(e)
}

func (l *DoublyLinkedList) remove(e *Element) {
	e.prev.next = e.next
	e.next.prev = e.prev
	e.prev = nil
	e.next = nil
	l.len--
}

func (l *DoublyLinkedList) Back() *Element {
	if l.len == 0 {
		return nil
	}
	return l.tail.prev
}

func (l *DoublyLinkedList) Len() int {
	return l.len
}

