package evm

import (
	"bytes"
	stdJSON "encoding/json"
	"fmt"

	js "github.com/dipdup-net/abi-indexer/internal/jsonschema"
	"github.com/ethereum/go-ethereum/accounts/abi"
	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type types struct {
	Type    string         `json:"type"`
	Inputs  *js.JSONSchema `json:"inputs,omitempty"`
	Outputs *js.JSONSchema `json:"outputs,omitempty"`
}

// JSONSchema -
func (vm *VirtualMachine) JSONSchema() ([]byte, error) {
	if vm.contractABI == nil {
		return nil, ErrNilABI
	}

	schema, err := vm.createEntrypointsSchema()
	if err != nil {
		return nil, err
	}

	data, err := json.Marshal(schema)
	if err != nil {
		return nil, err
	}
	var buf bytes.Buffer
	err = stdJSON.Compact(&buf, data)
	return buf.Bytes(), err
}

func (vm *VirtualMachine) createEntrypointsSchema() (map[string]types, error) {
	result := make(map[string]types)

	for _, event := range vm.contractABI.Events {
		var (
			typ = types{
				Type: "event",
				Inputs: &js.JSONSchema{
					Schema: js.Draft201909,
					Type:   js.ItemTypeObject,
				},
			}
		)

		inputsBody, err := getBodyByArgs(event.Inputs)
		if err != nil {
			return nil, err
		}
		if inputsBody != nil {
			typ.Inputs.ObjectItem = inputsBody
		}

		result[event.Name] = typ
	}

	for _, method := range vm.contractABI.Methods {

		var (
			typ = types{
				Type: "method",
				Inputs: &js.JSONSchema{
					Schema: js.Draft201909,
					Type:   js.ItemTypeObject,
				},
				Outputs: &js.JSONSchema{
					Schema: js.Draft201909,
					Type:   js.ItemTypeObject,
				},
			}
		)

		inputsBody, err := getBodyByArgs(method.Inputs)
		if err != nil {
			return nil, err
		}
		if inputsBody != nil {
			typ.Inputs.ObjectItem = inputsBody
		}

		outputsBody, err := getBodyByArgs(method.Outputs)
		if err != nil {
			return nil, err
		}
		if outputsBody != nil {
			typ.Outputs.ObjectItem = outputsBody
		}

		result[method.Name] = typ
	}

	return result, nil
}

func getBodyByArgs(args abi.Arguments) (*js.ObjectItem, error) {
	if len(args) == 0 {
		return nil, nil
	}
	body := &js.ObjectItem{
		Properties: make(map[string]js.JSONSchema),
		Required:   []string{},
	}
	for idx, arg := range args {
		argSchema, err := createSchemaItem(arg.Name, idx, &arg.Type)
		if err != nil {
			return nil, err
		}
		body.Properties[argSchema.Title] = argSchema
		body.Required = append(body.Required, argSchema.Title)
	}

	return body, nil
}

func createSchemaItem(name string, idx int, typ *abi.Type) (js.JSONSchema, error) {
	if name == "" {
		name = fmt.Sprintf("%s_%02d", typ.String(), idx)
	}
	switch typ.T {
	case abi.AddressTy, abi.StringTy:
		return js.JSONSchema{
			Type:    js.ItemTypeString,
			Title:   name,
			Comment: typ.String(),
		}, nil

	case abi.ArrayTy, abi.SliceTy:
		schema := js.JSONSchema{
			Type:  js.ItemTypeArray,
			Title: name,
		}

		elemName := fmt.Sprintf("%s_elem", name)
		elem, err := createSchemaItem(elemName, idx, typ.Elem)
		if err != nil {
			return elem, err
		}

		schema.ArrayItem = &js.ArrayItem{
			Items: []js.JSONSchema{
				elem,
			},
		}

		return schema, nil

	case abi.TupleTy:
		schema := js.JSONSchema{
			Type:  js.ItemTypeObject,
			Title: name,
			ObjectItem: &js.ObjectItem{
				Properties: make(map[string]js.JSONSchema),
				Required:   make([]string, 0),
			},
		}

		for compIdx, component := range typ.TupleElems {
			elem, err := createSchemaItem(typ.TupleRawNames[compIdx], compIdx, component)
			if err != nil {
				return elem, err
			}
			schema.ObjectItem.Properties[typ.TupleRawNames[compIdx]] = elem
			schema.ObjectItem.Required = append(schema.ObjectItem.Required, typ.TupleRawNames[compIdx])
		}

		return schema, nil
	case abi.BoolTy:
		return js.JSONSchema{
			Type:    js.ItemTypeBoolean,
			Title:   name,
			Comment: typ.String(),
		}, nil
	case abi.BytesTy, abi.FixedBytesTy, abi.FunctionTy:
		return js.JSONSchema{
			Type:    js.ItemTypeString,
			Title:   name,
			Comment: typ.String(),
		}, nil
	case abi.IntTy, abi.UintTy, abi.FixedPointTy:
		return js.JSONSchema{
			Type:    js.ItemTypeNumber,
			Title:   name,
			Comment: typ.String(),
		}, nil
	default:
		return js.JSONSchema{}, errors.Errorf("unknown argument type: %d", typ.T)
	}
}
