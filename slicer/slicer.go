// Package slicer generates typed slice.
package slicer

import (
	"fmt"
	"log"
	"path/filepath"
	"strings"

	"github.com/mh-cbon/astutil"
	"github.com/mh-cbon/ggt/opts"
	"github.com/mh-cbon/ggt/utils"
)

// Cmd slicer
type Cmd struct{}

// Run slicer command
func (c Cmd) Run(options *opts.Cli) {
	outPkg := options.OutPkg
	contract := options.Contract

	if options.Help {
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

		if err := processType(contract, todo, fileOut); err != nil {
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
	fmt.Printf(`ggt [options] slicer ...[FromTypeName:ToTypeName]

generates typed slice

[options]

    -c        Create a contract of the generated type.
    -p        Force out package name

...[FromTypeName:ToTypeName]

    A list of types such as src:dst.
    A type is defined by its package path and its type name,
    [pkgpath/]name.
    If the Package path is empty, it is set to the package name being generated.
    Name can be a valid type identifier such as TypeName, *TypeName, []TypeName

Example

    ggt -c slicer MySrcType:gen/*NewGenType
    ggt -c slicer myModule/*MySrcType:gen/NewGenType
`)
}

func processType(contract bool, todo utils.TransformArg, fileOut *utils.FileOut) error {
	dest := &fileOut.Body
	srcName := todo.FromTypeName
	destName := todo.ToTypeName
	srcConcrete := astutil.GetUnpointedType(srcName)
	destPointed := astutil.GetPointedType(destName)
	destConcrete := astutil.GetUnpointedType(destName)
	srcIsPointer := astutil.IsAPointedType(srcName)
	srcIsBasic := astutil.IsBasic(srcName)

	srcNameFq := srcName
	if todo.FromPkgPath != todo.ToPkgPath && !astutil.IsBasic(todo.FromTypeName) {
		srcNameFq = fmt.Sprintf("%v.%v", filepath.Base(todo.FromPkgPath), srcConcrete)
		if srcIsPointer {
			srcNameFq = "*" + srcNameFq
		}
		fileOut.AddImport(todo.FromPkgPath, todo.FromPkgID)
	}

	fmt.Fprintf(dest, `// %v implements a typed slice of %v`, destConcrete, srcNameFq)
	fmt.Fprintln(dest, "")
	fmt.Fprintf(dest, `type %v struct {items []%v}`, destConcrete, srcNameFq)
	fmt.Fprintln(dest, "")

	fmt.Fprintf(dest, `// New%v creates a new typed slice of %v`, destConcrete, srcNameFq)
	fmt.Fprintln(dest, "")
	fmt.Fprintf(dest, `func New%v() %v {
 return &%v{items: []%v{}}
}`, destConcrete, destPointed, destConcrete, srcNameFq)

	fmt.Fprintln(dest, "")

	fmt.Fprintf(dest, `// Push appends every %v`, srcNameFq)
	fmt.Fprintln(dest, "")
	fmt.Fprintf(dest, `func (t %v) Push(x ...%v) %v {
 t.items = append(t.items, x...)
 return t
}`, destPointed, srcNameFq, destPointed)

	fmt.Fprintln(dest, "")

	fmt.Fprintf(dest, `// Unshift prepends every %v`, srcNameFq)
	fmt.Fprintln(dest, "")
	fmt.Fprintf(dest, `func (t %v) Unshift(x ...%v) %v {
	t.items = append(x, t.items...)
	return t
}`, destPointed, srcNameFq, destPointed)

	fmt.Fprintln(dest, "")

	fmt.Fprintf(dest, `// Pop removes then returns the last %v.`, srcNameFq)
	fmt.Fprintln(dest, "")
	fmt.Fprintf(dest, `func (t %v) Pop() %v {
 var ret %v
 if len(t.items)>0 {
  ret = t.items[len(t.items)-1]
  t.items = append(t.items[:0], t.items[len(t.items)-1:]...)
 }
 return ret
}`, destPointed, srcNameFq, srcNameFq)

	fmt.Fprintln(dest, "")

	fmt.Fprintf(dest, `// Shift removes then returns the first %v.`, srcNameFq)
	fmt.Fprintln(dest, "")
	fmt.Fprintf(dest, `func (t %v) Shift() %v {
  var ret %v
  if len(t.items)>0 {
    ret = t.items[0]
    t.items = append(t.items[:0], t.items[1:]...)
  }
  return ret
}`, destPointed, srcNameFq, srcNameFq)

	fmt.Fprintln(dest, "")

	fmt.Fprintf(dest, `// Index of given %v. It must implements Ider interface.`, srcNameFq)
	fmt.Fprintln(dest, "")
	if srcIsBasic == false {
		fmt.Fprintf(dest, `func (t %v) Index(s %v) int {
	  ret := -1
	  for i,item:= range t.items {
			if s.GetID()==item.GetID() {
				ret = i
				break
			}
	  }
	  return ret
	}`, destPointed, srcNameFq)

	} else if srcIsPointer && srcIsBasic { // needed ?
		fmt.Fprintf(dest, `func (t %v) Index(s %v) int {
	  ret := -1
	  for i,item:= range t.items {
			if *s==*item {
				ret = i
				break
			}
	  }
	  return ret
	}`, destPointed, srcName)

	} else {
		fmt.Fprintf(dest, `func (t %v) Index(s %v) int {
	  ret := -1
	  for i,item:= range t.items {
			if s==item {
				ret = i
				break
			}
	  }
	  return ret
	}`, destPointed, srcNameFq)
	}
	fmt.Fprintln(dest, "")

	fmt.Fprintf(dest, `// Contains returns true if s in is t.`)
	fmt.Fprintln(dest, "")
	fmt.Fprintf(dest, `func (t %v) Contains(s %v) bool {
  return t.Index(s)>-1
}`, destPointed, srcNameFq)
	fmt.Fprintln(dest, "")

	fmt.Fprintf(dest, `// RemoveAt removes a %v at index i.`, srcNameFq)
	fmt.Fprintln(dest, "")
	fmt.Fprintf(dest, `func (t %v) RemoveAt(i int) bool {
  if i>=0 && i<len(t.items) {
    t.items = append(t.items[:i], t.items[i+1:]...)
		return true
  }
  return false
}`, destPointed)

	fmt.Fprintln(dest, "")

	fmt.Fprintf(dest, `// Remove removes given %v`, srcNameFq)
	fmt.Fprintln(dest, "")
	fmt.Fprintf(dest, `func (t %v) Remove(s %v) bool {
  if i := t.Index(s); i > -1 {
    t.RemoveAt(i)
		return true
  }
  return false
}`, destPointed, srcNameFq)

	fmt.Fprintln(dest, "")

	fmt.Fprintf(dest, `// InsertAt adds given %v at index i`, srcNameFq)
	fmt.Fprintln(dest, "")
	fmt.Fprintf(dest, `func (t %v) InsertAt(i int, s %v) %v {
	if i<0 || i >= len(t.items) {
		return t
	}
	res := []%v{}
	res = append(res, t.items[:0]...)
	res = append(res, s)
	res = append(res, t.items[i:]...)
	t.items = res
	return t
}`, destPointed, srcNameFq, destPointed, srcNameFq)

	fmt.Fprintln(dest, "")

	fmt.Fprintf(dest, `// Splice removes and returns a slice of %v, starting at start, ending at start+length.`, srcNameFq)
	fmt.Fprintln(dest, "")
	fmt.Fprintf(dest, `// If any s is provided, they are inserted in place of the removed slice.`)
	fmt.Fprintln(dest, "")
	fmt.Fprintf(dest, `func (t %v) Splice(start int, length int, s ...%v) []%v {
	var ret []%v
	for i := 0; i < len(t.items); i++ {
		if i >= start && i < start+length {
			ret = append(ret, t.items[i])
		}
	}
	if start >= 0 && start+length <= len(t.items) && start+length >= 0 {
		t.items = append(
			t.items[:start],
			append(s,
				t.items[start+length:]...,
			)...,
		)
	}
  return ret
}`, destPointed, srcNameFq, srcNameFq, srcNameFq)

	fmt.Fprintln(dest, "")

	fmt.Fprintf(dest, `// Slice returns a copied slice of %v, starting at start, ending at start+length.`, srcNameFq)
	fmt.Fprintln(dest, "")
	fmt.Fprintf(dest, `func (t %v) Slice(start int, length int) []%v {
  var ret []%v
	if start >= 0 && start+length <= len(t.items) && start+length >= 0 {
		ret = t.items[start:start+length]
	}
	return ret
}`, destPointed, srcNameFq, srcNameFq)

	fmt.Fprintln(dest, "")

	fmt.Fprintf(dest, `// Reverse the slice.`)
	fmt.Fprintln(dest, "")
	fmt.Fprintf(dest, `func (t %v) Reverse() %v {
  for i, j := 0, len(t.items)-1; i < j; i, j = i+1, j-1 {
    t.items[i], t.items[j] = t.items[j], t.items[i]
  }
  return t
}`, destPointed, destPointed)

	fmt.Fprintln(dest, "")

	fmt.Fprintf(dest, `// Len of the slice.`)
	fmt.Fprintln(dest, "")
	fmt.Fprintf(dest, `func (t %v) Len() int {
  return len(t.items)
}`, destPointed)

	fmt.Fprintln(dest, "")

	fmt.Fprintf(dest, `// Set the slice.`)
	fmt.Fprintln(dest, "")
	fmt.Fprintf(dest, `func (t %v) Set(x []%v) %v {
  t.items = append(t.items[:0], x...)
	return t
}`, destPointed, srcNameFq, destPointed)
	fmt.Fprintln(dest, "")

	fmt.Fprintf(dest, `// Get the slice.`)
	fmt.Fprintln(dest, "")
	fmt.Fprintf(dest, `func (t %v) Get() []%v {
	return t.items
}`, destPointed, srcNameFq)
	fmt.Fprintln(dest, "")

	fmt.Fprintf(dest, `// At return the item at index i.`)
	fmt.Fprintln(dest, "")
	fmt.Fprintf(dest, `func (t %v) At(i int) %v {
	return t.items[i]
}`, destPointed, srcNameFq)
	fmt.Fprintln(dest, "")

	fmt.Fprintf(dest, `// Filter return a new %v with all items satisfying f.`, destName)
	fmt.Fprintln(dest, "")
	fmt.Fprintf(dest, `func (t %v) Filter(filters ...func(%v) bool) %v {
	ret := New%v()
	for _, i := range t.items {
		ok := true
		for _, f := range filters {
			ok = ok && f(i)
			if ! ok {
				break
			}
		}
		if ok {
			ret.Push(i)
		}
	}
	return ret
}`, destPointed, srcNameFq, destPointed, destConcrete)
	fmt.Fprintln(dest, "")

	// todo: handle more cases like ArrayType etc.
	fmt.Fprintf(dest, `// Map return a new %v of each items modified by f.`, destName)
	fmt.Fprintln(dest, "")

	if srcIsPointer {
		fmt.Fprintf(dest, `func (t %v) Map(mappers ...func(%v) %v) %v {
		ret := New%v()
		for _, i := range t.items {
			val := i
			for _, m := range mappers {
				val = m(val)
				if val == nil {
					break
				}
			}
			if val != nil {
				ret.Push(val)
			}
		}
		return ret
	}`, destPointed, srcNameFq, srcNameFq, destPointed, destConcrete)

	} else {
		fmt.Fprintf(dest, `func (t %v) Map(mappers ...func(%v) %v) %v {
		ret := New%v()
		for _, i := range t.items {
			val := i
			for _, m := range mappers {
				val = m(val)
			}
			ret.Push(val)
		}
		return ret
	}`, destPointed, srcNameFq, srcNameFq, destPointed, destConcrete)
	}
	fmt.Fprintln(dest, "")

	fmt.Fprintf(dest, `// First returns the first value or default.`)
	fmt.Fprintln(dest, "")
	fmt.Fprintf(dest, `func (t %v) First() %v {
	var ret %v
	if len(t.items)>0 {
		ret = t.items[0]
	}
	return ret
}`, destPointed, srcNameFq, srcNameFq)
	fmt.Fprintln(dest, "")

	fmt.Fprintf(dest, `// Last returns the last value or default.`)
	fmt.Fprintln(dest, "")
	fmt.Fprintf(dest, `func (t %v) Last() %v {
	var ret %v
	if len(t.items)>0 {
		ret = t.items[len(t.items)-1]
	}
	return ret
}`, destPointed, srcNameFq, srcNameFq)
	fmt.Fprintln(dest, "")

	fmt.Fprintf(dest, `// Empty returns true if the slice is empty.`)
	fmt.Fprintln(dest, "")
	fmt.Fprintf(dest, `func (t %v) Empty() bool {
	return len(t.items)==0
}`, destPointed)
	fmt.Fprintln(dest, "")

	fmt.Fprintf(dest, `// Transact execute one op.`)
	fmt.Fprintln(dest, "")
	fmt.Fprintf(dest, `func (t %v) Transact(f func(%v))  {
	f(t)
}`, destPointed, destPointed)
	fmt.Fprintln(dest, "")

	fileOut.AddImport("encoding/json", "json")
	// Add marshalling capabilities
	fmt.Fprintf(dest, `
//UnmarshalJSON JSON unserializes %v
func (t %v) UnmarshalJSON(b []byte) error {
	var items []%v
	if err := json.Unmarshal(b, &items); err != nil {
		return err
	}
	t.items = items
	return nil
}
`, destConcrete, destPointed, srcNameFq)
	fmt.Fprintln(dest)

	fmt.Fprintf(dest, `
//MarshalJSON JSON serializes %v
func (t %v) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.items)
}
`, destConcrete, destPointed)
	fmt.Fprintln(dest)

	fmt.Fprintln(dest, "")

	if contract {

		fmt.Fprintf(dest, `
	// %vContract are the requirements of %v
	type %vContract interface {
	`, destConcrete, destConcrete, destConcrete)

		fmt.Fprintf(dest, `Push(x ...%v) %v
	`, srcNameFq, destPointed)

		fmt.Fprintf(dest, `Unshift(x ...%v) %v
	`, srcNameFq, destPointed)

		fmt.Fprintf(dest, `Pop() %v
	`, srcNameFq)

		fmt.Fprintf(dest, `Shift() %v
	`, srcNameFq)

		fmt.Fprintf(dest, `Index(s %v) int
	`, srcNameFq)

		fmt.Fprintf(dest, `Contains(s %v) bool
	`, srcNameFq)

		fmt.Fprintf(dest, `RemoveAt(i int) bool
	`)

		fmt.Fprintf(dest, `Remove(s %v) bool
	`, srcNameFq)

		fmt.Fprintf(dest, `InsertAt(i int,s %v) %v
	`, srcNameFq, destPointed)

		fmt.Fprintf(dest, `Splice(start int, length int, s ...%v) []%v
	`, srcNameFq, srcNameFq)

		fmt.Fprintf(dest, `Slice(start int, length int) []%v
	`, srcNameFq)

		fmt.Fprintf(dest, `Reverse() %v
	`, destPointed)

		fmt.Fprintf(dest, `Set(x []%v) %v
	`, srcNameFq, destPointed)

		fmt.Fprintf(dest, `Get() []%v
	`, srcNameFq)

		fmt.Fprintf(dest, `At(i int) %v
	`, srcNameFq)

		fmt.Fprintf(dest, `Filter(filters ...func(%v) bool) %v
	`, srcNameFq, destPointed)

		fmt.Fprintf(dest, `Map(mappers ...func(%v) %v) %v
	`, srcNameFq, srcNameFq, destPointed)

		fmt.Fprintf(dest, `First() %v
	`, srcNameFq)

		fmt.Fprintf(dest, `Last() %v
	`, srcNameFq)

		fmt.Fprintf(dest, `Transact(func(%v))
	`, destPointed)

		fmt.Fprint(dest, `Len() int
	Empty() bool
	`)

		fmt.Fprintf(dest, `}
	`)
	}

	if astutil.IsBasic(todo.FromTypeName) == false {
		if err := processFilter(todo, fileOut); err != nil {
			return err
		}
		if err := processSetter(todo, fileOut); err != nil {
			return err
		}
	}

	return nil
}

func processFilter(todo utils.TransformArg, fileOut *utils.FileOut) error {

	dest := &fileOut.Body
	srcName := todo.FromTypeName
	destName := todo.ToTypeName
	srcIsPointer := astutil.IsAPointedType(srcName)
	srcConcrete := astutil.GetUnpointedType(srcName)
	destConcrete := astutil.GetUnpointedType(destName)

	pkgToLoad := todo.FromPkgPath
	if pkgToLoad == "" {
		pkgToLoad = utils.GetPkgToLoad()
	}

	prog := astutil.GetProgramFast(pkgToLoad)
	pkg := prog.Package(pkgToLoad)

	foundStruct := astutil.GetStruct(pkg, astutil.GetUnpointedType(srcName))
	if foundStruct == nil {
		log.Println("Can not locate the type " + srcName)
		return nil
	}

	srcNameFq := srcName
	if todo.FromPkgPath != todo.ToPkgPath && !astutil.IsBasic(todo.FromTypeName) {
		srcNameFq = fmt.Sprintf("%v.%v", filepath.Base(todo.FromPkgPath), srcConcrete)
		if srcIsPointer {
			srcNameFq = "*" + srcNameFq
		}
	}

	props := astutil.StructPropsDeep(pkg, foundStruct)

	newStructProps := ""
	for _, prop := range props {
		//todo: find a way to detect if the type implements Eq or something like this.
		propType := prop["type"]
		if !astutil.IsArrayType(propType) {
			propName := prop["name"]
			newStructProps += fmt.Sprintf(`By%v func(...%v) func (%v) bool
			`, propName, propType, srcNameFq)
			newStructProps += fmt.Sprintf(`Not%v func(...%v) func (%v) bool
			`, propName, propType, srcNameFq)
			if strings.Index(prop["type"], ".") > 0 {
				pType := strings.Split(prop["type"], ".")[0]
				path := astutil.GetImportPath(pkg, pType)
				if path == "" {
					log.Printf("package path not found. id:%q path:%q\n", pType, path)
				} else {
					fileOut.AddImport(path, pType)
				}
			}
		}
	}

	if newStructProps != "" {
		fmt.Fprintf(dest, "// Filter%v provides filters for a struct.\n", destConcrete)
		fmt.Fprintf(dest, `var Filter%v = struct {`, destConcrete)
		fmt.Fprintln(dest)
		fmt.Fprintln(dest, newStructProps+"\n")
		fmt.Fprintln(dest, "}{")
		for _, prop := range props {
			//todo: find a way to detect if the type implements Eq or something like this.
			propType := prop["type"]
			if !astutil.IsArrayType(propType) {
				propName := prop["name"]
				fmt.Fprintf(dest, `By%v: func(all ...%v) func(%v) bool {
					return func(o %v) bool {
						for _, v := range all {
							if o.%v==v {
								return true
							}
						}
						return false
					}
				},
					`, propName, propType, srcNameFq, srcNameFq, propName)
				fmt.Fprintf(dest, `Not%v: func(all ...%v) func(%v) bool {
					return func(o %v) bool {
						for _, v := range all {
							if o.%v==v {
								return false
							}
						}
						return true
					}
				},
					`, propName, propType, srcNameFq, srcNameFq, propName)
			}
		}
		fmt.Fprintln(dest)
		fmt.Fprintln(dest, "}")
	}

	return nil
}

func processSetter(todo utils.TransformArg, fileOut *utils.FileOut) error {

	dest := &fileOut.Body
	srcName := todo.FromTypeName
	destName := todo.ToTypeName
	srcIsPointer := astutil.IsAPointedType(srcName)
	srcConcrete := astutil.GetUnpointedType(srcName)
	destConcrete := astutil.GetUnpointedType(destName)

	prog := astutil.GetProgramFast(todo.FromPkgPath)
	pkg := prog.Package(todo.FromPkgPath)

	foundStruct := astutil.GetStruct(pkg, astutil.GetUnpointedType(srcName))
	if foundStruct == nil {
		log.Println("Can not locate the type " + srcName)
		return nil
	}

	srcNameFq := srcName
	if todo.FromPkgPath != todo.ToPkgPath && !astutil.IsBasic(todo.FromTypeName) {
		srcNameFq = fmt.Sprintf("%v.%v", filepath.Base(todo.FromPkgPath), srcConcrete)
		if srcIsPointer {
			srcNameFq = "*" + srcNameFq
		}
	}

	props := astutil.StructPropsDeep(pkg, foundStruct)

	newStructProps := ""
	for _, prop := range props {
		//todo: find a way to detect if the type implements Eq or something like this.
		propType := prop["type"]
		if !astutil.IsArrayType(propType) {
			propName := prop["name"]
			newStructProps += fmt.Sprintf("Set%v func(%v) func (%v) %v", propName, propType, srcNameFq, srcNameFq)
			newStructProps += "\n"
		}
	}

	if newStructProps != "" {
		fmt.Fprintf(dest, "// Setter%v provides sets properties.\n", destConcrete)
		fmt.Fprintf(dest, `var Setter%v = struct {`, destConcrete)
		fmt.Fprintln(dest)
		fmt.Fprintln(dest, newStructProps+"\n")
		fmt.Fprintln(dest, "}{")
		for _, prop := range props {
			//todo: find a way to detect if the type implements Eq or something like this.
			propType := prop["type"]
			if !astutil.IsArrayType(propType) {
				propName := prop["name"]
				fmt.Fprintf(dest, `Set%v: func(v %v) func(%v) %v {`, propName, propType, srcNameFq, srcNameFq)
				fmt.Fprintf(dest, `	return func(o %v) %v {`, srcNameFq, srcNameFq)
				fmt.Fprintf(dest, `	 o.%v = v
														 return o`, propName)
				fmt.Fprintf(dest, `	}`)
				fmt.Fprintf(dest, `},`)
				fmt.Fprintln(dest, "")
			}
		}
		fmt.Fprintln(dest)
		fmt.Fprintln(dest, "}")
	}

	return nil
}
