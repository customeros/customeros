package user_details

import (
	"encoding/json"
	"strings"

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
	roles := strings.Join(filterNonEmpty(instance.Roles), ",")
	fullName := strings.Join(
		filterNonEmpty([]string{instance.FirstName, instance.LastName}),
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

func filterNonEmpty(arr []string) []string {
	var res []string
	for _, str := range arr {
		if str != "" {
			res = append(res, str)
		}
	}
	return res
}
