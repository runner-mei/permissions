package permissions

import (
	"database/sql"
	"time"

	"github.com/lib/pq"
)

type roles struct{}

func (self *roles) scan(scanner RowScanner) (*Role, error) {
	var value Role
	var nullDescription sql.NullString
	var nullPermissionKeys sql.NullString
	var nullCreatedAt pq.NullTime
	var nullUpdatedAt pq.NullTime

	e := scanner.Scan(
		&value.ID,
		&value.Name,
		&nullDescription,
		&nullPermissionKeys,
		&nullCreatedAt,
		&nullUpdatedAt)
	if nil != e {
		return nil, e
	}

	if nullDescription.Valid {
		value.Description = nullDescription.String
	}
	if nullPermissionKeys.Valid {
		value.PermissionKeys = nullPermissionKeys.String
	}
	if nullCreatedAt.Valid {
		value.CreatedAt = nullCreatedAt.Time
	}
	if nullUpdatedAt.Valid {
		value.UpdatedAt = nullUpdatedAt.Time
	}
	return &value, nil
}

const rolePrefix = "select id, name, description, permission_keys, created_at, updated_at from tpt_roles "

func (self *roles) QueryRowWith(db *sql.DB, queryString string, args ...interface{}) (*Role, error) {
	queryString, err := PlaceholderFormat(queryString)
	if err != nil {
		return nil, err
	}

	row := db.QueryRow(rolePrefix+queryString, args...)
	return self.scan(row)
}

func (self *roles) QueryWith(db *sql.DB, queryString string, args ...interface{}) ([]*Role, error) {
	queryString, err := PlaceholderFormat(queryString)
	if err != nil {
		return nil, err
	}

	rows, err := db.Query(rolePrefix+queryString, args...)
	if nil != err {
		return nil, err
	}
	defer rows.Close()

	results := make([]*Role, 0, 4)
	for rows.Next() {
		v, err := self.scan(rows)
		if nil != err {
			return nil, err
		}
		results = append(results, v)
	}
	return results, rows.Err()
}

func (self *roles) FindByID(db *sql.DB, id int64) (*Role, error) {
	return self.QueryRowWith(db, "WHERE id = ?", id)
}

func (self *roles) FindByName(db *sql.DB, name string) (*Role, error) {
	return self.QueryRowWith(db, "WHERE name = ?", name)
}

func (self *roles) FindByUserID(db *sql.DB, userID int64) ([]*Role, error) {
	return self.QueryWith(db, "WHERE EXISTS (SELECT * FROM tpt_user_roles WHERE user_id = ? AND tpt_roles.id = tpt_user_roles.role_id)", userID)
}

func (self *roles) FindByUserName(db *sql.DB, username string) ([]*Role, error) {
	return self.QueryWith(db, "WHERE EXISTS (SELECT * FROM tpt_user_roles WHERE tpt_roles.id = tpt_user_roles.role_id AND EXISTS (SELECT * FROM tpt_users WHERE name = ? AND tpt_user_roles.user_id = tpt_users.id))", username)
}

func (self *roles) CreateIt(db *sql.DB, value *Role) (int64, error) {
	sqlString := "INSERT INTO tpt_roles(name, description, permission_keys, created_at, updated_at) VALUES (?, ?, ?, ?, ?)"
	sqlString, err := PlaceholderFormat(sqlString)
	if err != nil {
		return 0, err
	}

	now := time.Now()
	if IsReturning {
		sqlString = sqlString + " RETURNING \"id\""

		err := db.QueryRow(sqlString,
			value.Name,
			value.Description,
			value.PermissionKeys,
			now,
			now).Scan(&value.ID)
		return value.ID, err
	}

	result, err := db.Exec(sqlString, value.Name, value.Description, value.PermissionKeys, now, now)
	if nil != err {
		return 0, err
	}
	return result.LastInsertId()
}

func (self *roles) UpdateIt(db *sql.DB, value *Role) error {
	if 0 == value.ID {
		return ThrowPrimaryKeyInvalid("tpt_roles")
	}

	updateString := "UPDATE tpt_roles SET name=?, description=?, permission_keys=?, updated_at=? WHERE id = ?"
	updateString, err := PlaceholderFormat(updateString)
	if err != nil {
		return err
	}

	result, err := db.Exec(updateString,
		value.Name,
		value.Description,
		value.PermissionKeys,
		time.Now(),
		value.ID)
	if nil != err {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if nil != err {
		return err
	}
	if 0 == rowsAffected {
		return ErrNotUpdated
	}
	return nil
}

func (self *roles) DeleteIt(db *sql.DB, value *Role) error {
	return self.DeleteByID(db, value.ID)
}

func (self *roles) DeleteByID(db *sql.DB, key int64) error {
	if 0 == key {
		return ThrowPrimaryKeyInvalid("tpt_roles")
	}

	deleteString := "DELETE FROM tpt_roles WHERE id = ?"
	deleteString, err := PlaceholderFormat(deleteString)
	if err != nil {
		return err
	}
	result, err := db.Exec(deleteString, key)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrNotDeleted
	}
	return nil
}

type users struct{}

func (self *users) scan(scanner RowScanner) (*User, error) {
	var value User
	var nullDescription sql.NullString
	var nullPassword sql.NullString
	var nullPhone sql.NullString
	var nullEmail sql.NullString
	var nullState sql.NullInt64
	//var nullAttributes sql.NullString
	var nullCreatedAt pq.NullTime
	var nullUpdatedAt pq.NullTime

	e := scanner.Scan(
		&value.ID,
		&value.Name,
		&nullDescription,
		&nullPassword,
		&nullPhone,
		&nullEmail,
		&nullState,
		//&nullAttributes,
		&nullCreatedAt,
		&nullUpdatedAt)
	if nil != e {
		return nil, e
	}

	if nullDescription.Valid {
		value.Description = nullDescription.String
	}
	if nullPassword.Valid {
		value.Password = nullPassword.String
	}
	if nullPhone.Valid {
		value.Phone = nullPhone.String
	}
	if nullEmail.Valid {
		value.Email = nullEmail.String
	}
	if nullState.Valid {
		value.State = nullState.Int64
	}
	if nullCreatedAt.Valid {
		value.CreatedAt = nullCreatedAt.Time
	}
	if nullUpdatedAt.Valid {
		value.UpdatedAt = nullUpdatedAt.Time
	}
	return &value, nil
}

const userPrefix = "select id, name, description, password, phone, email, state, created_at, updated_at from tpt_users "

func (self *users) QueryRowWith(db *sql.DB, queryString string, args ...interface{}) (*User, error) {
	queryString, err := PlaceholderFormat(queryString)
	if err != nil {
		return nil, err
	}

	row := db.QueryRow(userPrefix+queryString, args...)
	return self.scan(row)
}

func (self *users) QueryWith(db *sql.DB, queryString string, args ...interface{}) ([]*User, error) {
	queryString, err := PlaceholderFormat(queryString)
	if err != nil {
		return nil, err
	}

	rows, err := db.Query(userPrefix+queryString, args...)
	if nil != err {
		return nil, err
	}
	defer rows.Close()

	results := make([]*User, 0, 4)
	for rows.Next() {
		v, err := self.scan(rows)
		if nil != err {
			return nil, err
		}
		results = append(results, v)
	}
	return results, rows.Err()
}

func (self *users) FindByID(db *sql.DB, id int64) (*User, error) {
	return self.QueryRowWith(db, "WHERE id = ?", id)
}

func (self *users) FindByName(db *sql.DB, name string) (*User, error) {
	return self.QueryRowWith(db, "WHERE name = ?", name)
}

func (self *users) AddRole(db *sql.DB, userID, roleID int64) error {
	insertString := "INSERT INTO tpt_user_roles(user_id, role_id, created_at, updated_at) VALUES (?, ?, ?, ?)"
	insertString, err := PlaceholderFormat(insertString)
	if err != nil {
		return err
	}
	now := time.Now()
	_, err = db.Exec(insertString,
		userID,
		roleID,
		now,
		now)
	return err
}

func (self *users) RemoveRole(db *sql.DB, userID, roleID int64) error {
	deleteString := "DELETE FROM tpt_user_roles WHERE user_id = ? AND role_id = ?"
	deleteString, err := PlaceholderFormat(deleteString)
	if err != nil {
		return err
	}

	_, err = db.Exec(deleteString,
		userID,
		roleID)
	return err
}

func (self *users) ListRoles(db *sql.DB, userID int64) ([]*Role, error) {
	return Roles.FindByUserID(db, userID)
}

func (self *users) CreateIt(db *sql.DB, value *User) (int64, error) {
	sqlString := "INSERT INTO tpt_users(name, description, password, phone, email, state, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?)"
	sqlString, err := PlaceholderFormat(sqlString)
	if err != nil {
		return 0, err
	}

	now := time.Now()
	if IsReturning {
		sqlString = sqlString + " RETURNING \"id\""
		err := db.QueryRow(sqlString,
			value.Name,
			value.Description,
			value.Password,
			value.Phone,
			value.Email,
			value.State,
			now,
			now).Scan(&value.ID)

		return value.ID, err
	}

	result, err := db.Exec(sqlString,
		value.Name,
		value.Description,
		value.Password,
		value.Phone,
		value.Email,
		value.State,
		now,
		now)
	if nil != err {
		return 0, err
	}
	return result.LastInsertId()
}

func (self *users) UpdateIt(db *sql.DB, value *User) error {
	if 0 == value.ID {
		return ThrowPrimaryKeyInvalid("tpt_users")
	}

	updateString := "UPDATE tpt_users SET name=?, description=?, password=?, phone=?, email=?, state=?, updated_at=? WHERE id = ?"
	updateString, err := PlaceholderFormat(updateString)
	if err != nil {
		return err
	}

	result, err := db.Exec(updateString,
		value.Name,
		value.Description,
		value.Password,
		value.Phone,
		value.Email,
		value.State,
		time.Now(),
		value.ID)
	if nil != err {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if nil != err {
		return err
	}
	if 0 == rowsAffected {
		return ErrNotUpdated
	}
	return nil
}

func (self *users) DeleteIt(db *sql.DB, value *User) error {
	return self.DeleteByID(db, value.ID)
}

func (self *users) DeleteByID(db *sql.DB, key int64) error {
	if 0 == key {
		return ThrowPrimaryKeyInvalid("tpt_users")
	}

	deleteString := "DELETE FROM tpt_users WHERE id = ?"
	deleteString, err := PlaceholderFormat(deleteString)
	if err != nil {
		return err
	}
	result, err := db.Exec(deleteString, key)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrNotDeleted
	}
	return nil
}

type userProfiles struct{}

func (self *userProfiles) scan(scanner RowScanner) (*UserProfile, error) {
	var value UserProfile
	var nullUser sql.NullString
	var nullName sql.NullString
	var nullValue sql.NullString
	var nullCreatedAt pq.NullTime
	var nullUpdatedAt pq.NullTime

	e := scanner.Scan(
		&value.ID,
		&nullUser,
		&nullName,
		&nullValue,
		&nullCreatedAt,
		&nullUpdatedAt)
	if nil != e {
		return nil, e
	}

	if nullUser.Valid {
		value.User = nullUser.String
	}
	if nullName.Valid {
		value.Name = nullName.String
	}
	if nullValue.Valid {
		value.Value = nullValue.String
	}
	if nullCreatedAt.Valid {
		value.CreatedAt = nullCreatedAt.Time
	}
	if nullUpdatedAt.Valid {
		value.UpdatedAt = nullUpdatedAt.Time
	}
	return &value, nil
}

const userProfilePrefix = "select id, usr, name, value, created_at, updated_at from tpt_user_profiles "

func (self *userProfiles) QueryRowWith(db *sql.DB, queryString string, args ...interface{}) (*UserProfile, error) {
	queryString, err := PlaceholderFormat(queryString)
	if err != nil {
		return nil, err
	}

	row := db.QueryRow(userProfilePrefix+queryString, args...)
	return self.scan(row)
}

func (self *userProfiles) QueryWith(db *sql.DB, queryString string, args ...interface{}) ([]*UserProfile, error) {
	queryString, err := PlaceholderFormat(queryString)
	if err != nil {
		return nil, err
	}

	rows, err := db.Query(userProfilePrefix+queryString, args...)
	if nil != err {
		return nil, err
	}
	defer rows.Close()

	results := make([]*UserProfile, 0, 4)
	for rows.Next() {
		v, err := self.scan(rows)
		if nil != err {
			return nil, err
		}
		results = append(results, v)
	}
	return results, rows.Err()
}

func (self *userProfiles) FindByID(db *sql.DB, id int64) (*UserProfile, error) {
	return self.QueryRowWith(db, "WHERE id = ?", id)
}

func (self *userProfiles) CreateIt(db *sql.DB, value *UserProfile) (int64, error) {
	sqlString := "INSERT INTO tpt_user_profiles(usr, name, value, created_at, updated_at) VALUES (?, ?, ?, ?, ?)"
	sqlString, err := PlaceholderFormat(sqlString)
	if err != nil {
		return 0, err
	}

	now := time.Now()
	if IsReturning {
		sqlString = sqlString + " RETURNING \"id\""

		err := db.QueryRow(sqlString,
			value.User,
			value.Name,
			value.Value,
			now,
			now).Scan(&value.ID)
		return value.ID, err
	}

	result, err := db.Exec(sqlString,
		value.User,
		value.Name,
		value.Value,
		now,
		now)
	if nil != err {
		return 0, err
	}
	return result.LastInsertId()
}

func (self *userProfiles) UpdateIt(db *sql.DB, value *UserProfile) error {
	if 0 == value.ID {
		return ThrowPrimaryKeyInvalid("tpt_user_profiles")
	}

	updateString := "UPDATE tpt_user_profiles SET usr=?, name=?, value=?, updated_at=? WHERE id = ?"
	updateString, err := PlaceholderFormat(updateString)
	if err != nil {
		return err
	}

	result, err := db.Exec(updateString,
		value.User,
		value.Name,
		value.Value,
		time.Now(),
		value.ID)
	if nil != err {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if nil != err {
		return err
	}
	if 0 == rowsAffected {
		return ErrNotUpdated
	}
	return nil
}

func (self *userProfiles) DeleteIt(db *sql.DB, value *UserProfile) error {
	return self.DeleteByID(db, value.ID)
}

func (self *userProfiles) DeleteByID(db *sql.DB, key int64) error {
	if 0 == key {
		return ThrowPrimaryKeyInvalid("tpt_user_profiles")
	}

	deleteString := "DELETE FROM tpt_user_profiles WHERE id = ?"
	deleteString, err := PlaceholderFormat(deleteString)
	if err != nil {
		return err
	}
	result, err := db.Exec(deleteString, key)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrNotDeleted
	}
	return nil
}
