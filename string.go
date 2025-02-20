package zog

import (
	"net/url"
	"regexp"
	"strings"

	"github.com/Oudwins/zog/conf"
	p "github.com/Oudwins/zog/internals"
	"github.com/Oudwins/zog/zconst"
)

var (
	emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
)

type stringProcessor struct {
	preTransforms  []p.PreTransform
	tests          []p.Test
	postTransforms []p.PostTransform
	defaultVal     *string
	required       *p.Test
	catch          *string
}

func String() *stringProcessor {
	return &stringProcessor{
		tests: []p.Test{},
	}
}

func (v *stringProcessor) Parse(data any, dest *string, options ...ParsingOption) p.ZogErrList {
	errs := p.NewErrsList()
	ctx := p.NewParseCtx(errs, conf.ErrorFormatter)
	for _, opt := range options {
		opt(ctx)
	}
	path := p.PathBuilder("")

	v.process(data, dest, path, ctx)

	return errs.List
}

func (v *stringProcessor) process(val any, dest any, path p.PathBuilder, ctx ParseCtx) {
	primitiveProcessor(val, dest, path, ctx, v.preTransforms, v.tests, v.postTransforms, v.defaultVal, v.required, v.catch, conf.Coercers.String)
}

// Adds pretransform function to schema
func (v *stringProcessor) PreTransform(transform p.PreTransform) *stringProcessor {
	if v.preTransforms == nil {
		v.preTransforms = []p.PreTransform{}
	}
	v.preTransforms = append(v.preTransforms, transform)
	return v
}

// Adds posttransform function to schema
func (v *stringProcessor) PostTransform(transform p.PostTransform) *stringProcessor {
	if v.postTransforms == nil {
		v.postTransforms = []p.PostTransform{}
	}
	v.postTransforms = append(v.postTransforms, transform)
	return v
}

// ! MODIFIERS

// marks field as required
func (v *stringProcessor) Required(options ...TestOption) *stringProcessor {
	r := p.Required()
	for _, opt := range options {
		opt(&r)
	}
	v.required = &r
	return v
}

// marks field as optional
func (v *stringProcessor) Optional() *stringProcessor {
	v.required = nil
	return v
}

// sets the default value
func (v *stringProcessor) Default(val string) *stringProcessor {
	v.defaultVal = &val
	return v
}

// sets the catch value (i.e the value to use if the validation fails)
func (v *stringProcessor) Catch(val string) *stringProcessor {
	v.catch = &val
	return v
}

// ! VALIDATORS
// custom test function call it -> schema.Test(t z.Test, opts ...TestOption)
func (v *stringProcessor) Test(t p.Test, opts ...TestOption) *stringProcessor {
	for _, opt := range opts {
		opt(&t)
	}
	v.tests = append(v.tests, t)
	return v
}

// checks that the value is one of the enum values
func (v *stringProcessor) OneOf(enum []string, options ...TestOption) *stringProcessor {
	t := p.In(enum)
	for _, opt := range options {
		opt(&t)
	}
	v.tests = append(v.tests, t)
	return v
}

// checks that the value is at least n characters long
func (v *stringProcessor) Min(n int, options ...TestOption) *stringProcessor {
	t := p.LenMin[string](n)
	for _, opt := range options {
		opt(&t)
	}
	v.tests = append(v.tests, t)
	return v
}

// checks that the value is at most n characters long
func (v *stringProcessor) Max(n int, options ...TestOption) *stringProcessor {
	t := p.LenMax[string](n)
	for _, opt := range options {
		opt(&t)
	}
	v.tests = append(v.tests, t)
	return v
}

// checks that the value is exactly n characters long
func (v *stringProcessor) Len(n int, options ...TestOption) *stringProcessor {
	t := p.Len[string](n)
	for _, opt := range options {
		opt(&t)
	}
	v.tests = append(v.tests, t)
	return v
}

// checks that the value is a valid email address
func (v *stringProcessor) Email(options ...TestOption) *stringProcessor {
	t := p.Test{
		ErrCode: zconst.ErrCodeEmail,
		ValidateFunc: func(v any, ctx ParseCtx) bool {
			email, ok := v.(string)
			if !ok {
				return false
			}
			return emailRegex.MatchString(email)
		},
	}
	for _, opt := range options {
		opt(&t)
	}
	v.tests = append(v.tests, t)
	return v
}

func (v *stringProcessor) URL(options ...TestOption) *stringProcessor {
	t := p.Test{
		ErrCode: zconst.ErrCodeURL,
		ValidateFunc: func(v any, ctx ParseCtx) bool {
			s, ok := v.(string)
			if !ok {
				return false
			}
			u, err := url.Parse(s)
			return err == nil && u.Scheme != "" && u.Host != ""
		},
	}
	for _, opt := range options {
		opt(&t)
	}
	v.tests = append(v.tests, t)
	return v
}

func (v *stringProcessor) HasPrefix(s string, options ...TestOption) *stringProcessor {
	t := p.Test{
		ErrCode: zconst.ErrCodeHasPrefix,
		Params:  make(map[string]any, 1),
		ValidateFunc: func(v any, ctx ParseCtx) bool {
			val, ok := v.(string)
			if !ok {
				return false
			}
			return strings.HasPrefix(val, s)
		},
	}
	t.Params[zconst.ErrCodeHasPrefix] = s
	for _, opt := range options {
		opt(&t)
	}
	v.tests = append(v.tests, t)
	return v
}

func (v *stringProcessor) HasSuffix(s string, options ...TestOption) *stringProcessor {
	t := p.Test{
		ErrCode: zconst.ErrCodeHasSuffix,
		Params:  make(map[string]any, 1),
		ValidateFunc: func(v any, ctx ParseCtx) bool {
			val, ok := v.(string)
			if !ok {
				return false
			}
			return strings.HasSuffix(val, s)
		},
	}
	t.Params[zconst.ErrCodeHasSuffix] = s
	for _, opt := range options {
		opt(&t)
	}
	v.tests = append(v.tests, t)
	return v
}

func (v *stringProcessor) Contains(sub string, options ...TestOption) *stringProcessor {
	t := p.Test{
		ErrCode: zconst.ErrCodeContains,
		Params:  make(map[string]any, 1),
		ValidateFunc: func(v any, ctx ParseCtx) bool {
			val, ok := v.(string)
			if !ok {
				return false
			}
			return strings.Contains(val, sub)
		},
	}
	t.Params[zconst.ErrCodeContains] = sub
	for _, opt := range options {
		opt(&t)
	}
	v.tests = append(v.tests, t)
	return v
}

func (v *stringProcessor) ContainsUpper(options ...TestOption) *stringProcessor {
	t := p.Test{
		ErrCode: zconst.ErrCodeContainsUpper,
		ValidateFunc: func(v any, ctx ParseCtx) bool {
			val, ok := v.(string)
			if !ok {
				return false
			}
			for _, r := range val {
				if r >= 'A' && r <= 'Z' {
					return true
				}
			}
			return false
		},
	}
	for _, opt := range options {
		opt(&t)
	}
	v.tests = append(v.tests, t)
	return v
}

func (v *stringProcessor) ContainsDigit(options ...TestOption) *stringProcessor {
	t := p.Test{
		ErrCode: zconst.ErrCodeContainsDigit,
		ValidateFunc: func(v any, ctx ParseCtx) bool {
			val, ok := v.(string)
			if !ok {
				return false
			}
			for _, r := range val {
				if r >= '0' && r <= '9' {
					return true
				}
			}
			return false
		},
	}

	for _, opt := range options {
		opt(&t)
	}

	v.tests = append(v.tests, t)
	return v
}

func (v *stringProcessor) ContainsSpecial(options ...TestOption) *stringProcessor {
	t :=
		p.Test{
			ErrCode: zconst.ErrCodeContainsSpecial,
			ValidateFunc: func(v any, ctx ParseCtx) bool {
				val, ok := v.(string)
				if !ok {
					return false
				}
				for _, r := range val {
					if (r >= '!' && r <= '/') ||
						(r >= ':' && r <= '@') ||
						(r >= '[' && r <= '`') ||
						(r >= '{' && r <= '~') {
						return true
					}
				}
				return false
			},
		}
	for _, opt := range options {
		opt(&t)
	}
	v.tests = append(v.tests, t)
	return v
}
