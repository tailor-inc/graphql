package graphql

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/tailor-inc/graphql/language/ast"
	"github.com/tailor-inc/graphql/language/parser"
)

func TestAstAsSchemaConfig(t *testing.T) {
	parser := NewGraphqlParser(func(typeName string, fieldName string) FieldResolveFn {
		return DefaultResolveFn
	})

	TestType := NewScalar(ScalarConfig{
		Name: "Test",
	})

	config, err := parser.AstAsSchemaConfig([]ast.Node{ast.NewScalarDefinition(&ast.ScalarDefinition{
		Name: ast.NewName(&ast.Name{Value: "Test"}),
	})}, func(_type string) Type {
		return TestType
	})
	assert.NoError(t, err)

	assert.Equal(t, config.Types[0], TestType)

}

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

"""Description for TestQuery"""
input TestQuery {
	"""Description for i"""
	i: Int
	and: Query
}

type Hello {
	world: String
}

extend type Hello {
	extend: String
}

type Query {
	hello: Hello
    test: Test
	tests(page: Int!, limit: Int!, query: TestQuery): [Test!]!
}

extend type Query {
	testExtend: Test
}

type Mutation {
    create(name: String, need: Int!): Test
}

extend type Mutation {
    update(id: ID!, name: String, need: Int!): Test
}

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
	assert.Equal(t, "Description for TestQuery", schema.Type("TestQuery").Description())
	assert.Equal(t, "Description for i", schema.Type("TestQuery").(*InputObject).Fields()["i"].Description())
	// extend type
	assert.Equal(t, "String", schema.Type("Hello").(*Object).Fields()["extend"].Type.Name())
	assert.Equal(t, "Test", schema.Type("Query").(*Object).Fields()["test"].Type.Name())
	assert.Equal(t, "Test", schema.Type("Query").(*Object).Fields()["testExtend"].Type.Name())
	assert.Equal(t, "Test", schema.Type("Mutation").(*Object).Fields()["create"].Type.Name())
	assert.Equal(t, "Test", schema.Type("Mutation").(*Object).Fields()["update"].Type.Name())

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

func TestParseSDL_FederationDirectives(t *testing.T) {
	sdl := `
	directive @external on FIELD_DEFINITION
	directive @extends on OBJECT | INTERFACE
	
	type User @extends {
		id: ID! @external
		name: String
		email: String @external
	}
	
	type Product {
		id: ID!
		name: String
		user: User
	}
	`

	doc, err := parser.Parse(parser.ParseParams{
		Source: sdl,
	})
	assert.NoError(t, err)

	// Check that directives are parsed correctly
	var externalDirective, extendsDirective *ast.DirectiveDefinition
	var userType *ast.ObjectDefinition

	for _, def := range doc.Definitions {
		switch node := def.(type) {
		case *ast.DirectiveDefinition:
			if node.Name.Value == "external" {
				externalDirective = node
			}
			if node.Name.Value == "extends" {
				extendsDirective = node
			}
		case *ast.ObjectDefinition:
			if node.Name.Value == "User" {
				userType = node
			}
		}
	}

	assert.NotNil(t, externalDirective, "external directive should be parsed")
	assert.Equal(t, "external", externalDirective.Name.Value)
	assert.Equal(t, 1, len(externalDirective.Locations))
	assert.Equal(t, "FIELD_DEFINITION", externalDirective.Locations[0].Value)

	assert.NotNil(t, extendsDirective, "extends directive should be parsed")
	assert.Equal(t, "extends", extendsDirective.Name.Value)
	assert.Equal(t, 2, len(extendsDirective.Locations))

	assert.NotNil(t, userType, "User type should be parsed")
	assert.Equal(t, 1, len(userType.Directives), "User type should have extends directive")
	assert.Equal(t, "extends", userType.Directives[0].Name.Value)

	// Check fields with external directive
	var idField, emailField *ast.FieldDefinition
	for _, field := range userType.Fields {
		if field.Name.Value == "id" {
			idField = field
		}
		if field.Name.Value == "email" {
			emailField = field
		}
	}

	assert.NotNil(t, idField, "id field should exist")
	assert.Equal(t, 1, len(idField.Directives), "id field should have external directive")
	assert.Equal(t, "external", idField.Directives[0].Name.Value)

	assert.NotNil(t, emailField, "email field should exist")
	assert.Equal(t, 1, len(emailField.Directives), "email field should have external directive")
	assert.Equal(t, "external", emailField.Directives[0].Name.Value)
}
