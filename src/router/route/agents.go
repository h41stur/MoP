package route

import (
	"MoP/src/controllers"
	"net/http"
)

// carrega a URI, o método e cada função correspondente do pacote controllers
var agentsRoute = []Route{
	{
		URI:      "/agents",
		Method:   http.MethodPost,
		Function: controllers.CreateAgent,
	},
	{
		URI:      "/agents/{agentId}/com",
		Method:   http.MethodGet,
		Function: controllers.GetCommand,
	},
	{
		URI:      "/agents/{agentId}/com",
		Method:   http.MethodPost,
		Function: controllers.PostCommand,
	},
	{
		URI:      "/agents/{agentId}/fl",
		Method:   http.MethodGet,
		Function: controllers.GetFile,
	},
	{
		URI:      "/agents/{agentId}/fl",
		Method:   http.MethodPost,
		Function: controllers.PostFile,
	},
	{
		URI:      "/agents/{agentId}/fl/{fileId}",
		Method:   http.MethodDelete,
		Function: controllers.DeleteFile,
	},
	{
		URI:      "/agents/{agentId}/alive",
		Method:   http.MethodGet,
		Function: controllers.AgentAlive,
	},
}
