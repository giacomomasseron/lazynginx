package boxlayout

// Dimensions represents the coordinates of a box
type Dimensions struct {
	X0 int
	X1 int
	Y0 int
	Y1 int
}

// Direction defines how children are laid out
type Direction int

const (
	ROW    Direction = iota // Children stacked vertically
	COLUMN                  // Children arranged horizontally
)

// Box represents a layout box that can contain children or a window
type Box struct {
	// Direction decides how the children boxes are laid out
	Direction Direction

	// Children are the sub-boxes
	Children []*Box

	// Window refers to the name of the window this box represents
	Window string

	// Size is static size (height for ROW parent, width for COLUMN parent)
	Size int

	// Weight is dynamic size proportion (after static sizes allocated)
	Weight int
}

// ArrangeWindows calculates the dimensions for all windows in the box tree
func ArrangeWindows(root *Box, x0, y0, width, height int) map[string]Dimensions {
	children := root.Children
	if len(children) == 0 {
		// leaf node
		if root.Window != "" {
			dimensionsForWindow := Dimensions{
				X0: x0,
				Y0: y0,
				X1: x0 + width - 1,
				Y1: y0 + height - 1,
			}
			return map[string]Dimensions{root.Window: dimensionsForWindow}
		}
		return map[string]Dimensions{}
	}

	direction := root.Direction

	var availableSize int
	if direction == COLUMN {
		availableSize = width
	} else {
		availableSize = height
	}

	sizes := calcSizes(children, availableSize)

	result := map[string]Dimensions{}
	offset := 0
	for i, child := range children {
		boxSize := sizes[i]

		var resultForChild map[string]Dimensions
		if direction == COLUMN {
			resultForChild = ArrangeWindows(child, x0+offset, y0, boxSize, height)
		} else {
			resultForChild = ArrangeWindows(child, x0, y0+offset, width, boxSize)
		}

		result = mergeDimensionMaps(result, resultForChild)
		offset += boxSize
	}

	return result
}

// calcSizes determines the size of each box based on static sizes and weights
func calcSizes(boxes []*Box, availableSpace int) []int {
	totalWeight := 0
	reservedSpace := 0

	for _, box := range boxes {
		if box.isStatic() {
			reservedSpace += box.Size
		} else {
			totalWeight += box.Weight
		}
	}

	dynamicSpace := max(0, availableSpace-reservedSpace)

	unitSize := 0
	extraSpace := 0
	if totalWeight > 0 {
		unitSize = dynamicSpace / totalWeight
		extraSpace = dynamicSpace % totalWeight
	}

	result := make([]int, len(boxes))
	for i, box := range boxes {
		if box.isStatic() {
			result[i] = min(availableSpace, box.Size)
		} else {
			result[i] = unitSize * box.Weight
		}
	}

	// Distribute remainder across dynamic boxes
	for i := 0; i < len(boxes) && extraSpace > 0; i++ {
		if !boxes[i].isStatic() {
			result[i]++
			extraSpace--
		}
	}

	return result
}

func (b *Box) isStatic() bool {
	return b.Size > 0
}

func mergeDimensionMaps(a map[string]Dimensions, b map[string]Dimensions) map[string]Dimensions {
	result := map[string]Dimensions{}
	for k, v := range a {
		result[k] = v
	}
	for k, v := range b {
		result[k] = v
	}
	return result
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
