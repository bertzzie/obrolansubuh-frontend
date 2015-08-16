package controllers

import (
	"github.com/revel/revel"
	osc "obrolansubuh.com/modules/gorm/app/controllers"
)

func init() {
	revel.OnAppStart(osc.InitDB)

	revel.InterceptMethod((*osc.GormController).Begin, revel.BEFORE)
	revel.InterceptMethod((*osc.GormController).Commit, revel.AFTER)
	revel.InterceptMethod((*osc.GormController).RollBack, revel.FINALLY)
}
