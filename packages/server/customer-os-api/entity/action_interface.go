package entity

const (
	ActionName_PageView = "PageViewAction"
)

type Action interface {
	Action()
	ActionName() string
}

type ActionEntities []Action
