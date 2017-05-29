// Package consumer generates typed http client.
package consumer

import (
	"fmt"
	"log"
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

generates http client implementation of given type.

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
	// srcNameFq := srcName
	if todo.FromPkgPath != todo.ToPkgPath && !astutil.IsBasic(todo.FromTypeName) {
		// srcNameFq = fmt.Sprintf("%v.%v", filepath.Base(todo.FromPkgPath), srcConcrete)
		if srcIsPointer {
			// srcNameFq = "*" + srcNameFq
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
		client *http.Client
	}
			`, dstConcrete)

	} else {
		fmt.Fprintf(dest, `
	type %v struct{
		client *http.Client
	}
			`, dstConcrete)
	}

	// Make the constructor
	// should param *http.Client be an interface ?
	fmt.Fprintf(dest, `// New%v constructs an http-clienter of %v
`, dstConcrete, srcName)

	if mode == routeMode {
		fmt.Fprintf(dest, `func New%v(router *mux.Router, client *http.Client) *%v {
	if client == nil {
		client = http.DefaultClient
	}
	ret := &%v{
		router: router,
		client: client,
	}
  return ret
}
`, dstConcrete, dstConcrete, dstConcrete)

	} else {
		fmt.Fprintf(dest, `func New%v(client *http.Client) *%v {
	if client == nil {
		client = http.DefaultClient
	}
	ret := &%v{
		client: client,
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
		fileOut.AddImport("context", "")
		fileOut.AddImport("github.com/gorilla/mux", "")

		// cheat.
		for _, x := range []string{
			"bytes.MinRead",
			"fmt.Println",
			"url.PathEscape",
			"strings.ToUpper",
			"context.Canceled",
			"mux.Vars",
			"io.Copy",
			"http.StatusOK",
		} {
			fmt.Fprintf(dest, `var xx%v = %v
				`, utils.Hash(fileOut.Path+x), x)
		}
	} else {
		fileOut.AddImport("net/http", "")

		// cheat.
		for _, x := range []string{"http.StatusOK"} {
			fmt.Fprintf(dest, `var xx%v = %v
				`, utils.Hash(fileOut.Path+x), x)
		}
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
		lParams := astutil.MethodInputs(m)
		// paramNames := astutil.MethodParamNames(m)
		lParamNames := astutil.MethodInputNames(m)
		paramTypes := astutil.MethodInputTypes(m)

		retTypes := astutil.MethodReturnTypes(m)
		retVars := astutil.MethodReturnNamesNormalized(m)
		sRetVarNames := strings.TrimSpace(strings.Join(retVars, ", "))

		defVal := astutil.TypesToDefVal(retTypes)
		if i := astutil.TypesIndex(retTypes, "error"); i > -1 {
			defVal[i] = fmt.Sprintf(`errors.New(%q)`, "todo")
			fileOut.AddImport("errors", "")
		}
		defValRet := strings.Join(defVal, ", ")

		handleErr := func(errVarName string, subjects ...string) string {
			return fmt.Sprintf(`if %v!= nil {
		return %v
		}
		`, errVarName, defValRet)
		}

		importIDs := astutil.GetSignatureImportIdentifiers(m)
		for _, i := range importIDs {
			path := astutil.GetImportPath(pkg, i)
			if path == "" {
				log.Printf("package path not found. id:%q path:%q\n", i, path)
			} else {
				fileOut.AddImport(path, i)
			}
		}

		if mode == "rpc" {

			fileOut.AddImport("errors", "")
			fileOut.AddImport("bytes", "")

			fmt.Fprintf(dest, `// %v constructs a request to %v
		`, methodName, methodName)

			body := `var reqBody bytes.Buffer
		`
			if len(paramTypes) > 0 {
				fileOut.AddImport("encoding/json", "json")

				body += fmt.Sprintf(`
				{
					input := struct{
						%v
					}{
						%v
					}
					encErr := json.NewEncoder(&reqBody).Encode(&input)
					%v
				}
					`, mapParamsToStruct(paramTypes, false),
					mapParamsToStructValues(lParamNames),
					handleErr("encErr", "req", "json", "encode"))
			}

			// - create the request object
			// preferedMethod := getPreferredMethod(annotations)
			body += fmt.Sprintf(`finalURL := %q
			req, reqErr := http.NewRequest(%q, finalURL, &reqBody)
			%v
			`, "/"+methodName, "POST", handleErr("reqErr"))

			if len(retVars) == 0 {
				body += fmt.Sprintf(`_, resErr := t.client.Do(req)
			%v
			return
			`, handleErr("resErr"))

			} else {

				outputHandling := fmt.Sprintf(`
				output := struct{
					%v
				}{}
				{
					decErr := json.NewDecoder(res.Body).Decode(&output)
					%v
				}
					`, mapParamsToStruct(retTypes, false),
					handleErr("decErr", "res", "json", "decode"))

				retValues := ""
				for _, r := range mapParamsToStructValueNames(retVars, "output") {
					retValues += r + ", "
				}
				if retValues != "" {
					retValues = retValues[:len(retValues)-2]
				}

				body += fmt.Sprintf(`res, resErr := t.client.Do(req)
				%v
				%v
				return %v
				`, handleErr("resErr"), outputHandling, retValues)

			}

			sRetTypes := strings.TrimSpace(strings.Join(retTypes, ", "))
			fmt.Fprintf(dest, `func(t %v) %v(%v) (%v) {
			%v
		}
		`, destName, methodName, params, sRetTypes, body)

		} else if route, ok := annotations["route"]; ok && mode == routeMode {

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
				paramType := paramTypes[paramIndex]

				if !managedParamNames.Contains(paramName) {

					if paramName == "getValues" && paramType == mapStringString {

						getParams += fmt.Sprintf(`
							for k, v := range %v {
								reqURL.Query().Add(k,v)
							}
							`, paramName)

					} else if paramName == "getValues" && paramType == mapStringSliceString {

						getParams += fmt.Sprintf(`
							for k, vv := range %v {
								for _, v := range vv {
									reqURL.Query().Add(k,v)
								}
							}
							`, paramName)

					} else if strings.HasPrefix(paramName, "get") || strings.HasPrefix(paramName, "url") || strings.HasPrefix(paramName, "req") {
						k := strings.ToLower(paramName[3:])

						if astutil.IsArrayType(paramType) {
							getParams += fmt.Sprintf(`
								for _, item%v := range %v {
									var xx%v string
									%v
									reqURL.Query().Add(%q, xx%v)
								}
								`, paramName, paramName, paramName,
								convertToStr("item"+paramName, "xx"+paramName, astutil.GetUnslicedType(paramType), handleErr, destName, methodName, "get"),
								k, paramName)

						} else if astutil.IsAPointedType(paramType) {
							getParams += fmt.Sprintf(`
								if %v != nil {
									var xx%v string
									%v
									reqURL.Query().Add(%q, xx%v)
								}
								`, paramName, paramName,
								convertToStr(paramName, "xx"+paramName, paramType, handleErr, destName, methodName, "get"),
								k, paramName)

						} else {
							getParams += fmt.Sprintf(`var xx%v string
								%v
								reqURL.Query().Add(%q, xx%v)
								`, paramName, convertToStr(paramName, "xx"+paramName, paramType, handleErr, destName, methodName, "get"), k, paramName)
						}

					} else if strings.HasPrefix(paramName, "post") {
						k := strings.ToLower(paramName[4:])

						if astutil.IsAPointedType(paramType) {
							postParams += fmt.Sprintf(`form.Add(%q, *%v)
							`, k, paramName)
						} else {
							postParams += fmt.Sprintf(`form.Add(%q, %v)
							`, k, paramName)
						}
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
				body += fmt.Sprintf("_, resErr := t.client.Do(req)\n")
				body += handleErr("resErr")

			} else {

				fileOut.AddImport("encoding/json", "json")

				body += fmt.Sprintf(`
				{
					res, resErr := t.client.Do(req)
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

func convertStrTo(fromStrVarName, toVarName, toType string, errHandler func(string, ...string) string, subjects ...string) string {
	if astutil.GetUnpointedType(toType) == "string" {
		if astutil.IsAPointedType(toType) {
			return fmt.Sprintf("%v = &%v", toVarName, fromStrVarName)
		}
		return fmt.Sprintf("%v = %v", toVarName, fromStrVarName)

	} else if astutil.GetUnpointedType(toType) == "bool" {
		if astutil.IsAPointedType(toType) {
			return fmt.Sprintf(`{
				xxTmp := %v=="true"
					%v = &xxTmp
			}
	`, fromStrVarName, toVarName)
		}
		return fmt.Sprintf(`%v = %v=="true"
`, fromStrVarName, toVarName)

	} else if astutil.GetUnpointedType(toType) == "int" {
		if astutil.IsAPointedType(toType) {
			return fmt.Sprintf(`%v, err := strconv.Atoi(*%v)
		%v
	`, toVarName, fromStrVarName, errHandler("err", subjects...))
		}
		return fmt.Sprintf(`%v, err := strconv.Atoi(%v)
		%v
	`, toVarName, fromStrVarName, errHandler("err", subjects...))
	}
	return fmt.Sprintf(`// conversion not handled string to %v
	`, toType)
}

func convertToStr(fromVarName, toStrVarName, fromType string, errHandler func(string, ...string) string, subjects ...string) string {
	if astutil.GetUnpointedType(fromType) == "string" {
		if astutil.IsAPointedType(fromType) {
			return fmt.Sprintf("%v = *%v", toStrVarName, fromVarName)
		}
		return fmt.Sprintf("%v = %v", toStrVarName, fromVarName)

	} else if astutil.GetUnpointedType(fromType) == "bool" {
		if astutil.IsAPointedType(fromType) {
			return fmt.Sprintf(`%v = "false"
			if %v != nil && *%v {
				%v = "true"
			}
	`, toStrVarName, fromVarName, fromVarName, toStrVarName)
		}
		return fmt.Sprintf(`%v = "false"
			if %v {
				%v = "true"
			}
	`, toStrVarName, fromVarName, toStrVarName)

	} else if astutil.GetUnpointedType(fromType) == "int" {
		return fmt.Sprintf(`%v = fmt.Sprintf("%%v", %v)
	`, toStrVarName, fromVarName)

	}
	return fmt.Sprintf(`// conversion not handled  %v to string
	`, fromType)
}

var headersValueName = "headers"
var ctxCtx = "context.Context"
var httpRequest = "*http.Request"
var httpResponse = "http.ResponseWriter"
var httpCookie = "http.Cookie"
var ggtFile = "ggt.File"
var fileValues = "fileValues"
var cookieValues = "cookieValues"
var ioReader = "io.Reader"
var mapStringString = "map[string]string"
var mapStringSliceString = "map[string][]string"
var rpcMode = "rpc"
var routeMode = "route"

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

func mapParamsToStructValueNames(params []string, fromvarname string) []string {
	var ret []string
	for i := range params {
		x := fmt.Sprintf("Arg%v", i)
		if fromvarname != "" {
			x = fmt.Sprintf("%v.%v", fromvarname, x)
		}
		ret = append(ret, x)
	}
	return ret
}

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
