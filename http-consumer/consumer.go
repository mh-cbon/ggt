// Package consumer generates typed http client.
package consumer

import (
	"fmt"
	"log"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/mh-cbon/astutil"
	"github.com/mh-cbon/ggt/opts"
	"github.com/mh-cbon/ggt/utils"
)

// Cmd http-consumer
type Cmd struct{}

// Run http-consumer command
func (c Cmd) Run(options *opts.Cli) {
	outPkg := options.OutPkg
	mode := options.Mode

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

	if mode != "rpc" && mode != "route" {
		panic("wrong mode: " + mode)
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
	fmt.Printf(`ggt [options] http-consumer ...[FromTypeName:ToTypeName]

generates typed slice

[options]

    -p        Force out package name
    -mode     Thep referred generation mode (rpc|route)

...[FromTypeName:ToTypeName]

    A list of types such as src:dst.
    A type is defined by its package path and its type name,
    [pkgpath/]name.
    If the Package path is empty, it is set to the package name being generated.
    Name can be a valid type identifier such as TypeName, *TypeName, []TypeName

Example

    ggt http-consumer MySrcType:gen/*NewGenType
    ggt http-consumer myModule/*MySrcType:gen/NewGenType
    ggt -mode rpc http-consumer myModule/*MySrcType:gen/NewGenType
    ggt -mode route http-consumer myModule/*MySrcType:gen/NewGenType
`)
}

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
	structAnnotations := astutil.GetAnnotations(structComment, "@")

	dstConcrete := astutil.GetUnpointedType(destName)

	srcIsPointer := astutil.IsAPointedType(srcName)
	srcNameFq := srcName
	if todo.FromPkgPath != todo.ToPkgPath && !astutil.IsBasic(todo.FromTypeName) {
		srcNameFq = fmt.Sprintf("%v.%v", filepath.Base(todo.FromPkgPath), srcConcrete)
		if srcIsPointer {
			srcNameFq = "*" + srcNameFq
		}
		fileOut.AddImport(todo.FromPkgPath, todo.FromPkgID)
	}

	// Declare the new type

	fmt.Fprintf(dest, `
// %v is an http-clienter of %v.
%v`, dstConcrete, srcName, structComment)

	if mode == routeMode {
		fmt.Fprintf(dest, `
	type %v struct{
		router *mux.Router
	  Base string
		Client *http.Client
	}
			`, dstConcrete)

	} else {
		fmt.Fprintf(dest, `
	type %v struct{
	  Base string
		Client *http.Client
	}
			`, dstConcrete)
	}

	// Make the constructor
	// should param *http.Client be an interface ?
	fmt.Fprintf(dest, `// New%v constructs an http-clienter of %v
`, dstConcrete, srcName)

	if mode == routeMode {
		fmt.Fprintf(dest, `func New%v(router *mux.Router) *%v {
	ret := &%v{
		router: router,
		Client: http.DefaultClient,
	}
  return ret
}
`, dstConcrete, dstConcrete, dstConcrete)

	} else {
		fmt.Fprintf(dest, `func New%v() *%v {
	ret := &%v{
		Client: http.DefaultClient,
	}
  return ret
}
`, dstConcrete, dstConcrete, dstConcrete)
	}

	if mode == routeMode {
		fileOut.AddImport("bytes", "")
		fileOut.AddImport("fmt", "")
		fileOut.AddImport("io", "")
		fileOut.AddImport("net/http", "")
		fileOut.AddImport("net/url", "")
		fileOut.AddImport("strings", "")
		fileOut.AddImport("github.com/gorilla/mux", "")
	} else {
		fileOut.AddImport("net/http", "")
	}

	for _, m := range foundMethods[srcConcrete] {
		methodName := astutil.MethodName(m)

		if astutil.IsExported(methodName) == false {
			continue
		}
		if strings.HasSuffix(methodName, "Finalizer") {
			continue
		}

		comment := astutil.GetComment(prog, m.Pos())
		annotations := astutil.GetAnnotations(comment, "@")
		annotations = mergeAnnotations(structAnnotations, annotations)
		params := astutil.MethodParams(m)
		lParams := commaArgsToSlice(params)
		paramNames := astutil.MethodParamNames(m)
		lParamNames := commaArgsToSlice(paramNames)
		paramTypes := astutil.MethodInputTypes(m)

		retTypes := astutil.MethodReturnTypes(m)
		retVars := astutil.MethodReturnNamesNormalized(m)
		sRetVarNames := strings.TrimSpace(strings.Join(retVars, ", "))

		handleErr := func(errVarName string) string {
			return fmt.Sprintf(`if %v!= nil {
		return %v
		}
		`, errVarName, sRetVarNames)
		}

		if mode == "rpc" {

			importIDs := astutil.GetSignatureImportIdentifiers(m)
			for _, i := range importIDs {
				fileOut.AddImport(astutil.GetImportPath(pkg, i), i)
			}

			fileOut.AddImport("errors", "")

			fmt.Fprintf(dest, `// %v constructs a request to %v
		`, methodName, methodName)

			defVal := astutil.TypesToDefVal(retTypes)
			if i := astutil.TypesIndex(retTypes, "error"); i > -1 {
				defVal[i] = fmt.Sprintf(`errors.New(%q)`, "todo")
			}
			sRet := strings.Join(defVal, ", ")

			sRetTypes := strings.TrimSpace(strings.Join(retTypes, ", "))
			fmt.Fprintf(dest, `func(t %v) %v(%v) (%v) {
			return %v
		}
		`, destName, methodName, params, sRetTypes, sRet)

		} else if route, ok := annotations["route"]; ok {

			importIDs := astutil.GetSignatureImportIdentifiers(m)
			for _, i := range importIDs {
				fileOut.AddImport(astutil.GetImportPath(pkg, i), i)
			}

			getParams := ""
			postParams := ""
			routeName, _ := annotations["name"]

			// - look for every route params
			managedParamNames := NewStringSlice()
			routeParamsExpr := []string{}
			routeParamNames := getRouteParamNamesFromRoute(mode, route)
			for _, p := range routeParamNames {
				routeParamsExpr = append(routeParamsExpr, fmt.Sprintf("%q", p))
				routeParamsExpr = append(routeParamsExpr, p)
				methodParam := getMethodParamForRouteParam(mode, lParamNames, p, managedParamNames)
				if methodParam == "" {
					log.Println("route param not identified into the method parameters " + p)
					continue
				}
				managedParamNames.Push(methodParam)
			}

			// - look for url/req/post params, not already managed by the route params
			for paramIndex, paramName := range lParamNames {
				// paramType := lParamTypes[i]
				if paramName == reqBodyVarName {
					continue
				}

				if !managedParamNames.Contains(paramName) {

					if strings.HasPrefix(paramName, "get") || strings.HasPrefix(paramName, "url") || strings.HasPrefix(paramName, "req") {
						k := strings.ToLower(paramName[3:])
						if astutil.IsAPointedType(paramTypes[paramIndex]) {
							getParams += fmt.Sprintf("url.Query().Add(%q, *%v)", k, paramName)
						} else {
							getParams += fmt.Sprintf("url.Query().Add(%q, %v)", k, paramName)
						}
						managedParamNames.Push(paramName)

					} else if strings.HasPrefix(paramName, "post") {
						k := strings.ToLower(paramName[4:])
						if astutil.IsAPointedType(paramTypes[paramIndex]) {
							postParams += fmt.Sprintf("form.Add(%q, *%v)", k, paramName)
						} else {
							postParams += fmt.Sprintf("form.Add(%q, %v)", k, paramName)
						}
						managedParamNames.Push(paramName)
					}
				}
			}

			// - forge url from the router using the route name
			urlHandling := ""
			if routeName != "" {
				k := ""
				if len(routeParamsExpr) > 0 {
					k = strings.Join(routeParamsExpr, ", ")
					k = k[:len(k)-2]
				}
				urlHandling = fmt.Sprintf(`reqURL, URLerr := t.router.Get(%q).URL(%v)
									`, routeName, k)
			} else {
				// - a route without name neeeds a jit update.
				urlHandling += fmt.Sprintf(`sReqURL := %q
										`, route)

				routeParams := getRouteParamsFromRoute(mode, route)
				managedParamNames = NewStringSlice()
				for i, routeParam := range routeParams {
					routeParamName := routeParamNames[i]
					methodParam := getMethodParamForRouteParam(mode, lParamNames, routeParamName, managedParamNames)
					if methodParam == "" {
						log.Println("route param not identified into the method parameters " + routeParam)
						continue
					}
					urlHandling += fmt.Sprintf(`sReqURL = strings.Replace(sReqURL, "%v", fmt.Sprintf("%%v", %v), 1)
													`, routeParam, methodParam)
					managedParamNames.Push(methodParam)
				}

				urlHandling += fmt.Sprintf(`reqURL, URLerr := url.ParseRequestURI(sReqURL)
									`)
			}
			urlHandling += handleErr("URLerr")

			// - if any GET params, handle them
			if getParams != "" {
				urlHandling += fmt.Sprintf("%v\n", getParams)
			}
			// - if any POST params, handle them
			if postParams != "" {
				urlHandling += fmt.Sprintf("form := url.Values{}\n%v\n", postParams)
			}

			urlHandling += fmt.Sprint("finalURL := reqURL.String()\n")

			// - build the final url
			if base, ok := annotations["base"]; ok {
				urlHandling += fmt.Sprintf("finalURL = fmt.Sprint(%q, %q, finalURL)\n", "%v%v", base)
			}
			urlHandling += fmt.Sprintf("finalURL = fmt.Sprintf(%q, t.Base, finalURL)\n", "%v%v")

			// modify method params to transform a reqBody ? to reqBody io.Reader
			methodParams := changeParamType(lParams, "reqBody", "io.Reader")

			body := ""
			// - handle the request body
			if hasReqBody(lParamNames) {
				body += fmt.Sprintf(`
				var body io.ReadWriter
				{
					var b bytes.Buffer
					body = &b
					encErr := json.NewEncoder(body).Encode(jsonReqBody)
					%v
				}
					`, handleErr("encErr"))
			}

			// - create the request object
			preferedMethod := getPreferredMethod(annotations)
			body += fmt.Sprintf("%v\n", urlHandling)
			if postParams != "" {
				body += fmt.Sprintf(" req, reqErr := http.NewRequest(%q, finalURL, strings.NewReader(form.Encode()))\n", preferedMethod)
			} else if hasReqBody(lParamNames) {
				body += fmt.Sprintf(" req, reqErr := http.NewRequest(%q, finalURL, body)\n", preferedMethod)
			} else {
				body += fmt.Sprintf(" req, reqErr := http.NewRequest(%q, finalURL, nil)\n", preferedMethod)
			}
			body += handleErr("reqErr")

			if !hasResBody(retVars) {
				body += fmt.Sprintf("_, resErr := t.Client.Do(req)\n")
				body += handleErr("resErr")

			} else {

				fileOut.AddImport("encoding/json", "json")

				body += fmt.Sprintf(`
				{
					res, resErr := t.Client.Do(req)
					%v
					decErr := json.NewDecoder(res.Body).Decode(jsonResBody)
					%v
				}
					`, handleErr("resErr"), handleErr("decErr"))
			}

			sRetVars := ""
			for i := range retVars {
				sRetVars += fmt.Sprintf(`%v %v, `, retVars[i], retTypes[i])
			}
			if sRetVars != "" {
				sRetVars = sRetVars[:len(sRetVars)-2]
				sRetVars = "(" + sRetVars + ")"
			}

			// - print the method
			fmt.Fprintf(dest, "// %v constructs a request to %v\n", methodName, route)
			fmt.Fprintf(dest, `func(t %v) %v(%v) %v {
					        %v
					        return %v
					      }
					      `, destName, methodName, strings.Join(methodParams, ","), sRetVars, body, sRetVarNames)
		}

	}

	return nil
}

func getMethodParamForRouteParam(mode string,
	methodParamNames []string,
	routeParamName string,
	managed *StringSlice) string {
	for _, methodParamName := range methodParamNames {

		if strings.HasPrefix(methodParamName, "route") {

			valueName := methodParamName[5:]
			if strings.ToLower(valueName) == strings.ToLower(routeParamName) {
				return methodParamName
			}
		} else if strings.HasPrefix(methodParamName, "get") ||
			strings.HasPrefix(methodParamName, "url") ||
			strings.HasPrefix(methodParamName, "req") {

			valueName := methodParamName[3:]
			if strings.ToLower(valueName) == strings.ToLower(routeParamName) {
				return methodParamName
			}
		}
	}
	return ""
}

//
// func handleErr(errVarName string) string {
// 	return fmt.Sprintf(`if %v!= nil {
// return nil, %v
// }
// `, errVarName, errVarName)
// }

func hasReqBody(paramNames []string) bool {
	for _, p := range paramNames {
		if p == "jsonReqBody" {
			return true
		}
	}
	return false
}

func hasResBody(paramNames []string) bool {
	for _, p := range paramNames {
		if p == "jsonResBody" {
			return true
		}
	}
	return false
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
		}
	}
	return false
}

func changeParamType(lParams []string, name, t string) []string {
	ret := []string{}
	for _, p := range lParams {
		p = strings.TrimSpace(p)
		if strings.Index(p, name) == 0 {
			p = name + " " + t
		}
		ret = append(ret, p)
	}
	return ret
}

var re = regexp.MustCompile(`({[^}]+})`)

func getRouteParamsFromRoute(mode, route string) []string {
	ret := []string{}
	if mode == routeMode {
		//todo: find a better way.
		res := re.FindAllStringSubmatch(route, -1)
		for _, r := range res {
			if len(r) > 0 {
				k := strings.TrimSpace(r[0])
				if len(k) > 0 {
					ret = append(ret, k)
				}
			}
		}
	}
	return ret
}

func getRouteParamNamesFromRoute(mode, route string) []string {
	ret := []string{}
	if mode == routeMode {
		//todo: find a better way.
		res := re.FindAllStringSubmatch(route, -1)
		for _, r := range res {
			if len(r) > 0 {
				k := strings.TrimSpace(r[0])
				if len(k) > 2 { // there is braces inside
					k = k[1 : len(k)-1]
					if strings.Index(k, ":") > -1 {
						j := strings.Split(k, ":")
						ret = append(ret, j[0])
					} else {
						ret = append(ret, k)
					}
				}
			}
		}
	}
	return ret
}

var routeMode = "route"
var stdMode = "std"
var reqBodyVarName = "reqBody"

func getPreferredMethod(annotations map[string]string) string {
	preferedMethod := "GET"
	if m, ok := annotations["metods"]; ok {
		methods := commaArgsToSlice(m)
		if len(methods) > 0 {
			preferedMethod = strings.ToUpper(methods[0])
		}
	}
	return preferedMethod
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

func commaArgsToSlice(s string) []string {
	ret := []string{}
	for _, l := range strings.Split(s, ",") {
		l = strings.TrimSpace(l)
		if l != "" {
			ret = append(ret, l)
		}
	}
	return ret
}
