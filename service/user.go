// Copyright 2019 tree xie
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

package service

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/lib/pq"
	"github.com/vicanso/elton"
	session "github.com/vicanso/elton-session"
	"github.com/vicanso/hes"
	lruTTL "github.com/vicanso/lru-ttl"
	"github.com/vicanso/origin/cs"
	"github.com/vicanso/origin/helper"
	"github.com/vicanso/origin/util"
	"gorm.io/gorm"

	"go.uber.org/zap"
)

const (
	// UserSessionInfoKey user session info
	UserSessionInfoKey = "user-session-info"
	// UserLoginToken user login token
	UserLoginToken = "loginToken"

	errUserCategory = "user"
)

var (
	errAccountOrPasswordInvalid = &hes.Error{
		Message:    "账户或者密码错误",
		StatusCode: http.StatusBadRequest,
		Category:   errUserCategory,
	}
)

var (
	// 用户角色
	userRolesMap map[string]string
	// 用户分组
	userGroupsMap map[string]string

	// userNameCache 用户名字缓存
	userNameCache *lruTTL.Cache
)

type (
	// UserSessionInfo user session info
	UserSessionInfo struct {
		Account   string
		ID        uint
		Roles     []string
		Groups    []string
		LoginedAt string
	}
	// UserSession user session struct
	UserSession struct {
		se   *session.Session
		info *UserSessionInfo
	}
	Users []*User
	// User user
	User struct {
		helper.Model

		Account  string `json:"account,omitempty" gorm:"type:varchar(20);not null;unique_index:idx_users_account"`
		Password string `json:"-,omitempty" gorm:"type:varchar(128);not null"`
		Name     string `json:"name,omitempty"`

		// 用户角色
		Roles pq.StringArray `json:"roles,omitempty" gorm:"type:text[]"`
		// 用户角色描述
		RolesDesc []string `json:"rolesDesc,omitempty" gorm:"-"`

		// 用户群组
		Groups pq.StringArray `json:"groups,omitempty" gorm:"type:text[]"`
		// 用户群组描述
		GroupsDesc []string `json:"groupsDesc,omitempty" gorm:"-"`

		// 用户销售分组
		MarketingGroup string `json:"marketingGroup,omitempty"`

		// 用户状态
		Status     int    `json:"status,omitempty"`
		StatusDesc string `json:"statusDesc,omitempty" gorm:"-"`
		Email      string `json:"email,omitempty"`
		Mobile     string `json:"mobile,omitempty"`
		// 推荐人
		Recommender     uint   `json:"recommender,omitempty"`
		RecommenderName string `json:"recommenderName,omitempty" gorm:"-"`
	}
	// UserRole user role
	UserRole struct {
		Name  string `json:"name,omitempty"`
		Value string `json:"value,omitempty"`
	}
	// UserGroup user group
	UserGroup struct {
		Name  string `json:"name,omitempty"`
		Value string `json:"value,omitempty"`
	}

	UserLoginRecords []*UserLoginRecord
	// UserLoginRecord user login
	UserLoginRecord struct {
		helper.Model

		Account       string `json:"account,omitempty" gorm:"type:varchar(20);not null;index:idx_user_logins_account"`
		UserAgent     string `json:"userAgent,omitempty"`
		IP            string `json:"ip,omitempty" gorm:"type:varchar(64);not null"`
		TrackID       string `json:"trackId,omitempty" gorm:"type:varchar(64);not null"`
		SessionID     string `json:"sessionId,omitempty" gorm:"type:varchar(64);not null"`
		XForwardedFor string `json:"xForwardedFor,omitempty" gorm:"type:varchar(128)"`
		Country       string `json:"country,omitempty" gorm:"type:varchar(64)"`
		Province      string `json:"province,omitempty" gorm:"type:varchar(64)"`
		City          string `json:"city,omitempty" gorm:"type:varchar(64)"`
		ISP           string `json:"isp,omitempty" gorm:"type:varchar(64)"`

		Width         int    `json:"width,omitempty"`
		Height        int    `json:"height,omitempty"`
		PixelRatio    int    `json:"pixelRatio,omitempty"`
		Platform      string `json:"platform,omitempty"`
		UUID          string `json:"uuid,omitempty"`
		SystemVersion string `json:"systemVersion,omitempty"`
		Brand         string `json:"brand,omitempty"`
		Version       string `json:"version,omitempty"`
		BuildNumber   string `json:"buildNumber,omitempty"`
	}
	// UserTrackRecord user track record
	UserTrackRecord struct {
		helper.Model

		TrackID   string `json:"trackId,omitempty" gorm:"type:varchar(64);not null;index:idx_user_track_id"`
		UserAgent string `json:"userAgent,omitempty"`
		IP        string `json:"ip,omitempty" gorm:"type:varchar(64);not null"`
		Country   string `json:"country,omitempty" gorm:"type:varchar(64)"`
		Province  string `json:"province,omitempty" gorm:"type:varchar(64)"`
		City      string `json:"city,omitempty" gorm:"type:varchar(64)"`
		ISP       string `json:"isp,omitempty" gorm:"type:varchar(64)"`
	}
	// UserSrv user service
	UserSrv struct{}

	// UserAmount user amount
	UserAmount struct {
		UpdatedAt     *time.Time `json:"updatedAt,omitempty"`
		EnabledAmount float64    `json:"enabledAmount,omitempty"`
		FrozenAmount  float64    `json:"frozenAmount,omitempty"`
		Amount        float64    `json:"amount,omitempty"`
	}
)

func init() {
	err := pgGetClient().AutoMigrate(
		&User{},
		&UserLoginRecord{},
		&UserTrackRecord{},
	)
	if err != nil {
		panic(err)
	}

	userRolesMap = map[string]string{
		cs.UserRoleNormal: "普通用户",
		cs.UserRoleAdmin:  "管理员",
		cs.UserRoleSu:     "超级用户",
	}
	userGroupsMap = map[string]string{
		cs.UserGroupIT:        "研发部",
		cs.UserGroupMarketing: "市场部",
		cs.UserGroupFinance:   "财务部",
		cs.UserGroupLogistics: "物流部",
	}

	ttl := 300 * time.Second
	// 本地开发环境，设置缓存为1秒
	if util.IsDevelopment() {
		ttl = time.Second
	}
	userNameCache = lruTTL.New(100, ttl)
}

// AfterCreate after create hook
func (u *User) AfterCreate(tx *gorm.DB) (err error) {
	// 首次创建账号，设置su权限
	if u.ID == 1 {
		tx.Model(u).Updates(User{
			Roles: []string{
				cs.UserRoleSu,
			},
		})
	}
	return
}

// BeforeCreate before create hook
func (u *User) BeforeCreate(_ *gorm.DB) (err error) {
	if len(u.Roles) == 0 {
		// 自动添加用户角色
		u.Roles = []string{
			cs.UserRoleNormal,
		}
	}
	return
}

func (u *User) AfterFind(_ *gorm.DB) (err error) {
	u.StatusDesc = getStatusDesc(u.Status)

	userRolesDesc := make([]string, 0)
	for _, role := range u.Roles {
		value, ok := userRolesMap[role]
		if ok {
			userRolesDesc = append(userRolesDesc, value)
		}
	}
	u.RolesDesc = userRolesDesc

	userGroupsDesc := make([]string, 0)
	for _, group := range u.Groups {
		value, ok := userGroupsMap[group]
		if ok {
			userGroupsDesc = append(userGroupsDesc, value)
		}
	}
	u.GroupsDesc = userGroupsDesc

	if u.Recommender != 0 {
		rcmder, _ := userSrv.FindByID(u.Recommender)
		u.RecommenderName = rcmder.Name
		if len(rcmder.Name) != 0 {
			u.RecommenderName = rcmder.Account
		}
	}

	return
}

func (us Users) AfterFind(tx *gorm.DB) (err error) {
	for _, u := range us {
		err = u.AfterFind(tx)
		if err != nil {
			return
		}
	}
	return
}

// ListRoles list all user roles
func (srv *UserSrv) ListRole() []*UserRole {
	userRoles := make([]*UserRole, 0)
	for key, value := range userRolesMap {
		userRoles = append(userRoles, &UserRole{
			Name:  value,
			Value: key,
		})
	}
	return userRoles
}

// ListGroups list all user group
func (srv *UserSrv) ListGroup() []*UserGroup {
	userGroups := make([]*UserGroup, 0)
	for key, value := range userGroupsMap {
		userGroups = append(userGroups, &UserGroup{
			Name:  value,
			Value: key,
		})
	}
	return userGroups
}

// createByID create a user model by id
func (srv *UserSrv) createByID(id uint) *User {
	u := &User{}
	u.Model.ID = id
	return u
}

// createLoginRecordByID cerate login record by id
func (srv *UserSrv) createLoginRecordByID(id uint) *UserLoginRecord {
	ulr := &UserLoginRecord{}
	ulr.Model.ID = id
	return ulr
}

// Add add user
func (srv *UserSrv) Add(data User) (u *User, err error) {
	u = &data
	if u.Status == 0 {
		u.Status = cs.StatusEnabled
	}

	err = pgCreate(u)
	return
}

// Login user login
func (srv *UserSrv) Login(account, password, token string) (u *User, err error) {
	u = &User{}
	err = pgGetClient().Where("account = ?", account).First(u).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			err = errAccountOrPasswordInvalid
		}
		return
	}
	pwd := util.Sha256(u.Password + token)
	// 用于自动化测试使用
	if util.IsDevelopment() && password == "fEqNCco3Yq9h5ZUglD3CZJT4lBsfEqNCco31Yq9h5ZUB" {
		pwd = password
	}
	if pwd != password {
		err = errAccountOrPasswordInvalid
		return
	}
	return
}

// UpdateByID update user by id
func (srv *UserSrv) UpdateByID(id uint, user User) (err error) {
	err = pgGetClient().Model(srv.createByID(id)).Updates(user).Error
	return
}

// UpdateByAccount update user by account
func (srv *UserSrv) UpdateByAccount(account string, value interface{}) (err error) {
	err = pgGetClient().Model(&User{}).Where("account = ?", account).Updates(value).Error
	return
}

// FindByID find user by id
func (srv *UserSrv) FindByID(id uint) (user *User, err error) {
	user = &User{}
	err = pgGetClient().First(user, "id = ?", id).Error
	return
}

// FindOneByAccount find one by account
func (srv *UserSrv) FindOneByAccount(account string) (user *User, err error) {
	user = &User{}
	err = pgGetClient().First(user, "account = ?", account).Error
	return
}

// UpdateLoginRecordByID update login record by id
func (srv *UserSrv) UpdateLoginRecordByID(id uint, value interface{}) (err error) {
	err = pgGetClient().Model(srv.createLoginRecordByID(id)).Updates(value).Error
	return
}

// AddLoginRecord add user login record
func (srv *UserSrv) AddLoginRecord(r *UserLoginRecord, c *elton.Context) (err error) {
	err = pgCreate(r)
	if r.ID != 0 {
		id := r.ID
		ip := r.IP
		go func() {
			lo, err := GetLocationByIP(ip, c)
			if err != nil {
				logger.Error("get location by ip fail",
					zap.String("ip", ip),
					zap.Error(err),
				)
				return
			}
			_ = srv.UpdateLoginRecordByID(id, map[string]string{
				"country":  lo.Country,
				"province": lo.Province,
				"city":     lo.City,
				"isp":      lo.ISP,
			})
		}()
	}
	return
}

// AddTrackRecord add track record
func (srv *UserSrv) AddTrackRecord(r *UserTrackRecord, c *elton.Context) (err error) {
	err = pgCreate(r)
	if r.ID != 0 {
		id := r.ID
		ip := r.IP
		go func() {
			lo, err := GetLocationByIP(ip, c)
			if err != nil {
				logger.Error("get location by ip fail",
					zap.String("ip", ip),
					zap.Error(err),
				)
				return
			}
			_ = srv.UpdateLoginRecordByID(id, map[string]string{
				"country":  lo.Country,
				"province": lo.Province,
				"city":     lo.City,
				"isp":      lo.ISP,
			})
		}()
	}
	return
}

// List list users
func (srv *UserSrv) List(params PGQueryParams, args ...interface{}) (result Users, err error) {
	result = make(Users, 0)
	err = pgQuery(params, args...).Find(&result).Error
	if err != nil {
		return
	}
	return
}

// Count count the users
func (srv *UserSrv) Count(args ...interface{}) (count int64, err error) {
	return pgCount(&User{}, args...)
}

// ListLoginRecord list login record
func (srv *UserSrv) ListLoginRecord(params PGQueryParams, args ...interface{}) (result UserLoginRecords, err error) {
	result = make(UserLoginRecords, 0)
	err = pgQuery(params, args...).Find(&result).Error
	return
}

// CountLoginRecord count login record
func (srv *UserSrv) CountLoginRecord(args ...interface{}) (count int64, err error) {
	return pgCount(&UserLoginRecord{}, args...)
}

// GetNameFromCache get user's name from cache
func (srv *UserSrv) GetNameFromCache(id uint) (name string, err error) {
	if id == 0 {
		return
	}
	value, ok := userNameCache.Get(id)
	if ok {
		return value.(string), nil
	}
	user, err := srv.FindByID(id)
	if err != nil {
		return
	}
	name = user.Name
	userNameCache.Add(id, name)
	return
}

// GetAmount get user's amount
func (*UserSrv) GetAmount(id uint) (userAmout *UserAmount, err error) {
	commissions, err := orderCommissionSrv.ListAll(PGQueryParams{
		Fields: "commissionAmount,updatedAt",
	}, "recommender = ?", id)
	if err != nil {
		return
	}
	amount := 0.0
	frozenAmount := 0.0
	initTime := time.Unix(0, 0)
	updatedAt := &initTime
	for _, commission := range commissions {
		if commission.UpdatedAt.Unix() > updatedAt.Unix() {
			updatedAt = commission.UpdatedAt
		}
		amount += commission.CommissionAmount
	}
	userAmout = &UserAmount{
		UpdatedAt:     updatedAt,
		EnabledAmount: amount - frozenAmount,
		Amount:        amount,
		FrozenAmount:  frozenAmount,
	}
	return
}

// GetUesrInfo get user info
func (us *UserSession) GetInfo() (info *UserSessionInfo, err error) {
	info = us.info
	if info != nil {
		return
	}
	data := us.se.GetString(UserSessionInfoKey)
	info = new(UserSessionInfo)
	err = json.Unmarshal([]byte(data), info)
	if err != nil {
		return
	}
	us.info = info
	return
}

// MustGetInfo get user info, if not exists, it will panic
func (us *UserSession) MustGetInfo() (info *UserSessionInfo) {
	info, err := us.GetInfo()
	if err != nil {
		panic(err)
	}
	if info == nil {
		panic(errors.New("get user info fail"))
	}
	return info
}

// IsLogined check user is logined
func (us *UserSession) IsLogined() bool {
	info, err := us.GetInfo()
	if err != nil || info == nil {
		return false
	}
	return info.Account != ""
}

// SetInfo set user session info
func (us *UserSession) SetInfo(data User) (err error) {
	info := UserSessionInfo{
		Account:   data.Account,
		ID:        data.ID,
		Roles:     data.Roles,
		Groups:    data.Groups,
		LoginedAt: util.NowString(),
	}
	buf, err := json.Marshal(&info)
	if err != nil {
		return
	}
	err = us.se.Set(UserSessionInfoKey, string(buf))
	if err != nil {
		return
	}
	return
}

// GetAccount get the account
func (us *UserSession) GetAccount() string {
	info := us.MustGetInfo()
	return info.Account
}

// GetID get user id
func (us *UserSession) GetID() uint {
	info := us.MustGetInfo()
	return info.ID
}

// SetLoginToken get user login token
func (us *UserSession) SetLoginToken(token string) error {
	return us.se.Set(UserLoginToken, token)
}

// GetLoginToken get user login token
func (us *UserSession) GetLoginToken() string {
	return us.se.GetString(UserLoginToken)
}

// GetRoles get user roles
func (us *UserSession) GetRoles() []string {
	info := us.MustGetInfo()
	return info.Roles
}

// GetGroups get user groups
func (us *UserSession) GetGroups() []string {
	info := us.MustGetInfo()
	return info.Groups
}

func (us *UserSession) GetLoginedAt() string {
	info := us.MustGetInfo()
	return info.LoginedAt
}

// Destroy destroy user session
func (us *UserSession) Destroy() error {
	return us.se.Destroy()
}

// Refresh refresh user session
func (us *UserSession) Refresh() error {
	return us.se.Refresh()
}

// ClearSessionID clear session id
func (us *UserSession) ClearSessionID() {
	us.se.ID = ""
}

// NewUserSession create a user session
func NewUserSession(c *elton.Context) *UserSession {
	v, ok := c.Get(session.Key)
	if !ok {
		return nil
	}
	data, ok := c.Get(cs.UserSession)
	if ok {
		us, ok := data.(*UserSession)
		if ok {
			return us
		}
	}
	us := &UserSession{
		se: v.(*session.Session),
	}
	c.Set(cs.UserSession, us)

	return us
}
