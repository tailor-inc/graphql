package graphql

const (
	// Operations
	DirectiveLocationQuery              = "QUERY"
	DirectiveLocationMutation           = "MUTATION"
	DirectiveLocationSubscription       = "SUBSCRIPTION"
	DirectiveLocationField              = "FIELD"
	DirectiveLocationFragmentDefinition = "FRAGMENT_DEFINITION"
	DirectiveLocationFragmentSpread     = "FRAGMENT_SPREAD"
	DirectiveLocationInlineFragment     = "INLINE_FRAGMENT"

	// Schema Definitions
	DirectiveLocationSchema               = "SCHEMA"
	DirectiveLocationScalar               = "SCALAR"
	DirectiveLocationObject               = "OBJECT"
	DirectiveLocationFieldDefinition      = "FIELD_DEFINITION"
	DirectiveLocationArgumentDefinition   = "ARGUMENT_DEFINITION"
	DirectiveLocationInterface            = "INTERFACE"
	DirectiveLocationUnion                = "UNION"
	DirectiveLocationEnum                 = "ENUM"
	DirectiveLocationEnumValue            = "ENUM_VALUE"
	DirectiveLocationInputObject          = "INPUT_OBJECT"
	DirectiveLocationInputFieldDefinition = "INPUT_FIELD_DEFINITION"
)

// DefaultDeprecationReason Constant string used for default reason for a deprecation.
const DefaultDeprecationReason = "No longer supported"

// SpecifiedRules The full list of specified directives.
var SpecifiedDirectives = []*Directive{
	IncludeDirective,
	SkipDirective,
	DeprecatedDirective,
}

// Directive structs are used by the GraphQL runtime as a way of modifying execution
// behavior. Type system creators will usually not create these directly.
type Directive struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Locations   []string    `json:"locations"`
	Args        []*Argument `json:"args"`

	err error
}

// DirectiveConfig options for creating a new GraphQLDirective
type DirectiveConfig struct {
	Name        string              `json:"name"`
	Description string              `json:"description"`
	Locations   []string            `json:"locations"`
	Args        FieldConfigArgument `json:"args"`
}

func NewDirective(config DirectiveConfig) *Directive {
	dir := &Directive{}

	// Ensure directive is named
	if dir.err = invariant(config.Name != "", "Directive must be named."); dir.err != nil {
		return dir
	}

	// Ensure directive name is valid
	if dir.err = assertValidName(config.Name); dir.err != nil {
		return dir
	}

	// Ensure locations are provided for directive
	if dir.err = invariant(len(config.Locations) > 0, "Must provide locations for directive."); dir.err != nil {
		return dir
	}

	args := []*Argument{}

	for argName, argConfig := range config.Args {
		if dir.err = assertValidName(argName); dir.err != nil {
			return dir
		}
		args = append(args, &Argument{
			PrivateName:        argName,
			PrivateDescription: argConfig.Description,
			Type:               argConfig.Type,
			DefaultValue:       argConfig.DefaultValue,
		})
	}

	dir.Name = config.Name
	dir.Description = config.Description
	dir.Locations = config.Locations
	dir.Args = args
	return dir
}

// IncludeDirective is used to conditionally include fields or fragments.
var IncludeDirective = NewDirective(DirectiveConfig{
	Name: "include",
	Description: "Directs the executor to include this field or fragment only when " +
		"the `if` argument is true.",
	Locations: []string{
		DirectiveLocationField,
		DirectiveLocationFragmentSpread,
		DirectiveLocationInlineFragment,
	},
	Args: FieldConfigArgument{
		"if": &ArgumentConfig{
			Type:        NewNonNull(Boolean),
			Description: "Included when true.",
		},
	},
})

// SkipDirective Used to conditionally skip (exclude) fields or fragments.
var SkipDirective = NewDirective(DirectiveConfig{
	Name: "skip",
	Description: "Directs the executor to skip this field or fragment when the `if` " +
		"argument is true.",
	Args: FieldConfigArgument{
		"if": &ArgumentConfig{
			Type:        NewNonNull(Boolean),
			Description: "Skipped when true.",
		},
	},
	Locations: []string{
		DirectiveLocationField,
		DirectiveLocationFragmentSpread,
		DirectiveLocationInlineFragment,
	},
})

// DeprecatedDirective  Used to declare element of a GraphQL schema as deprecated.
var DeprecatedDirective = NewDirective(DirectiveConfig{
	Name:        "deprecated",
	Description: "Marks an element of a GraphQL schema as no longer supported.",
	Args: FieldConfigArgument{
		"reason": &ArgumentConfig{
			Type: String,
			Description: "Explains why this element was deprecated, usually also including a " +
				"suggestion for how to access supported similar data. Formatted" +
				"in [Markdown](https://daringfireball.net/projects/markdown/).",
			DefaultValue: DefaultDeprecationReason,
		},
	},
	Locations: []string{
		DirectiveLocationFieldDefinition,
		DirectiveLocationEnumValue,
	},
})

// ExternalDirective
// directive @external on FIELD_DEFINITION
var ExternalDirective = NewDirective(DirectiveConfig{
	Name: "external",
	Locations: []string{
		DirectiveLocationFieldDefinition,
	},
})

// RequiresDirective The @requires directive is used to annotate the required input fieldset from a base type for a resolver. It is used to develop a query plan where the required fields may not be needed by the client, but the service may need additional information from other services
var RequiresDirective = NewDirective(DirectiveConfig{
	Name:        "requires",
	Description: "The @requires directive is used to annotate the required input fieldset from a base type for a resolver. It is used to develop a query plan where the required fields may not be needed by the client, but the service may need additional information from other services",
	Locations: []string{
		DirectiveLocationFieldDefinition,
	},
	Args: FieldConfigArgument{
		"fields": &ArgumentConfig{
			Type: NewNonNull(FieldSet),
		},
	},
})

// ProvidesDirective The @provides directive is used to annotate the expected returned fieldset from a field on a base type that is guaranteed to be selectable by the gateway
var ProvidesDirective = NewDirective(DirectiveConfig{
	Name:        "provides",
	Description: "The @provides directive is used to annotate the expected returned fieldset from a field on a base type that is guaranteed to be selectable by the gateway",
	Locations: []string{
		DirectiveLocationFieldDefinition,
	},
	Args: FieldConfigArgument{
		"fields": &ArgumentConfig{
			Type: NewNonNull(FieldSet),
		},
	},
})

// KeyDirective The @key directive is used to indicate a combination of fields that can be used to uniquely identify and fetch an object or interface.
var KeyDirective = NewDirective(DirectiveConfig{
	Name:        "key",
	Description: "The @key directive is used to indicate a combination of fields that can be used to uniquely identify and fetch an object or interface.",
	Locations: []string{
		DirectiveLocationObject,
		DirectiveLocationInterface,
	},
	Args: FieldConfigArgument{
		"fields": &ArgumentConfig{
			Type: String,
		},
	},
})

// Link__Purpose /
var Link__Purpose = NewEnum(EnumConfig{
	Name: "link__Purpose",
	Values: EnumValueConfigMap{
		"SECURITY": &EnumValueConfig{
			Value:       "SECURITY",
			Description: "`SECURITY` features provide metadata necessary to securely resolve fields.",
		},
		"EXECUTION": &EnumValueConfig{
			Value:       "EXECUTION",
			Description: "`EXECUTION` features provide metadata necessary for operation execution.",
		},
	},
})

// LinkDirective The @link directive links definitions within the document to external schemas.
var LinkDirective = NewDirective(DirectiveConfig{
	Name:        "link",
	Description: "The @link directive links definitions within the document to external schemas.",
	Locations: []string{
		DirectiveLocationObject,
		DirectiveLocationInterface,
	},
	Args: FieldConfigArgument{
		"url": &ArgumentConfig{
			Type: String,
		},
		"as": &ArgumentConfig{
			Type: String,
		},
		"for": &ArgumentConfig{
			Type: Link__Purpose,
		},
		"import": &ArgumentConfig{
			Type: NewList(link__Import),
		},
	},
})

// ShareableDirective The @shareable directive is used to indicate that a field can be resolved by multiple subgraphs. Any subgraph that includes a shareable field can potentially resolve a query for that field. To successfully compose, a field must have the same shareability mode (either shareable or non-shareable) across all subgraphs.
var ShareableDirective = NewDirective(DirectiveConfig{
	Name:        "shareable",
	Description: "The @shareable directive is used to indicate that a field can be resolved by multiple subgraphs. Any subgraph that includes a shareable field can potentially resolve a query for that field. To successfully compose, a field must have the same shareability mode (either shareable or non-shareable) across all subgraphs.",
	Locations: []string{
		DirectiveLocationObject,
		DirectiveLocationFieldDefinition,
	},
})

// OverrideDirective The @override directive is used to indicate that the current subgraph is taking responsibility for resolving the marked field away from the subgraph specified in the from argument.
var OverrideDirective = NewDirective(DirectiveConfig{
	Name:        "override",
	Description: "The @override directive is used to indicate that the current subgraph is taking responsibility for resolving the marked field away from the subgraph specified in the from argument.",
	Locations: []string{
		DirectiveLocationFieldDefinition,
	},
})

// InaccessibleDirective The @composeDirective directive is used to indicate to composition that the custom directive specified should be preserved in the supergraph
var InaccessibleDirective = NewDirective(DirectiveConfig{
	Name:        "inaccessible",
	Description: "The @composeDirective directive is used to indicate to composition that the custom directive specified should be preserved in the supergraph",
	Locations: []string{
		DirectiveLocationFieldDefinition,
		DirectiveLocationInterface,
		DirectiveLocationObject,
		DirectiveLocationUnion,
		DirectiveLocationArgumentDefinition,
		DirectiveLocationScalar,
		DirectiveLocationEnum,
		DirectiveLocationEnumValue,
		DirectiveLocationInputObject,
		DirectiveLocationInputFieldDefinition,
	},
})

// ComposeDirective The @composeDirective directive is used to indicate to composition that the custom directive specified should be preserved in the supergraph
var ComposeDirective = NewDirective(DirectiveConfig{
	Name:        "composeDirective",
	Description: "The @composeDirective directive is used to indicate to composition that the custom directive specified should be preserved in the supergraph",
	Locations: []string{
		DirectiveLocationSchema,
	},
})
