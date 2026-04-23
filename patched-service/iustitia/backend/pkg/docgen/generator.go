package docgen

import (
	"strings"
)

type DocumentContext struct {
	CaseID    string
	Defendant string
	Verdict   string
	Details   string
}

// PATCH - Vuln 3 (Server-Side Template Injection + arbitrary file read).
// Было: text/template.Parse(userTemplate) на пользовательский ввод, плюс
// метод (*DocumentContext).ReadClassified(path) был доступен из темплейта
// как {{.ReadClassified "secrets/classified.txt"}} - давал чтение любого
// файла, включая secrets/internal.txt (SECRET_MARKER_S / SECRET_MARKER_F).
// Стало: ReadClassified удалён. Шаблон больше не парсится Go text/template.
// Вместо этого - узкая подстановка только whitelisted-плейсхолдеров
// {{case_id}} / {{defendant}} / {{verdict}} / {{details}}. Любой другой
// {{...}} остаётся как есть (плейн-текст) - ни FuncMap, ни method
// invocation, ни pipeline'ов.
func Generate(userTemplate string, ctx *DocumentContext) (string, error) {
	replacements := []struct {
		placeholder string
		value       string
	}{
		{"{{case_id}}", ctx.CaseID},
		{"{{defendant}}", ctx.Defendant},
		{"{{verdict}}", ctx.Verdict},
		{"{{details}}", ctx.Details},
	}
	out := userTemplate
	for _, r := range replacements {
		out = strings.ReplaceAll(out, r.placeholder, r.value)
	}
	return out, nil
}
