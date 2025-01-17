package property

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
)

func TestEnum(t *testing.T) {
	type testcase struct {
		name        string
		oldProperty *apiextensionsv1.JSONSchemaProps
		newProperty *apiextensionsv1.JSONSchemaProps
		err         error
		handled     bool
		enum        *Enum
	}

	for _, tc := range []testcase{
		{
			name: "no diff, no error, handled",
			oldProperty: &apiextensionsv1.JSONSchemaProps{
				Enum: []apiextensionsv1.JSON{
					{
						Raw: []byte("foo"),
					},
				},
			},
			newProperty: &apiextensionsv1.JSONSchemaProps{
				Enum: []apiextensionsv1.JSON{
					{
						Raw: []byte("foo"),
					},
				},
			},
			err:     nil,
			handled: true,
			enum:    &Enum{},
		},
		{
			name: "new enum constraint, error, handled",
			oldProperty: &apiextensionsv1.JSONSchemaProps{
				Enum: []apiextensionsv1.JSON{},
			},
			newProperty: &apiextensionsv1.JSONSchemaProps{
				Enum: []apiextensionsv1.JSON{
					{
						Raw: []byte("foo"),
					},
				},
			},
			err:     errors.New("enum constraints [foo] added when there were no restrictions previously"),
			handled: true,
			enum:    &Enum{},
		},
		{
			name: "new enum constraint, addition enforcement set to None, no error, handled",
			oldProperty: &apiextensionsv1.JSONSchemaProps{
				Enum: []apiextensionsv1.JSON{},
			},
			newProperty: &apiextensionsv1.JSONSchemaProps{
				Enum: []apiextensionsv1.JSON{
					{
						Raw: []byte("foo"),
					},
				},
			},
			err:     nil,
			handled: true,
			enum: &Enum{
				AdditionEnforcement: EnumValidationAdditionEnforcementNone,
			},
		},
		{
			name: "new allowed enum value added, addition enforcement set to IfPreviouslyConstrained, no error, handled",
			oldProperty: &apiextensionsv1.JSONSchemaProps{
				Enum: []apiextensionsv1.JSON{
					{
						Raw: []byte("foo"),
					},
				},
			},
			newProperty: &apiextensionsv1.JSONSchemaProps{
				Enum: []apiextensionsv1.JSON{
					{
						Raw: []byte("foo"),
					},
					{
						Raw: []byte("bar"),
					},
				},
			},
			err:     nil,
			handled: true,
			enum: &Enum{
				AdditionEnforcement: EnumValidationAdditionEnforcementIfPreviouslyConstrained,
			},
		},
		{
			name: "new allowed enum value added, addition enforcement set to None, no error, handled",
			oldProperty: &apiextensionsv1.JSONSchemaProps{
				Enum: []apiextensionsv1.JSON{
					{
						Raw: []byte("foo"),
					},
				},
			},
			newProperty: &apiextensionsv1.JSONSchemaProps{
				Enum: []apiextensionsv1.JSON{
					{
						Raw: []byte("foo"),
					},
					{
						Raw: []byte("bar"),
					},
				},
			},
			err:     nil,
			handled: true,
			enum: &Enum{
				AdditionEnforcement: EnumValidationAdditionEnforcementNone,
			},
		},
		{
			name: "new allowed enum value added, addition enforcement set to Strict, error, handled",
			oldProperty: &apiextensionsv1.JSONSchemaProps{
				Enum: []apiextensionsv1.JSON{
					{
						Raw: []byte("foo"),
					},
				},
			},
			newProperty: &apiextensionsv1.JSONSchemaProps{
				Enum: []apiextensionsv1.JSON{
					{
						Raw: []byte("foo"),
					},
					{
						Raw: []byte("bar"),
					},
				},
			},
			err:     errors.New("enums [bar] added to the set of allowed values"),
			handled: true,
			enum: &Enum{
				AdditionEnforcement: EnumValidationAdditionEnforcementStrict,
			},
		},
		{
			name: "remove enum value, error, handled",
			oldProperty: &apiextensionsv1.JSONSchemaProps{
				Enum: []apiextensionsv1.JSON{
					{
						Raw: []byte("foo"),
					},
					{
						Raw: []byte("bar"),
					},
				},
			},
			newProperty: &apiextensionsv1.JSONSchemaProps{
				Enum: []apiextensionsv1.JSON{
					{
						Raw: []byte("bar"),
					},
				},
			},
			err:     errors.New("enums [foo] removed from the set of previously allowed values"),
			handled: true,
			enum:    &Enum{},
		},
		{
			name: "remove enum value, removal enforcement set to None, no error, handled",
			oldProperty: &apiextensionsv1.JSONSchemaProps{
				Enum: []apiextensionsv1.JSON{
					{
						Raw: []byte("foo"),
					},
					{
						Raw: []byte("bar"),
					},
				},
			},
			newProperty: &apiextensionsv1.JSONSchemaProps{
				Enum: []apiextensionsv1.JSON{
					{
						Raw: []byte("bar"),
					},
				},
			},
			err:     nil,
			handled: true,
			enum: &Enum{
				RemovalEnforcement: EnumValidationRemovalEnforcementNone,
			},
		},
		{
			name: "different field changed, no error, not handled",
			oldProperty: &apiextensionsv1.JSONSchemaProps{
				ID: "foo",
			},
			newProperty: &apiextensionsv1.JSONSchemaProps{
				ID: "bar",
			},
			err:     nil,
			handled: false,
			enum:    &Enum{},
		},
		{
			name: "different field changed with enum, no error, not handled",
			oldProperty: &apiextensionsv1.JSONSchemaProps{
				ID: "foo",
				Enum: []apiextensionsv1.JSON{
					{
						Raw: []byte("foo"),
					},
				},
			},
			newProperty: &apiextensionsv1.JSONSchemaProps{
				ID: "bar",
				Enum: []apiextensionsv1.JSON{
					{
						Raw: []byte("foo"),
					},
				},
			},
			err:     nil,
			handled: false,
			enum:    &Enum{},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			_, handled, err := tc.enum.Validate(NewDiff(tc.oldProperty, tc.newProperty))
			require.Equal(t, tc.err, err)
			require.Equal(t, tc.handled, handled)
		})
	}
}
