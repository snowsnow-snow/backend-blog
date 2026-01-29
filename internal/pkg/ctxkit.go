package pkg

type ctxTraceKey struct{}
type ctxIPKey struct{}

var TraceKey = ctxTraceKey{}
var IPKey = ctxIPKey{}
