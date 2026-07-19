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

// Node 表示 A* 演算法中的節點
type Node struct {
	Position protocol.Position
	G        float64 // 從起點到當前節點的實際距離
	H        float64 // 從當前節點到終點的啟發式距離
	F        float64 // G + H
	Parent   *Node
	Index    int // heap 索引
}

// NodeHeap 實現 heap.Interface 用於優先佇列
type NodeHeap []*Node

func (h NodeHeap) Len() int           { return len(h) }
func (h NodeHeap) Less(i, j int) bool { return h[i].F < h[j].F }
func (h NodeHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
	h[i].Index = i
	h[j].Index = j
}

func (h *NodeHeap) Push(x interface{}) {
	n := len(*h)
	node := x.(*Node)
	node.Index = n
	*h = append(*h, node)
}

func (h *NodeHeap) Pop() interface{} {
	old := *h
	n := len(old)
	node := old[n-1]
	node.Index = -1
	*h = old[0 : n-1]
	return node
}

// AStar 使用 A* 演算法尋找路徑（新增 maxNodeCount 限制）
func AStar(world bot.World, start, goal mgl64.Vec3, maxNodeCount int) ([]mgl64.Vec3, error) {
	// 將浮點數座標轉換為區塊整數座標
	startPos := protocol.Position{int32(math.Floor(start.X())), int32(math.Floor(start.Y())), int32(math.Floor(start.Z()))}
	goalPos := protocol.Position{int32(math.Floor(goal.X())), int32(math.Floor(goal.Y())), int32(math.Floor(goal.Z()))}

	// 如果終點本身就不可通行，直接防呆返回
	if !isWalkable(world, goalPos) {
		return nil, nil
	}

	openSet := &NodeHeap{}
	heap.Init(openSet)

	closedSet := make(map[protocol.Position]bool)
	allNodes := make(map[protocol.Position]*Node)

	// 初始化起點節點
	startNode := &Node{
		Position: startPos,
		G:        0,
		H:        heuristic(startPos, goalPos),
		Index:    0,
	}
	startNode.F = startNode.G + startNode.H

	heap.Push(openSet, startNode)
	allNodes[startPos] = startNode

	// 追蹤走過的節點數量，以及目前最靠近終點的節點（防禦性設計）
	nodesExplored := 0
	bestNode := startNode

	for openSet.Len() > 0 {
		// 檢查是否超過最大節點搜尋限制
		if maxNodeCount > 0 && nodesExplored >= maxNodeCount {
			// 選擇 1：返回目前為止最接近終點的路徑（推薦，機器人不會卡死）
			return reconstructPath(bestNode), ErrMaxNodesExceeded
		}

		current := heap.Pop(openSet).(*Node)
		nodesExplored++

		// 更新目前最接近終點的節點（根據啟發式距離 H，越小代表越接近終點）
		if current.H < bestNode.H {
			bestNode = current
		}

		// 找到終點，開始回溯路徑
		if current.Position == goalPos {
			return reconstructPath(current), nil
		}

		closedSet[current.Position] = true

		// 檢查 6 個方向的相鄰節點
		for _, neighbor := range getNeighbors(current.Position) {
			if closedSet[neighbor] {
				continue
			}

			// 檢查該位置機器人是否容納得下
			if !isWalkable(world, neighbor) {
				continue
			}

			tentativeG := current.G + distance(current.Position, neighbor)

			neighborNode, exists := allNodes[neighbor]
			if !exists {
				neighborNode = &Node{
					Position: neighbor,
					G:        math.Inf(1),
					H:        heuristic(neighbor, goalPos),
					Index:    -1,
				}
				allNodes[neighbor] = neighborNode
			}

			// 如果這條路徑比之前找到的更好，更新節點資訊
			if tentativeG < neighborNode.G {
				neighborNode.Parent = current
				neighborNode.G = tentativeG
				neighborNode.F = neighborNode.G + neighborNode.H

				if neighborNode.Index == -1 {
					heap.Push(openSet, neighborNode)
				} else {
					heap.Fix(openSet, neighborNode.Index)
				}
			}
		}
	}

	return nil, nil // 找不到路徑
}

// heuristic 計算啟發式距離（曼哈頓距離）
func heuristic(a, b protocol.Position) float64 {
	return math.Abs(float64(a[0]-b[0])) + math.Abs(float64(a[1]-b[1])) + math.Abs(float64(a[2]-b[2]))
}

// distance 計算兩點間的實際距離
func distance(a, b protocol.Position) float64 {
	dx := float64(a[0] - b[0])
	dy := float64(a[1] - b[1])
	dz := float64(a[2] - b[2])
	return math.Sqrt(dx*dx + dy*dy + dz*dz)
}

// getNeighbors 獲取相鄰節點
func getNeighbors(pos protocol.Position) []protocol.Position {
	neighbors := []protocol.Position{
		{pos[0] + 1, pos[1], pos[2]}, // 東
		{pos[0] - 1, pos[1], pos[2]}, // 西
		{pos[0], pos[1], pos[2] + 1}, // 南
		{pos[0], pos[1], pos[2] - 1}, // 北
		{pos[0], pos[1] + 1, pos[2]}, // 上
		{pos[0], pos[1] - 1, pos[2]}, // 下
	}
	return neighbors
}

// isWalkable 檢查位置是否可通行
func isWalkable(world bot.World, pos protocol.Position) bool {
	// 檢查腳部位置
	footBlock, err := world.GetBlock(pos)
	if err != nil {
		return false
	}

	// 檢查頭部位置
	headPos := protocol.Position{pos[0], pos[1] + 1, pos[2]}
	headBlock, err := world.GetBlock(headPos)
	if err != nil {
		return false
	}

	supportPos := protocol.Position{pos[0], pos[1] - 1, pos[2]}
	supportBlock, err := world.GetBlock(supportPos)
	if err != nil {
		return false
	}

	return block.IsAirBlock(footBlock) && block.IsAirBlock(headBlock) && !block.IsAirBlock(supportBlock)
}

// reconstructPath 重建路徑
func reconstructPath(node *Node) []mgl64.Vec3 {
	var path []mgl64.Vec3
	current := node

	for current != nil {
		pos := mgl64.Vec3{
			float64(current.Position[0]),
			float64(current.Position[1]),
			float64(current.Position[2]),
		}
		path = append([]mgl64.Vec3{pos}, path...)
		current = current.Parent
	}

	return path
}
