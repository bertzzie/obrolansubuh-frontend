package controllers

import (
	"github.com/revel/revel"
	osc "obrolansubuh.com/modules/gorm/app/controllers"
)

type App struct {
	osc.GormController
}

func (c App) Index() revel.Result {
	return c.Render()
}
