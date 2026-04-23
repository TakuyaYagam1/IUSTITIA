package usecase

import (
	"fmt"
	"strconv"

	"github.com/TakuyaYagam1/iustitia/pkg/docgen"
)

// VerdictRenderer адаптирует pkg/docgen под интерфейс docgenRenderer,
// объявленный в usecase/case.go.
type VerdictRenderer struct{}

func NewVerdictRenderer() *VerdictRenderer {
	return &VerdictRenderer{}
}

// RenderVerdict рендерит verdict-шаблон. caseSeq подставляется в CaseID
// как "№N/Δ" - внешний формат, который видит пользователь в UI.
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
