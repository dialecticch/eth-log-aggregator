// Copyright 2019 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package pkg

import "github.com/ethereum/go-ethereum/accounts/abi"

// tmplData is the data structure required to fill the binding template.
type tmplData struct {
	Package   string                   // Name of the package to place the generated file in
	Contracts map[string]*tmplContract // List of contracts to generate into this file
	Libraries map[string]string        // Map the bytecode's link pattern to the library name
	Structs   map[string]*tmplStruct   // Contract struct type definitions
}

// tmplContract contains the data needed to generate an individual contract binding.
type tmplContract struct {
	Type        string                 // Type name of the main contract binding
	InputABI    string                 // JSON ABI used as the input to generate the binding from
	InputBin    string                 // Optional EVM bytecode used to generate deploy code from
	FuncSigs    map[string]string      // Optional map: string signature -> 4-byte signature
	Constructor abi.Method             // Contract constructor for deploy parametrization
	Calls       map[string]*tmplMethod // Contract calls that only read state data
	Transacts   map[string]*tmplMethod // Contract calls that write state data
	Fallback    *tmplMethod            // Additional special fallback function
	Receive     *tmplMethod            // Additional special receive function
	Events      map[string]*tmplEvent  // Contract events accessors
	Libraries   map[string]string      // Same as tmplData, but filtered to only keep what the contract needs
	Library     bool                   // Indicator whether the contract is a library
}

// tmplMethod is a wrapper around an abi.Method that contains a few preprocessed
// and cached data fields.
type tmplMethod struct {
	Original   abi.Method // Original method as parsed by the abi package
	Normalized abi.Method // Normalized version of the parsed method (capitalized names, non-anonymous args/returns)
	Structured bool       // Whether the returns should be accumulated into a struct
}

// tmplEvent is a wrapper around an abi.Event that contains a few preprocessed
// and cached data fields.
type tmplEvent struct {
	Original   abi.Event // Original event as parsed by the abi package
	Normalized abi.Event // Normalized version of the parsed fields
}

// tmplField is a wrapper around a struct field with binding language
// struct type definition and relative filed name.
type tmplField struct {
	Type    string   // Field type representation depends on target binding language
	Name    string   // Field name converted from the raw user-defined field name
	SolKind abi.Type // Raw abi type information
}

// tmplStruct is a wrapper around an abi.tuple and contains an auto-generated
// struct name.
type tmplStruct struct {
	Name   string       // Auto-generated struct name(before solidity v0.5.11) or raw name.
	Fields []*tmplField // Struct fields definition depends on the binding language.
}

const tmplSource = `
// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package {{.Package}}

import (
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/event"
)

{{$structs := .Structs}}
{{range $contract := .Contracts}}

	// Watcher watches for harvested events and pipes them to an output channel.
	type {{.Type}}Watcher struct {	
		contracts []{{.Type}}
	
		{{range .Events}}
			subscriptionsFor{{.Normalized.Name}} []event.Subscription
			sinkFor{{.Normalized.Name}} <-chan *{{$contract.Type}}{{.Normalized.Name}}
		{{end}}
	}
	
	// NewWatcher returns a watcher.
	func New{{.Type}}Watcher(contracts []{{.Type}}) *{{.Type}}Watcher {
		return &{{.Type}}Watcher{
			contracts: contracts,
		{{range .Events}}
			subscriptionsFor{{.Normalized.Name}}: make([]event.Subscription, 0),
		{{end}}
		}
	}

	{{range .Events}}
		// Watch is a Harvested log subscription binding to a set of contracts.
		func (w *{{$contract.Type}}Watcher) Watch{{.Normalized.Name}}(opts *bind.WatchOpts, sink chan<- *{{$contract.Type}}{{.Normalized.Name}}{{range .Normalized.Inputs}}{{if .Indexed}}, {{.Name}} []{{bindtype .Type $structs}}{{end}}{{end}}) (event.Subscription, error) {
			for _, c := range w.contracts {
				sub, err := c.Watch{{.Normalized.Name}}(
					opts,
					sink,
					{{range .Normalized.Inputs}}{{if .Indexed}}
					{{.Name}},
					{{end}}{{end}}
				)
	
				if err != nil {
					w.unsubscribe{{.Normalized.Name}}()
					return nil, err
				}
	
				w.subscriptionsFor{{.Normalized.Name}} = append(w.subscriptionsFor{{.Normalized.Name}}, sub)
			}
		
			return event.NewSubscription(func(quit <-chan struct{}) error {
				for {
					select {
					case evt := <-w.sinkFor{{.Normalized.Name}}:
						sink <- evt
					case <-quit:
						w.unsubscribe{{.Normalized.Name}}()
						return nil
					default:
						for _, sub := range w.subscriptionsFor{{.Normalized.Name}} {
							select {
							case err := <-sub.Err():
								return err
							default:
								continue
							}
						}
					}
				}
			}), nil
		}
		
		func (w *{{$contract.Type}}Watcher) unsubscribe{{.Normalized.Name}}() {
			for _, sub := range w.subscriptionsFor{{.Normalized.Name}} {
				sub.Unsubscribe()
			}
		}
	{{end}}

{{end}}
`
