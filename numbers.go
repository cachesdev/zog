package zog

import (
	"github.com/Oudwins/zog/conf"
	p "github.com/Oudwins/zog/internals"
)

type Numeric interface {
	~int | ~float64
}

type numberProcessor[T Numeric] struct {
	preTransforms  []p.PreTransform
	tests          []p.Test
	postTransforms []p.PostTransform
	defaultVal     *T
	required       *p.Test
	catch          *T
}

// creates a new float64 processor
func Float() *numberProcessor[float64] {
	return &numberProcessor[float64]{}
}

// creates a new int processor
func Int() *numberProcessor[int] {
	return &numberProcessor[int]{}
}

// parses the value and stores it in the destination
func (v *numberProcessor[T]) Parse(data any, dest *T, options ...ParsingOption) p.ZogErrList {
	errs := p.NewErrsList()
	ctx := p.NewParseCtx(errs, conf.ErrorFormatter)
	for _, opt := range options {
		opt(ctx)
	}

	path := p.PathBuilder("")

	v.process(data, dest, path, ctx)

	return errs.List
}

func (v *numberProcessor[T]) process(val any, dest any, path p.PathBuilder, ctx ParseCtx) {

	var coercer conf.CoercerFunc
	switch any(dest).(type) {
	case *float64:
		coercer = conf.Coercers.Float64
	case *int:
		coercer = conf.Coercers.Int
	}

	primitiveProcessor(val, dest, path, ctx, v.preTransforms, v.tests, v.postTransforms, v.defaultVal, v.required, v.catch, coercer)
}

// GLOBAL METHODS

func (v *numberProcessor[T]) PreTransform(transform p.PreTransform) *numberProcessor[T] {
	if v.preTransforms == nil {
		v.preTransforms = []p.PreTransform{}
	}
	v.preTransforms = append(v.preTransforms, transform)
	return v
}

// Adds posttransform function to schema
func (v *numberProcessor[T]) PostTransform(transform p.PostTransform) *numberProcessor[T] {
	if v.postTransforms == nil {
		v.postTransforms = []p.PostTransform{}
	}
	v.postTransforms = append(v.postTransforms, transform)
	return v
}

// ! MODIFIERS

// marks field as required
func (v *numberProcessor[T]) Required(options ...TestOption) *numberProcessor[T] {
	r := p.Required()
	for _, opt := range options {
		opt(&r)
	}
	v.required = &r
	return v
}

// marks field as optional
func (v *numberProcessor[T]) Optional() *numberProcessor[T] {
	v.required = nil
	return v
}

// sets the default value
func (v *numberProcessor[T]) Default(val T) *numberProcessor[T] {
	v.defaultVal = &val
	return v
}

// sets the catch value (i.e the value to use if the validation fails)
func (v *numberProcessor[T]) Catch(val T) *numberProcessor[T] {
	v.catch = &val
	return v
}

// custom test function call it -> schema.Test(test, options)
func (v *numberProcessor[T]) Test(t p.Test, opts ...TestOption) *numberProcessor[T] {
	for _, opt := range opts {
		opt(&t)
	}
	v.tests = append(v.tests, t)
	return v
}

// UNIQUE METHODS

// Check that the value is one of the enum values
func (v *numberProcessor[T]) OneOf(enum []T, options ...TestOption) *numberProcessor[T] {
	t := p.In(enum)
	for _, opt := range options {
		opt(&t)
	}
	v.tests = append(v.tests, t)
	return v
}

// checks for equality
func (v *numberProcessor[T]) EQ(n T, options ...TestOption) *numberProcessor[T] {
	t := p.EQ(n)
	for _, opt := range options {
		opt(&t)
	}
	v.tests = append(v.tests, t)
	return v
}

// checks for lesser or equal
func (v *numberProcessor[T]) LTE(n T, options ...TestOption) *numberProcessor[T] {
	t := p.LTE(n)
	for _, opt := range options {
		opt(&t)
	}
	v.tests = append(v.tests, t)
	return v
}

// checks for greater or equal
func (v *numberProcessor[T]) GTE(n T, options ...TestOption) *numberProcessor[T] {
	t := p.GTE(n)
	for _, opt := range options {
		opt(&t)
	}
	v.tests = append(v.tests, t)
	return v
}

// checks for lesser
func (v *numberProcessor[T]) LT(n T, options ...TestOption) *numberProcessor[T] {
	t := p.LT(n)
	for _, opt := range options {
		opt(&t)
	}
	v.tests = append(v.tests, t)
	return v
}

// checks for greater
func (v *numberProcessor[T]) GT(n T, options ...TestOption) *numberProcessor[T] {
	t := p.GT(n)
	for _, opt := range options {
		opt(&t)
	}
	v.tests = append(v.tests, t)
	return v
}
