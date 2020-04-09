package xlsx

import (
	"fmt"
)

// Row represents a single Row in the current Sheet.
type Row struct {
	Hidden       bool    // Hidden determines whether this Row is hidden or not.
	Sheet        *Sheet  // Sheet is a reference back to the Sheet that this Row is within.
	Height       float64 // Height is the current height of the Row in PostScript Points
	OutlineLevel uint8   // OutlineLevel contains the outline level of this Row.  Used for collapsing.
	isCustom     bool    // isCustom is a flag that is set to true when the Row has been modified
	num          int     // Num hold the positional number of the Row in the Sheet
	cellCount    int     // The current number of cells
	cells        []*Cell // the cells
}

// SetHeight sets the height of the Row in PostScript Points
func (r *Row) SetHeight(ht float64) {
	r.Height = ht
	r.isCustom = true
}

// SetHeightCM sets the height of the Row in centimetres, inherently converting it to PostScript points.
func (r *Row) SetHeightCM(ht float64) {
	r.Height = ht * 28.3464567 // Convert CM to postscript points
	r.isCustom = true
}

// AddCell adds a new Cell to the Row
func (r *Row) AddCell() *Cell {
	cell := newCell(r, r.cellCount)
	r.cellCount++
	r.cells = append(r.cells, cell)
	return cell
}

// CopyCell adds a new Cell to the Row from copy
func (r *Row) CopyCell(cell *Cell) {
	r.cellCount++
	r.cells = append(r.cells, cell)
	return
}

func (r *Row) makeCellKey(colIdx int) string {
	return fmt.Sprintf("%s:%06d:%06d", r.Sheet.Name, r.num, colIdx)
}

func (r *Row) key() string {
	return r.makeCellKeyRowPrefix()
}

func (r *Row) makeCellKeyRowPrefix() string {
	return fmt.Sprintf("%s:%06d", r.Sheet.Name, r.num)
}

func (r *Row) growCellsSlice(newSize int) {
	capacity := cap(r.cells)
	if newSize >= capacity {
		newCap := 2 * capacity
		if newSize > newCap {
			newCap = newSize
		}
		newSlice := make([]*Cell, newCap, newCap)
		copy(newSlice, r.cells)
		r.cells = newSlice
	}
}

// GetCell returns the Cell at a given column index, creating it if it doesn't exist.
func (r *Row) GetCell(colIdx int) *Cell {
	if colIdx >= len(r.cells) {
		cell := newCell(r, colIdx)
		r.growCellsSlice(colIdx + 1)

		r.cells[colIdx] = cell
		return cell
	}

	cell := r.cells[colIdx]
	if cell == nil {
		cell = newCell(r, colIdx)
		r.cells[colIdx] = cell
	}
	return cell
}

// ForEachCell will call the provided CellVisitorFunc for each
// currently defined cell in the Row.
func (r *Row) ForEachCell(cvf CellVisitorFunc) error {
	fn := func(c *Cell) error {
		if c != nil {
			c.Row = r
			return cvf(c)
		}
		return nil
	}

	for _, cell := range r.cells {
		err := fn(cell)
		if err != nil {
			return err
		}
	}

	return nil
}
