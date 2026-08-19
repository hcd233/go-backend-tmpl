package constant

const (
	// ProjectName 项目名，用作 JWT Issuer/Audience、日志标识等。
	ProjectName = "aris-api-tmpl"

	// ErrorModelTemplate model.Error 的 Error() 输出模板。
	ErrorModelTemplate = "code: %d, message: %s"
)

const (
	// LocaleEmbedDir i18n 翻译文件内嵌目录。
	LocaleEmbedDir = "locales"

	// LocaleFileExt 翻译文件扩展名。
	LocaleFileExt = ".json"

	// LocaleAcceptSeparator Accept-Language 主语言分隔符。
	LocaleAcceptSeparator = "-"

	// LocalePrimaryZH 中文主语言标识。
	LocalePrimaryZH = "zh"

	// LocalePrimaryJA 日文主语言标识。
	LocalePrimaryJA = "ja"

	// HTTPHeaderAcceptLanguage Accept-Language 请求头名。
	HTTPHeaderAcceptLanguage = "Accept-Language"
)
