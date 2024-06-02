package mysql

import (
	"datasupervision/internal/service"
	"fmt"
)

func (db *DB) AddRole(rolename string) error {
	_, err := db.DBConn.Exec("CREATE ROLE " + rolename)
	if err != nil {
		return err
	}
	return nil
}

func (db *DB) DeleteRole(rolename string) error {
	_, err := db.DBConn.Exec("DROP ROLE " + rolename)
	if err != nil {
		return err
	}
	return nil
}

func (db *DB) ListRoles() ([]service.Role, error) {
	rows, err := db.DBConn.Query("SELECT rolname FROM pg_roles")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var roles []service.Role
	for rows.Next() {
		var role service.Role
		if err := rows.Scan(&role.Rolename); err != nil {
			return nil, err
		}
		roles = append(roles, role)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return roles, nil
}

func (db *DB) ListUsers() ([]service.User, error) {
	rows, err := db.DBConn.Query("SELECT usename FROM pg_user")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []service.User
	for rows.Next() {
		var user service.User
		if err := rows.Scan(&user.Username); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, rows.Err()
}

func (db *DB) CreateUser(username, password string) error {
	_, err := db.DBConn.Exec(fmt.Sprintf("CREATE USER %s WITH PASSWORD '%s'", username, password))
	return err
}

func (db *DB) DeleteUser(username string) error {
	_, err := db.DBConn.Exec(fmt.Sprintf("DROP USER %s", username))
	return err
}

func (db *DB) AssignRoleToUser(username, rolename string) error {
	query := fmt.Sprintf("GRANT %s TO %s", rolename, username)
	println(query)

	_, err := db.DBConn.Exec(query)
	if err != nil {
		return err
	}
	return nil
}

func (db *DB) ListUsersForRole(role string) ([]service.User, error) {
	query := `
    SELECT u.usename
    FROM pg_user u
    JOIN pg_auth_members m ON u.usesysid = m.member
    JOIN pg_roles r ON r.oid = m.roleid
    WHERE r.rolname = $1`
	rows, err := db.DBConn.Query(query, role)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []service.User
	for rows.Next() {
		var user service.User
		if err := rows.Scan(&user.Username); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, rows.Err()
}
