package userdetails

import (
	"encoding/json"
	"strings"

	utils "github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	"github.com/gin-gonic/gin"
)

type UserDetails struct {
	Id        string   `json:"id"`
	Email     string   `json:"email"`
	FirstName string   `json:"firstName"`
	LastName  string   `json:"lastName"`
	Roles     []string `json:"roles"`
}

func (instance *UserDetails) ToHeaders(c *gin.Context) {
	roles := strings.Join(utils.RemoveEmpties(instance.Roles), ",")
	fullName := strings.Join(
		utils.RemoveEmpties([]string{instance.FirstName, instance.LastName}),
		" ",
	)

	c.Header("X-User-Id", instance.Id)
	c.Header("X-User-Email", instance.Email)
	c.Header("X-User-Name", fullName)
	c.Header("X-User-Roles", roles)
}

func (instance *UserDetails) Marshal() ([]byte, error) {
	return json.Marshal(instance)
}
