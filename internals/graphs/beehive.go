package graphs

import (
	"errors"
	"math"
)

type Point struct {
	X, Y float64
}

type Direction bool

const (
	Prev Direction = false
	Next Direction = true
)

func visitToGatherXBasedColliderData(dln *Node, visitedDln *Node, direction Direction, minDistance float64, xBasedPossibleColliders []float64) (interface{}, error) {
	if visitedDln == nil {
		return nil, errors.New("No nearest possible collider found")
	}

	if direction == Prev && dln.value-visitedDln.value > minDistance {
		return nil, errors.New("No nearest possible collider found")
	}

	if direction == Next && visitedDln.value-dln.value > minDistance {
		return nil, errors.New("No nearest possible collider found")
	}

	xBasedPossibleColliders = append(xBasedPossibleColliders, visitedDln.id)

	if direction == Prev {
		visitToGatherXBasedColliderData(dln, visitedDln.prev, direction, minDistance, xBasedPossibleColliders)
	} else {
		visitToGatherXBasedColliderData(dln, visitedDln.next, direction, minDistance, xBasedPossibleColliders)
	}

	return xBasedPossibleColliders, nil

}

func findNearestPossibleCollider(dln *Node, visitedDln *Node, direction Direction, minDistance float64) (float64, error) {
	if visitedDln == nil {
		return -1, errors.New("No nearest possible collider found")
	}

	if direction == Prev && dln.value-visitedDln.value > minDistance {
		return -1, errors.New("No nearest possible collider found")
	}

	if direction == Next && visitedDln.value-dln.value > minDistance {
		return -1, errors.New("No nearest possible collider found")
	}

	if visitedDln.free != math.Inf(-1) {
		return visitedDln.free, nil
	}

	if direction == Prev {
		return findNearestPossibleCollider(dln, visitedDln.prev, direction, minDistance)
	} else {
		return findNearestPossibleCollider(dln, visitedDln.next, direction, minDistance)
	}

}

func Beehive(xValues []float64, yValues []float64, r float64) {

	minDistance := 2 * r
	minSquaredDistance := minDistance * minDistance

	xBasedDataManager := DoublyLinkedList{}

	for i := 0; i < len(xValues); i++ {
		xBasedDataManager.Add(&Node{id: float64(i), value: xValues[i], free: math.Inf(-1), datum: Point{X: xValues[i], Y: yValues[i]}})
	}

	xBasedColliderManager := DoublyLinkedList{}
	yBasedColliderManager := DoublyLinkedList{}

	for i := 0; i < len(xValues); i++ {
		bestYPosition := math.Inf(-1)
		xBasedPossibleColliders := []float64{}
		dln := xBasedDataManager.Find(i)

		nearestXPrevAlreadyArrangedData, errPrev := findNearestPossibleCollider(dln, dln.prev, Prev, minDistance)
		nearestXNextAlreadyArrangedData, errNext := findNearestPossibleCollider(dln, dln.next, Next, minDistance)

		if errPrev == nil {
			dln = xBasedDataManager.Find(nearestXPrevAlreadyArrangedData)
			visitToGatherXBasedColliderData(dln, dln.prev, Prev, minDistance, xBasedPossibleColliders)
		}

		if errNext == nil {
			dln = xBasedDataManager.Find(nearestXNextAlreadyArrangedData)
			visitToGatherXBasedColliderData(dln, dln.next, Next, minDistance, xBasedPossibleColliders)
		}

		if len(xBasedPossibleColliders) == 0 {
			bestYPosition = yValues[i]
		}

	}

}
