package helper

import (
	"fmt"
	"github.com/savioruz/mikti-task/internal/domain/model"
)

func (h *ContextHelper) BuildCacheKey(opts model.TodoQueryOptions) string {
	var cacheKey string

	if opts.IsAdmin {
		if opts.UserID != nil {
			cacheKey = fmt.Sprintf("todos:admin:user:%s:page:%d:size:%d", *opts.UserID, opts.Page, opts.Size)
		} else {
			cacheKey = fmt.Sprintf("todos:admin:all:page:%d:size:%d", opts.Page, opts.Size)
		}
	} else {
		cacheKey = fmt.Sprintf("todos:user:%s:page:%d:size:%d", *opts.UserID, opts.Page, opts.Size)
	}

	if opts.Title != nil {
		cacheKey = fmt.Sprintf("%s:title:%s", cacheKey, *opts.Title)
	}

	return cacheKey
}
