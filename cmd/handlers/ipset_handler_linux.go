package handlers

import (
	"ipsetsv/common"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/nadoo/ipset"
)

type IPSet struct {
	Token   string
	Timeout uint32
}

func (ctrl IPSet) SyncIPSet(c echo.Context) error {
	var req common.IPSetReq
	if err := c.Bind(&req); err != nil {
		c.NoContent(http.StatusBadRequest)
		return err
	}
	if req.Token != ctrl.Token {
		c.NoContent(http.StatusUnauthorized)
		return nil
	}
	ipset.Create(req.SetName, ipset.OptTimeout(ctrl.Timeout))
	for _, ip := range req.IPList {
		ipset.Add(req.SetName, ip, ipset.OptTimeout(req.Timeout))
	}
	return nil
}
