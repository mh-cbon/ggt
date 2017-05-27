package capable

//go:generate ggt -mode route http-provider Get:RestGet
//go:generate ggt -mode rpc http-provider Get:RPCGet

//go:generate ggt -mode route http-provider Post:RestPost
//go:generate ggt -mode rpc http-provider Post:RPCPost

//go:generate ggt -mode route http-provider Route:RestRoute
//go:generate ggt -mode rpc http-provider Route:RPCRoute

//go:generate ggt -mode route http-provider URL:RestURL
//go:generate ggt -mode rpc http-provider URL:RPCURL

//go:generate ggt -mode route http-provider Req:RestReq
//go:generate ggt -mode rpc http-provider Req:RPReq

//go:generate ggt -mode route http-provider JSON:RestJSON
//go:generate ggt -mode rpc http-provider JSON:RPCJSON

//go:generate ggt -mode route http-provider Error:RestError
//go:generate ggt -mode rpc http-provider Error:RPCError
