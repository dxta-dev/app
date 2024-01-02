package graphs

type Node struct {
	id	float64
	datum interface{}
	free float64
	prev *Node
	next *Node
	value float64
}

type DoublyLinkedList struct {
	head *Node
	tail *Node
}

func (list *DoublyLinkedList) First() *Node {
	return list.head
}

func (list *DoublyLinkedList) Last() *Node {
	return list.tail
}

func (list *DoublyLinkedList) Add(node *Node) {
	if list.head == nil {
		list.head = node
		list.tail = node
	} else {
		list.tail.next = node
		node.prev = list.tail
		list.tail = node
	}
}

func (list *DoublyLinkedList) Find(id float64) *Node {
	node := list.head
	for node != nil {
		if node.id == id {
			return node
		}
		node = node.next
	}
	return nil
}
