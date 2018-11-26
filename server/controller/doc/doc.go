package doc

import (
	"github.com/ktpswjz/httpserver/document"
	"github.com/ktpswjz/om/server/controller"
)

type doc struct {
	controller.Controller

	Document document.Document
}
