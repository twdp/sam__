package controllers

import (
	"github.com/pkg/errors"
	"tianwei.pro/business/controller"
)

type TestController struct {
	controller.RestfulController
}


// @router /a [get]
func (t *TestController) Test() {
	panic(errors.New("----"))
}