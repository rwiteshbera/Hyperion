package blockchain

func enqueue[T comparable](queue []*T, element *T) []*T {
	queue = append(queue, element)
	return queue
}

// func dequeue[T comparable](queue []*T) (*T, []*T, error) {
// 	if len(queue) == 0 {
// 		return nil, queue, errors.New("cannot dequeue from an empty queue")
// 	}
// 	removed := queue[0]
// 	return removed, queue[1:], nil
// }

// func peakQueue[T comparable](queue []*T) *T {
// 	return queue[0]
// }

// func isQueueEmpty[T comparable](queue []*T) bool {
// 	return len(queue) == 0
// }
