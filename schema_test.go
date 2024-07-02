package graphql

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_typeMapReducer(t *testing.T) {
	type args struct {
		schema     *Schema
		typeMap    TypeMap
		objectType Type
	}
	tests := []struct {
		name    string
		args    args
		want    TypeMap
		wantErr bool
	}{
		{
			name: "failure - interface isn't nil but the concrete value is nil",
			args: args{
				schema:     &Schema{},
				typeMap:    TypeMap{},
				objectType: nil,
			},
			want:    TypeMap{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := typeMapReducer(tt.args.schema, tt.args.typeMap, tt.args.objectType)
			if tt.wantErr {
				assert.Error(t, err)
			}
			assert.Equalf(t, tt.want, got, "typeMapReducer(%v, %v, %v)", tt.args.schema, tt.args.typeMap, tt.args.objectType)
		})
	}
}
