// Copyright © 2023 OpenIM. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package api

import (
	"github.com/gin-gonic/gin"

	"github.com/OpenIMSDK/protocol/constant"
	"github.com/OpenIMSDK/protocol/msggateway"
	user "github.com/OpenIMSDK/protocol/user"
	"github.com/OpenIMSDK/tools/a2r"
	"github.com/OpenIMSDK/tools/apiresp"
	"github.com/OpenIMSDK/tools/errs"
	"github.com/OpenIMSDK/tools/log"

	"github.com/openimsdk/open-im-server/v3/pkg/common/config"
	"github.com/openimsdk/open-im-server/v3/pkg/rpcclient"
)

type UserApi rpcclient.User

func NewUserApi(client rpcclient.User) UserApi {
	return UserApi(client)
}

// UserRegisterReq
// @Summary User registration
// @Description User registration
// @Tags User
// @ID UserRegister
// @Accept json
// @Param OperationId header string true "Operation Id"
// @Param UserInfo body user.UserRegisterReq true "Secret is the Openim key. For details, see the server Config.yaml Secret field."
// @Produce json
// @Success 200 {object} user.UserRegisterResp
// @Failure 500 {object} error "ERRCODE is 500 generally an internal error of the server"
// @Failure 400 {object} error "Errcode is 400, which is generally a parameter input error.""
// @Router /user/user_register [post]
func (u *UserApi) UserRegister(c *gin.Context) {
	a2r.Call(user.UserClient.UserRegister, u.Client, c)
}

// @Summary Modify user information
// @Description Modify user information Userid Faceurl, etc.
// @Tags User
// @ID UpdateUserInfo
// @Accept json
// @Param OperationId header string true "Operation Id"
// @Param req body user.UpdateUserInfoReq true "Request"
// @Produce json
// @Success 200 {object} user.UpdateUserInfoResp
// @Failure 500 {object} error "ERRCODE is 500 generally an internal error of the server"
// @Failure 400 {object} error "Errcode is 400, which is generally a parameter input error."
// @Security	ApiKeyAuth
// @Router /user/update_user_info [post]
func (u *UserApi) UpdateUserInfo(c *gin.Context) {
	a2r.Call(user.UserClient.UpdateUserInfo, u.Client, c)
}

// @Summary Set the overall disturbance
// @Description Set the overall disturbance
// @Tags User
// @ID SetGlobalRecvMessageOpt
// @Accept json
// @Param OperationId header string true "Operation Id"
// @Param req body user.SetGlobalRecvMessageOptReq true "GlobalRecvmsGopt is the global disturbance setting 0 to turn off 1 to open"
// @Produce json
// @Success 200 {object} user.SetGlobalRecvMessageOptResp
// @Failure 500 {object} error "ERRCODE is 500 generally an internal error of the server"
// @Failure 400 {object} error "Errcode is 400, which is generally a parameter input error."
// @Security	ApiKeyAuth
// @Router /user/set_global_msg_recv_opt [post]
func (u *UserApi) SetGlobalRecvMessageOpt(c *gin.Context) {
	a2r.Call(user.UserClient.SetGlobalRecvMessageOpt, u.Client, c)
}

//// @Success 0 {object} user.GetDesignateUsersResp{usersInfo=[]sdkws.UserInfo}

// @Summary Get user information
// @Description Obtain user information in batches according to the user list
// @Tags User
// @ID GetUsersInfo
// @Accept json
// @Param OperationId header string true "Operation Id"
// @Param req body user.GetDesignateUsersReq true "Request"
// @Produce json
// @Success 200 {object} user.GetDesignateUsersResp
// @Failure 500 {object} error "Errcode is 500 一For the internal error of the server"
// @Failure 400 {object} error "errcode is 400 一Input errors in the parameter, token is not brought up"
// @Security	ApiKeyAuth
// @Router /user/get_users_info [post]
func (u *UserApi) GetUsersPublicInfo(c *gin.Context) {
	a2r.Call(user.UserClient.GetDesignateUsers, u.Client, c)
}

func (u *UserApi) GetAllUsersID(c *gin.Context) {
	a2r.Call(user.UserClient.GetAllUserID, u.Client, c)
}

func (u *UserApi) AccountCheck(c *gin.Context) {
	a2r.Call(user.UserClient.AccountCheck, u.Client, c)
}

func (u *UserApi) GetUsers(c *gin.Context) {
	a2r.Call(user.UserClient.GetPaginationUsers, u.Client, c)
}

// @Summary Get user online status
// @Description Get user online status
// @Tags User
// @ID GetUsersOnlineStatus
// @Accept json
// @Param OperationId header string true "Operation Id"
// @Param req body msggateway.GetUsersOnlineStatusReq true "Request"
// @Produce json
// @Success 200 {object} apiresp.ApiResponse{data=[]msggateway.GetUsersOnlineStatusResp_SuccessResult}
// @Failure 500 {object} error "ERRCODE is 500 generally an internal error of the server"
// @Failure 400 {object} error "Errcode is 400, which is generally a parameter input error."
// @Security	ApiKeyAuth
// @Router /user/get_users_online_status [post]
func (u *UserApi) GetUsersOnlineStatus(c *gin.Context) {
	var req msggateway.GetUsersOnlineStatusReq
	if err := c.BindJSON(&req); err != nil {
		apiresp.GinError(c, errs.ErrArgs.WithDetail(err.Error()).Wrap())
		return
	}
	conns, err := u.Discov.GetConns(c, config.Config.RpcRegisterName.OpenImMessageGatewayName)
	if err != nil {
		apiresp.GinError(c, err)
		return
	}

	var wsResult []*msggateway.GetUsersOnlineStatusResp_SuccessResult
	var respResult []*msggateway.GetUsersOnlineStatusResp_SuccessResult
	flag := false

	// Online push message
	for _, v := range conns {
		msgClient := msggateway.NewMsgGatewayClient(v)
		reply, err := msgClient.GetUsersOnlineStatus(c, &req)
		if err != nil {
			log.ZWarn(c, "GetUsersOnlineStatus rpc err", err)

			parseError := apiresp.ParseError(err)
			if parseError.ErrCode == errs.NoPermissionError {
				apiresp.GinError(c, err)
				return
			}
		} else {
			wsResult = append(wsResult, reply.SuccessResult...)
		}
	}
	// Traversing the userIDs in the api request body
	for _, v1 := range req.UserIDs {
		flag = false
		res := new(msggateway.GetUsersOnlineStatusResp_SuccessResult)
		// Iterate through the online results fetched from various gateways
		for _, v2 := range wsResult {
			// If matches the above description on the line, and vice versa
			if v2.UserID == v1 {
				flag = true
				res.UserID = v1
				res.Status = constant.OnlineStatus
				res.DetailPlatformStatus = append(res.DetailPlatformStatus, v2.DetailPlatformStatus...)
				break
			}
		}
		if !flag {
			res.UserID = v1
			res.Status = constant.OfflineStatus
		}
		respResult = append(respResult, res)
	}
	apiresp.GinSuccess(c, respResult)
}

func (u *UserApi) UserRegisterCount(c *gin.Context) {
	a2r.Call(user.UserClient.UserRegisterCount, u.Client, c)
}

// GetUsersOnlineTokenDetail Get user online token details.
func (u *UserApi) GetUsersOnlineTokenDetail(c *gin.Context) {
	var wsResult []*msggateway.GetUsersOnlineStatusResp_SuccessResult
	var respResult []*msggateway.SingleDetail
	flag := false
	var req msggateway.GetUsersOnlineStatusReq
	if err := c.BindJSON(&req); err != nil {
		apiresp.GinError(c, errs.ErrArgs.WithDetail(err.Error()).Wrap())
		return
	}
	conns, err := u.Discov.GetConns(c, config.Config.RpcRegisterName.OpenImMessageGatewayName)
	if err != nil {
		apiresp.GinError(c, err)
		return
	}
	// Online push message
	for _, v := range conns {
		msgClient := msggateway.NewMsgGatewayClient(v)
		reply, err := msgClient.GetUsersOnlineStatus(c, &req)
		if err != nil {
			log.ZWarn(c, "GetUsersOnlineStatus rpc  err", err)
			continue
		} else {
			wsResult = append(wsResult, reply.SuccessResult...)
		}
	}

	for _, v1 := range req.UserIDs {
		m := make(map[string][]string, 10)
		flag = false
		temp := new(msggateway.SingleDetail)
		for _, v2 := range wsResult {
			if v2.UserID == v1 {
				flag = true
				temp.UserID = v1
				temp.Status = constant.OnlineStatus
				for _, status := range v2.DetailPlatformStatus {
					if v, ok := m[status.Platform]; ok {
						m[status.Platform] = append(v, status.Token)
					} else {
						m[status.Platform] = []string{status.Token}
					}
				}
			}
		}
		for p, tokens := range m {
			t := new(msggateway.SinglePlatformToken)
			t.Platform = p
			t.Token = tokens
			t.Total = int32(len(tokens))
			temp.SinglePlatformToken = append(temp.SinglePlatformToken, t)
		}

		if flag {
			respResult = append(respResult, temp)
		}
	}

	apiresp.GinSuccess(c, respResult)
}

// SubscriberStatus Presence status of subscribed users.
func (u *UserApi) SubscriberStatus(c *gin.Context) {
	a2r.Call(user.UserClient.SubscribeOrCancelUsersStatus, u.Client, c)
}

// UnSubscriberStatus Unsubscribe a user's presence.
func (u *UserApi) UnSubscriberStatus(c *gin.Context) {
	a2r.Call(user.UserClient.SubscribeOrCancelUsersStatus, u.Client, c)
}

func (u *UserApi) GetUserStatus(c *gin.Context) {
	a2r.Call(user.UserClient.GetUserStatus, u.Client, c)
}

// GetSubscribeUsersStatus Get the online status of subscribers.
func (u *UserApi) GetSubscribeUsersStatus(c *gin.Context) {
	a2r.Call(user.UserClient.GetSubscribeUsersStatus, u.Client, c)
}
