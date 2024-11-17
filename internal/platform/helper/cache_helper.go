package helper

import (
	"fmt"
	"github.com/savioruz/mikti-task/internal/domain/model"
)

func (h *ContextHelper) BuildCacheKey(opts model.TodoQueryOptions) string {
	if opts.IsAdmin {
		if opts.UserID != nil {
			return fmt.Sprintf("todos:admin:user:%s:page:%d:size:%d", *opts.UserID, opts.Page, opts.Size)
		}
		return fmt.Sprintf("todos:admin:all:page:%d:size:%d", opts.Page, opts.Size)
	}
	return fmt.Sprintf("todos:user:%s:page:%d:size:%d", *opts.UserID, opts.Page, opts.Size)
}
