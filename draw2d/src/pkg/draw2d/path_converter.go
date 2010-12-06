package draw2d

import (
	"math"
)
type PathConverter struct {
	converter VertexConverter
	ApproximationScale, AngleTolerance, CuspLimit float
	startX, startY, x, y float
}

func NewPathConverter(converter VertexConverter) (*PathConverter) {
	return &PathConverter{converter, 1, 0, 0, 0, 0, 0, 0}
}

func (c *PathConverter) Convert(path *Path) {
	j := 0
	for _, cmd := range path.commands {
		j = j + c.ConvertCommand(cmd, path.vertices[j:]...)
	}
	c.converter.NextCommand(VertexStopCommand)
}


func (c *PathConverter) ConvertCommand(cmd PathCmd, vertices... float) int {
	switch cmd {
		case MoveTo:
			c.MoveTo(vertices[0], vertices[1])
			return 2
		case LineTo:
			c.LineTo(vertices[0], vertices[1])
			return 2
		case QuadCurveTo:
			c.QuadCurveTo(vertices[0], vertices[1], vertices[2], vertices[3])
			return 4
		case CubicCurveTo:
			c.CubicCurveTo(vertices[0], vertices[1], vertices[2], vertices[3], vertices[4], vertices[5])
			return 6
		case ArcTo:
			c.x, c.y = arc(c.converter, vertices[0], vertices[1], vertices[2], vertices[3], vertices[4], vertices[5], c.ApproximationScale)
			if(c.startX == c.x && c.startY== c.y) {
				c.converter.NextCommand(VertexCloseCommand)
			} 
			c.converter.Vertex(c.x, c.y)
			return 6
		case Close:
			c.Close()
			return 0
		}
		return 0
}

func (c *PathConverter) MoveTo(x, y float) *PathConverter {
	c.x, c.y = x, y
	c.startX, c.startY = c.x, c.y
	c.converter.NextCommand(VertexStopCommand)
	c.converter.NextCommand(VertexStartCommand)
	c.converter.Vertex(c.x, c.y)
	return c
}

func (c *PathConverter) RMoveTo( dx, dy float) *PathConverter {
	c.MoveTo(c.x+dx, c.y+dy)
	return c
}

func (c *PathConverter) LineTo( x, y float) *PathConverter {
	c.x, c.y = x, y
	if(c.startX == c.x && c.startY== c.y) {
		c.converter.NextCommand(VertexCloseCommand)
	} 
	c.converter.Vertex(c.x, c.y)
	c.converter.NextCommand(VertexJoinCommand)
	return c
}

func (c *PathConverter) RLineTo( dx, dy float) *PathConverter {
	c.LineTo(c.x+dx, c.y+dy)
	return c
}

func (c *PathConverter) QuadCurveTo( cx, cy, x, y float) *PathConverter {
	quadraticBezier(c.converter, c.x, c.y, cx, cy, x, y, c.ApproximationScale, c.AngleTolerance)
	c.x, c.y = x, y
	if(c.startX == c.x && c.startY== c.y) {
		c.converter.NextCommand(VertexCloseCommand)
	} 
	c.converter.Vertex(c.x, c.y)
	return c
}

func (c *PathConverter) RQuadCurveTo( dcx, dcy, dx, dy float) *PathConverter {
	c.QuadCurveTo(c.x+dcx, c.y+dcy, c.x+dx, c.y+dy)
	return c
}

func (c *PathConverter) CubicCurveTo( cx1, cy1, cx2, cy2, x, y float) *PathConverter {
	cubicBezier(c.converter, c.x, c.y, cx1, cy1, cx2, cy2, x, y, c.ApproximationScale, c.AngleTolerance, c.CuspLimit)
	c.x, c.y = x, y
	if(c.startX == c.x && c.startY== c.y) {
		c.converter.NextCommand(VertexCloseCommand)
	} 
	c.converter.Vertex(c.x, c.y)
	return c
}

func (c *PathConverter) RCubicCurveTo( dcx1, dcy1, dcx2, dcy2, dx, dy float) *PathConverter {
	c.CubicCurveTo(c.x+dcx1, c.y+dcy1, c.x+dcx2, c.y+dcy2, c.x+dx, c.y+dy)
	return c
}

func (c *PathConverter) ArcTo( cx, cy, rx, ry, startAngle, angle float) *PathConverter {
	endAngle := startAngle + angle
	clockWise := true
	if angle < 0 {
		clockWise = false
	}
	// normalize
	if clockWise {
		for endAngle < startAngle {
			endAngle += math.Pi * 2.0
		}
	} else {
		for startAngle < endAngle {
			startAngle += math.Pi * 2.0
		}
	}
	startX := cx + cos(startAngle)*rx
	startY := cy + sin(startAngle)*ry
	c.MoveTo(startX, startY)
	c.x, c.y = arc(c.converter, cx, cy, rx, ry, startAngle, angle, c.ApproximationScale)
	if(c.startX == c.x && c.startY== c.y) {
		c.converter.NextCommand(VertexCloseCommand)
	} 
	c.converter.Vertex(c.x, c.y)
	return c
}

func (c *PathConverter) RArcTo( dcx, dcy, rx, ry, startAngle, angle float) *PathConverter {
	c.ArcTo(c.x+dcx, c.y+dcy, rx, ry, startAngle, angle)
	return c
}

func (c *PathConverter) Close() *PathConverter {
	c.converter.NextCommand(VertexCloseCommand)
	c.converter.Vertex(c.startX, c.startY)
	return c
}
