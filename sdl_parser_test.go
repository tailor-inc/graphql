package graphql

import (
	"github.com/tailor-inc/graphql/language/parser"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseSDL(t *testing.T) {

	sdl := `
enum TestEnum {
	enum1
	enum2
	enum3
}

type Foo {
	foo: String
}
type Var {
	var: String
}

union TestUnion = Foo | Var

type Test {
	id: ID
	i: Int
	f: Float
	b: Boolean
	s: String
	u: TestUnion
	arr: [String]
	ni: Int!
	parent: Test
	enum: TestEnum
}

input TestQuery {
	i: Int
	and: Query
}

type Hello {
	world: String
}

type Query {
	hello: Hello
    test: Test
	tests(page: Int!, limit: Int!, query: TestQuery): [Test!]!
}

type Mutation {
    create(name: String, need: Int!): Test
}
`
	schema, err := ParseSDL(sdl, func(name, field string) FieldResolveFn {
		switch name {
		case "Query":
			switch field {
			case "hello":
				return func(p ResolveParams) (interface{}, error) {
					return map[string]interface{}{
						"world": "hello",
					}, nil
				}
			}
		case "Hello":
			switch field {
			case "world":
				return func(p ResolveParams) (interface{}, error) {
					return "world", nil
				}
			}
		}
		return nil
	})

	assert.NoError(t, err)

	assert.NotNil(t, t, schema.QueryType())
	assert.Len(t, schema.TypeMap(), 22)

	query := `
query Example {
	hello {
        world
    }
}
`
	astDoc, err := parser.Parse(parser.ParseParams{
		Source: query,
		Options: parser.ParseOptions{
			// include source, for error reporting
			NoSource: false,
		},
	})
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	args := map[string]interface{}{
		"size": 100,
	}
	data := map[string]interface{}{}
	operationName := "Example"
	ep := ExecuteParams{
		Schema:        *schema,
		Root:          data,
		AST:           astDoc,
		OperationName: operationName,
		Args:          args,
	}
	result := Execute(ep)
	if len(result.Errors) > 0 {
		t.Fatalf("wrong result, unexpected errors: %v", result.Errors)
	}
	t.Log(result)

}
