// Based on https://github.com/nsf/jsondiff
package jsondiff

import (
	"bytes"
	"encoding/json"
	"reflect"
	"sort"
	"strconv"
)

type Difference int

const (
	FullMatch Difference = iota
	SupersetMatch
	NoMatch
	ValueDiff
	FirstArgIsInvalidJson
	SecondArgIsInvalidJson
	BothArgsAreInvalidJson
)

func (d Difference) String() string {
	switch d {
	case FullMatch:
		return "FullMatch"
	case ValueDiff:
		return "ValueDiff"
	case SupersetMatch:
		return "SupersetMatch"
	case NoMatch:
		return "NoMatch"
	case FirstArgIsInvalidJson:
		return "FirstArgIsInvalidJson"
	case SecondArgIsInvalidJson:
		return "SecondArgIsInvalidJson"
	case BothArgsAreInvalidJson:
		return "BothArgsAreInvalidJson"
	}
	return "Invalid"
}

type Tag struct {
	Begin string
	End   string
}

type Options struct {
	Normal           Tag
	Added            Tag
	Removed          Tag
	Changed          Tag
	Skipped          Tag
	Prefix           string
	Indent           string
	PrintTypes       bool
	ChangedSeparator string
	// When provided, this function will be used to compare two numbers. By default numbers are compared using their
	// literal representation byte by byte.
	CompareNumbers func(a, b json.Number) bool
	SkipValueDiff  bool
}

// Provides a set of options in JSON format that are fully parseable.
func DefaultJSONOptions() Options {
	return Options{
		Added:            Tag{Begin: "\"missing\":{", End: "}"},
		Removed:          Tag{Begin: "\"wrong\":{", End: "}"},
		Changed:          Tag{Begin: "{\"changed\":[", End: "]}"},
		ChangedSeparator: ", ",
		Indent:           "    ",
	}
}

// Provides a set of options that are well suited for console output. Options
// use ANSI foreground color escape sequences to highlight changes.
func DefaultConsoleOptions() Options {
	return Options{
		Added:            Tag{Begin: "\033[0;32m", End: "\033[0m"},
		Removed:          Tag{Begin: "\033[0;31m", End: "\033[0m"},
		Changed:          Tag{Begin: "\033[0;33m", End: "\033[0m"},
		Skipped:          Tag{Begin: "\033[0;90m", End: "\033[0m"},
		ChangedSeparator: " => ",
		Indent:           "    ",
	}
}

// Provides a set of options that are well suited for HTML output. Works best
// inside <pre> tag.
func DefaultHTMLOptions() Options {
	return Options{
		Added:            Tag{Begin: `<span style="background-color: #8bff7f">`, End: `</span>`},
		Removed:          Tag{Begin: `<span style="background-color: #fd7f7f">`, End: `</span>`},
		Changed:          Tag{Begin: `<span style="background-color: #fcff7f">`, End: `</span>`},
		Skipped:          Tag{Begin: `<span style="color: rgba(0, 0, 0, 0.3)">`, End: `</span>`},
		ChangedSeparator: " => ",
		Indent:           "    ",
	}
}

type context struct {
	opts    *Options
	level   int
	lastTag *Tag
	diff    Difference
}

func (ctx *context) compareNumbers(a, b json.Number) bool {
	if ctx.opts.CompareNumbers != nil {
		return ctx.opts.CompareNumbers(a, b)
	} else {
		return a == b
	}
}

func (ctx *context) terminateTag(buf *bytes.Buffer) {
	if ctx.lastTag != nil {
		buf.WriteString(ctx.lastTag.End)
		ctx.lastTag = nil
	}
}

func (ctx *context) newline(buf *bytes.Buffer, s string) {
	buf.WriteString(s)
	if ctx.lastTag != nil {
		buf.WriteString(ctx.lastTag.End)
	}
	buf.WriteString("\n")
	buf.WriteString(ctx.opts.Prefix)
	for i := 0; i < ctx.level; i++ {
		buf.WriteString(ctx.opts.Indent)
	}
	if ctx.lastTag != nil {
		buf.WriteString(ctx.lastTag.Begin)
	}
}

func (ctx *context) key(buf *bytes.Buffer, k string) {
	buf.WriteString(strconv.Quote(k))
	buf.WriteString(": ")
}

func (ctx *context) writeValue(buf *bytes.Buffer, v interface{}, full bool) {
	switch vv := v.(type) {
	case bool:
		buf.WriteString(strconv.FormatBool(vv))
	case json.Number:
		buf.WriteString(string(vv))
	case string:
		buf.WriteString(strconv.Quote(vv))
	case []interface{}:
		if full {
			if len(vv) == 0 {
				buf.WriteString("[")
			} else {
				ctx.level++
				ctx.newline(buf, "[")
			}
			for i, v := range vv {
				ctx.writeValue(buf, v, true)
				if i != len(vv)-1 {
					ctx.newline(buf, ",")
				} else {
					ctx.level--
					ctx.newline(buf, "")
				}
			}
			buf.WriteString("]")
		} else {
			buf.WriteString("[]")
		}
	case map[string]interface{}:
		if full {
			if len(vv) == 0 {
				buf.WriteString("{")
			} else {
				ctx.level++
				ctx.newline(buf, "{")
			}

			keys := make([]string, 0, len(vv))
			for key := range vv {
				keys = append(keys, key)
			}
			sort.Strings(keys)

			i := 0
			for _, k := range keys {
				v := vv[k]
				ctx.key(buf, k)
				ctx.writeValue(buf, v, true)
				if i != len(vv)-1 {
					ctx.newline(buf, ",")
				} else {
					ctx.level--
					ctx.newline(buf, "")
				}
				i++
			}
			buf.WriteString("}")
		} else {
			buf.WriteString("{}")
		}
	default:
		buf.WriteString("null")
	}

	ctx.writeTypeMaybe(buf, v)
}

func (ctx *context) writeTypeMaybe(buf *bytes.Buffer, v interface{}) {
	if ctx.opts.PrintTypes {
		buf.WriteString(" ")
		ctx.writeType(buf, v)
	}
}

func (ctx *context) writeType(buf *bytes.Buffer, v interface{}) {
	switch v.(type) {
	case bool:
		buf.WriteString("(boolean)")
	case json.Number:
		buf.WriteString("(number)")
	case string:
		buf.WriteString("(string)")
	case []interface{}:
		buf.WriteString("(array)")
	case map[string]interface{}:
		buf.WriteString("(object)")
	default:
		buf.WriteString("(null)")
	}
}

func (ctx *context) writeMismatch(buf *bytes.Buffer, a, b interface{}) {
	ctx.writeValue(buf, a, false)
	buf.WriteString(ctx.opts.ChangedSeparator)
	ctx.writeValue(buf, b, false)
}

func (ctx *context) tag(buf *bytes.Buffer, tag *Tag) {
	if ctx.lastTag == tag {
		return
	} else if ctx.lastTag != nil {
		buf.WriteString(ctx.lastTag.End)
	}
	buf.WriteString(tag.Begin)
	ctx.lastTag = tag
}

func (ctx *context) result(d Difference) {
	if d == NoMatch {
		ctx.diff = NoMatch
	} else if d == SupersetMatch && ctx.diff != NoMatch {
		ctx.diff = SupersetMatch
	} else if d == ValueDiff && ctx.diff != NoMatch && ctx.diff != SupersetMatch {
		ctx.diff = ValueDiff
	} else if ctx.diff != NoMatch && ctx.diff != SupersetMatch && ctx.diff != ValueDiff {
		ctx.diff = FullMatch
	}
}

func (ctx *context) printMismatch(buf *bytes.Buffer, a, b interface{}) {
	ctx.tag(buf, &ctx.opts.Changed)
	ctx.writeMismatch(buf, a, b)
}

func (ctx *context) printSkipped(buf *bytes.Buffer, n *int, strfunc func(n int) string, last bool) {
	if *n == 0 || strfunc == nil {
		return
	}
	ctx.tag(buf, &ctx.opts.Skipped)
	buf.WriteString(strfunc(*n))
	if !last {
		ctx.tag(buf, &ctx.opts.Normal)
		ctx.newline(buf, ",")
	}
	*n = 0
}

func (ctx *context) finalize(buf *bytes.Buffer) string {
	ctx.terminateTag(buf)
	return buf.String()
}

type collectionConfig struct {
	open    string
	close   string
	skipped func(n int) string
	value   interface{}
}

type dualIterator interface {
	clone() dualIterator
	count() int
	next() (a interface{}, aOK bool, b interface{}, bOK bool, i int)
	key(buf *bytes.Buffer)
}

type dualSliceIterator struct {
	a       []interface{}
	b       []interface{}
	max     int
	current int
}

func (it *dualSliceIterator) clone() dualIterator {
	copy := *it
	return &copy
}

func (it *dualSliceIterator) count() int {
	return it.max
}

func (it *dualSliceIterator) next() (a interface{}, aOK bool, b interface{}, bOK bool, i int) {
	it.current++
	i = it.current
	if i <= it.max {
		if i < len(it.a) {
			a = it.a[i]
			aOK = true
		}
		if i < len(it.b) {
			b = it.b[i]
			bOK = true
		}
	} else {
		i = -1
	}
	return
}

func (it *dualSliceIterator) key(buf *bytes.Buffer) {
	// noop
}

type dualMapIterator struct {
	a       map[string]interface{}
	b       map[string]interface{}
	keys    []string
	current int
}

func (it *dualMapIterator) clone() dualIterator {
	copy := *it
	return &copy
}

func (it *dualMapIterator) count() int {
	return len(it.keys)
}

func (it *dualMapIterator) next() (a interface{}, aOK bool, b interface{}, bOK bool, i int) {
	it.current++
	i = it.current
	if i < len(it.keys) {
		key := it.keys[i]
		a, aOK = it.a[key]
		b, bOK = it.b[key]
	} else {
		i = -1
	}
	return
}

func (it *dualMapIterator) key(buf *bytes.Buffer) {
	key := it.keys[it.current]
	buf.WriteString(strconv.Quote(key))
	buf.WriteString(": ")
}

func makeDualMapIterator(a, b map[string]interface{}) dualIterator {
	keysMap := make(map[string]struct{})
	for k := range a {
		keysMap[k] = struct{}{}
	}
	for k := range b {
		keysMap[k] = struct{}{}
	}
	keys := make([]string, 0, len(keysMap))
	for k := range keysMap {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return &dualMapIterator{
		a:       a,
		b:       b,
		keys:    keys,
		current: -1,
	}
}

func makeDualSliceIterator(a, b []interface{}) dualIterator {
	max := len(a)
	if len(b) > max {
		max = len(b)
	}
	return &dualSliceIterator{
		a:       a,
		b:       b,
		max:     max,
		current: -1,
	}
}

func (ctx *context) collectDiffs(it dualIterator) (diffs []string, last int) {
	ctx.level++
	last = -1
	for {
		a, aok, b, bok, i := it.next()
		if i == -1 {
			break
		}
		var diff string
		if aok && bok {
			diff = ctx.printDiff(a, b)
		}
		if len(diff) > 0 || aok != bok {
			last = i
		}
		diffs = append(diffs, diff)
	}
	ctx.level--
	return
}

func (ctx *context) printCollectionDiff(cfg *collectionConfig, it dualIterator) string {
	var buf bytes.Buffer
	diffs, lastDiff := ctx.collectDiffs(it.clone())
	if lastDiff == -1 {
		// no diffs
		return ""
	}

	// some diffs or empty collection
	ctx.tag(&buf, &ctx.opts.Normal)
	if it.count() == 0 {
		buf.WriteString(cfg.open)
		buf.WriteString(cfg.close)
		ctx.writeTypeMaybe(&buf, cfg.value)
		return ctx.finalize(&buf)
	} else {
		ctx.level++
		ctx.newline(&buf, cfg.open)
	}

	noDiffSpan := 0
	for {
		va, aok, vb, bok, i := it.next()
		equals := true
		if aok && bok {
			diff := diffs[i]
			if len(diff) > 0 {
				equals = false
				ctx.printSkipped(&buf, &noDiffSpan, cfg.skipped, false)
				it.key(&buf)
				buf.WriteString(diff)
			}
		} else if aok {
			equals = false
			ctx.printSkipped(&buf, &noDiffSpan, cfg.skipped, false)
			ctx.tag(&buf, &ctx.opts.Removed)
			it.key(&buf)
			ctx.writeValue(&buf, va, true)
			ctx.result(SupersetMatch)
		} else if bok {
			equals = false
			ctx.printSkipped(&buf, &noDiffSpan, cfg.skipped, false)
			ctx.tag(&buf, &ctx.opts.Added)
			it.key(&buf)
			ctx.writeValue(&buf, vb, true)
			ctx.result(NoMatch)
		}
		if equals {
			noDiffSpan++
		}

		wroteItem := !equals
		willWriteMoreItems := (i < lastDiff) || (cfg.skipped != nil && lastDiff < it.count()-1)

		if wroteItem && willWriteMoreItems {
			ctx.tag(&buf, &ctx.opts.Normal)
			ctx.newline(&buf, ",")
		}
		if i == it.count()-1 {
			// we're done
			ctx.printSkipped(&buf, &noDiffSpan, cfg.skipped, true)
			ctx.level--
			ctx.tag(&buf, &ctx.opts.Normal)
			ctx.newline(&buf, "")
			break
		}
	}

	buf.WriteString(cfg.close)
	ctx.writeTypeMaybe(&buf, cfg.value)
	return ctx.finalize(&buf)
}

func Compare(b, a []byte, opts *Options) (Difference, string) {
	var av, bv interface{}
	da := json.NewDecoder(bytes.NewReader(a))
	da.UseNumber()
	db := json.NewDecoder(bytes.NewReader(b))
	db.UseNumber()
	errA := da.Decode(&av)
	errB := db.Decode(&bv)
	if errA != nil && errB != nil {
		return BothArgsAreInvalidJson, "both arguments are invalid json"
	}
	if errA != nil {
		return FirstArgIsInvalidJson, "first argument is invalid json"
	}
	if errB != nil {
		return SecondArgIsInvalidJson, "second argument is invalid json"
	}

	var buf bytes.Buffer

	ctx := context{opts: opts}
	buf.WriteString(ctx.printDiff(av, bv))
	return ctx.diff, buf.String()
}

func (ctx *context) printDiff(a, b interface{}) string {
	var buf bytes.Buffer

	if a == nil || b == nil {
		if a == nil && b == nil {
			// noop
		} else if !ctx.opts.SkipValueDiff {
			// mismatch
			ctx.printMismatch(&buf, a, b)
			ctx.result(ValueDiff)
		}
		return ctx.finalize(&buf)
	}

	ka := reflect.TypeOf(a).Kind()
	kb := reflect.TypeOf(b).Kind()
	if ka != kb {
		// Go type does not match
		ctx.printMismatch(&buf, a, b)
		ctx.result(NoMatch)
		return ctx.finalize(&buf)
	}

	// big switch here handles type-specific mismatches and returns if that's the case
	// buf if control flow goes past through this switch, it's a match
	// NOTE: ka == kb at this point
	switch ka {
	case reflect.Bool:
		if a.(bool) != b.(bool) {
			if !ctx.opts.SkipValueDiff {
				ctx.printMismatch(&buf, a, b)
				ctx.result(ValueDiff)
				return ctx.finalize(&buf)
			}
		}
	case reflect.String:
		// string can be a json.Number here too (because it's a string type)
		switch aa := a.(type) {
		case json.Number:
			bb, ok := b.(json.Number)
			if !ok || !ctx.compareNumbers(aa, bb) {
				if !ctx.opts.SkipValueDiff {
					ctx.printMismatch(&buf, a, b)
					ctx.result(ValueDiff)
					return ctx.finalize(&buf)
				}
			}
		case string:
			bb, ok := b.(string)
			if !ok || aa != bb {
				if !ctx.opts.SkipValueDiff {
					ctx.printMismatch(&buf, a, b)
					ctx.result(ValueDiff)
					return ctx.finalize(&buf)
				}
			}
		}
	case reflect.Slice:
		sa, sb := a.([]interface{}), b.([]interface{})
		return ctx.printCollectionDiff(&collectionConfig{
			open:  "[",
			close: "]",
			value: a,
		}, makeDualSliceIterator(sa, sb))
	case reflect.Map:
		ma, mb := a.(map[string]interface{}), b.(map[string]interface{})
		return ctx.printCollectionDiff(&collectionConfig{
			open:  "{",
			close: "}",
			value: a,
		}, makeDualMapIterator(ma, mb))
	}
	return ctx.finalize(&buf)
}
