package openwrt

import "github.com/Blaze757/routeflux/internal/backend/xray"

// NewXrayController returns an init.d backed Xray controller.
func NewXrayController() xray.InitdController {
	return xray.InitdController{ScriptPath: XrayServicePath()}
}
