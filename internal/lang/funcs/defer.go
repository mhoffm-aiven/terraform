package funcs

import (
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/function"
)

var DeferFunc = function.New(&function.Spec{
	Params: []function.Parameter{
		{
			Name:             "value",
			Type:             cty.DynamicPseudoType,
			AllowUnknown:     true,
			AllowNull:        true,
			AllowMarked:      true,
			AllowDynamicType: true,
		},
	},
	Type: func(args []cty.Value) (cty.Type, error) {
		return args[0].Type(), nil
	},
	Impl: func(args []cty.Value, retType cty.Type) (ret cty.Value, err error) {
		arg := args[0]

		if arg.IsNull() {
			return cty.UnknownVal(arg.Type()), nil
		}
		return arg, nil
	},
})

func Defer(v cty.Value) (cty.Value, error) {
	return DeferFunc.Call([]cty.Value{v})
}
