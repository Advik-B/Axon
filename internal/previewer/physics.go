package previewer

import (
	"math"

	"github.com/Advik-B/Axon/pkg/axon"
)

// Physics constants for the simulation. Tweak these to change the "feel".
const (
	stiffness      = 0.01  // Spring stiffness
	damping        = 0.95  // Velocity damping factor per frame
	repulsion      = 50000 // Force pushing nodes apart
	mass           = 5.0   // Node mass
	minRepelDist   = float64(nodeWidth * 2)
	maxRepelDistSq = minRepelDist * minRepelDist
)

// Vec2 represents a 2D vector for physics calculations.
type Vec2 struct {
	X, Y float64
}

// PhysicsNode wraps a LayoutNode with physics properties.
type PhysicsNode struct {
	*LayoutNode
	Position Vec2
	Velocity Vec2
	Force    Vec2
}

// simulatePhysics runs one tick of the physics simulation.
func simulatePhysics(nodes map[string]*PhysicsNode, edges []*axon.DataEdge, execEdges []*axon.ExecEdge, draggedNode *PhysicsNode) {
	// 1. Reset forces on all nodes
	for _, n := range nodes {
		n.Force = Vec2{}
	}

	// 2. Calculate spring forces from data and exec edges
	for _, edge := range edges {
		applySpringForce(nodes[edge.FromNodeId], nodes[edge.ToNodeId], float64(nodeWidth+hSpacing))
	}
	for _, edge := range execEdges {
		applySpringForce(nodes[edge.FromNodeId], nodes[edge.ToNodeId], float64(nodeWidth+hSpacing))
	}

	// 3. Calculate repulsion forces between all pairs of nodes
	nodeSlice := make([]*PhysicsNode, 0, len(nodes))
	for _, n := range nodes {
		nodeSlice = append(nodeSlice, n)
	}
	for i := 0; i < len(nodeSlice); i++ {
		for j := i + 1; j < len(nodeSlice); j++ {
			applyRepulsionForce(nodeSlice[i], nodeSlice[j])
		}
	}

	// 4. Update position and velocity based on accumulated forces
	for _, n := range nodes {
		if n == draggedNode {
			n.Velocity = Vec2{} // Stop physics movement if being dragged
			continue
		}

		// Apply force to velocity (F = ma, so a = F/m)
		n.Velocity.X += n.Force.X / mass
		n.Velocity.Y += n.Force.Y / mass

		// Apply damping
		n.Velocity.X *= damping
		n.Velocity.Y *= damping

		// Apply velocity to position
		n.Position.X += n.Velocity.X
		n.Position.Y += n.Velocity.Y

		// Update the visual rectangle
		n.updateRect()
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