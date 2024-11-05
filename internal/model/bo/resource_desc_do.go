package bo

import "backend-blog/internal/model/entity"

type ResourceDescBo struct {
	File         entity.File
	BlogImage    entity.BlogImage
	BlogVideo    entity.BlogVideo
	BlogMarkdown entity.BlogMarkdown
}
