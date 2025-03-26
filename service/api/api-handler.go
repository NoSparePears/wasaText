package api

import (
	"net/http"
)

// Handler returns an instance of httprouter.Router that handle APIs registered here
func (rt *_router) Handler() http.Handler {
	// Special routes
	rt.router.GET("/liveness", rt.liveness)

	rt.router.POST("/session", rt.wrap(rt.doLogin, false))

	rt.router.GET("/profiles", rt.wrap(rt.searchUsers, true))

	rt.router.PUT("/profiles/:userID/username", rt.wrap(rt.setMyUserName, true))

	rt.router.PUT("/profiles/:userID/photo", rt.wrap(rt.setMyPhoto, true))

	rt.router.GET("/profiles/:userID/conversations", rt.wrap(rt.getMyConversations, true))

	rt.router.PUT("/profiles/:userID/conversations/:destID", rt.wrap(rt.createConversation, true))

	rt.router.GET("/profiles/:userID/conversations/:destID", rt.wrap(rt.getConversation, true))

	rt.router.DELETE("/profiles/:userID/conversations/:destID", rt.wrap(rt.deleteConversation, true))

	rt.router.GET("/profiles/:userID/conversations/:destID/messages", rt.wrap(rt.getMyMessages, true))

	rt.router.POST("/profiles/:userID/conversations/:destID/messages", rt.wrap(rt.sendMessage, true))

	rt.router.DELETE("/profiles/:userID/conversations/:destID/messages/:msgID", rt.wrap(rt.deleteMessage, true))

	rt.router.POST("/profiles/:userID/conversations/:destID/messages/:msgID", rt.wrap(rt.forwardMessage, true))

	rt.router.PUT("/profiles/:userID/conversations/:destID/messages/:msgID/comments", rt.wrap(rt.commentMessage, true))

	rt.router.DELETE("/profiles/:userID/conversations/:destID/messages/:msgID/comments/:commID", rt.wrap(rt.uncommmentMessage, true))

	rt.router.POST("/profiles/:userID/groups", rt.wrap(rt.createGroup, true))

	rt.router.GET("/profiles/:userID/groups", rt.wrap(rt.getMyGroups, true))

	rt.router.GET("/profiles/:userID/groups/:groupID", rt.wrap(rt.getGroup, true))

	rt.router.GET("/profiles/:userID/groups/:groupID/messages", rt.wrap(rt.getGroupMessages, true))

	rt.router.POST("/profiles/:userID/groups/:groupID/messages", rt.wrap(rt.sendGroupMessage, true))

	rt.router.DELETE("/profiles/:userID/groups/:groupID/messages/:msgID", rt.wrap(rt.deleteGroupMessage, true))

	rt.router.PUT("/profiles/:userID/groups/:groupID/members", rt.wrap(rt.addToGroup, true))

	rt.router.GET("/profiles/:userID/groups/:groupID/members", rt.wrap(rt.GetGroupMembers, true))

	rt.router.DELETE("/profiles/:userID/groups/:groupID", rt.wrap(rt.leaveGroup, true))

	rt.router.PUT("/profiles/:userID/groups/:groupID/g_name", rt.wrap(rt.setGroupName, true))

	rt.router.PUT("/profiles/:userID/groups/:groupID/g_photo", rt.wrap(rt.setGroupPhoto, true))

	return rt.router
}
