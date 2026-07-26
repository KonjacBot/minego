package player

import (
	"container/heap"
	"errors"
	"math"

	"github.com/go-gl/mathgl/mgl64"

	"github.com/KonjacBot/go-mc/level/block"

	"github.com/KonjacBot/minego/pkg/bot"
	"github.com/KonjacBot/minego/pkg/protocol"
)

var ErrMaxNodesExceeded = errors.New("a* pathfinding exceeded max node count")

const fastAStarHeuristicWeight = 1.25

type Node struct {
	Position protocol.Position
	G        float64
	H        float64
	F        float64
	Parent   *Node
	Index    int
}

type NodeHeap []*Node

func (h NodeHeap) Len() int { return len(h) }
func (h NodeHeap) Less(left, right int) bool {
	if h[left].F != h[right].F {
		return h[left].F < h[right].F
	}
	if h[left].H != h[right].H {
		return h[left].H < h[right].H
	}
	return positionLess(h[left].Position, h[right].Position)
}
func (h NodeHeap) Swap(left, right int) {
	h[left], h[right] = h[right], h[left]
	h[left].Index = left
	h[right].Index = right
}
func (h *NodeHeap) Push(value any) {
	node := value.(*Node)
	node.Index = len(*h)
	*h = append(*h, node)
}
func (h *NodeHeap) Pop() any {
	old := *h
	last := len(old) - 1
	node := old[last]
	old[last] = nil
	node.Index = -1
	*h = old[:last]
	return node
}

type pathSearchOptions struct {
	maxNodes        int
	goalDistance    float64
	heuristicWeight float64
	flight          bool
}

type walkabilityCache struct {
	world  bot.World
	flight bool
	values map[protocol.Position]bool
}

func newWalkabilityCache(world bot.World, flight ...bool) *walkabilityCache {
	cache := &walkabilityCache{world: world, values: make(map[protocol.Position]bool)}
	if len(flight) > 0 {
		cache.flight = flight[0]
	}
	return cache
}

func (cache *walkabilityCache) walkable(pos protocol.Position) bool {
	if value, ok := cache.values[pos]; ok {
		return value
	}
	value := isWalkable(cache.world, pos)
	if cache.flight {
		value = canFlyThrough(cache.world, pos)
	}
	cache.values[pos] = value
	return value
}

// AStar keeps the original exact-goal contract while using the optimized
// Minecraft neighbor model and voxel cache.
func AStar(world bot.World, start, goal mgl64.Vec3, maxNodeCount int) ([]mgl64.Vec3, error) {
	return findPath(world, start, goal, pathSearchOptions{
		maxNodes: maxNodeCount, heuristicWeight: 1,
	})
}

// FastAStarWithin uses weighted A* and succeeds at the closest walkable cell
// within goalDistance. The range is measured from the player's feet to the
// supplied goal vector.
func FastAStarWithin(world bot.World, start, goal mgl64.Vec3, maxNodeCount int, goalDistance float64) ([]mgl64.Vec3, error) {
	return findPath(world, start, goal, pathSearchOptions{
		maxNodes: maxNodeCount, goalDistance: max(0, goalDistance), heuristicWeight: fastAStarHeuristicWeight,
	})
}

// FastFlightAStarWithin finds a bounded six-direction path through spaces that
// fit the player's feet and head. Unlike walking paths, it does not require a
// supporting block below every waypoint.
func FastFlightAStarWithin(world bot.World, start, goal mgl64.Vec3, maxNodeCount int, goalDistance float64) ([]mgl64.Vec3, error) {
	return findPath(world, start, goal, pathSearchOptions{
		maxNodes: maxNodeCount, goalDistance: max(0, goalDistance), heuristicWeight: fastAStarHeuristicWeight, flight: true,
	})
}

func findPath(world bot.World, start, goal mgl64.Vec3, options pathSearchOptions) ([]mgl64.Vec3, error) {
	startPos := vectorCell(start)
	goalPos := vectorCell(goal)
	if options.maxNodes <= 0 {
		options.maxNodes = 1
	}
	if options.heuristicWeight < 1 {
		options.heuristicWeight = 1
	}

	walkable := newWalkabilityCache(world, options.flight)
	if reachedGoal(startPos, goal, options.goalDistance) {
		return []mgl64.Vec3{cellVector(startPos)}, nil
	}
	if options.goalDistance == 0 && !walkable.walkable(goalPos) {
		return nil, nil
	}

	startNode := &Node{
		Position: startPos,
		H:        goalHeuristic(startPos, goal, options.goalDistance, options.flight),
		Index:    0,
	}
	startNode.F = startNode.H * options.heuristicWeight
	open := NodeHeap{startNode}
	heap.Init(&open)
	nodes := map[protocol.Position]*Node{startPos: startNode}
	closed := make(map[protocol.Position]struct{}, min(options.maxNodes, 256))
	best := startNode

	for open.Len() > 0 {
		if len(closed) >= options.maxNodes {
			return reconstructPath(best), ErrMaxNodesExceeded
		}
		current := heap.Pop(&open).(*Node)
		if _, alreadyClosed := closed[current.Position]; alreadyClosed {
			continue
		}
		if current.H < best.H {
			best = current
		}
		if reachedGoal(current.Position, goal, options.goalDistance) {
			return reconstructPath(current), nil
		}
		closed[current.Position] = struct{}{}

		neighbors := minecraftNeighbors(current.Position, walkable.walkable)
		if options.flight {
			neighbors = flightNeighbors(current.Position, walkable.walkable)
		}
		for _, neighbor := range neighbors {
			if _, visited := closed[neighbor]; visited {
				continue
			}
			candidateG := current.G + cellDistance(current.Position, neighbor)
			node, exists := nodes[neighbor]
			if !exists {
				node = &Node{
					Position: neighbor,
					G:        math.Inf(1),
					H:        goalHeuristic(neighbor, goal, options.goalDistance, options.flight),
					Index:    -1,
				}
				nodes[neighbor] = node
			}
			if candidateG >= node.G {
				continue
			}
			node.Parent = current
			node.G = candidateG
			node.F = candidateG + node.H*options.heuristicWeight
			if node.Index < 0 {
				heap.Push(&open, node)
			} else {
				heap.Fix(&open, node.Index)
			}
		}
	}
	return nil, nil
}

func flightNeighbors(pos protocol.Position, canFit func(protocol.Position) bool) []protocol.Position {
	directions := [...]protocol.Position{
		{1, 0, 0}, {-1, 0, 0}, {0, 0, 1}, {0, 0, -1}, {0, 1, 0}, {0, -1, 0},
	}
	neighbors := make([]protocol.Position, 0, len(directions))
	for _, direction := range directions {
		candidate := protocol.Position{pos[0] + direction[0], pos[1] + direction[1], pos[2] + direction[2]}
		if canFit(candidate) {
			neighbors = append(neighbors, candidate)
		}
	}
	return neighbors
}

func minecraftNeighbors(pos protocol.Position, walkable func(protocol.Position) bool) []protocol.Position {
	directions := [...]protocol.Position{
		{1, 0, 0}, {-1, 0, 0}, {0, 0, 1}, {0, 0, -1},
		{1, 0, 1}, {1, 0, -1}, {-1, 0, 1}, {-1, 0, -1},
	}
	neighbors := make([]protocol.Position, 0, len(directions))
	for _, direction := range directions {
		sameLevel := protocol.Position{pos[0] + direction[0], pos[1], pos[2] + direction[2]}
		diagonal := direction[0] != 0 && direction[2] != 0
		if walkable(sameLevel) && (!diagonal || diagonalClear(pos, direction, walkable)) {
			neighbors = append(neighbors, sameLevel)
			continue
		}
		if diagonal {
			continue
		}
		stepUp := protocol.Position{sameLevel[0], sameLevel[1] + 1, sameLevel[2]}
		if walkable(stepUp) {
			neighbors = append(neighbors, stepUp)
			continue
		}
		stepDown := protocol.Position{sameLevel[0], sameLevel[1] - 1, sameLevel[2]}
		if walkable(stepDown) {
			neighbors = append(neighbors, stepDown)
		}
	}
	return neighbors
}

func diagonalClear(pos, direction protocol.Position, walkable func(protocol.Position) bool) bool {
	return walkable(protocol.Position{pos[0] + direction[0], pos[1], pos[2]}) &&
		walkable(protocol.Position{pos[0], pos[1], pos[2] + direction[2]})
}

func vectorCell(value mgl64.Vec3) protocol.Position {
	return protocol.Position{
		int32(math.Floor(value.X())),
		int32(math.Floor(value.Y())),
		int32(math.Floor(value.Z())),
	}
}

func cellVector(pos protocol.Position) mgl64.Vec3 {
	return mgl64.Vec3{float64(pos[0]), float64(pos[1]), float64(pos[2])}
}

func cellCenter(pos protocol.Position) mgl64.Vec3 {
	return mgl64.Vec3{float64(pos[0]) + 0.5, float64(pos[1]), float64(pos[2]) + 0.5}
}

func reachedGoal(pos protocol.Position, goal mgl64.Vec3, goalDistance float64) bool {
	if goalDistance == 0 {
		return pos == vectorCell(goal)
	}
	return cellCenter(pos).Sub(goal).LenSqr() <= goalDistance*goalDistance
}

func goalHeuristic(pos protocol.Position, goal mgl64.Vec3, goalDistance float64, flight bool) float64 {
	delta := cellCenter(pos).Sub(goal)
	if flight {
		estimate := math.Abs(delta.X()) + math.Abs(delta.Y()) + math.Abs(delta.Z())
		return max(0, estimate-goalDistance)
	}
	horizontalMin := min(math.Abs(delta.X()), math.Abs(delta.Z()))
	horizontalMax := max(math.Abs(delta.X()), math.Abs(delta.Z()))
	octile := horizontalMax + (math.Sqrt2-1)*horizontalMin
	estimate := octile + math.Abs(delta.Y())
	return max(0, estimate-goalDistance)
}

func cellDistance(left, right protocol.Position) float64 {
	dx := float64(left[0] - right[0])
	dy := float64(left[1] - right[1])
	dz := float64(left[2] - right[2])
	return math.Sqrt(dx*dx + dy*dy + dz*dz)
}

func positionLess(left, right protocol.Position) bool {
	if left[1] != right[1] {
		return left[1] < right[1]
	}
	if left[0] != right[0] {
		return left[0] < right[0]
	}
	return left[2] < right[2]
}

func isWalkable(world bot.World, pos protocol.Position) bool {
	footBlock, err := world.GetBlock(pos)
	if err != nil {
		return false
	}
	headBlock, err := world.GetBlock(protocol.Position{pos[0], pos[1] + 1, pos[2]})
	if err != nil {
		return false
	}
	supportBlock, err := world.GetBlock(protocol.Position{pos[0], pos[1] - 1, pos[2]})
	if err != nil {
		return false
	}
	return block.IsAirBlock(footBlock) && block.IsAirBlock(headBlock) && !block.IsAirBlock(supportBlock)
}

func canFlyThrough(world bot.World, pos protocol.Position) bool {
	footBlock, err := world.GetBlock(pos)
	if err != nil {
		return false
	}
	headBlock, err := world.GetBlock(protocol.Position{pos[0], pos[1] + 1, pos[2]})
	if err != nil {
		return false
	}
	return block.IsAirBlock(footBlock) && block.IsAirBlock(headBlock)
}

func reconstructPath(last *Node) []mgl64.Vec3 {
	length := 0
	for node := last; node != nil; node = node.Parent {
		length++
	}
	path := make([]mgl64.Vec3, length)
	for index, node := length-1, last; node != nil; index, node = index-1, node.Parent {
		path[index] = cellVector(node.Position)
	}
	return path
}
