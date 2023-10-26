package mconsts

// package name
const (
	ContextPkg     = "context"
	HertzPkg       = "github.com/cloudwego/hertz"
	HertzServerPkg = HertzPkg + "/pkg/app/server"
	HertzAppPkg    = HertzPkg + "/pkg/app"
	HertzUtils     = HertzPkg + "/pkg/common/utils"
	GinPkg         = "github.com/gin-gonic/gin"
)

// param name
const (
	GinCtx    = "gin.Context"
	HertzCtx  = "app.RequestContext"
	NormalCtx = "context.Context"
)
