package z3

// #include <stdlib.h>
// #include "go-z3.h"
import "C"
import "unsafe"

// AST represents an AST value in Z3.
//
// AST memory management is automatically managed by the Context it
// is contained within. When the Context is freed, so are the AST nodes.
type AST struct {
	rawCtx C.Z3_context
	rawAST C.Z3_ast
}

// String returns a human-friendly string version of the AST.
func (a *AST) String() string {
	return C.GoString(C.Z3_ast_to_string(a.rawCtx, a.rawAST))
}

// DeclName returns the name of a declaration. The AST value must be a
// func declaration for this to work.
func (a *AST) DeclName() *Symbol {
	return &Symbol{
		rawCtx: a.rawCtx,
		rawSymbol: C.Z3_get_decl_name(
			a.rawCtx, C.Z3_to_func_decl(a.rawCtx, a.rawAST)),
	}
}

//-------------------------------------------------------------------
// func creation & application
//-------------------------------------------------------------------

type FDECL struct {
	rawCtx  C.Z3_context
	rawFunc C.Z3_func_decl
}

func (c *Context) FuncDecl(s *Symbol, domain []*Sort, rangeSort *Sort) *FDECL {
	l := len(domain)
	dom := make([]C.Z3_sort, l)
	for i, d := range domain {
		dom[i] = (*d).rawSort
	}

	return &FDECL{
		rawCtx: c.raw,
		rawFunc: C.Z3_mk_func_decl(
			c.raw,
			s.rawSymbol,
			C.uint(l),
			(*C.Z3_sort)(unsafe.Pointer(&dom[0])),
			rangeSort),
	}
}

func (c *Context) App(fdecl *FDECL, args []*AST) *AST {
	l := len(args)
	arg := make([]C.Z3_ast, l)
	for i, a := range args {
		arg[i] = (*a).rawAST
	}

	return &AST{
		rawCtx: c.raw,
		rawAST: C.mk_app(
			c.raw,
			fdecl.rawFunc,
			C.uint(l),
			(*C.Z3_ast)(unsafe.Pointer(&arg[0]))),
	}
}

//-------------------------------------------------------------------
// Var, Literal Creation
//-------------------------------------------------------------------

// Const declares a variable. It is called "Const" since internally
// this is equivalent to create a function that always returns a constant
// value. From an initial user perspective this may be confusing but go-z3
// is following identical naming convention.
func (c *Context) Const(s *Symbol, typ *Sort) *AST {
	return &AST{
		rawCtx: c.raw,
		rawAST: C.Z3_mk_const(c.raw, s.rawSymbol, typ.rawSort),
	}
}

// Int creates an integer type.
//
// Maps: Z3_mk_int
func (c *Context) Int(v int, typ *Sort) *AST {
	return &AST{
		rawCtx: c.raw,
		rawAST: C.Z3_mk_int(c.raw, C.int(v), typ.rawSort),
	}
}

// True creates the value "true".
//
// Maps: Z3_mk_true
func (c *Context) True() *AST {
	return &AST{
		rawCtx: c.raw,
		rawAST: C.Z3_mk_true(c.raw),
	}
}

// False creates the value "false".
//
// Maps: Z3_mk_false
func (c *Context) False() *AST {
	return &AST{
		rawCtx: c.raw,
		rawAST: C.Z3_mk_false(c.raw),
	}
}

//-------------------------------------------------------------------
// Value Readers
//-------------------------------------------------------------------

// Int gets the integer value of this AST. The value must be able to fit
// into a machine integer.
func (a *AST) Int() int {
	var dst C.int
	C.Z3_get_numeral_int(a.rawCtx, a.rawAST, &dst)
	return int(dst)
}
