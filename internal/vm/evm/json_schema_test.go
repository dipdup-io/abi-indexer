package evm

import (
	"log"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/stretchr/testify/assert"
)

func TestVirtualMachine_JSONSchema(t *testing.T) {
	tests := []struct {
		name        string
		contractABI string
		want        []byte
		wantErr     bool
	}{
		{
			name:        "error",
			contractABI: "",
			wantErr:     true,
		}, {
			name: "test #1",
			contractABI: `[
				{
					"inputs": [],
					"name": "getCount",
					"outputs": [
						{
							"internalType": "uint256",
							"name": "",
							"type": "uint256"
						}
					],
					"stateMutability": "view",
					"type": "function"
				},
				{
					"inputs": [],
					"name": "increment",
					"outputs": [],
					"stateMutability": "nonpayable",
					"type": "function"
				}
			]`,
			want: []byte(`{"getCount":{"type":"method","inputs":{"$schema":"http://json-schema.org/draft/2019-09/schema#","type":"object"},"outputs":{"$schema":"http://json-schema.org/draft/2019-09/schema#","type":"object","properties":{"uint256_00":{"type":"number","$comment":"uint256","title":"uint256_00"}},"required":["uint256_00"]}},"increment":{"type":"method","inputs":{"$schema":"http://json-schema.org/draft/2019-09/schema#","type":"object"},"outputs":{"$schema":"http://json-schema.org/draft/2019-09/schema#","type":"object"}}}`),
		}, {
			name: "test #2",
			contractABI: `[{
				"type":"event",
				"inputs": [{"name":"a","type":"uint256","indexed":true},{"name":"b","type":"bytes32","indexed":false}],
				"name":"Event"
			}]`,
			want: []byte(`{"Event":{"type":"event","inputs":{"$schema":"http://json-schema.org/draft/2019-09/schema#","type":"object","properties":{"a":{"type":"number","$comment":"uint256","title":"a"},"b":{"type":"string","$comment":"bytes32","title":"b"}},"required":["a","b"]}}}`),
		}, {
			name: "test #3",
			contractABI: `[
				{
				  "name": "f",
				  "type": "function",
				  "inputs": [
					{
					  "name": "s",
					  "type": "tuple",
					  "components": [
						{
						  "name": "a",
						  "type": "uint256"
						},
						{
						  "name": "b",
						  "type": "uint256[]"
						},
						{
						  "name": "c",
						  "type": "tuple[]",
						  "components": [
							{
							  "name": "x",
							  "type": "uint256"
							},
							{
							  "name": "y",
							  "type": "uint256"
							}
						  ]
						}
					  ]
					},
					{
					  "name": "t",
					  "type": "tuple",
					  "components": [
						{
						  "name": "x",
						  "type": "uint256"
						},
						{
						  "name": "y",
						  "type": "uint256"
						}
					  ]
					},
					{
					  "name": "a",
					  "type": "uint256"
					}
				  ],
				  "outputs": []
				}
			  ]`,
			want: []byte(`{
				"f": {
					"type": "method",
					"inputs": {
						"$schema": "http://json-schema.org/draft/2019-09/schema#",
						"type": "object",
						"properties": {
							"a": {
								"type": "number",
								"$comment": "uint256",
								"title": "a"
							},
							"t": {
								"title": "t",
								"type": "object",
								"properties": {
									"x": {
										"$comment": "uint256",
										"title": "x",
										"type": "number"
									},
									"y": {
										"$comment": "uint256",
										"title": "y",
										"type": "number"
									}
								},
								"required": [
									"x",
									"y"
								]
							},
							"s": {
								"title": "s",
								"type": "object",
								"properties": {
									"a": {
										"type": "number",
										"$comment": "uint256",
										"title": "a"
									},
									"b": {
										"type": "array",
										"title": "b",
										"items": [
											{
												"type": "number",
												"$comment": "uint256",
												"title": "b_elem"
											}
										],
										"uniqueItems": false
									},
									"c": {
										"type": "array",
										"title": "c",
										"items": [
											{
												"title": "c_elem",
												"type": "object",
												"properties": {
													"x": {
														"$comment": "uint256",
														"title": "x",
														"type": "number"
													},
													"y": {
														"$comment": "uint256",
														"title": "y",
														"type": "number"
													}
												},
												"required": [
													"x",
													"y"
												]
											}
										],
										"uniqueItems": false
									}
								},
								"required": ["a","b","c"]
							}
						},
						"required": ["s", "t", "a"]
					},
					"outputs": {
						"$schema": "http://json-schema.org/draft/2019-09/schema#",
						"type": "object"
					}
				}
			}`),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var contractABI *abi.ABI
			if tt.contractABI != "" {
				contractABI = new(abi.ABI)
				if err := json.Unmarshal([]byte(tt.contractABI), contractABI); err != nil {
					t.Errorf("json.Unmarshal error = %v, wantErr %v", err, tt.wantErr)
					return
				}
			}
			vm := &VirtualMachine{
				contractABI: contractABI,
			}
			got, err := vm.JSONSchema()
			if (err != nil) != tt.wantErr {
				t.Errorf("VirtualMachine.JSONSchema() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.contractABI != "" {
				log.Print(string(got))
				assert.JSONEq(t, string(tt.want), string(got))
			}
		})
	}
}
