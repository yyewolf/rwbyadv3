package main

import (
	_ "embed"
	"fmt"
	"log"
	"reflect"
	"strings"

	"entgo.io/ent/entc"
	"entgo.io/ent/entc/gen"
	"github.com/yyewolf/entvis"
)

//go:embed edgetags.tmpl
var edgeTags string

func main() {
	err := entc.Generate("./ent/schema", &gen.Config{
		Features: []gen.Feature{
			gen.FeatureUpsert,
			gen.FeatureSnapshot,
			gen.FeatureIntercept,
		},
		Hooks: []gen.Hook{
			AddOmitzero(),
		},
		Templates: []*gen.Template{
			gen.MustParse(gen.NewTemplate("edgeTags").Parse(edgeTags)),
			gen.MustParse(gen.NewTemplate("utils").ParseDir("./ent/schema/templates")),
		},
	},
		entc.Extensions(entvis.NewViewExtension()),
	)
	if err != nil {
		log.Fatalf("running ent codegen: %v", err)
	}
}

func doStructTagReplace(tag reflect.StructTag) string {
	if tag.Get("json") == "" || tag.Get("json") == "-" {
		return string(tag)
	}

	oldTag := fmt.Sprintf(`json:"%s"`, tag.Get("json"))
	newTag := fmt.Sprintf(`json:"%s,omitzero"`, tag.Get("json"))

	// Replace the old tag with the new tag
	tag = reflect.StructTag(strings.ReplaceAll(string(tag), oldTag, newTag))
	return string(tag)
}

func AddOmitzero() gen.Hook {
	return func(next gen.Generator) gen.Generator {
		return gen.GenerateFunc(func(g *gen.Graph) error {
			for _, node := range g.Nodes {
				fields := node.Fields
				if node.ID != nil {
					fields = append(fields, node.ID)
				}

				for _, field := range fields {
					field.StructTag = doStructTagReplace(reflect.StructTag(field.StructTag))
				}

				for _, edge := range node.Edges {
					edge.StructTag = doStructTagReplace(reflect.StructTag(edge.StructTag))
				}
			}
			return next.Generate(g)
		})
	}
}
