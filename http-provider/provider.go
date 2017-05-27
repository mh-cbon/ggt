// Package provider generates typed http provider.
package provider

import (
	"fmt"
	"log"
	"path/filepath"
	"strings"

	"github.com/mh-cbon/astutil"
	"github.com/mh-cbon/ggt/opts"
	"github.com/mh-cbon/ggt/utils"
)

var name = "http-provider"
var version = "0.0.0"

// Cmd mutexer
type Cmd struct{}

// Run mutexer command
func (c Cmd) Run(options *opts.Cli) {

	outPkg := options.OutPkg
	help := options.Help
	mode := options.Mode

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

		if err := processType(mode, todo, fileOut); err != nil {
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
	fmt.Printf(`ggt [options] http-provider ...[FromTypeName:ToTypeName]

generates http oriented implementation of given type.

[options]

    -p        Force out package name
    -mode     Generation mode (rpc|route).

...[FromTypeName:ToTypeName]

    A list of types such as src:dst.
    A type is defined by its package path and its type name,
    [pkgpath/]name.
    If the Package path is empty, it is set to the package name being generated.
    Name can be a valid type identifier such as TypeName, *TypeName, []TypeName

Example

    ggt -c http-provider MySrcType:gen/*NewGenType
    ggt -c http-provider myModule/*MySrcType:gen/NewGenType
`)
}

var httpCookieType = "http.Cookie"

func processType(mode string, todo utils.TransformArg, fileOut *utils.FileOut) error {

	dest := &fileOut.Body
	srcName := todo.FromTypeName
	destName := todo.ToTypeName

	prog := astutil.GetProgramFast(todo.FromPkgPath)
	pkg := prog.Package(todo.FromPkgPath)
	foundMethods := astutil.FindMethods(pkg)

	srcConcrete := astutil.GetUnpointedType(srcName)
	// the json input must provide a key/value for each params.
	structType := astutil.FindStruct(pkg, srcConcrete)
	structComment := astutil.GetComment(prog, structType.Pos())
	// todo: might do better to send only annotations or do other improvemenets.
	structComment = makeCommentLines(structComment)
	dstStar := astutil.GetPointedType(destName)
	structAnnotations := astutil.GetAnnotations(structComment, "@")

	srcIsPointer := astutil.IsAPointedType(srcName)
	srcNameFq := srcName
	if todo.FromPkgPath != todo.ToPkgPath && !astutil.IsBasic(todo.FromTypeName) {
		srcNameFq = fmt.Sprintf("%v.%v", filepath.Base(todo.FromPkgPath), srcConcrete)
		if srcIsPointer {
			srcNameFq = "*" + srcNameFq
		}
		fileOut.AddImport(todo.FromPkgPath, todo.FromPkgID)
	}

	addlog := func(receiver, errName string, subjects ...string) string {
		s := fmt.Sprintf("%q, ", "")
		if len(subjects) > 0 {
			subjects = append(subjects, destName)
			s = ""
			for _, sub := range subjects {
				if len(sub) > 0 {
					if strings.HasPrefix(sub, "!") {
						s += fmt.Sprintf(`%v, `, sub[1:])
					} else {
						s += fmt.Sprintf(`%q, `, sub)
					}
				}
			}
			s = s[:len(s)-2]
		}
		return fmt.Sprintf(`%v.Log.Handle(nil, nil, %v, %v)`, receiver, errName, s)
	}

	fileOut.AddImport("io", "")
	fileOut.AddImport("net/http", "")
	fileOut.AddImport("strconv", "")
	fileOut.AddImport("github.com/mh-cbon/ggt/lib", "ggt")

	// cheat.
	for _, x := range []string{"strconv.Atoi", "io.Copy", "http.StatusOK"} {
		fmt.Fprintf(dest, `var xx%v = %v
			`, utils.Hash(fileOut.Path+x), x)
	}

	// Declare the new type
	fmt.Fprintf(dest, `
// %v is an httper of %v.
%v
		`, destName, srcName, structComment)
	fmt.Fprintf(dest, `type %v struct{
	embed %v
	Log ggt.HTTPLogger
	Session ggt.SessionStoreProvider
}
		`, destName, srcNameFq)

	// Make the constructor
	fmt.Fprintf(dest, `// New%v constructs an httper of %v
`, destName, srcName)

	fmt.Fprintf(dest, `func New%v(embed %v) *%v {
	ret := &%v{
		embed: embed,
		Log: &ggt.VoidLog{},
		Session: &ggt.VoidSession{},
	}
	%v
  return ret
}
`, destName, srcNameFq, destName, destName, addlog("ret", "nil", "constructor"))

	// wrap each method
	for _, m := range foundMethods[srcConcrete] {
		methodName := astutil.MethodName(m)

		// ensure it is desired to facade this method.
		if astutil.IsExported(methodName) == false {
			continue
		}
		if strings.HasPrefix(methodName, "Finalize") {
			continue
		}

		params := astutil.MethodInputs(m)
		lParamNames := astutil.MethodInputNames(m)
		paramNames := astutil.MethodParamNames(m)
		paramTypes := astutil.MethodParamTypes(m)
		lParamTypes := astutil.MethodInputTypes(m)
		comment := astutil.GetComment(prog, m.Pos())
		comment = makeCommentLines(comment)
		// annotations := astutil.GetAnnotations(comment, "@")
		// annotations = mergeAnnotations(structAnnotations, annotations)
		hasEllipse := astutil.MethodHasEllipse(m)
		// _, hasRoute := annotations["route"]

		addlog = func(receiver, errName string, subjects ...string) string {
			s := fmt.Sprintf("%q, ", "")
			if len(subjects) > 0 {
				subjects = append(subjects, destName, methodName)
				s = ""
				for _, sub := range subjects {
					if len(sub) > 0 {
						if strings.HasPrefix(sub, "!") {
							s += fmt.Sprintf(`%v, `, sub[1:])
						} else {
							s += fmt.Sprintf(`%q, `, sub)
						}
					}
				}
				s = s[:len(s)-2]
			}
			return fmt.Sprintf(`%v.Log.Handle(w, r, %v, %v)`, receiver, errName, s)
		}

		errHandler := func(errName string, subjects ...string) string {
			subjects = append(subjects, "error")
			var ret string
			if astutil.HasMethod(pkg, srcConcrete, methodName+"Finalizer") {
				ret = fmt.Sprintf(`
					%v
					t.embed.%vFinalizer(w, r, %v)
				`, addlog("t", errName, subjects...), methodName, errName)
			} else if astutil.HasMethod(pkg, srcConcrete, "Finalizer") {
				ret = fmt.Sprintf(`
					%v
					t.embed.Finalizer(w, r, %v)
				`, addlog("t", errName, subjects...), errName)
			} else {
				ret = fmt.Sprintf(`
					%v
					http.Error(w, %v.Error(), http.StatusInternalServerError)
				`, addlog("t", errName, subjects...), errName)
			}
			if ret != "" {
				ret = fmt.Sprintf(`
				if %v != nil {
					%v
					return
				}
				`, errName, ret)
			}
			return ret
		}

		if mode == rpcMode {
			importIDs := astutil.GetSignatureImportIdentifiers(m)
			for _, i := range importIDs {
				path := astutil.GetImportPath(pkg, i)
				if path == "" {
					log.Printf("package path not found. id:%q path:%q\n", i, path)
				} else {
					fileOut.AddImport(path, i)
				}
			}
		}

		bodyFunc := ""

		if hasPostParam(lParamNames) {
			bodyFunc += fmt.Sprintf(`
			{
				err := r.ParseForm()
				%v
			}
			`, errHandler("err", "parseform"))
		}

		if mode == "route" {
			if hasRouteParam(lParamNames) {
				fileOut.AddImport("github.com/gorilla/mux", "")
				bodyFunc += fmt.Sprintf(`
					xxRouteVars := mux.Vars(r)
				`)
			}
			if hasGetParam(lParamNames) {
				bodyFunc += fmt.Sprintf(`
					xxURLValues := r.URL.Query()
				`)
			}
		}

		if mode == rpcMode {
			if len(paramTypes) > 0 {
				fileOut.AddImport("encoding/json", "json")
				bodyFunc += fmt.Sprintf(`input := struct{
						%v
					}{}
					decErr := json.NewDecoder(r.Body).Decode(&input)
					%v
				`, mapParamsToStruct(params, hasEllipse),
					errHandler("decErr", "req", "json", "decode"))

			}

		} else {
			for i, paramName := range lParamNames {
				paramType := lParamTypes[i]

				if paramName == "postValues" && paramType == mapStringSliceString {
					bodyFunc += fmt.Sprintf(`%v := r.PostForm
							`, paramName)

				} else if paramName == "postValues" && paramType == mapStringString {
					bodyFunc += fmt.Sprintf(`%v := map[string]string{}
					{
						for k,v := range r.PostForm {
							if len(v)>0{
								%v[k] = v[0]
							}
						}
					}
					`, paramName, paramName)

				} else if paramName == "getValues" && paramType == mapStringSliceString {
					bodyFunc += fmt.Sprintf(`%v := xxURLValues
					`, paramName)

				} else if paramName == "getValues" && paramType == mapStringString {
					bodyFunc += fmt.Sprintf(`%v := map[string]string{}
					{
						for k,v := range xxURLValues {
							if len(v)>0{
								%v[k] = v[0]
							}
						}
					}
				`, paramName, paramName)

				} else if paramName == "reqValues" && paramType == mapStringSliceString {
					bodyFunc += fmt.Sprintf(`
						var %v map[string][]string
						{
							%v = map[string][]string{}
							xxTempValue := %v
							for k,v := range xxRouteVars {
								if _, ok := xxTempValue[k]; ok {
									xxTempValue[k] = append(xxTempValue[k], v)
								} else {
									xxTempValue[k] = []string{v}
								}
							}
							for k,v := range xxURLValues {
								if _, ok := xxTempValue[k]; ok {
									xxTempValue[k] = append(xxTempValue[k], v...)
								} else {
									xxTempValue[k] = v
								}
							}
							for k,v := range r.Form {
								if _, ok := xxTempValue[k]; ok {
									xxTempValue[k] = append(xxTempValue[k], v...)
								} else {
									xxTempValue[k] = v
								}
							}
					}
					`, paramName, paramName, paramName)

				} else if paramName == "reqValues" && paramType == mapStringString {
					bodyFunc += fmt.Sprintf(`
						var %v map[string]string
						{
							%v = map[string]string{}
							xxTempValue := %v
							for k,v := range xxRouteVars {
								if len(v)>0{
									xxTempValue[k] = v
								}
							}
							for k,v := range xxURLValues {
								if len(v)>0{
									if _, ok := xxTempValue[k]; ok {
										for _, vv := range v{
											if len(vv)>0 {
												xxTempValue[k] = vv
												break
											}
										}
									} else {
										for _, vv := range v{
											if len(vv)>0 {
												xxTempValue[k] = vv
												break
											}
										}
									}
								}
							}
							for k,v := range r.Form {
								if len(v)>0{
									if _, ok := xxTempValue[k]; ok {
										for _, vv := range v{
											if len(vv)>0 {
												xxTempValue[k] = vv
												break
											}
										}
									} else {
										for _, vv := range v{
											if len(vv)>0 {
												xxTempValue[k] = vv
												break
											}
										}
									}
								}
							}
					}
					`, paramName, paramName, paramName)

				} else if paramName == "urlValues" && paramType == mapStringSliceString {
					bodyFunc += fmt.Sprintf(`
							var %v map[string][]string
							{
								%v = map[string][]string{}
								xxTempValue := %v
								for k,v := range xxRouteVars {
									if _, ok := xxTempValue[k]; ok {
										xxTempValue[k] = append(xxTempValue[k], v)
									} else {
										xxTempValue[k] = []string{v}
									}
								}
								for k,v := range xxURLValues {
									if _, ok := xxTempValue[k]; ok {
										xxTempValue[k] = append(xxTempValue[k], v...)
									} else {
										xxTempValue[k] = v
									}
								}
						}
						`, paramName, paramName, paramName)

				} else if paramName == "urlValues" && paramType == mapStringString {
					bodyFunc += fmt.Sprintf(`
						var %v map[string]string
						{
							%v = map[string]string{}
							xxTempValue := %v
							for k,v := range xxRouteVars {
								if len(v)>0{
									xxTempValue[k] = v
								}
							}
							for k,v := range xxURLValues {
								if len(v)>0{
									if _, ok := xxTempValue[k]; ok {
										for _, vv := range v{
											if len(vv)>0 {
												xxTempValue[k] = vv
												break
											}
										}
									} else {
										for _, vv := range v{
											if len(vv)>0 {
												xxTempValue[k] = vv
												break
											}
										}
									}
								}
							}
						}
					`, paramName, paramName, paramName)

				} else if paramName == "routeValues" && paramType == mapStringString {
					bodyFunc += fmt.Sprintf(`%v := xxRouteVars
					`, paramName)

				} else if paramName == "cookieValues" && paramType == mapStringString {
					bodyFunc += fmt.Sprintf(`var %v map[string]string
					{
						for _, v := range r.Cookies() {
							%v[v.Name] = v.Value
						}
					}
					`, paramName, paramName)

				} else if paramName == "cookieValues" && paramType == "[]*http.Cookie" {
					bodyFunc += fmt.Sprintf(`%v := r.Cookies()
					`, paramName)

				} else if strings.HasPrefix(paramName, "session") {
					k := strings.ToLower(paramName[7:])
					bodyFunc += fmt.Sprintf(`var %v %v
					`, paramName, paramType)
					bodyFunc += fmt.Sprintf(`{
					`)
					bodyFunc += fmt.Sprintf(`
						store%v, storeErr := t.Session.Get(r, %q)
						%v
					`, paramName, k, errHandler("storeErr", "session", "store", "get", "error", paramName))
					bodyFunc += fmt.Sprintf(`
						defer func() {
							saveErr := store%v.Save(r,w)
							%v
						}()
					`, paramName, errHandler("saveErr", "session", "save", "error", paramName))
					bodyFunc += fmt.Sprintf(`
						val%v, getErr := store%v.Get()
						%v
					`, paramName, paramName, errHandler("getErr", "session", "read", "error", paramName))
					bodyFunc += fmt.Sprintf(`
						%v = val%v
					`, paramName, paramName)
					bodyFunc += fmt.Sprintf(`
					}
					`)

				} else if strings.HasPrefix(paramName, "get") {
					k := strings.ToLower(paramName[3:])
					bodyFunc += fmt.Sprintf("var %v %v", paramName, paramType)

					if astutil.IsArrayType(paramType) {
						bodyFunc += fmt.Sprintf(`
										if _, ok := xxURLValues[%q]; ok {
											xxTmp%v := xxURLValues[%q]
											%v
											%v
										}
									`, k, paramName, k,
							addlog("t", "nil", "input", "get", k /*, "!xxTmp"+paramName*/),
							convertStrTo("xxTmp"+paramName, paramName, paramType, errHandler, destName, methodName, "get"))

					} else {
						bodyFunc += fmt.Sprintf(`
									if _, ok := xxURLValues[%q]; ok {
										xxTmp%v := xxURLValues.Get(%q)
										%v
										%v
									}
								`, k, paramName, k,
							addlog("t", "nil", "input", "get", k, "!xxTmp"+paramName),
							convertStrTo("xxTmp"+paramName, paramName, paramType, errHandler, destName, methodName, "get"))
					}

				} else if strings.HasPrefix(paramName, "post") {
					k := strings.ToLower(paramName[4:])
					bodyFunc += fmt.Sprintf("var %v %v", paramName, paramType)

					if astutil.IsArrayType(paramType) {
						bodyFunc += fmt.Sprintf(`
							if _, ok := r.Form[%q]; ok {
										xxTmp%v := r.Form[%q]
										%v
										%v
									}
								`, k, paramName, k,
							addlog("t", "nil", "input", "form", k /*, "!xxTmp"+paramName*/),
							convertStrTo("xxTmp"+paramName, paramName, paramType, errHandler, destName, methodName, "form"))

					} else {
						bodyFunc += fmt.Sprintf(`
									if _, ok := r.Form[%q]; ok {
										xxTmp%v := r.FormValue(%q)
										%v
										%v
									}
								`, k, paramName, k,
							addlog("t", "nil", "input", "form", k, "!xxTmp"+paramName),
							convertStrTo("xxTmp"+paramName, paramName, paramType, errHandler, "post"))
					}

				} else if strings.HasPrefix(paramName, "route") {
					k := strings.ToLower(paramName[5:])
					bodyFunc += fmt.Sprintf("var %v %v", paramName, paramType)
					bodyFunc += fmt.Sprintf(`
							if _, ok := xxRouteVars[%q]; ok {
								xxTmp%v := xxRouteVars[%q]
								%v
								%v
							}
						`, k, paramName, k,
						addlog("t", "nil", "input", "route", k, "!xxTmp"+paramName),
						convertStrTo("xxTmp"+paramName, paramName, paramType, errHandler, "route"))

				} else if strings.HasPrefix(paramName, "url") {
					k := strings.ToLower(paramName[3:])
					bodyFunc += fmt.Sprintf("var %v %v", paramName, paramType)
					bodyFunc += fmt.Sprintf(`
						if _, ok := xxRouteVars[%q]; ok {
							xxTmp%v := xxRouteVars[%q]
							%v
						}`,
						k, paramName, k, convertStrTo("xxTmp"+paramName, paramName, paramType, errHandler, "url", "route"))

					bodyFunc += fmt.Sprintf(`else if _, ok := xxURLValues[%q]; ok {
							xxTmp%v := xxURLValues.Get(%q)
								%v
							}
							`, k, paramName, k, convertStrTo("xxTmp"+paramName, paramName, paramType, errHandler, "url", "get"))

				} else if strings.HasPrefix(paramName, "req") {
					k := strings.ToLower(paramName[3:])
					bodyFunc += fmt.Sprintf("var %v %v", paramName, paramType)
					bodyFunc += fmt.Sprintf(`
						if _, ok := xxRouteVars[%q]; ok {
							xxTmp%v := xxRouteVars[%q]
							%v
						}`,
						k, paramName, k, convertStrTo("xxTmp"+paramName, paramName, paramType, errHandler, "route"))

					bodyFunc += fmt.Sprintf(`else if _, ok := xxURLValues[%q]; ok {
						xxTmp%v := xxURLValues.Get(%q)
						%v
						}`, k, paramName, k, convertStrTo("xxTmp"+paramName, paramName, paramType, errHandler, "get"))

					bodyFunc += fmt.Sprintf(`else if _, ok := r.Form[%q]; ok {
								xxTmp%v := r.FormValue(%q)
								%v
							}
						`, k, paramName, k, convertStrTo("xxTmp"+paramName, paramName, paramType, errHandler, "form"))

				} else if strings.HasPrefix(paramName, "cookie") {
					k := strings.ToLower(paramName[6:])
					bodyFunc += fmt.Sprintf("var %v %v", paramName, paramType)
					if paramType == "http.Cookie" {
						bodyFunc += fmt.Sprintf(`
							{
								c, cookieErr := r.Cookie(%q)
								%v
								if c != nil {
									%v = *c
								}
							}
							`, k, errHandler("cookieErr", "req", "cookie", "error"), paramName)
					} else if paramType == "*http.Cookie" {
						bodyFunc += fmt.Sprintf(`
							{
								c, cookieErr := r.Cookie(%q)
								%v
								%v = c
							}
							`, k, errHandler("cookieErr", "req", "cookie", "error"), paramName)
					} else {
						bodyFunc += fmt.Sprintf(`
							{
								c, cookieErr := r.Cookie(%q)
								%v
								if c != nil {
									%v
								}
							}
							`, k, errHandler("cookieErr", "req", "cookie", "error"), convertStrTo("c.Value", paramName, paramType, errHandler, "route"))
					}

				} else if paramName == "jsonReqBody" {

					paramPkgID := astutil.GetPkgID(paramType)
					if paramPkgID != "" {
						fileOut.AddImport(astutil.GetImportPath(pkg, paramPkgID), paramPkgID)
					}

					bodyFunc += fmt.Sprintf("var %v %v", paramName, paramType)
					bodyFunc += fmt.Sprintf(`
							{
								jsonReqBody = %v
								decErr := json.NewDecoder(r.Body).Decode(jsonReqBody)
								%v
						    defer r.Body.Close()
							}
						`, astutil.GetTypeToStructInit(paramType), errHandler("decErr", "req", "json", "decode"))

				} else if paramName == "headers" {
					bodyFunc += fmt.Sprintf("%v := r.Header\n", paramName)

				} else if paramType == "*"+httpCookieType {
					bodyFunc += fmt.Sprintf("%v := %v\n", paramName, "r")

					bodyFunc += fmt.Sprintf("")
					bodyFunc += fmt.Sprintf(`
						var %v %v
						{
							c, cookieErr := r.Cookie(%q)
							%v
							%v = c
						}
						`, paramName, paramType,
						paramName, errHandler("cookieErr", "req", "cookie", "error"), paramName)

				} else if paramType == httpCookieType {
					bodyFunc += fmt.Sprintf("%v := %v\n", paramName, "r")

					bodyFunc += fmt.Sprintf("")
					bodyFunc += fmt.Sprintf(`
						var %v %v
						{
							c, cookieErr := r.Cookie(%q)
							%v
							%v = *c
						}
						`, paramName, paramType,
						paramName, errHandler("cookieErr", "req", "cookie", "error"), paramName)

				} else if paramType == "*http.Request" && paramName != "r" {
					bodyFunc += fmt.Sprintf("%v := %v\n", paramName, "r")

				} else if paramType == "http.ResponseWriter" && paramName != "w" {
					bodyFunc += fmt.Sprintf("%v := %v\n", paramName, "w")

				} else if paramType == "context.Context" {
					bodyFunc += fmt.Sprintf("%v := %v.Context()\n", paramName, "r")

				} else {
					bodyFunc += fmt.Sprintf("var %v %v\n", paramName, paramType)
				}
			}

		}

		retTypes := astutil.MethodReturnTypes(m)
		retVars := astutil.MethodReturnNamesNormalized(m)
		sRetVars := strings.TrimSpace(strings.Join(retVars, ", "))
		// hasErr := astutil.MethodReturnError(m)

		// proceed to the method call on embed
		if mode == rpcMode {
			if sRetVars == "" {
				bodyFunc += fmt.Sprintf(`
						t.embed.%v(%v)
						`, methodName, mapParamNamesToStructProps(lParamNames, hasEllipse))
			} else {
				bodyFunc += fmt.Sprintf(`
					%v := t.embed.%v(%v)
					`, sRetVars, methodName, mapParamNamesToStructProps(lParamNames, hasEllipse))
			}
		}

		if mode == routeMode {

			if sRetVars == "" {
				bodyFunc += fmt.Sprintf(`
							t.embed.%v(%v)
							w.WriteHeader(200)
						`, methodName, paramNames)

			} else {
				bodyFunc += fmt.Sprintf(`
								%v := t.embed.%v(%v)
							`, sRetVars, methodName, paramNames)

				for i, retVar := range retVars {
					if retTypes[i] == "error" {
						bodyFunc += errHandler(retVar, "business")
					}
				}
			}
		}

		if mode == rpcMode {
			if sRetVars != "" {

				mappedParams := mapParamsToStruct(retTypes, false)
				// mappedParamNames := mapParamsToStructNames(retTypes)
				mappedParamValues := mapParamsToStructValues(retVars)
				bodyFunc += fmt.Sprintf(`output := struct{
			%v
		}{
			%v
		}
		`, mappedParams, mappedParamValues)

				fileOut.AddImport("encoding/json", "json")

				bodyFunc += fmt.Sprintf(`
			{
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(200)
				encErr := json.NewEncoder(w).Encode(output)
				%v
			}
				`, errHandler("encErr", "res", "json", "encode"))
			}
		} else {

			if sRetVars != "" {
				for _, retVar := range retVars {
					if retVar == "jsonResBody" {

						fileOut.AddImport("encoding/json", "json")

						bodyFunc += fmt.Sprintf(`
							{
								w.Header().Set("Content-Type", "application/json")
								w.WriteHeader(200)
								encErr := json.NewEncoder(w).Encode(jsonResBody)
								%v
							}
								`, errHandler("encErr", "res", "json", "encode"))

					}
				}
			}
		}

		for i, retVar := range retVars {
			if strings.HasPrefix(retVar, "header") {
				k := strings.ToLower(retVar[5:])
				bodyFunc += fmt.Sprintf(`
				w.Header().Set(%q %v)
				`, k, retVar)

			} else if retTypes[i] == "*"+httpCookieType {
				fileOut.AddImport("time", "")
				bodyFunc += fmt.Sprintf(`
					if %v == nil {
						http.SetCookie(w, &http.Cookie{
							Name: %q,
							Expires: time.Now().Add(-time.Hour * 24 * 100),
						})
					} else {
						http.SetCookie(w, %v)
					}
						`, retVar, retVar, retVar)

			} else if retTypes[i] == httpCookieType {
				bodyFunc += fmt.Sprintf(`http.SetCookie(w, &%v)
			`, retVar)

			} else if strings.HasPrefix(retVar, "cookie") {
				bodyFunc += fmt.Sprintf(`
				http.SetCookie(w, %v)
				`, retVar)
			}
		}

		fmt.Fprintf(dest, `// %v invoke %v.%v using the request body as a json payload.
			%v
		`, methodName, srcName, methodName, comment)

		fmt.Fprintf(dest, `func (t %v) %v(w http.ResponseWriter, r *http.Request) {
			%v
		  %v
			%v
		}

		`, dstStar, methodName, addlog("t", "nil", "begin"), bodyFunc, addlog("t", "nil", "end"))
	}

	// write the method set for the binder
	fileOut.AddImport("net/http", "")

	// declare the descriptor type
	fmt.Fprintf(dest, `// %vDescriptor describe a %v
			`, destName, dstStar)
	fmt.Fprintf(dest, `type %vDescriptor struct {
		ggt.TypeDescriptor
		about %v
			`, destName, dstStar)

	// write type properties
	for _, m := range foundMethods[srcConcrete] {
		methodName := astutil.MethodName(m)

		// ensure it is desired to facade this method.
		if astutil.IsExported(methodName) == false {
			continue
		}
		if strings.HasPrefix(methodName, "Finalize") {
			continue
		}
		fmt.Fprintf(dest, `method%v *ggt.MethodDescriptor
				`, methodName)
	}

	fmt.Fprint(dest, `}
		`)

	// write the constructor
	fmt.Fprintf(dest, `// New%vDescriptor describe a %v
			`, destName, dstStar)
	fmt.Fprintf(dest, `func New%vDescriptor (about %v) %vDescriptor {
		ret := &%vDescriptor{about: about}
			`, destName, dstStar, dstStar, destName)

	// write the setters
	for _, m := range foundMethods[srcConcrete] {
		methodName := astutil.MethodName(m)

		// ensure it is desired to facade this method.
		if astutil.IsExported(methodName) == false {
			continue
		}
		if strings.HasPrefix(methodName, "Finalize") {
			continue
		}
		comment := astutil.GetComment(prog, m.Pos())
		comment = makeCommentLines(comment)
		annotations := astutil.GetAnnotations(comment, "@")
		annotations = mergeAnnotations(structAnnotations, annotations)

		methods := fmt.Sprint("[]string{}")
		route := methodName
		if r, ok := annotations["route"]; ok {
			route = r
		}
		if m, ok := annotations["methods"]; ok {
			methods = fmt.Sprintf("[]string{%v}", stringifyList(m))
		}
		fmt.Fprintf(dest, `ret.method%v = &ggt.MethodDescriptor{
				Name     : %q,
				Handler  : about.%v,
				Route    : %q,
				Methods  : %v,
		}
				`, methodName, methodName, methodName, route, methods)
		fmt.Fprintf(dest, `ret.TypeDescriptor.Register(ret.method%v)
				`, methodName)
	}
	fmt.Fprint(dest, `return ret
		}
	`)

	// write the getters
	for _, m := range foundMethods[srcConcrete] {
		methodName := astutil.MethodName(m)

		// ensure it is desired to facade this method.
		if astutil.IsExported(methodName) == false {
			continue
		}
		if strings.HasPrefix(methodName, "Finalize") {
			continue
		}

		fmt.Fprintf(dest, `// %v returns a MethodDescriptor
				`, methodName)
		fmt.Fprintf(dest, `func (t %vDescriptor) %v() *ggt.MethodDescriptor { return t.method%v }
				`, dstStar, methodName, methodName)
	}

	return nil
}

var mapStringString = "map[string]string"
var mapStringSliceString = "map[string][]string"
var rpcMode = "rpc"
var routeMode = "route"

func stringifyList(s string) string {
	var ret []string
	for _, l := range strings.Split(s, ",") {
		l = strings.TrimSpace(l)
		if len(l) > 0 {
			ret = append(ret, fmt.Sprintf("%q", l))
		}
	}
	return strings.Join(ret, ", ")
}

func split(s string, o string) []string {
	ret := strings.Split(s, o)
	for i, p := range ret {
		ret[i] = strings.TrimSpace(p)
	}
	return ret
}

func hasPostParam(paramNames []string) bool {
	for _, paramName := range paramNames {
		if strings.HasPrefix(paramName, "post") {
			return true
		} else if strings.HasPrefix(paramName, "req") {
			return true
		}
	}
	return false
}

func hasRouteParam(paramNames []string) bool {
	for _, paramName := range paramNames {
		if strings.HasPrefix(paramName, "route") {
			return true
		} else if strings.HasPrefix(paramName, "url") {
			return true
		} else if strings.HasPrefix(paramName, "req") {
			return true
		}
	}
	return false
}

func hasGetParam(paramNames []string) bool {
	for _, paramName := range paramNames {
		if strings.HasPrefix(paramName, "get") {
			return true
		} else if strings.HasPrefix(paramName, "url") {
			return true
		} else if strings.HasPrefix(paramName, "req") {
			return true
		}
	}
	return false
}

func makeCommentLines(s string) string {
	s = strings.TrimSpace(s)
	comment := ""
	for _, k := range strings.Split(s, "\n") {
		comment += "// " + k + "\n"
	}
	comment = strings.TrimSpace(comment)
	if comment == "" {
		comment = "//"
	}
	return comment
}

func convertStrTo(fromStrVarName, toVarName, toType string, errHandler func(string, ...string) string, subjects ...string) string {
	if astutil.GetUnpointedType(toType) == "string" {
		if astutil.IsAPointedType(toType) {
			return fmt.Sprintf("%v = &%v", toVarName, fromStrVarName)
		}
		return fmt.Sprintf("%v = %v", toVarName, fromStrVarName)

	} else if astutil.GetUnpointedType(toType) == "bool" {
		if astutil.IsAPointedType(toType) {
			return fmt.Sprintf(`{
				xxTmpValue, err := strconv.ParseBool(%v)
				%v
				%v = &xxTmpValue
			}
	`, fromStrVarName, errHandler("err", subjects...), toVarName)
		}
		return fmt.Sprintf(`{
			var err error
			%v, err = strconv.ParseBool(%v)
			%v
		}
	`, toVarName, fromStrVarName, errHandler("err", subjects...))

	} else if astutil.GetUnpointedType(toType) == "int" {
		if astutil.IsAPointedType(toType) {
			return fmt.Sprintf(`{
				xxTmpValue, err := strconv.Atoi(%v)
				%v
				%v = &xxTmpValue
			}
	`, fromStrVarName, errHandler("err", subjects...), toVarName)
		}
		return fmt.Sprintf(`{
			var err error
			%v, err = strconv.Atoi(%v)
			%v
		}
	`, toVarName, fromStrVarName, errHandler("err", subjects...))

	} else if astutil.IsArrayType(toType) {
		toType = toType[2:]
		return fmt.Sprintf(`
			for _, xxValueTemp := range %v {
				var xxNewValueTemp %v
				%v
				%v = append(%v, xxNewValueTemp)
			}
		`, fromStrVarName, toType, convertStrTo("xxValueTemp", "xxNewValueTemp", toType, errHandler), toVarName, toVarName)

	}
	return fmt.Sprintf(`// conversion not handled string to %v
	`, toType)
}

func paramType(params string) string {
	x := strings.Split(params, ",")
	return x[len(x)-1]
}

func mergeAnnotations(structAnnot, methodAnnot map[string]string) map[string]string {
	ret := map[string]string{}
	for k, v := range methodAnnot {
		ret[k] = v
	}
	for k, v := range structAnnot {
		if _, ok := ret[k]; !ok {
			ret[k] = v
		}
	}
	return ret
}

func mapParamNamesToStructProps(params []string, hasEllipse bool) string {
	ret := ""
	if len(params) > 0 {
		for i := range params {
			if i == len(params)-1 && hasEllipse {
				ret += fmt.Sprintf("input.Arg%v..., ", i)
			} else {
				ret += fmt.Sprintf("input.Arg%v, ", i)
			}
		}
		ret = ret[:len(ret)-2]
	}
	return ret
}

func mapParamsToStruct(params []string, hasEllipse bool) string {
	ret := ""
	if len(params) > 0 {
		for i, p := range params {
			p = strings.TrimSpace(p)
			y := strings.Split(p, " ")
			t := strings.Replace(y[0], "...", "", -1)
			if len(y) > 1 {
				t = strings.Replace(y[1], "...", "", -1)
			}
			if i == len(params)-1 && hasEllipse {
				ret += fmt.Sprintf("Arg%v []%v\n", i, t)
			} else {
				ret += fmt.Sprintf("Arg%v %v\n", i, t)
			}
		}
	}
	return ret
}

func mapParamsToStructValues(params []string) string {
	ret := ""
	if len(params) > 0 {
		for i, p := range params {
			p = strings.TrimSpace(p)
			y := strings.Split(p, " ")
			ret += fmt.Sprintf("Arg%v: %v,\n", i, y[0])
		}
	}
	return ret
}

func mapParamsToStructNames(params []string) []string {
	var ret []string
	if len(params) > 0 {
		for i := range params {
			ret = append(ret, fmt.Sprintf(`Arg%v`, i))
		}
	}
	return ret
}
