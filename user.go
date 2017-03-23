package permissions

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"
)

// ErrNotUpdated - 表示没有更新任何记录
var ErrNotUpdated = errors.New("no record is updated")

// ErrNotDeleted - 表示没有删除任何记录
var ErrNotDeleted = errors.New("no record is deleted")

// ThrowPrimaryKeyInvalid 返回一个 主键无效的错误
func ThrowPrimaryKeyInvalid(tableName string) error {
	return errors.New("primary key of '" + tableName + "' is invalid")
}

// RowScanner is the interface that wraps the Scan method.
//
// Scan behaves like database/sql.Row.Scan.
type RowScanner interface {
	Scan(...interface{}) error
}

var (
	// Question is a PlaceholderFormat instance that leaves placeholders as
	// question marks.
	Question = func(sql string) (string, error) {
		return sql, nil
	}

	// Dollar is a PlaceholderFormat instance that replaces placeholders with
	// dollar-prefixed positional placeholders (e.g. $1, $2, $3).
	Dollar = func(sql string) (string, error) {
		buf := &bytes.Buffer{}
		i := 0
		for {
			p := strings.Index(sql, "?")
			if p == -1 {
				break
			}

			if len(sql[p:]) > 1 && sql[p:p+2] == "??" { // escape ?? => ?
				buf.WriteString(sql[:p])
				buf.WriteString("?")
				if len(sql[p:]) == 1 {
					break
				}
				sql = sql[p+2:]
			} else {
				i++
				buf.WriteString(sql[:p])
				fmt.Fprintf(buf, "$%d", i)
				sql = sql[p+1:]
			}
		}

		buf.WriteString(sql)
		return buf.String(), nil
	}

	// PlaceholderFormat takes a SQL statement and replaces each question mark
	// placeholder with a (possibly different) SQL placeholder.
	PlaceholderFormat = Question

	// IsReturning use returning case in the insert statement.
	IsReturning bool
)

// Role 代表一个用户角色
type Role struct {
	ID             int64     `json:"id,omitempty"`
	Name           string    `json:"name,omitempty"`
	Description    string    `json:"description,omitempty"`
	PermissionKeys string    `json:"permission_keys,omitempty"`
	UpdatedAt      time.Time `json:"updated_at,omitempty"`
	CreatedAt      time.Time `json:"created_at,omitempty"`
}

func (role *Role) Keys() {
	var keys []string
	if err := json.Unmarshal([]byte(role.PermissionKeys), &keys); err != nil {

	}
}

// User 代表一个用户
type User struct {
	ID          int64     `json:"id,omitempty"`
	Name        string    `json:"name,omitempty"`
	Description string    `json:"description,omitempty"`
	Password    string    `json:"password,omitempty"`
	Phone       string    `json:"phone,omitempty"`
	Email       string    `json:"email,omitempty"`
	State       int64     `json:"state,omitempty"`
	CreatedAt   time.Time `json:"created_at,omitempty"`
	UpdatedAt   time.Time `json:"updated_at,omitempty"`
}

func (user *User) CreateIt(db *sql.DB) (int64, error) {
	return Users.CreateIt(db, user)
}

func (user *User) UpdateIt(db *sql.DB) error {
	return Users.UpdateIt(db, user)
}

func (user *User) DeleteIt(db *sql.DB) error {
	return Users.DeleteIt(db, user)
}

// UserProfile 代表用户的属性
type UserProfile struct {
	ID        int64     `json:"id,omitempty"`
	User      string    `json:"usr,omitempty"`
	Name      string    `json:"name,omitempty"`
	Value     string    `json:"value,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

func (userProfile *UserProfile) CreateIt(db *sql.DB) (int64, error) {
	return UserProfiles.CreateIt(db, userProfile)
}

func (userProfile *UserProfile) UpdateIt(db *sql.DB) error {
	return UserProfiles.UpdateIt(db, userProfile)
}

func (userProfile *UserProfile) DeleteIt(db *sql.DB) error {
	return UserProfiles.DeleteIt(db, userProfile)
}

func (role *Role) CreateIt(db *sql.DB) (int64, error) {
	return Roles.CreateIt(db, role)
}

func (role *Role) UpdateIt(db *sql.DB) error {
	return Roles.UpdateIt(db, role)
}

func (role *Role) DeleteIt(db *sql.DB) error {
	return Roles.DeleteIt(db, role)
}

var (
	Roles        = roles{}
	Users        = users{}
	UserProfiles = userProfiles{}
)

// User 代表一个用户
type UserRBAC struct {
	User User
	//Role []Role
}

func (self *UserRBAC) Name() string {
	return self.User.Name
}

func (self *UserRBAC) IsAdmin() bool {
	return self.User.Name == "admin" || self.User.Name == "administrator"
}

func (self *UserRBAC) HasPermission(key string) bool {
	panic("not implemented")
}

func (self *UserRBAC) Data() interface{} {
	panic("not implemented")
}

func QueryUserRBAC(db *sql.DB, userName string) (*UserRBAC, error) {
	user, err := Users.FindByName(db, userName)
	if err != nil {
		return nil, errors.New("load user fial, " + err.Error())
	}
	roles, err := Users.ListRoles(db, user.ID)
	if err != nil {
		return nil, errors.New("load roles fial, " + err.Error())
	}

	rbac := &UserRBAC{
		User: *user,
	}
	for _, r := range roles {
		fmt.Println(r.Name)
		panic("not implemented")
	}

	return rbac, nil
}
