package mongo

import (
	"strconv"

	"github.com/zelic91/repogen/internal/code"
	"github.com/zelic91/repogen/internal/codegen"
	"github.com/zelic91/repogen/internal/spec"
)

func (g RepositoryGenerator) generateFindBody(
	operation spec.FindOperation) (codegen.FunctionBody, error) {

	return findBodyGenerator{
		baseMethodGenerator: g.baseMethodGenerator,
		operation:           operation,
	}.generate()
}

type findBodyGenerator struct {
	baseMethodGenerator
	operation spec.FindOperation
}

func (g findBodyGenerator) generate() (codegen.FunctionBody, error) {
	querySpec, err := g.convertQuerySpec(g.operation.Query)
	if err != nil {
		return nil, err
	}

	sortsCode, err := g.generateSortMap()
	if err != nil {
		return nil, err
	}

	if g.operation.Mode == spec.QueryModeOne {
		return g.generateFindOneBody(querySpec, sortsCode), nil
	}

	return g.generateFindManyBody(querySpec, sortsCode), nil
}

func (g findBodyGenerator) generateFindOneBody(querySpec querySpec,
	sortsCode codegen.MapStatement) codegen.FunctionBody {

	return codegen.FunctionBody{
		codegen.DeclStatement{
			Name: "entity",
			Type: code.SimpleType(g.structModel.Name),
		},
		codegen.IfBlock{
			Condition: []codegen.Statement{
				codegen.DeclAssignStatement{
					Vars: []string{"err"},
					Values: codegen.StatementList{
						codegen.NewChainBuilder("r").
							Chain("collection").
							Call("FindOne",
								codegen.Identifier("arg0"),
								querySpec.Code(),
								codegen.NewChainBuilder("options").
									Call("FindOne").
									Call("SetSort", sortsCode).
									Build(),
							).
							Call("Decode",
								codegen.RawStatement("&entity"),
							).Build(),
					},
				},
				errOccurred,
			},
			Statements: []codegen.Statement{
				returnNilErr,
			},
		},
		codegen.ReturnStatement{
			codegen.RawStatement("&entity"),
			codegen.Identifier("nil"),
		},
	}
}

func (g findBodyGenerator) generateFindManyBody(querySpec querySpec,
	sortsCode codegen.MapStatement) codegen.FunctionBody {

	return codegen.FunctionBody{
		codegen.DeclAssignStatement{
			Vars: []string{"cursor", "err"},
			Values: codegen.StatementList{
				codegen.NewChainBuilder("r").
					Chain("collection").
					Call("Find",
						codegen.Identifier("arg0"),
						querySpec.Code(),
						g.findManyOptions(sortsCode),
					).Build(),
			},
		},
		ifErrReturnNilErr,
		codegen.DeclStatement{
			Name: "entities",
			Type: code.ArrayType{
				ContainedType: code.PointerType{
					ContainedType: code.SimpleType(g.structModel.Name),
				},
			},
		},
		codegen.IfBlock{
			Condition: []codegen.Statement{
				codegen.DeclAssignStatement{
					Vars: []string{"err"},
					Values: codegen.StatementList{
						codegen.NewChainBuilder("cursor").
							Call("All",
								codegen.Identifier("arg0"),
								codegen.RawStatement("&entities"),
							).Build(),
					},
				},
				errOccurred,
			},
			Statements: []codegen.Statement{
				returnNilErr,
			},
		},
		codegen.ReturnStatement{
			codegen.Identifier("entities"),
			codegen.Identifier("nil"),
		},
	}
}

func (g findBodyGenerator) findManyOptions(
	sortsCode codegen.MapStatement) codegen.Statement {

	optionsBuilder := codegen.NewChainBuilder("options").
		Call("Find").
		Call("SetSort", sortsCode)
	if g.operation.Limit > 0 {
		optionsBuilder = optionsBuilder.Call("SetLimit",
			codegen.Identifier(strconv.Itoa(g.operation.Limit)),
		)
	}

	return optionsBuilder.Build()
}

func (g findBodyGenerator) generateSortMap() (
	codegen.MapStatement, error) {

	sortsCode := codegen.MapStatement{
		Type: "bson.M",
	}

	for _, s := range g.operation.Sorts {
		bsonFieldReference, err := g.bsonFieldReference(s.FieldReference)
		if err != nil {
			return codegen.MapStatement{}, err
		}

		sortValueIdentifier := codegen.Identifier("1")
		if s.Ordering == spec.OrderingDescending {
			sortValueIdentifier = codegen.Identifier("-1")
		}

		sortsCode.Pairs = append(sortsCode.Pairs, codegen.MapPair{
			Key:   bsonFieldReference,
			Value: sortValueIdentifier,
		})
	}

	return sortsCode, nil
}
