package presenter

import (
	"github.com/kutsuzawa/slim-load-recorder/interactor"
	"github.com/line/line-bot-sdk-go/linebot"
)

// LineModel is line image message object referenced by
// https://developers.line.me/ja/reference/messaging-api/#anchor-66fcb7a0e9f4b2deb9eb2a0ace36ca6a4cc70e41
type LineModel struct {
	linebot.ImagemapMessage
}

type linePresenter struct {
	viewer ConsoleViewer
}

// NewLinePresenter init linePresenter
func NewLinePresenter(viewer ConsoleViewer) interactor.SlimLoadResponder {
	return &linePresenter{
		viewer: viewer,
	}
}

func (sp *linePresenter) Present(response []interactor.Response) {

}
