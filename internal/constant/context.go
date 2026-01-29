package constant

type ctxTxKey struct{}

// DBTxKey 是用于在 Context 中存取事务对象的唯一标识
var DBTxKey = ctxTxKey{}

type ctxUserKey struct{}

// ContextUserKey 用于在 Context 中存取用户名
var ContextUserKey = ctxUserKey{}
