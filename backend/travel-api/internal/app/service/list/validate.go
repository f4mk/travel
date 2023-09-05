package list

import "github.com/f4mk/travel/backend/travel-api/internal/pkg/web"

func (nl NewList) Validate() error {
	return web.Check(nl)
}

func (ni NewItem) Validate() error {
	return web.Check(ni)
}

func (np NewPoint) Validate() error {
	return web.Check(np)
}

func (ul UpdateList) Validate() error {
	return web.Check(ul)
}

func (ui UpdateItem) Validate() error {
	return web.Check(ui)
}

func (up UpdatePoint) Validate() error {
	return web.Check(up)
}
