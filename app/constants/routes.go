package constants

// route constants
const (
	RouteApi        = "/api"
	RouteLoggedOnly = "/loggedonly"
	RouteAdminOnly  = "/adminonly"

	RouteRegister             = "/register"
	RouteAuthorize            = "/authorize"
	RouteVerifyEmail          = "/verifyEmail/{username}/{activationToken}"
	RouteRequestResetPassword = "/requestResetPassword"
	RouteResetPassword        = "/resetPassword"
	RouteLogin                = "/login"
	RouteLogout               = "/logout"
	RouteKeepAlive            = "/keepalive"
	RoutePoll                 = "/poll"
	RouteJoinChat             = "/joinChat"
	RouteLeaveChat            = "/leaveChat"
	RouteGetChats             = "/getChats"
	RouteGetUsers             = "/getUsers"
	RouteGetLastMessages      = "/getLastMessages"
	RoutePostMessage          = "/postMessage"
	RouteCreateChat           = "/createChat"
	RouteRemoveChat           = "/removeChat"
)
