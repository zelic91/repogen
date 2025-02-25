package main

import (
	"errors"
	"flag"
	"fmt"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path/filepath"

	"github.com/zelic91/repogen/internal/code"
	"github.com/zelic91/repogen/internal/generator"
	"github.com/zelic91/repogen/internal/spec"
)

const usageText = `repogen generates MongoDB repository implementation from repository interface

  Find more information at: https://github.com/zelic91/repogen

Supported options:`

// version indicates the version of repogen.
const version = "v0.3-next"

func main() {
	flag.Usage = printUsage

	versionPtr := flag.Bool("version", false, "print version of repogen")
	pkgDirPtr := flag.String("pkg", ".", "package directory to scan for model struct and repository interface")
	destPtr := flag.String("dest", "", "destination file")
	modelPtr := flag.String("model", "", "model struct name")
	repoPtr := flag.String("repo", "", "repository interface name")

	flag.Parse()

	if *versionPtr {
		printVersion()
		return
	}

	if *modelPtr == "" {
		printUsage()
		log.Fatal("-model flag required")
	}
	if *repoPtr == "" {
		printUsage()
		log.Fatal("-repo flag required")
	}

	code, err := generateFromRequest(*pkgDirPtr, *modelPtr, *repoPtr)
	if err != nil {
		panic(err)
	}

	dest := os.Stdout
	if *destPtr != "" {
		if err := os.MkdirAll(filepath.Dir(*destPtr), os.ModePerm); err != nil {
			panic(err)
		}
		file, err := os.Create(*destPtr)
		if err != nil {
			panic(err)
		}
		defer file.Close()
		dest = file
	}

	if _, err := dest.WriteString(code); err != nil {
		panic(err)
	}
}

func printUsage() {
	fmt.Println(usageText)
	flag.PrintDefaults()
}

func printVersion() {
	fmt.Println(version)
}

func generateFromRequest(pkgDir, structModelName, repositoryInterfaceName string) (string, error) {
	fset := token.NewFileSet()
	dir, err := parser.ParseDir(fset, pkgDir, nil, parser.ParseComments)
	if err != nil {
		return "", err
	}

	pkg, err := code.ParsePackage(dir)
	if err != nil {
		return "", err
	}

	return generateRepository(pkg, structModelName, repositoryInterfaceName)
}

var (
	errStructNotFound    = errors.New("struct not found")
	errInterfaceNotFound = errors.New("interface not found")
)

func generateRepository(pkg code.Package, structModelName, repositoryInterfaceName string) (string, error) {
	structModel, ok := pkg.Structs[structModelName]
	if !ok {
		return "", errStructNotFound
	}

	intf, ok := pkg.Interfaces[repositoryInterfaceName]
	if !ok {
		return "", errInterfaceNotFound
	}

	var methodSpecs []spec.MethodSpec
	for _, method := range intf.Methods {
		methodSpec, err := spec.ParseInterfaceMethod(pkg.Structs, structModel, method)
		if err != nil {
			return "", err
		}
		methodSpecs = append(methodSpecs, methodSpec)
	}

	return generator.GenerateRepository(pkg.Name, structModel, intf.Name, methodSpecs)
}
