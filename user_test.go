package permissions

import (
	"database/sql"
	"flag"
	"testing"
)

var driverName = flag.String("dbDrv", "postgres", "")
var dataSourceName = flag.String("dbURL", "host=127.0.0.1 dbname=tpt_data_test user=xxx password=xxx sslmode=disable", "")

func dbTest(t *testing.T, cb func(db *sql.DB)) {
	conn, err := sql.Open(*driverName, *dataSourceName)
	if err != nil {
		t.Error(err)
		return
	}
	defer conn.Close()

	_, err = conn.Exec(`
DROP TABLE IF EXISTS tpt_user_roles;
DROP TABLE IF EXISTS tpt_users;
DROP TABLE IF EXISTS tpt_roles;
DROP TABLE IF EXISTS tpt_user_profiles;

CREATE TABLE tpt_roles
(
  id serial,
  name character varying(50),
  permission_keys character varying(40000),
  description character varying(200),
  created_at timestamp with time zone DEFAULT now(),
  updated_at timestamp with time zone DEFAULT now(),
  CONSTRAINT tpt_roles_pkey PRIMARY KEY (id),
  CONSTRAINT tpt_roles_name_uq UNIQUE (name)
);

CREATE TABLE tpt_users
(
  id serial,
  name character varying(50),
  password character varying(200),
  phone character varying(50),
  email character varying(100),
  description character varying(200),
  created_at timestamp with time zone DEFAULT now(),
  updated_at timestamp with time zone DEFAULT now(),
  state integer NOT NULL DEFAULT 0,
  CONSTRAINT tpt_users_pkey PRIMARY KEY (id),
  CONSTRAINT tpt_users_name_uq UNIQUE (name)
);


CREATE TABLE tpt_user_profiles
(
  id serial,
  usr character varying(50),
  name character varying(50),
  value character varying(10000) NOT NULL,
  created_at timestamp without time zone,
  updated_at timestamp without time zone,
  CONSTRAINT tpt_user_profiles_pkey PRIMARY KEY (id)
);

CREATE TABLE tpt_user_roles
(
  id serial,
  user_id bigint NOT NULL,
  role_id bigint NOT NULL,
  created_at timestamp without time zone,
  updated_at timestamp without time zone,
  CONSTRAINT tpt_user_roles_pkey PRIMARY KEY (id),
  CONSTRAINT tpt_user_roles_role_id_fkey FOREIGN KEY (role_id)
      REFERENCES public.tpt_roles (id) MATCH SIMPLE
      ON UPDATE NO ACTION ON DELETE CASCADE,
  CONSTRAINT tpt_user_roles_user_id_fkey FOREIGN KEY (user_id)
      REFERENCES public.tpt_users (id) MATCH SIMPLE
      ON UPDATE NO ACTION ON DELETE CASCADE
);`)
	if err != nil {
		t.Error(err)
		return
	}
	PlaceholderFormat = Dollar
	IsReturning = true

	cb(conn)
}

func TestRoleDao(t *testing.T) {
	dbTest(t, func(db *sql.DB) {
		role1 := &Role{
			Name:           "a",
			Description:    "a_descr",
			PermissionKeys: "k1,k2",
		}

		id, err := role1.CreateIt(db)
		if err != nil {
			t.Error(err)
			return
		}

		role2, err := Roles.FindByID(db, id)
		if err != nil {
			t.Error(err)
			return
		}

		role3, err := Roles.FindByName(db, role1.Name)
		if err != nil {
			t.Error(err)
			return
		}

		assertRole := func(oldRole, newRole *Role) {
			if oldRole.Name != newRole.Name {
				t.Error(oldRole.Name, newRole.Name)
			}

			if oldRole.Description != newRole.Description {
				t.Error(oldRole.Description, newRole.Description)
			}
			if oldRole.PermissionKeys != newRole.PermissionKeys {
				t.Error(oldRole.PermissionKeys, newRole.PermissionKeys)
			}
			if newRole.CreatedAt.IsZero() {
				t.Error("newRole.CreatedAt.IsZero()")
			}
			if newRole.UpdatedAt.IsZero() {
				t.Error("newRole.UpdatedAt.IsZero()")
			}
		}
		for _, newRole := range []*Role{role2, role3} {
			assertRole(role1, newRole)
		}

		role2.Name = "aaa"
		role2.Description = "aaa_descr"
		role2.PermissionKeys = "k3,k4"
		if err := role2.UpdateIt(db); err != nil {
			t.Error(err)
			return
		}

		role4, err := Roles.FindByID(db, id)
		if err != nil {
			t.Error(err)
			return
		}
		assertRole(role2, role4)
	})
}

func TestUserDao(t *testing.T) {
	dbTest(t, func(db *sql.DB) {
		user1 := &User{
			Name:        "a",
			Description: "a_descr",
			Password:    "a_pwd",
			Phone:       "123",
			Email:       "a@h.com",
			State:       23,
		}

		id, err := user1.CreateIt(db)
		if err != nil {
			t.Error(err)
			return
		}

		user2, err := Users.FindByID(db, id)
		if err != nil {
			t.Error(err)
			return
		}

		user3, err := Users.FindByName(db, user1.Name)
		if err != nil {
			t.Error(err)
			return
		}

		assertUser := func(oldUser, newUser *User) {
			if oldUser.Name != newUser.Name {
				t.Error(oldUser.Name, newUser.Name)
			}

			if oldUser.Description != newUser.Description {
				t.Error(oldUser.Description, newUser.Description)
			}
			if oldUser.Password != newUser.Password {
				t.Error(oldUser.Password, newUser.Password)
			}
			if oldUser.Phone != newUser.Phone {
				t.Error(oldUser.Phone, newUser.Phone)
			}
			if oldUser.Email != newUser.Email {
				t.Error(oldUser.Email, newUser.Email)
			}
			if oldUser.State != newUser.State {
				t.Error(oldUser.State, newUser.State)
			}
			if newUser.CreatedAt.IsZero() {
				t.Error("newUser.CreatedAt.IsZero()")
			}
			if newUser.UpdatedAt.IsZero() {
				t.Error("newUser.UpdatedAt.IsZero()")
			}
		}
		for _, newUser := range []*User{user2, user3} {
			assertUser(user1, newUser)
		}

		user2.Name = "aaa"
		user2.Description = "aaa_descr"
		user2.Password = "aaa_pwd"
		user2.Phone = "23"
		user2.Email = "a1@h.com"
		user2.State = 123
		if err := user2.UpdateIt(db); err != nil {
			t.Error(err)
			return
		}

		user4, err := Users.FindByID(db, id)
		if err != nil {
			t.Error(err)
			return
		}
		assertUser(user2, user4)
	})
}

func TestUserProfileDao(t *testing.T) {
	dbTest(t, func(db *sql.DB) {
		userProfile1 := &UserProfile{
			User:  "u1",
			Name:  "a",
			Value: "a_descr",
		}

		id, err := userProfile1.CreateIt(db)
		if err != nil {
			t.Error(err)
			return
		}

		userProfile2, err := UserProfiles.FindByID(db, id)
		if err != nil {
			t.Error(err)
			return
		}

		assertUserProfile := func(oldUserProfile, newUserProfile *UserProfile) {
			if oldUserProfile.Name != newUserProfile.Name {
				t.Error(oldUserProfile.Name, newUserProfile.Name)
			}

			if oldUserProfile.User != newUserProfile.User {
				t.Error(oldUserProfile.User, newUserProfile.User)
			}
			if oldUserProfile.Value != newUserProfile.Value {
				t.Error(oldUserProfile.Value, newUserProfile.Value)
			}
			if newUserProfile.CreatedAt.IsZero() {
				t.Error("newUserProfile.CreatedAt.IsZero()")
			}
			if newUserProfile.UpdatedAt.IsZero() {
				t.Error("newUserProfile.UpdatedAt.IsZero()")
			}
		}
		for _, newUserProfile := range []*UserProfile{userProfile2} {
			assertUserProfile(userProfile1, newUserProfile)
		}

		userProfile2.User = "u3"
		userProfile2.Name = "aaa"
		userProfile2.Value = "aaa_descr"
		if err := userProfile2.UpdateIt(db); err != nil {
			t.Error(err)
			return
		}

		userProfile4, err := UserProfiles.FindByID(db, id)
		if err != nil {
			t.Error(err)
			return
		}
		assertUserProfile(userProfile2, userProfile4)
	})
}

func TestUserRoleDao(t *testing.T) {
	dbTest(t, func(db *sql.DB) {
		role1 := &Role{
			Name:           "a1",
			Description:    "a_descr",
			PermissionKeys: "k1,k2",
		}
		role2 := &Role{
			Name:           "a2",
			Description:    "a_descr",
			PermissionKeys: "k1,k2",
		}
		role3 := &Role{
			Name:           "a3",
			Description:    "a_descr",
			PermissionKeys: "k1,k2",
		}
		user1 := &User{
			Name:        "user1",
			Description: "a_descr",
			Password:    "a_pwd",
			Phone:       "123",
			Email:       "a@h.com",
			State:       23,
		}

		r1, err := role1.CreateIt(db)
		if err != nil {
			t.Error(err)
			return
		}
		r2, err := role2.CreateIt(db)
		if err != nil {
			t.Error(err)
			return
		}
		r3, err := role3.CreateIt(db)
		if err != nil {
			t.Error(err)
			return
		}
		u1, err := user1.CreateIt(db)
		if err != nil {
			t.Error(err)
			return
		}

		err = Users.AddRole(db, u1, r1)
		if err != nil {
			t.Error(err)
			return
		}

		err = Users.AddRole(db, u1, r2)
		if err != nil {
			t.Error(err)
			return
		}

		assertUserRole := func(idList []int64) {
			roleList1, err := Users.ListRoles(db, u1)
			if err != nil {
				t.Error(err)
				return
			}

			roleList2, err := Roles.FindByUserName(db, user1.Name)
			if err != nil {
				t.Error(err)
				return
			}

			for _, roleList := range [][]*Role{roleList1, roleList2} {
				if len(roleList) != len(idList) {
					t.Error("len(roleList) != len(idList)", len(roleList), len(idList), idList)
					for _, r := range roleList {
						t.Log(r.ID)
					}
				}

				for _, id := range idList {
					found := false
					for _, r := range roleList {
						if id == r.ID {
							found = true
							break
						}
					}
					if !found {
						t.Error(id, "isn't found")
						return
					}
				}
			}
		}

		t.Log("====", 1)
		assertUserRole([]int64{r1, r2})
		t.Log("====", 1, "END")

		err = Users.AddRole(db, u1, r3)
		if err != nil {
			t.Error(err)
			return
		}
		t.Log("====", 2)
		assertUserRole([]int64{r1, r2, r3})
		t.Log("====", 2, "END")

		err = Users.RemoveRole(db, u1, r1)
		if err != nil {
			t.Error(err)
			return
		}

		t.Log("====", 3)
		assertUserRole([]int64{r2, r3})
		t.Log("====", 3, "END")
	})
}
