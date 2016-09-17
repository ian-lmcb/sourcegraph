package ctags

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	opentracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"

	"sourcegraph.com/sourcegraph/sourcegraph/lang/ctags/parser"
	"sourcegraph.com/sourcegraph/sourcegraph/pkg/jsonrpc2"
	"sourcegraph.com/sourcegraph/sourcegraph/pkg/lsp"
)

var nameToSymbolKind = map[string]lsp.SymbolKind{
	"file":        lsp.SKFile,
	"module":      lsp.SKModule,
	"namespace":   lsp.SKNamespace,
	"package":     lsp.SKPackage,
	"class":       lsp.SKClass,
	"method":      lsp.SKMethod,
	"property":    lsp.SKProperty,
	"field":       lsp.SKField,
	"constructor": lsp.SKConstructor,
	"enum":        lsp.SKEnum,
	"interface":   lsp.SKInterface,
	"function":    lsp.SKFunction,
	"variable":    lsp.SKVariable,
	"constant":    lsp.SKConstant,
	"string":      lsp.SKString,
	"number":      lsp.SKNumber,
	"boolean":     lsp.SKBoolean,
	"array":       lsp.SKArray,
}

func (h *Handler) handleSymbol(ctx context.Context, req *jsonrpc2.Request, params lsp.WorkspaceSymbolParams) (symbols []lsp.SymbolInformation, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ctags.handleSymbol")
	if params.Query != "" {
		span.SetTag("query", params.Query)
	}
	defer func() {
		if err != nil {
			ext.Error.Set(span, true)
			span.LogEvent(fmt.Sprintf("error: %v", err))
		}
		span.SetTag("returned symbols count", len(symbols))
		span.Finish()
	}()

	rootDir := h.init.RootPath
	vslog("Requesting workspace symbols for ", rootDir)
	p, err := parser.Parse(ctx, rootDir, nil)
	if err != nil {
		return nil, err
	}

	tags := p.Tags()
	span.SetTag("tags count", len(tags))
	vslog("Definitions found: ", strconv.Itoa(len(tags)))

	if params.Query != "" {
		span, _ := opentracing.StartSpanFromContext(ctx, "filter tags")

		q := strings.ToLower(params.Query)
		exact, prefix, contains := []parser.Tag{}, []parser.Tag{}, []parser.Tag{}
		for _, t := range tags {
			name := strings.ToLower(t.Name)
			if name == q {
				exact = append(exact, t)
			} else if strings.HasPrefix(name, q) {
				prefix = append(prefix, t)
			} else if strings.Contains(name, q) {
				contains = append(contains, t)
			}
		}
		tags = append(append(exact, prefix...), contains...) // Basic ranking

		span.Finish()
	}

	symbols = make([]lsp.SymbolInformation, 0, len(tags))
	for _, tag := range tags {
		nameIdx := strings.Index(tag.DefLinePrefix, tag.Name)
		if nameIdx < 0 {
			// Drop this tag if we couldn't find the name in the def line prefix.
			continue
		}
		kind := nameToSymbolKind[tag.Kind]
		if kind == 0 {
			kind = lsp.SKVariable
		}
		symbols = append(symbols, lsp.SymbolInformation{
			Name: tag.Name,
			Kind: kind,
			Location: lsp.Location{
				URI: "file://" + rootDir + "/" + tag.File,
				Range: lsp.Range{
					Start: lsp.Position{Line: tag.Line - 1, Character: nameIdx},
					End:   lsp.Position{Line: tag.Line - 1, Character: nameIdx + len(tag.Name)},
				},
			},
		})
	}

	// Limit the amount of symbols we serve to the client. Allowing an
	// excessively large amount to be returned will generate a huge response
	// object, which slows down the performance of the pipeline significantly.
	const limit = 100
	if len(symbols) > limit {
		symbols = symbols[:limit]
	}

	return symbols, nil
}
