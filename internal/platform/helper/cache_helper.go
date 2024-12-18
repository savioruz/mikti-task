package helper

import (
	"fmt"
	"github.com/savioruz/mikti-task/internal/domain/model"
)

func (h *ContextHelper) BuildCacheKey(opts model.TodoQueryOptions) string {
	var cacheKey string

	if opts.IsAdmin {
		if opts.UserID != nil {
			cacheKey = fmt.Sprintf("todos:admin:user:%s:page:%d:size:%dsort:%sorder:%s", *opts.UserID, opts.Page, opts.Size, opts.Sort, opts.Order)
		} else {
			cacheKey = fmt.Sprintf("todos:admin:all:page:%d:size:%dsort:%sorder:%s", opts.Page, opts.Size, opts.Sort, opts.Order)
		}
	} else {
		cacheKey = fmt.Sprintf("todos:user:%s:page:%d:size:%dsort:%sorder:%s", *opts.UserID, opts.Page, opts.Size, opts.Sort, opts.Order)
	}

	if opts.Title != nil {
		cacheKey = fmt.Sprintf("%s:title:%s", cacheKey, *opts.Title)
	}

	return cacheKey
}
