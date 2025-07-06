package previewer

import (
	"math"

	"github.com/Advik-B/Axon/pkg/axon"
)

const (
	stiffness      = 0.01
	damping        = 0.92
	repulsion      = 60000
	mass           = 5.0
	attraction     = 0.002
	minRepelDist   = float64(nodeWidth * 2)
	maxRepelDistSq = minRepelDist * minRepelDist
)

type Vec2 struct {
	X, Y float64
}

type PhysicsNode struct {
	*LayoutNode
	Position       Vec2
	TargetPosition Vec2
	Velocity       Vec2
	Force          Vec2
}

// simulatePhysics runs one tick of the physics simulation.
// **THE FIX**: It now accepts the current orientation to pass to updateRect.
func simulatePhysics(nodes map[string]*PhysicsNode, edges []*axon.DataEdge, execEdges []*axon.ExecEdge, draggedNode *PhysicsNode, orientation LayoutOrientation) {
	for _, n := range nodes {
		n.Force = Vec2{}
	}

	for _, edge := range edges {
		applySpringForce(nodes[edge.FromNodeId], nodes[edge.ToNodeId], float64(nodeWidth+hSpacing))
	}
	for _, edge := range execEdges {
		applySpringForce(nodes[edge.FromNodeId], nodes[edge.ToNodeId], float64(nodeWidth+hSpacing))
	}

	nodeSlice := make([]*PhysicsNode, 0, len(nodes))
	for _, n := range nodes {
		nodeSlice = append(nodeSlice, n)
	}
	for i := 0; i < len(nodeSlice); i++ {
		for j := i + 1; j < len(nodeSlice); j++ {
			applyRepulsionForce(nodeSlice[i], nodeSlice[j])
		}
	}

	for _, n := range nodes {
		applyAttractionForce(n)
	}

	for _, n := range nodes {
		if n == draggedNode {
			n.Velocity = Vec2{}
			continue
		}

		n.Velocity.X += n.Force.X / mass
		n.Velocity.Y += n.Force.Y / mass
		n.Velocity.X *= damping
		n.Velocity.Y *= damping
		n.Position.X += n.Velocity.X
		n.Position.Y += n.Velocity.Y

		// **THE FIX**: Pass the required orientation argument.
		n.updateRect(orientation)
	}
}

func applySpringForce(n1, n2 *PhysicsNode, restLength float64) {
	if n1 == nil || n2 == nil {
		return
	}
	dx := n2.Position.X - n1.Position.X
	dy := n2.Position.Y - n1.Position.Y
	dist := math.Sqrt(dx*dx + dy*dy)
	if dist == 0 {
		return
	}
	displacement := dist - restLength
	forceMag := stiffness * displacement
	forceX := (dx / dist) * forceMag
	forceY := (dy / dist) * forceMag
	n1.Force.X += forceX
	n1.Force.Y += forceY
	n2.Force.X -= forceX
	n2.Force.Y -= forceY
}

func applyRepulsionForce(n1, n2 *PhysicsNode) {
	dx := n2.Position.X - n1.Position.X
	dy := n2.Position.Y - n1.Position.Y
	distSq := dx*dx + dy*dy
	if distSq == 0 || distSq > maxRepelDistSq {
		return
	}
	dist := math.Sqrt(distSq)
	forceMag := repulsion / distSq
	forceX := (dx / dist) * forceMag
	forceY := (dy / dist) * forceMag
	n1.Force.X -= forceX
	n1.Force.Y -= forceY
	n2.Force.X += forceX
	n2.Force.Y += forceY
}

func applyAttractionForce(n *PhysicsNode) {
	dx := n.TargetPosition.X - n.Position.X
	dy := n.TargetPosition.Y - n.Position.Y
	n.Force.X += dx * attraction * mass
	n.Force.Y += dy * attraction * mass
}
