// Package chaner generates channed type.
package chaner

import (
	"fmt"
	"go/ast"
	"go/parser"
	"log"
	"strings"

	"github.com/mh-cbon/astutil"
	"github.com/mh-cbon/ggt/opts"
	"github.com/mh-cbon/ggt/utils"
)

// Cmd chaner
type Cmd struct{}

// Run chaner command
func (c Cmd) Run(options *opts.Cli) {

	outPkg := options.OutPkg
	help := options.Help

	if help {
		showHelp()
		return
	}

	args := options.Args
	if len(args) < 1 {
		wrongInput("not enough type to trasnform")
		return
	}

	out := ""
	if args[0] == "-" {
		args = args[1:]
		out = "-"
	}

	todos, err := utils.NewTransformsArgs("").Parse(args)
	if err != nil {
		panic(err)
	}

	filesOut := utils.NewFilesOut("github.com/mh-cbon/ggt")

	for _, todo := range todos.Args {

		fileOut := filesOut.Get(todo.ToPath)

		fileOut.PkgName = outPkg
		if fileOut.PkgName == "" {
			fileOut.PkgName = utils.FindOutPkg(todo)
		}

		if err := processType(todo, fileOut); err != nil {
			log.Println(err)
		}
	}

	filesOut.Write(out)
}

func wrongInput(format string, a ...interface{}) {
	showHelp()
	fmt.Printf(`
    wrong input: %v
    `, fmt.Sprintf(format, a...))
}
func showHelp() {
	fmt.Printf(`ggt [options] chaner ...[FromTypeName:ToTypeName]

generates race protected type

[options]

    -p        Force out package name

...[FromTypeName:ToTypeName]

    A list of types such as src:dst.
    A type is defined by its package path and its type name,
    [pkgpath/]name.
    If the Package path is empty, it is set to the package name being generated.
    Name can be a valid type identifier such as TypeName, *TypeName, []TypeName

Example

    ggt -c chaner MySrcType:gen/*NewGenType
    ggt -c chaner myModule/*MySrcType:gen/NewGenType
`)
}

func processType(todo utils.TransformArg, fileOut *utils.FileOut) error {
	dest := &fileOut.Body
	srcName := todo.FromTypeName
	destName := todo.ToTypeName
	srcConcrete := astutil.GetUnpointedType(srcName)
	dstConcrete := astutil.GetUnpointedType(destName)

	prog := astutil.GetProgramFast(todo.FromPkgPath)
	pkg := prog.Package(todo.FromPkgPath)

	foundTypes := astutil.FindTypes(pkg)
	foundMethods := astutil.FindMethods(pkg)
	foundCtors := astutil.FindCtors(pkg, foundTypes)

	// srcIsPointer := astutil.IsAPointedType(srcName)
	// srcNameFq := srcName
	if todo.FromPkgPath != todo.ToPkgPath && !astutil.IsBasic(todo.FromTypeName) {
		// srcNameFq = fmt.Sprintf("%v.%v", filepath.Base(todo.FromPkgPath), srcConcrete)
		// if srcIsPointer {
		// 	srcNameFq = "*" + srcNameFq
		// }
		// fileOut.AddImport(todo.FromPkgPath, todo.FromPkgID)
		//todo: it does not handle cyclic import.
	}

	fmt.Fprintf(dest, `
// %v is channeled.
type %v struct{
	embed %v
	ops chan func()
	stop chan bool
	tick chan bool
}
		`, dstConcrete, dstConcrete, srcName)

	ctorParams := ""
	ctorParamsInvokation := ""
	ctorName := ""
	ctorIsPointer := false
	if x, ok := foundCtors[srcConcrete]; ok {
		withEllipse := astutil.MethodHasEllipse(x)
		ctorParamsInvokation = astutil.MethodParamNamesInvokation(x, withEllipse)
		ctorParams = astutil.MethodParams(x)
		ctorIsPointer = astutil.MethodReturnPointer(x)
		ctorName = "New" + srcConcrete
	}

	if !(astutil.IsAPointedType(srcName) == ctorIsPointer) {
		ctorParams = ""
	}

	fmt.Fprintf(dest, `// New%v constructs a channeled version of %v
func New%v(%v) *%v {
	ret := &%v{
		ops: make(chan func()),
		tick: make(chan bool),
		stop: make(chan bool),
	}
`,
		dstConcrete, srcName, dstConcrete, ctorParams, dstConcrete, dstConcrete)

	if ctorName != "" && astutil.IsAPointedType(srcName) == ctorIsPointer {
		fmt.Fprintf(dest, "	ret.embed = %v(%v)\n", ctorName, ctorParamsInvokation)
	}
	fmt.Fprintf(dest, "	go ret.Start()\n")
	fmt.Fprintf(dest, "	return ret\n")
	fmt.Fprintf(dest, "}\n")

	receiverName := "t"

	for _, m := range foundMethods[srcConcrete] {
		withEllipse := astutil.MethodHasEllipse(m)
		paramNames := astutil.MethodParamNamesInvokation(m, withEllipse)
		receiverName = astutil.ReceiverName(m)
		methodName := astutil.MethodName(m)

		if methodName == "Transact" {
			continue
		}

		varExpr := ""
		assignExpr := ""
		callExpr := fmt.Sprintf("%v.embed.%v(%v)", receiverName, methodName, paramNames)
		returnExpr := ""
		methodReturnTypes := astutil.MethodReturnTypes(m)
		if len(methodReturnTypes) > 0 {
			retVars := astutil.MethodReturnVars(m)
			for i, r := range retVars {
				varExpr += fmt.Sprintf("var %v %v\n", r, methodReturnTypes[i])
			}
			varExpr = varExpr[:len(varExpr)-1]
			assignExpr = fmt.Sprintf("%v = ", strings.Join(retVars, ", "))
			returnExpr = fmt.Sprintf(`
				return %v
				`, strings.Join(retVars, ", "))
		}
		sExpr := fmt.Sprintf(`
	%v
	%v.ops<-func() {%v%v}
	<-t.tick
	%v

`, varExpr, receiverName, assignExpr, callExpr, returnExpr)

		sExpr = fmt.Sprintf(`func(){%v}`, sExpr)
		expr, err := parser.ParseExpr(sExpr)
		if err != nil {
			panic(err)
		}
		// astutil.SetReceiverName(m, "t")
		astutil.SetReceiverTypeName(m, dstConcrete)
		astutil.SetReceiverPointer(m, true)
		m.Body = expr.(*ast.FuncLit).Body
		fmt.Fprintf(dest, "// %v is channeled\n", methodName)
		m.Doc = nil // clear the doc.
		fmt.Fprintf(dest, "%v\n", astutil.Print(m))

		mIds := astutil.GetSignatureImportIdentifiers(m)
		mImports := astutil.GetImportPaths(pkg, mIds)
		for i := range mImports {
			fileOut.AddImport(mImports[i], mIds[i])
		}
	}

	fmt.Fprintf(dest, `// Transact execute one op.`)
	fmt.Fprintln(dest, "")
	fmt.Fprintf(dest, `func (t *%v) Transact(f func(*%v))  {
	ref := &t.embed
	f(ref)
	t.embed = *ref
}`, destName, srcName)
	fmt.Fprintln(dest, "")

	fmt.Fprintf(dest, `// Start the main loop
	func (%v *%v) Start(){
		for {
			select{
			case op:=<-%v.ops:
				op()
				%v.tick<-true
			case <-%v.stop:
				return
			}
		}
	}
	`, receiverName, dstConcrete, receiverName, receiverName, receiverName)

	fmt.Fprintln(dest)
	fmt.Fprintf(dest, `// Stop the main loop
	func (%v *%v) Stop(){
		%v.stop <- true
	}
	`, receiverName, dstConcrete, receiverName)
	fmt.Fprintln(dest)

	return nil
}
