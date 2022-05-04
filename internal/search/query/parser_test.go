package query

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/hexops/autogold"

	"github.com/sourcegraph/sourcegraph/lib/errors"
)

func collectLabels(nodes []Node) (result labels) {
	for _, node := range nodes {
		switch v := node.(type) {
		case Operator:
			result |= v.Annotation.Labels
			result |= collectLabels(v.Operands)
		case Pattern:
			result |= v.Annotation.Labels
		case Parameter:
			result |= v.Annotation.Labels
		}
	}
	return result
}

func labelsToString(nodes []Node) string {
	labels := collectLabels(nodes)
	return strings.Join(labels.String(), ",")
}

func TestParseParameterList(t *testing.T) {
	type value struct {
		Result       string
		ResultLabels string
		ResultRange  string
	}

	test := func(input string) value {
		parser := &parser{buf: []byte(input), heuristics: parensAsPatterns | allowDanglingParens}
		result, err := parser.parseLeaves(Regexp)
		if err != nil {
			t.Fatal(fmt.Sprintf("Unexpected error: %s", err))
		}
		resultNode := result[0]
		got, _ := json.Marshal(resultNode)

		var gotRange string
		switch n := resultNode.(type) {
		case Pattern:
			gotRange = n.Annotation.Range.String()
		case Parameter:
			gotRange = n.Annotation.Range.String()
		}

		var gotLabels string
		if _, ok := resultNode.(Pattern); ok {
			gotLabels = labelsToString([]Node{resultNode})
		}

		return value{
			Result:       string(got),
			ResultLabels: gotLabels,
			ResultRange:  gotRange,
		}
	}

	autogold.Want("Normal field:value", value{
		Result:      `{"field":"file","value":"README.md","negated":false}`,
		ResultRange: `{"start":{"line":0,"column":0},"end":{"line":0,"column":14}}`,
	}).Equal(t, test(`file:README.md`))

	autogold.Want("Normal field:value with trailing space", value{
		Result:      `{"field":"file","value":"README.md","negated":false}`,
		ResultRange: `{"start":{"line":0,"column":0},"end":{"line":0,"column":14}}`,
	}).Equal(t, test(`file:README.md    `))

	autogold.Want("First char is colon", value{
		Result: `{"value":":foo","negated":false}`, ResultLabels: "Regexp",
		ResultRange: `{"start":{"line":0,"column":0},"end":{"line":0,"column":4}}`,
	}).Equal(t, test(`:foo`))

	autogold.Want("Last char is colon", value{
		Result: `{"value":"foo:","negated":false}`, ResultLabels: "Regexp",
		ResultRange: `{"start":{"line":0,"column":0},"end":{"line":0,"column":4}}`,
	}).Equal(t, test(`foo:`))

	autogold.Want("Match first colon", value{
		Result:      `{"field":"file","value":"bar:baz","negated":false}`,
		ResultRange: `{"start":{"line":0,"column":0},"end":{"line":0,"column":12}}`,
	}).Equal(t, test(`file:bar:baz`))

	autogold.Want("No field, start with minus", value{
		Result: `{"value":"-:foo","negated":false}`, ResultLabels: "Regexp",
		ResultRange: `{"start":{"line":0,"column":0},"end":{"line":0,"column":5}}`,
	}).Equal(t, test(`-:foo`))

	autogold.Want("Minus prefix on field", value{
		Result:      `{"field":"file","value":"README.md","negated":true}`,
		ResultRange: `{"start":{"line":0,"column":0},"end":{"line":0,"column":15}}`,
	}).Equal(t, test(`-file:README.md`))

	autogold.Want("NOT prefix on file", value{
		Result:      `{"field":"file","value":"README.md","negated":true}`,
		ResultRange: `{"start":{"line":0,"column":0},"end":{"line":0,"column":18}}`,
	}).Equal(t, test(`NOT file:README.md`))

	autogold.Want("NOT prefix on unsupported key-value pair", value{
		Result: `{"value":"foo:bar","negated":true}`, ResultLabels: "Regexp",
		ResultRange: `{"start":{"line":0,"column":0},"end":{"line":0,"column":11}}`,
	}).Equal(t, test(`NOT foo:bar`))
	autogold.Want("NOT prefix on content", value{
		Result:      `{"field":"content","value":"bar","negated":true}`,
		ResultRange: `{"start":{"line":0,"column":0},"end":{"line":0,"column":15}}`,
	}).Equal(t, test(`NOT content:bar`))

	autogold.Want("Double NOT", value{
		Result: `{"value":"NOT","negated":true}`, ResultLabels: "Regexp",
		ResultRange: `{"start":{"line":0,"column":0},"end":{"line":0,"column":7}}`,
	}).Equal(t, test(`NOT NOT`))

	autogold.Want("Double minus prefix on field", value{
		Result:       `{"value":"--foo:bar","negated":false}`,
		ResultLabels: "Regexp",
		ResultRange:  `{"start":{"line":0,"column":0},"end":{"line":0,"column":9}}`,
	}).Equal(t, test(`--foo:bar`))

	autogold.Want("Minus in the middle is not a valid field", value{
		Result:       `{"value":"fie-ld:bar","negated":false}`,
		ResultLabels: "Regexp",
		ResultRange:  `{"start":{"line":0,"column":0},"end":{"line":0,"column":10}}`,
	}).Equal(t, test(`fie-ld:bar`))

	autogold.Want("Preserve escaped whitespace", value{
		Result:       `{"value":"a\\ pattern","negated":false}`,
		ResultLabels: "Regexp",
		ResultRange:  `{"start":{"line":0,"column":0},"end":{"line":0,"column":10}}`,
	}).Equal(t, test(`a\ pattern`))

	autogold.Want("Quoted", value{
		Result: `{"value":"quoted","negated":false}`, ResultLabels: "Literal,Quoted",
		ResultRange: `{"start":{"line":0,"column":0},"end":{"line":0,"column":8}}`,
	}).Equal(t, test(`"quoted"`))

	autogold.Want("Escaped quote", value{
		Result: `{"value":"'","negated":false}`, ResultLabels: "Literal,Quoted",
		ResultRange: `{"start":{"line":0,"column":0},"end":{"line":0,"column":4}}`,
	}).Equal(t, test(`'\''`))

	autogold.Want("Regexp syntax with unbalanced paren", value{
		Result:       `{"value":"foo.*bar(","negated":false}`,
		ResultLabels: "HeuristicDanglingParens,Regexp",
		ResultRange:  `{"start":{"line":0,"column":0},"end":{"line":0,"column":9}}`,
	}).Equal(t, test(`foo.*bar(`))

	autogold.Want("Regexp delimiters", value{
		Result:       `{"value":"a regex pattern","negated":false}`,
		ResultLabels: "Regexp",
		ResultRange:  `{"start":{"line":0,"column":0},"end":{"line":0,"column":17}}`,
	}).Equal(t, test(`/a regex pattern/`))

	autogold.Want("Regexp group", value{
		Result:       `{"value":"Search()\\(","negated":false}`,
		ResultLabels: "Regexp",
		ResultRange:  `{"start":{"line":0,"column":0},"end":{"line":0,"column":10}}`,
	}).Equal(t, test(`Search()\(`))

	autogold.Want("Regexp non-empty group", value{
		Result:       `{"value":"Search(xxx)\\\\(","negated":false}`,
		ResultLabels: "HeuristicDanglingParens,Regexp",
		ResultRange:  `{"start":{"line":0,"column":0},"end":{"line":0,"column":14}}`,
	}).Equal(t, test(`Search(xxx)\\(`))

	autogold.Want("Regexp non-empty /.../", value{
		Result: `{"value":"book","negated":false}`, ResultLabels: "Regexp",
		ResultRange: `{"start":{"line":0,"column":0},"end":{"line":0,"column":6}}`,
	}).Equal(t, test(`/book/`))

	autogold.Want("Regexp empty /.../", value{
		Result: `{"value":"//","negated":false}`, ResultLabels: "Literal",
		ResultRange: `{"start":{"line":0,"column":0},"end":{"line":0,"column":2}}`,
	}).Equal(t, test(`//`))
}

func TestScanPredicate(t *testing.T) {
	type value struct {
		Result       string
		ResultLabels string
	}

	test := func(input string) value {
		parser := &parser{buf: []byte(input), heuristics: parensAsPatterns | allowDanglingParens}
		result, err := parser.parseLeaves(Regexp)
		if err != nil {
			t.Fatal(fmt.Sprintf("Unexpected error: %s", err))
		}
		resultNode := result[0]
		got, _ := json.Marshal(resultNode)
		gotLabels := labelsToString([]Node{resultNode})

		return value{
			Result:       string(got),
			ResultLabels: gotLabels,
		}
	}

	autogold.Want("Repo contains predicate", value{
		Result:       `{"field":"repo","value":"contains(file:test)","negated":false}`,
		ResultLabels: "IsPredicate",
	}).Equal(t, test(`repo:contains(file:test)`))

	autogold.Want("Repo contains commit after predicate", value{
		Result:       `{"field":"repo","value":"contains.commit.after(last thursday)","negated":false}`,
		ResultLabels: "IsPredicate",
	}).Equal(t, test(`repo:contains.commit.after(last thursday)`))

	autogold.Want("Repo contains commit before predicate does not exist", value{
		Result:       `{"field":"repo","value":"contains.commit.before(yesterday)","negated":false}`,
		ResultLabels: "None",
	}).Equal(t, test(`repo:contains.commit.before(yesterday)`))

	autogold.Want("Predicate contains escaped paranthesis", value{
		Result:       `{"field":"repo","value":"contains(\\()","negated":false}`,
		ResultLabels: "IsPredicate",
	}).Equal(t, test(`repo:contains(\()`))

	autogold.Want("Repo contains not predicate", value{
		Result:       `{"field":"repo","value":"contains","negated":false}`,
		ResultLabels: "None",
	}).Equal(t, test(`repo:contains`))

	autogold.Want("Repo with something that looks kinda like predicate", value{
		Result:       `{"Kind":1,"Operands":[{"field":"repo","value":"nopredicate","negated":false},{"value":"(file:foo","negated":false}],"Annotation":{"labels":0,"range":{"start":{"line":0,"column":0},"end":{"line":0,"column":0}}}}`,
		ResultLabels: "HeuristicDanglingParens,Regexp",
	}).Equal(t, test(`repo:nopredicate(file:foo or file:bar)`))

	autogold.Want("Pattern looks like predicate", value{
		Result:       `{"Kind":2,"Operands":[{"value":"abc","negated":false},{"value":"contains(file:test)","negated":false}],"Annotation":{"labels":0,"range":{"start":{"line":0,"column":0},"end":{"line":0,"column":0}}}}`,
		ResultLabels: "HeuristicDanglingParens,Regexp",
	}).Equal(t, test(`abc contains(file:test)`))

	autogold.Want("Resolve field aliases for predicates", value{
		Result:       `{"field":"r","value":"contains.file(sup)","negated":false}`,
		ResultLabels: "IsPredicate",
	}).Equal(t, test(`r:contains.file(sup)`))
}

func TestScanField(t *testing.T) {
	type value struct {
		Field   string
		Negated bool
		Advance int
	}

	test := func(input string) string {
		gotField, gotNegated, gotAdvance := ScanField([]byte(input))
		v, _ := json.Marshal(value{gotField, gotNegated, gotAdvance})
		return string(v)
	}

	autogold.Want("repo:foo", `{"Field":"repo","Negated":false,"Advance":5}`).Equal(t, test("repo:foo"))
	autogold.Want("RepO:foo", `{"Field":"RepO","Negated":false,"Advance":5}`).Equal(t, test("RepO:foo"))
	autogold.Want("after:", `{"Field":"after","Negated":false,"Advance":6}`).Equal(t, test("after:"))
	autogold.Want("-repo:", `{"Field":"repo","Negated":true,"Advance":6}`).Equal(t, test("-repo:"))
	autogold.Want("", `{"Field":"","Negated":false,"Advance":0}`).Equal(t, test(""))
	autogold.Want("-", `{"Field":"","Negated":false,"Advance":0}`).Equal(t, test("-"))
	autogold.Want("-:", `{"Field":"","Negated":false,"Advance":0}`).Equal(t, test("-:"))
	autogold.Want(":", `{"Field":"","Negated":false,"Advance":0}`).Equal(t, test(":"))
	autogold.Want("??:foo", `{"Field":"","Negated":false,"Advance":0}`).Equal(t, test("??:foo"))
	autogold.Want("repo", `{"Field":"","Negated":false,"Advance":0}`).Equal(t, test("repo"))
	autogold.Want("-repo", `{"Field":"","Negated":false,"Advance":0}`).Equal(t, test("-repo"))
	autogold.Want("--repo:", `{"Field":"","Negated":false,"Advance":0}`).Equal(t, test("--repo:"))
	autogold.Want(":foo", `{"Field":"","Negated":false,"Advance":0}`).Equal(t, test(":foo"))
}

func parseAndOrGrammar(in string) ([]Node, error) {
	if strings.TrimSpace(in) == "" {
		return nil, nil
	}
	parser := &parser{
		buf:        []byte(in),
		leafParser: SearchTypeRegex,
	}
	nodes, err := parser.parseOr()
	if err != nil {
		return nil, err
	}
	if parser.balanced != 0 {
		return nil, errors.New("unbalanced expression: unmatched closing parenthesis )")
	}
	return newOperator(nodes, And), nil
}

func TestParse(t *testing.T) {
	type value struct {
		Grammar   string
		Heuristic string
	}

	test := func(input string) value {
		var queryGrammar, queryHeuristic []Node
		var err error
		var resultGrammar, resultHeuristic string
		queryGrammar, err = parseAndOrGrammar(input) // Parse without heuristic.
		if err != nil {
			resultGrammar = err.Error()
		} else {
			resultGrammar = toString(queryGrammar)
		}

		queryHeuristic, err = Parse(input, SearchTypeRegex)
		if err != nil {
			resultHeuristic = err.Error()
		} else {
			resultHeuristic = toString(queryHeuristic)
		}

		if resultHeuristic == resultGrammar {
			resultHeuristic = "Same"
		}

		return value{
			Grammar:   resultGrammar,
			Heuristic: resultHeuristic,
		}
	}

	autogold.Want("Empty string", value{Grammar: "", Heuristic: "Same"}).Equal(t, test(""))
	autogold.Want("Whitespace", value{Grammar: "", Heuristic: "Same"}).Equal(t, test("             "))
	autogold.Want("Single", value{Grammar: `"a"`, Heuristic: "Same"}).Equal(t, test("a"))
	autogold.Want("Whitespace basic", value{Grammar: `(concat "a" "b")`, Heuristic: "Same"}).Equal(t, test("a b"))
	autogold.Want("Basic", value{Grammar: `(and "a" "b" "c")`, Heuristic: "Same"}).Equal(t, test("a and b and c"))

	autogold.Want("(f(x)oo((a|b))bar)", value{
		Grammar:   `(concat "f" "x" "oo" "a|b" "bar")`,
		Heuristic: `"(f(x)oo((a|b))bar)"`,
	}).Equal(t, test("(f(x)oo((a|b))bar)"))

	autogold.Want("aorb", value{Grammar: `"aorb"`, Heuristic: "Same"}).Equal(t, test("aorb"))
	autogold.Want("aANDb", value{Grammar: `"aANDb"`, Heuristic: "Same"}).Equal(t, test("aANDb"))
	autogold.Want("a oror b", value{Grammar: `(concat "a" "oror" "b")`, Heuristic: "Same"}).Equal(t, test("a oror b"))

	autogold.Want("Reduced complex query mixed caps", value{
		Grammar:   `(or (and "a" "b" "c") (and "d" (concat (or "e" "f") "g" "h" "i")) "j")`,
		Heuristic: "Same",
	}).Equal(t, test("a and b AND c or d and (e OR f) g h i or j"))

	autogold.Want("Basic reduced complex query", value{
		Grammar:   `(or (and "a" "b") (and "c" "d") "e")`,
		Heuristic: "Same",
	}).Equal(t, test("a and b or c and d or e"))

	autogold.Want("Reduced complex query, reduction over parens", value{
		Grammar:   `(or (and "a" "b") (and "c" "d") "e")`,
		Heuristic: "Same",
	}).Equal(t, test("(a and b or c and d) or e"))

	autogold.Want("Reduced complex query, nested 'or' trickles up", value{Grammar: `(or (and "a" "b") "c" "d")`, Heuristic: "Same"}).Equal(t, test("(a and b or c) or d"))

	autogold.Want("Reduced complex query, nested nested 'or' trickles up", value{
		Grammar:   `(or (and "a" "b") (and "c" "d") "f" "e")`,
		Heuristic: "Same",
	}).Equal(t, test("(a and b or (c and d or f)) or e"))

	autogold.Want("No reduction on precedence defined by parens", value{
		Grammar:   `(or (and "a" (or "b" "c") "d") "e")`,
		Heuristic: "Same",
	}).Equal(t, test("(a and (b or c) and d) or e"))

	autogold.Want("Paren reduction over operators", value{Grammar: `(and (concat "a" "b" "c") "d")`, Heuristic: `(and "(((a b c)))" "d")`}).Equal(t, test("(((a b c))) and d"))

	// Partition parameters and concatenated patterns.
	autogold.Want("a (b and c) d", value{Grammar: `(concat "a" (and "b" "c") "d")`, Heuristic: "Same"}).Equal(t, test("a (b and c) d"))

	autogold.Want("(a b c) and (d e f) and (g h i)", value{
		Grammar:   `(and (concat "a" "b" "c") (concat "d" "e" "f") (concat "g" "h" "i"))`,
		Heuristic: `(and "(a b c)" "(d e f)" "(g h i)")`,
	}).Equal(t, test("(a b c) and (d e f) and (g h i)"))

	autogold.Want("(a) repo:foo (b)", value{
		Grammar:   `(and "repo:foo" (concat "a" "b"))`,
		Heuristic: `(and "repo:foo" (concat "(a)" "(b)"))`,
	}).Equal(t, test("(a) repo:foo (b)"))

	autogold.Want("repo:foo func( or func(.*)", value{Grammar: "expected operand at 15", Heuristic: `(and "repo:foo" (or "func(" "func(.*)"))`}).Equal(t, test("repo:foo func( or func(.*)"))

	autogold.Want("repo:foo main { and bar {", value{
		Grammar:   `(and (and "repo:foo" (concat "main" "{")) (concat "bar" "{"))`,
		Heuristic: `(and "repo:foo" (concat "main" "{") (concat "bar" "{"))`,
	}).Equal(t, test("repo:foo main { and bar {"))

	autogold.Want("a b (repo:foo c d)", value{
		Grammar:   `(concat "a" "b" (and "repo:foo" (concat "c" "d")))`,
		Heuristic: "Same",
	}).Equal(t, test("a b (repo:foo c d)"))

	autogold.Want("a b (c d repo:foo)", value{
		Grammar:   `(concat "a" "b" (and "repo:foo" (concat "c" "d")))`,
		Heuristic: "Same",
	}).Equal(t, test("a b (c d repo:foo)"))

	autogold.Want("a repo:b repo:c (d repo:e repo:f)", value{
		Grammar:   `(and "repo:b" "repo:c" (concat "a" (and "repo:e" "repo:f" "d")))`,
		Heuristic: "Same",
	}).Equal(t, test("a repo:b repo:c (d repo:e repo:f)"))

	autogold.Want("a repo:b repo:c (repo:e repo:f (repo:g repo:h))", value{
		Grammar:   `(and "repo:b" "repo:c" "repo:e" "repo:f" "repo:g" "repo:h" "a")`,
		Heuristic: "Same",
	}).Equal(t, test("a repo:b repo:c (repo:e repo:f (repo:g repo:h))"))

	autogold.Want("a repo:b repo:c (repo:e repo:f (repo:g repo:h)) b", value{
		Grammar:   `(and "repo:b" "repo:c" "repo:e" "repo:f" "repo:g" "repo:h" (concat "a" "b"))`,
		Heuristic: "Same",
	}).Equal(t, test("a repo:b repo:c (repo:e repo:f (repo:g repo:h)) b"))
	autogold.Want("a repo:b repo:c (repo:e repo:f (repo:g repo:h b)) ", value{
		Grammar:   `(and "repo:b" "repo:c" (concat "a" (and "repo:e" "repo:f" "repo:g" "repo:h" "b")))`,
		Heuristic: "Same",
	}).Equal(t, test("a repo:b repo:c (repo:e repo:f (repo:g repo:h b)) "))

	autogold.Want("(repo:foo a (repo:bar b (repo:qux c)))", value{
		Grammar:   `(and "repo:foo" (concat "a" (and "repo:bar" (concat "b" (and "repo:qux" "c")))))`,
		Heuristic: "Same",
	}).Equal(t, test("(repo:foo a (repo:bar b (repo:qux c)))"))

	autogold.Want("a repo:b repo:c (d repo:e repo:f e)", value{
		Grammar:   `(and "repo:b" "repo:c" (concat "a" (and "repo:e" "repo:f" (concat "d" "e"))))`,
		Heuristic: "Same",
	}).Equal(t, test("a repo:b repo:c (d repo:e repo:f e)"))

	// Errors.
	autogold.Want("Unbalanced", value{
		Grammar:   "unbalanced expression: unmatched closing parenthesis )",
		Heuristic: `(concat "(foo)" "(bar")`,
	}).Equal(t, test("(foo) (bar"))
	autogold.Want("Illegal expression on the right", value{Grammar: "expected operand at 5", Heuristic: "Same"}).Equal(t, test("a or or b"))
	autogold.Want("Illegal expression on the right, mixed operators", value{Grammar: `(and "a" "OR")`, Heuristic: "Same"}).Equal(t, test("a and OR"))
	autogold.Want("paren reduction with ands", value{Grammar: `(and "a" "b" "c" "d")`, Heuristic: "Same"}).Equal(t, test("(a and b) and (c and d)"))
	autogold.Want("paren reduction with ors", value{Grammar: `(or "a" "b" "c" "d")`, Heuristic: "Same"}).Equal(t, test("(a or b) or (c or d)"))
	autogold.Want("nested paren reduction with whitespace", value{Grammar: `(concat "a" "b" "c" "d")`, Heuristic: `(concat "(((a b c)))" "d")`}).Equal(t, test("(((a b c))) d"))
	autogold.Want("left paren reduction with whitespace", value{Grammar: `(concat "a" "b" "c" "d")`, Heuristic: `(concat "(a b)" "c" "d")`}).Equal(t, test("(a b) c d"))
	autogold.Want("right paren reduction with whitespace", value{Grammar: `(concat "a" "b" "c" "d")`, Heuristic: `(concat "a" "b" "(c d)")`}).Equal(t, test("a b (c d)"))
	autogold.Want("grouped paren reduction with whitespace", value{Grammar: `(concat "a" "b" "c" "d")`, Heuristic: `(concat "(a b)" "(c d)")`}).Equal(t, test("(a b) (c d)"))

	// Escaping.
	autogold.Want("multiple grouped paren reduction with whitespace", value{Grammar: `(concat "a" "b" "c" "d" "e" "f")`, Heuristic: `(concat "(a b)" "(c d)" "(e f)")`}).Equal(t, test("(a b) (c d) (e f)"))

	autogold.Want("interpolated grouped paren reduction", value{Grammar: `(concat "a" "b" "c" "d" "e" "f")`, Heuristic: `(concat "(a b)" "c" "d" "(e f)")`}).Equal(t, test("(a b) c d (e f)"))

	autogold.Want("mixed interpolated grouped paren reduction", value{
		Grammar:   `(and "a" "b" (or "z" "q") "c" "d" "e" "f")`,
		Heuristic: "Same",
	}).Equal(t, test("(a and b and (z or q)) and (c and d) and (e and f)"))

	autogold.Want("empty paren", value{Grammar: `""`, Heuristic: `"()"`}).Equal(t, test("()"))
	autogold.Want("paren inside contiguous string", value{Grammar: `(concat "foo" "bar")`, Heuristic: `"foo()bar"`}).Equal(t, test("foo()bar"))
	autogold.Want("paren inside contiguous string with and", value{
		Grammar:   `(and "x" (concat "regex" "s" "?"))`,
		Heuristic: `(and "x" "regex(s)?")`,
	}).Equal(t, test("(x and regex(s)?)"))

	autogold.Want("paren containing whitespace inside contiguous string", value{Grammar: `(concat "foo" "bar")`, Heuristic: `"foo(   )bar"`}).Equal(t, test("foo(   )bar"))
	autogold.Want("nested empty paren", value{Grammar: `"x"`, Heuristic: `"(x())"`}).Equal(t, test("(x())"))
	autogold.Want("interpolated nested empty paren", value{Grammar: `"x"`, Heuristic: `"(()x(  )(())())"`}).Equal(t, test("(()x(  )(())())"))
	autogold.Want("empty paren on or", value{Grammar: `""`, Heuristic: `(or "()" "()")`}).Equal(t, test("() or ()"))
	autogold.Want("empty left paren on or", value{Grammar: `"x"`, Heuristic: `(or "()" "(x)")`}).Equal(t, test("() or (x)"))
	autogold.Want("complex interpolated nested empty paren", value{Grammar: `(concat "x" (or "y" "f"))`, Heuristic: `(concat "()" "x" "()" (or "y" "()" "(f)") "()")`}).Equal(t, test("(()x(  )(y or () or (f))())"))
	autogold.Want("disable parens as patterns heuristic if containing recognized operator", value{Grammar: `""`, Heuristic: `(or "()" "()")`}).Equal(t, test("(() or ())"))

	autogold.Want("NOT expression inside parentheses", value{
		Grammar:   `(and "r:foo" (concat "a/foo" (not ".svg")))`,
		Heuristic: "Same",
	}).Equal(t, test("r:foo (a/foo not .svg)"))

	autogold.Want("NOT expression at begining of group", value{Grammar: `(and "r:foo" (not ".svg"))`, Heuristic: "Same"}).Equal(t, test("r:foo (not .svg)"))

	// Escaping
	autogold.Want(`\(\)`, value{Grammar: `"\\(\\)"`, Heuristic: "Same"}).Equal(t, test(`\(\)`))
	autogold.Want(`\( \) ()`, value{Grammar: `(concat "\\(" "\\)")`, Heuristic: `(concat "\\(" "\\)" "()")`}).Equal(t, test(`\( \) ()`))
	autogold.Want(`\ `, value{Grammar: `"\\ "`, Heuristic: "Same"}).Equal(t, test(`\ `))
	autogold.Want(`\  \ `, value{Grammar: `(concat "\\ " "\\ ")`, Heuristic: "Same"}).Equal(t, test(`\  \ `))

	// Dangling parentheses heuristic.
	autogold.Want(`(`, value{Grammar: "expected operand at 1", Heuristic: `"("`}).Equal(t, test(`(`))
	autogold.Want(`)(())(`, value{
		Grammar:   "unsupported expression. The combination of parentheses in the query have an unclear meaning. Try using the content: filter to quote patterns that contain parentheses",
		Heuristic: "Same",
	}).Equal(t, test(`)(())(`))
	autogold.Want(`foo( and bar(`, value{Grammar: "expected operand at 5", Heuristic: `(and "foo(" "bar(")`}).Equal(t, test(`foo( and bar(`))
	autogold.Want(`repo:foo foo( or bar(`, value{Grammar: "expected operand at 14", Heuristic: `(and "repo:foo" (or "foo(" "bar("))`}).Equal(t, test(`repo:foo foo( or bar(`))
	autogold.Want(`(a or (b and )) or d)`, value{
		Grammar:   "unsupported expression. The combination of parentheses in the query have an unclear meaning. Try using the content: filter to quote patterns that contain parentheses",
		Heuristic: "Same",
	}).Equal(t, test(`(a or (b and )) or d)`))

	// Quotes and escape sequences.
	autogold.Want(`"`, value{Grammar: `"\""`, Heuristic: "Same"}).Equal(t, test(`"`))
	autogold.Want(`repo:foo' bar'`, value{Grammar: `(and "repo:foo'" "bar'")`, Heuristic: "Same"}).Equal(t, test(`repo:foo' bar'`))
	autogold.Want(`repo:'foo' 'bar'`, value{Grammar: `(and "repo:foo" "bar")`, Heuristic: "Same"}).Equal(t, test(`repo:'foo' 'bar'`))
	autogold.Want(`repo:"foo" "bar"`, value{Grammar: `(and "repo:foo" "bar")`, Heuristic: "Same"}).Equal(t, test(`repo:"foo" "bar"`))
	autogold.Want(`repo:"foo bar" "foo bar"`, value{Grammar: `(and "repo:foo bar" "foo bar")`, Heuristic: "Same"}).Equal(t, test(`repo:"foo bar" "foo bar"`))
	autogold.Want(`repo:"fo\"o" "bar"`, value{Grammar: `(and "repo:fo\"o" "bar")`, Heuristic: "Same"}).Equal(t, test(`repo:"fo\"o" "bar"`))
	autogold.Want(`repo:foo /b\/ar/`, value{Grammar: `(and "repo:foo" "b/ar")`, Heuristic: "Same"}).Equal(t, test(`repo:foo /b\/ar/`))
	autogold.Want(`repo:foo /a/file/path`, value{Grammar: `(and "repo:foo" "/a/file/path")`, Heuristic: "Same"}).Equal(t, test(`repo:foo /a/file/path`))
	autogold.Want(`repo:foo /a/file/path/`, value{Grammar: `(and "repo:foo" "/a/file/path/")`, Heuristic: "Same"}).Equal(t, test(`repo:foo /a/file/path/`))

	autogold.Want(`repo:foo /a/ /another/path/`, value{
		Grammar:   `(and "repo:foo" (concat "a" "/another/path/"))`,
		Heuristic: "Same",
	}).Equal(t, test(`repo:foo /a/ /another/path/`))

	autogold.Want(`repo:foo /\s+b\d+ar/ `, value{Grammar: `(and "repo:foo" "\\s+b\\d+ar")`, Heuristic: "Same"}).Equal(t, test(`repo:foo /\s+b\d+ar/ `))
	autogold.Want(`repo:foo /bar/ `, value{Grammar: `(and "repo:foo" "bar")`, Heuristic: "Same"}).Equal(t, test(`repo:foo /bar/ `))
	autogold.Want(`\t\r\n`, value{Grammar: `"\\t\\r\\n"`, Heuristic: "Same"}).Equal(t, test(`\t\r\n`))
	autogold.Want(`repo:foo\ bar \:\\`, value{Grammar: `(and "repo:foo\\ bar" "\\:\\\\")`, Heuristic: "Same"}).Equal(t, test(`repo:foo\ bar \:\\`))

	autogold.Want(`a file:\.(ts(?:(?:)|x)|js(?:(?:)|x))(?m:$)`, value{
		Grammar:   `(and "file:\\.(ts(?:(?:)|x)|js(?:(?:)|x))(?m:$)" "a")`,
		Heuristic: "Same",
	}).Equal(t, test(`a file:\.(ts(?:(?:)|x)|js(?:(?:)|x))(?m:$)`))

	autogold.Want(`(file:(a) file:(b))`, value{Grammar: `(and "file:(a)" "file:(b)")`, Heuristic: "Same"}).Equal(t, test(`(file:(a) file:(b))`))
	autogold.Want(`(repohascommitafter:"7 days")`, value{Grammar: `"repohascommitafter:7 days"`, Heuristic: "Same"}).Equal(t, test(`(repohascommitafter:"7 days")`))

	autogold.Want(`(foo repohascommitafter:"7 days")`, value{
		Grammar:   `(and "repohascommitafter:7 days" "foo")`,
		Heuristic: "Same",
	}).Equal(t, test(`(foo repohascommitafter:"7 days")`))

	// Fringe tests cases at the boundary of heuristics and invalid syntax.
	autogold.Want(`(0(F)(:())(:())(<0)0()`, value{
		Grammar:   "unbalanced expression: unmatched closing parenthesis )",
		Heuristic: `"(0(F)(:())(:())(<0)0()"`,
	}).Equal(t, test(`(0(F)(:())(:())(<0)0()`))

	// The space-looking character below is U+00A0.
	autogold.Want(`00 (000)`, value{Grammar: `(concat "00" "000")`, Heuristic: `(concat "00" "(000)")`}).Equal(t, test(`00 (000)`))

}

func TestScanDelimited(t *testing.T) {
	type value struct {
		Result string
		Count  int
		ErrMsg string
	}

	test := func(input string, delimiter rune) string {
		result, count, err := ScanDelimited([]byte(input), true, delimiter)
		var errMsg string
		if err != nil {
			errMsg = err.Error()
		}
		v, _ := json.Marshal(value{result, count, errMsg})
		return string(v)
	}

	autogold.Want(`""`, `{"Result":"","Count":2,"ErrMsg":""}`).Equal(t, test(`""`, '"'))
	autogold.Want(`"a"`, `{"Result":"a","Count":3,"ErrMsg":""}`).Equal(t, test(`"a"`, '"'))
	autogold.Want(`"\""`, `{"Result":"\"","Count":4,"ErrMsg":""}`).Equal(t, test(`"\""`, '"'))
	autogold.Want(`"\\""`, `{"Result":"\\","Count":4,"ErrMsg":""}`).Equal(t, test(`"\\""`, '"'))
	autogold.Want(`"\\\"`, `{"Result":"","Count":5,"ErrMsg":"unterminated literal: expected \""}`).Equal(t, test(`"\\\"`, '"'))
	autogold.Want(`"\\\""`, `{"Result":"\\\"","Count":6,"ErrMsg":""}`).Equal(t, test(`"\\\""`, '"'))
	autogold.Want(`"a`, `{"Result":"","Count":2,"ErrMsg":"unterminated literal: expected \""}`).Equal(t, test(`"a`, '"'))
	autogold.Want(`"\?"`, `{"Result":"","Count":3,"ErrMsg":"unrecognized escape sequence"}`).Equal(t, test(`"\?"`, '"'))
	autogold.Want(`/\//`, `{"Result":"/","Count":4,"ErrMsg":""}`).Equal(t, test(`/\//`, '/'))

	// The next invocation of test needs to panic.
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("expected panic for ScanDelimited")
		}
	}()
	_ = test(`a"`, '"')
}

func TestMergePatterns(t *testing.T) {
	test := func(input string) string {
		p := &parser{buf: []byte(input), heuristics: parensAsPatterns}
		nodes, err := p.parseLeaves(Regexp)
		got := nodes[0].(Pattern).Annotation.Range.String()
		if err != nil {
			t.Error(err)
		}
		return got
	}

	autogold.Want("foo()bar", `{"start":{"line":0,"column":0},"end":{"line":0,"column":8}}`).Equal(t, test("foo()bar"))
	autogold.Want("()bar", `{"start":{"line":0,"column":0},"end":{"line":0,"column":5}}`).Equal(t, test("()bar"))
}

func TestMatchUnaryKeyword(t *testing.T) {
	test := func(input string, pos int) string {
		p := &parser{buf: []byte(input), pos: pos}
		return fmt.Sprintf("%t", p.matchUnaryKeyword("NOT"))
	}

	autogold.Want("NOT bar", "true").Equal(t, test("NOT bar", 0))
	autogold.Want("foo NOT bar", "true").Equal(t, test("foo NOT bar", 4))
	autogold.Want("foo NOT", "false").Equal(t, test("foo NOT", 4))
	autogold.Want("fooNOT bar", "false").Equal(t, test("fooNOT bar", 3))
	autogold.Want("NOTbar", "false").Equal(t, test("NOTbar", 0))
	autogold.Want("(not bar)", "true").Equal(t, test("(not bar)", 1))
}
