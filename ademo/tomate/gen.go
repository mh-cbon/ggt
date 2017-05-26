package tomate

//go:generate ggt -c slicer *Tomate:Tomates
//go:generate ggt chaner Tomates:TomatesSync

//go:generate ggt -mode route http-provider Controller:RestController
//go:generate ggt -mode rpc http-provider Controller:RpcController

//go:generate ggt -mode route http-consumer Controller:RestClient
//go:generate ggt -mode rpc http-consumer Controller:RPCClient
