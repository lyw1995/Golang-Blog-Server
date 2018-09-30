package common

import "fmt"

//
//const (
//	Ok_Code           = 200 //get put ok
//	PostOk_Code       = 201
//	DelOk_Code        = 204
//	RequsetErr_Code   = 400 //参数错误等
//	UnAuthorized_Code = 401
//	Forbidden_Code    = 403
//	NotFound_Code     = 404
//	InternalErr_Code  = 500
//)

type Err struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func (err Err) Error() string {
	return fmt.Sprintf("ErrCode:%d , ErrMsg:%s", err.Code, err.Msg)
}

//返回结果中的错误码表示了用户调用API 的结果。其中，code 为公共错误码，其适用于所有模块的 API 接口。
//若 code 为 0，表示调用成功，否则，表示调用失败。当调用失败后，用户可以根据下表确定错误原因并采取相应措施。
var (
	ErrClientParams   = Err{Code: 4000, Msg: "缺少必要参数，或者参数值/路径格式不正确。"}
	ErrAuthorized     = Err{Code: 4100, Msg: "签名鉴权失败，请参考文档中鉴权部分。"}
	ErrAccoutDeny     = Err{Code: 4200, Msg: "帐号被封禁，或者不在接口针对的用户范围内等。"}
	ErrNotFound       = Err{Code: 4300, Msg: "资源不存在，或者访问了其他用户的资源。"}
	ErrMethodNotAllow = Err{Code: 4400, Msg: "协议不支持，请参考文档说明。"}
	ErrSignParams     = Err{Code: 4500, Msg: "签名错误，请参考文档说明。"}
	ErrInternal       = Err{Code: 5000, Msg: "服务器内部出现错误，请稍后重试。"}
	ErrApiClose       = Err{Code: 6200, Msg: "当前接口处于停服维护状态，请稍后重试。"}
)

//模块错误码
var (
	ErrUserExist          = Err{Code: 10000, Msg: "用户已存在,请修改后重试。"}
	ErrUserLogin          = Err{Code: 10010, Msg: "用户名或密码错误,请检查后重试。"}
	ErrUserNameUnExist    = Err{Code: 10020, Msg: "用户名不存在,请检查后重试。"}
	ErrUserNoExist        = Err{Code: 10030, Msg: "用户不存在,请检查后重试。"}
	ErrUserNameFormat     = Err{Code: 10040, Msg: fmt.Sprintf("用户名需字母开头长度(%d~%d)位字母数字_。", LenUserNameMin, LenUserNameMax)}
	ErrUserPwdFormat      = Err{Code: 10050, Msg: fmt.Sprintf("密码需字母开头长度(%d~%d)位字母数字_。", LenUserNameMin, LenPasswordMax)}
	ErrUserDescLen        = Err{Code: 10060, Msg: fmt.Sprintf("描述长度不能超过%d位,请改正后重试。", LenDesc)}
	ErrUserAddrLen        = Err{Code: 10070, Msg: fmt.Sprintf("地址长度不能超过%d位,请改正后重试。", LenAddr)}
	ErrUserEmailFormat    = Err{Code: 10080, Msg: "邮箱格式不正确,请检查后重试。"}
	ErrUserNickNameFormat = Err{Code: 10090, Msg: fmt.Sprintf("昵称长度在(%d~%d)之间,请改正后重试。", LenUserNameMin, LenUserNameMax)}
	ErrUpdateParams       = Err{Code: 10100, Msg: "修改用户信息,参数必填其一。"}
	ErrUserLinksNoExist   = Err{Code: 10110, Msg: "用户友链数据不存在,请检查后重试。"}
	ErrCategoryNoExist    = Err{Code: 10120, Msg: "该分类不存在,请检查后重试。"}
	ErrCategoryExist      = Err{Code: 10130, Msg: "该分类已存在,请修改后重试。"}
	ErrArticleNoExist     = Err{Code: 10140, Msg: "该文章不存在,请修改后重试。"}
	ErrUploadLenNotAllow  = Err{Code: 10150, Msg: "图片上传个数不允许,请修改后重试。"}
	ErrUploadExtNotAllow  = Err{Code: 10160, Msg: "仅支持jpg,jpeg,png格式图片,请修改后重试。"}
	ErrCollectSource            = Err{Code: 10170, Msg: "文章采集失败,请检查采集Url。"}
)
