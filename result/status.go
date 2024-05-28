package result

import "errors"

// 错误码规则:
// (1) 错误码需为 > 0 的数;
//
// (2) 错误码为 5 位数:
//              ----------------------------------------------------------
//                  第1位               2、3位                  4、5位
//              ----------------------------------------------------------
//                服务级错误码          模块级错误码	         具体错误码
//              ----------------------------------------------------------

var (
	OK = BuildSuccessResultWithMsg // 成功
	//WrongParameter   = BuildFailResultWithMsg("Wrong parameter")          // 参数错误
	//Err              = BuildFailResult()                                  // 错误
	MissingJWT       = BuildFailResultWithMsg("Missing or malformed JWT") // 无 JWT 认证
	SignatureInvalid = BuildFailResultWithMsg("Signature invalid")        // 签名无效
	NotFoundUser     = BuildFailResultWithMsg("Not found user")           // 无此用户
)
var (
	WrongParameter = errors.New("wrong parameter")
	Err            = errors.New("error")
	SaveVideoErr   = errors.New("save video error")
	SaveImgErr     = errors.New("save img error")
	DeleteFileErr  = errors.New("delete file error")
)
