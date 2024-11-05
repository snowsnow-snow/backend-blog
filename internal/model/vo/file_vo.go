package vo

type FileVo struct {
	ImageVos    []*ImageVo    `json:"imageVos"`
	VideoVos    []*VideoVo    `json:"videoVos"`
	MarkdownVos []*MarkdownVo `json:"markdownVos"`
}
