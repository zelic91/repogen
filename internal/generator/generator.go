package generator

import (
	"github.com/zelic91/repogen/internal/code"
	"github.com/zelic91/repogen/internal/codegen"
	"github.com/zelic91/repogen/internal/mongo"
	"github.com/zelic91/repogen/internal/spec"
)

// GenerateRepository generates repository implementation code from repository
// interface specification.
func GenerateRepository(packageName string, structModel code.Struct,
	interfaceName string, methodSpecs []spec.MethodSpec) (string, error) {

	generator := mongo.NewGenerator(structModel, interfaceName)

	codeBuilder := codegen.NewBuilder(
		"repogen",
		packageName,
		generator.Imports(),
	)

	constructorBuilder, err := generator.GenerateConstructor()
	if err != nil {
		return "", err
	}

	codeBuilder.AddImplementer(constructorBuilder)
	codeBuilder.AddImplementer(generator.GenerateStruct())

	for _, method := range methodSpecs {
		methodBuilder, err := generator.GenerateMethod(method)
		if err != nil {
			return "", err
		}
		codeBuilder.AddImplementer(methodBuilder)
	}

	return codeBuilder.Build()
}
