package z3

// #include "go-z3.h"
import "C"

// Sort represents a sort in Z3.
type Sort struct {
	rawCtx  C.Z3_context
	rawSort C.Z3_sort
}

// BoolSort returns the boolean type.
func (c *Context) BoolSort() *Sort {
	return &Sort{
		rawCtx:  c.raw,
		rawSort: C.Z3_mk_bool_sort(c.raw),
	}
}

// IntSort returns the int type.
func (c *Context) IntSort() *Sort {
	return &Sort{
		rawCtx:  c.raw,
		rawSort: C.Z3_mk_int_sort(c.raw),
	}
}

// UninterpretedSort returns the uninterpreted type.
func (c *Context) UninterpretedSort(symbol *Symbol) *Sort {
	return &Sort{
		rawCtx:  c.raw,
		rawSort: C.Z3_mk_uninterpreted_sort(c.raw, symbol.rawSymbol),
	}
}
