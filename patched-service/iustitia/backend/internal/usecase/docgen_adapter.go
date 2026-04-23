package usecase

import (
	"fmt"
	"strconv"

	"github.com/TakuyaYagam1/iustitia/pkg/docgen"
)

// VerdictRenderer адаптирует pkg/docgen под интерфейс docgenRenderer,
// объявленный в usecase/case.go. Делает Case usecase не-знающим о
// конкретной реализации шаблонного движка - это важно для F3 (SSTI):
// pkg/docgen держит whitelist-плейсхолдеры и canonical registry, а
// usecase оперирует только input-данными приговора.
type VerdictRenderer struct{}

func NewVerdictRenderer() *VerdictRenderer {
	return &VerdictRenderer{}
}

// RenderVerdict рендерит canonical verdict-шаблон из docgen.ResolveTemplate
// с whitelist-плейсхолдерами {{case_id}} / {{defendant}} / {{verdict}} /
// {{details}}. caseSeq подставляется в CaseID как "№N/Δ" - внешний формат,
// который видит пользователь в UI.
func (r *VerdictRenderer) RenderVerdict(caseSeq int64, defendant, crime, details, verdict string) (string, error) {
	body, ok := docgen.ResolveTemplate(docgen.TemplateVerdict)
	if !ok {
		return "", fmt.Errorf("docgen: canonical verdict template not found")
	}
	ctx := &docgen.DocumentContext{
		CaseID:    "№" + strconv.FormatInt(caseSeq, 10) + "/Δ",
		Defendant: defendant,
		Verdict:   verdict,
		Details:   details,
	}
	return docgen.Generate(body, ctx)
}
