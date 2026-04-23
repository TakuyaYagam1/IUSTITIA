package request

import (
	"github.com/TakuyaYagam1/iustitia/internal/domain"
	"github.com/TakuyaYagam1/iustitia/internal/openapi"
	"github.com/TakuyaYagam1/iustitia/internal/usecase"
)

func VerdictRequestToParams(req *openapi.VerdictRequest) domain.Verdict {
	return domain.Verdict(req.Verdict)
}

func VerdictRequestFull(req *openapi.VerdictRequest) (domain.Verdict, *string, string) {
	return domain.Verdict(req.Verdict), req.Sentence, req.Reasoning
}

func CaseListParams(limit, offset *int) (int, int) {
	l := 50
	if limit != nil {
		l = *limit
	}
	o := 0
	if offset != nil {
		o = *offset
	}
	return l, o
}

func SearchRequestToParams(req *openapi.SearchRequest) usecase.SearchRequest {
	out := usecase.SearchRequest{Q: req.Q}
	if req.OrderBy != nil {
		out.OrderBy = *req.OrderBy
	}
	if req.Direction != nil {
		out.Direction = string(*req.Direction)
	}
	if req.Limit != nil {
		out.Limit = *req.Limit
	}
	if req.Offset != nil {
		out.Offset = *req.Offset
	}
	return out
}
