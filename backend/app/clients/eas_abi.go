// Package clients contains the minimal ABI definitions for EAS contracts.
//
// Only the functions used by Miora are included:
//   - EAS.attest() — publish attestation
//   - EAS.getAttestation() — query attestation by UID
//   - EAS.isAttestationValid() — check if attestation exists
package clients

// EASABIJSON is the minimal ABI for the EAS contract (0x4200000000000000000000000000000000000021).
const EASABIJSON = `[
	{
		"name": "attest",
		"type": "function",
		"stateMutability": "payable",
		"inputs": [
			{
				"name": "request",
				"type": "tuple",
				"components": [
					{ "name": "schema", "type": "bytes32" },
					{
						"name": "data",
						"type": "tuple",
						"components": [
							{ "name": "recipient", "type": "address" },
							{ "name": "expirationTime", "type": "uint64" },
							{ "name": "revocable", "type": "bool" },
							{ "name": "refUID", "type": "bytes32" },
							{ "name": "data", "type": "bytes" },
							{ "name": "value", "type": "uint256" }
						]
					}
				]
			}
		],
		"outputs": [{ "name": "", "type": "bytes32" }]
	},
	{
		"name": "getAttestation",
		"type": "function",
		"stateMutability": "view",
		"inputs": [{ "name": "uid", "type": "bytes32" }],
		"outputs": [
			{
				"name": "",
				"type": "tuple",
				"components": [
					{ "name": "uid", "type": "bytes32" },
					{ "name": "schema", "type": "bytes32" },
					{ "name": "time", "type": "uint64" },
					{ "name": "expirationTime", "type": "uint64" },
					{ "name": "revocationTime", "type": "uint64" },
					{ "name": "refUID", "type": "bytes32" },
					{ "name": "recipient", "type": "address" },
					{ "name": "attester", "type": "address" },
					{ "name": "revocable", "type": "bool" },
					{ "name": "data", "type": "bytes" }
				]
			}
		]
	},
	{
		"name": "isAttestationValid",
		"type": "function",
		"stateMutability": "view",
		"inputs": [{ "name": "uid", "type": "bytes32" }],
		"outputs": [{ "name": "", "type": "bool" }]
	}
]`

// SchemaRegistryABIJSON is the minimal ABI for the SchemaRegistry contract (0x4200000000000000000000000000000000000020).
const SchemaRegistryABIJSON = `[
	{
		"name": "register",
		"type": "function",
		"stateMutability": "nonpayable",
		"inputs": [
			{ "name": "schema", "type": "string" },
			{ "name": "resolver", "type": "address" },
			{ "name": "revocable", "type": "bool" }
		],
		"outputs": [{ "name": "", "type": "bytes32" }]
	}
]`
