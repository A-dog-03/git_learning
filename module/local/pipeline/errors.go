package pipeline

import "../../../errors"

// genError 用于生成条目处理错误值。
func genError(errMsg string) error {
	return errors.NewCrawlerError(errors.ERROR_TYPE_PIPELINE,
		errMsg)
}

// genParameterError 用于生成条目处理参数错误值。
func genParameterError(errMsg string) error {
	return errors.NewCrawlerErrorBy(errors.ERROR_TYPE_PIPELINE,
		errors.NewIllegalParameterError(errMsg))
}
