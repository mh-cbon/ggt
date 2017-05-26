package login

//go:generate ggt -c slicer *HashedUser:HashedUsers
//go:generate ggt chaner HashedUsers:HashedUsersSync

//go:generate ggt -mode route http-provider Controller:RestController
//go:generate ggt -mode route http-consumer Controller:RestClient

//go:generate ggt -mode rpc http-provider Controller:RPCController
//go:generate ggt -mode rpc http-consumer Controller:RPCClient
