package evm

import (
	"github.com/dipdup-net/abi-indexer/internal/storage"
	"github.com/ethereum/go-ethereum/accounts/abi"
)

// VirtualMachine -
type VirtualMachine struct {
	contractABI *abi.ABI
}

// NewVM -
func NewVM(data []byte) (*VirtualMachine, error) {
	var contractABI abi.ABI
	if err := json.Unmarshal(data, &contractABI); err != nil {
		return nil, err
	}

	return &VirtualMachine{
		contractABI: &contractABI,
	}, nil
}

// Methods -
func (vm *VirtualMachine) Methods() ([]storage.Method, error) {
	if vm.contractABI == nil {
		return nil, ErrNilABI
	}

	methods := make([]storage.Method, 0)
	for name, method := range vm.contractABI.Methods {
		methods = append(methods, storage.Method{
			Name:        name,
			Signature:   method.Sig,
			SignatureID: method.ID,
			IsConst:     method.Constant,
			IsPayable:   method.Payable,
			Type:        int(method.Type),
			Mutability:  method.StateMutability,
		})
	}

	return methods, nil
}

// Events -
func (vm *VirtualMachine) Events() ([]storage.Event, error) {
	if vm.contractABI == nil {
		return nil, ErrNilABI
	}

	events := make([]storage.Event, 0)
	for name, event := range vm.contractABI.Events {
		events = append(events, storage.Event{
			Name:        name,
			Signature:   event.Sig,
			SignatureID: event.ID.Bytes(),
			Anonymous:   event.Anonymous,
		})
	}
	return events, nil
}
