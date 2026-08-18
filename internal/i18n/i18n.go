// Package i18n 多语言翻译支持：内嵌 locales/*.json 翻译表，
// 提供 Accept-Language 检测、翻译查询与 ctx 取值。
// 与 middleware.LocaleMiddleware + model.Error.Localize 形成
// "请求带 locale → 错误消息按语言翻译"的闭环。
package i18n

import (
	"context"
	"embed"
	"strings"
	"sync"

	"github.com/bytedance/sonic"
	"github.com/hcd233/aris-api-tmpl/internal/common/constant"
	"github.com/hcd233/aris-api-tmpl/internal/enum"
)

//go:embed locales/*.json
var localeFiles embed.FS

var (
	translations = make(map[enum.Locale]map[string]string)
	loadOnce     sync.Once
)

func loadLocales() {
	entries, err := localeFiles.ReadDir(constant.LocaleEmbedDir)
	if err != nil {
		return
	}
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), constant.LocaleFileExt) {
			continue
		}
		locale := enum.Locale(strings.TrimSuffix(entry.Name(), constant.LocaleFileExt))
		data, err := localeFiles.ReadFile(constant.LocaleEmbedDir + "/" + entry.Name())
		if err != nil {
			return
		}
		var m map[string]string
		if err := sonic.Unmarshal(data, &m); err != nil {
			return
		}
		translations[locale] = m
	}
}

func init() {
	loadOnce.Do(loadLocales)
}

// DetectLocale 从 Accept-Language 头检测语言环境，未识别时兜底为英语。
func DetectLocale(acceptLanguage string) enum.Locale {
	raw := strings.TrimSpace(acceptLanguage)
	if raw == "" {
		return enum.LocaleEN
	}
	parts := strings.SplitN(raw, ",", 2)
	primary := strings.TrimSpace(parts[0])
	if idx := strings.IndexAny(primary, constant.LocaleAcceptSeparator); idx > 0 {
		primary = primary[:idx]
	}
	primary = strings.ToLower(primary)
	switch primary {
	case constant.LocalePrimaryZH:
		return enum.LocaleZH
	default:
		return enum.LocaleEN
	}
}

// Translate 按 locale 翻译 key，缺失时回退英文表，再缺失时返回 fallback。
func Translate(locale enum.Locale, key, fallback string) string {
	if m, ok := translations[locale]; ok {
		if msg, ok := m[key]; ok {
			return msg
		}
	}
	if m, ok := translations[enum.LocaleEN]; ok {
		if msg, ok := m[key]; ok {
			return msg
		}
	}
	return fallback
}

// FromCtx 从 context 读取语言环境，未注入时兜底为英语。
func FromCtx(ctx context.Context) enum.Locale {
	if v, ok := ctx.Value(constant.CtxKeyLocale).(enum.Locale); ok {
		return v
	}
	return enum.LocaleEN
}
