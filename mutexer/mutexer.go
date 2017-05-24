// Package mutexer generates mutexed type.
package mutexer

import (
	"fmt"
	"go/ast"
	"go/parser"
	"log"

	"github.com/mh-cbon/astutil"
	"github.com/mh-cbon/ggt/opts"
	"github.com/mh-cbon/ggt/utils"
)

var name = "mutexer"
var version = "0.0.0"

// Cmd mutexer
type Cmd struct{}

// Run mutexer command
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
	fmt.Printf(``)
}

func processType(todo utils.TransformArg, fileOut *utils.FileOut) error {

	dest := &fileOut.Body
	srcName := todo.FromTypeName
	destName := todo.ToTypeName
	srcConcrete := astutil.GetUnpointedType(srcName)

	prog := astutil.GetProgramFast(todo.FromPkgPath)
	pkg := prog.Package(todo.FromPkgPath)

	foundTypes := astutil.FindTypes(pkg)
	foundMethods := astutil.FindMethods(pkg)
	foundCtors := astutil.FindCtors(pkg, foundTypes)

	dstConcrete := astutil.GetUnpointedType(destName)
	dstStar := astutil.GetPointedType(destName)

	// srcIsPointer := astutil.IsAPointedType(srcName)
	// srcNameFq := srcName
	if todo.FromPkgPath != todo.ToPkgPath && !astutil.IsBasic(todo.FromTypeName) {
		// srcNameFq = fmt.Sprintf("%v.%v", filepath.Base(todo.FromPkgPath), srcConcrete)
		// if srcIsPointer {
		// 	srcNameFq = "*" + srcNameFq
		// }
		fileOut.AddImport(todo.FromPkgPath, todo.FromPkgID)
	}

	fmt.Fprintf(dest, "// %v mutexes a %v\n", dstConcrete, srcConcrete)
	fmt.Fprintf(dest, `type %v struct{
		 embed %v
		 mutex *sync.Mutex
		 }
		 `, dstConcrete, srcName)

	ctorParams := ""
	ctorName := ""
	ctorIsPointer := false
	if x, ok := foundCtors[srcConcrete]; ok {
		ctorParams = astutil.MethodParams(x)
		ctorIsPointer = astutil.MethodReturnPointer(x)
		ctorName = "New" + srcConcrete
	}

	fmt.Fprintf(dest, `// New%v constructs a new %v
		`, dstConcrete, destName)
	fmt.Fprintf(dest, `func New%v(%v) %v {
		`, dstConcrete, ctorParams, dstStar)
	fmt.Fprintf(dest, `ret := &%v{}
		`, dstConcrete)
	if ctorName != "" {
		fmt.Fprintf(dest, "	embed := %v(%v)\n", ctorName, ctorParams)
		if !ctorIsPointer && astutil.IsAPointedType(srcName) {
			fmt.Fprintf(dest, "	ret.embed = *embed\n")
		} else {
			fmt.Fprintf(dest, "	ret.embed = embed\n")
		}
	}
	fmt.Fprintf(dest, `	ret.mutex = &sync.Mutex{}
		return ret
		}
	`)

	receiverName := "t"

	for _, m := range foundMethods[srcConcrete] {
		withEllipse := astutil.MethodHasEllipse(m)
		paramNames := astutil.MethodParamNamesInvokation(m, withEllipse)
		receiverName = astutil.ReceiverName(m)
		methodName := astutil.MethodName(m)

		if methodName == "Transact" {
			continue
		}

		callExpr := fmt.Sprintf("%v.embed.%v(%v)", receiverName, methodName, paramNames)
		sExpr := fmt.Sprintf(`
  %v.mutex.Lock()
  defer %v.mutex.Unlock()
	return %v`,
			receiverName, receiverName, callExpr)
		sExpr = fmt.Sprintf("func(){%v\n}", sExpr)
		expr, err := parser.ParseExpr(sExpr)
		if err != nil {
			panic(err)
		}
		astutil.SetReceiverTypeName(m, destName)
		astutil.SetReceiverPointer(m, true)
		m.Body = expr.(*ast.FuncLit).Body
		m.Doc = nil // clear the doc.
		fmt.Fprintf(dest, "// %v is mutexed\n", methodName)
		fmt.Fprintf(dest, "%v\n", astutil.Print(m))
	}

	if !astutil.HasMethod(pkg, srcConcrete, "UnmarshalJSON") ||
		!astutil.HasMethod(pkg, srcConcrete, "MarshalJSON") {
		fileOut.AddImport("encoding/json", "")
	}

	fmt.Fprintf(dest, `// Transact execute one op.`)
	fmt.Fprintln(dest, "")
	fmt.Fprintf(dest, `func (t *%v) Transact(f func(*%v))  {
	ref := &t.embed
	f(ref)
	t.embed = *ref
}`, destName, srcName)
	fmt.Fprintln(dest, "")

	if astutil.HasMethod(pkg, srcConcrete, "UnmarshalJSON") == false {
		// Add marshalling capabilities
		fmt.Fprintf(dest, `
		//UnmarshalJSON JSON unserializes %v
		func (%v %v) UnmarshalJSON(b []byte) error {
			t.mutex.Lock()
			defer t.mutex.Unlock()
			var items []%v
			if err := json.Unmarshal(b, &items); err != nil {
				return err
			}
			t.items = items
			return nil
		}
		`, dstConcrete, receiverName, dstStar, srcName)
		fmt.Fprintln(dest)
	}
	if astutil.HasMethod(pkg, srcConcrete, "UnmarshalJSON") == false {
		fmt.Fprintf(dest, `
		//MarshalJSON JSON serializes %v
		func (%v %v) MarshalJSON() ([]byte, error) {
			t.mutex.Lock()
			defer t.mutex.Unlock()
			return json.Marshal(t.items)
		}
		`, dstConcrete, receiverName, dstStar)
		fmt.Fprintln(dest)
	}

	return nil
}
