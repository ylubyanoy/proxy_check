package postgres

import (
	"context"
	"time"
	"ylubyanoy/proxy_check/internal/store"

	"github.com/jackc/pgx/v4"
)

// ProxyRepository ...
type ProxyRepository struct {
	pgstore *PGStore
}

// Update ...
func (r *ProxyRepository) Update(pID, pingCheck int, dateCheck time.Time) error {
	return r.pgstore.conn.QueryRow(
		context.Background(),
		"UPDATE proxy_api_proxy SET ping_check=$1, date_check=$2 WHERE id=$3 RETURNING id",
		pingCheck,
		dateCheck,
		pID,
	).Scan(&pID)
}

// Count ...
func (r *ProxyRepository) Count() (int, error) {
	var cnt int
	err := r.pgstore.conn.QueryRow(
		context.Background(),
		"SELECT count(id) FROM proxy_api_proxy",
	).Scan(&cnt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return 0, store.ErrRecordNotFound
		}
		return 0, err
	}
	return cnt, nil
}

// Find ...
func (r *ProxyRepository) Find(url string) ([]store.Proxy, error) {
	var proxies []store.Proxy

	rows, err := r.pgstore.conn.Query(
		context.Background(),
		"SELECT id, ip_addr, port, host_url, type_proxy, country, ping_check, date_check, item_datetime FROM proxy_api_proxy WHERE ip_addr=$1",
		url,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}
		return nil, err
	}
	for rows.Next() {
		var p store.Proxy
		err := rows.Scan(
			&p.ID,
			&p.IPAddr,
			&p.Port,
			&p.HostURL,
			&p.TypeProxy,
			&p.Country,
			&p.PingCheck,
			&p.DateCheck,
			&p.ItemDatetime,
		)
		if err != nil {
			return nil, err
		}
		proxies = append(proxies, p)
	}

	return proxies, nil
}

// GetProxy ...
func (r *ProxyRepository) GetList(limit, offset int) ([]store.Proxy, error) {
	var proxies []store.Proxy
	rows, err := r.pgstore.conn.Query(
		context.Background(),
		"SELECT id, ip_addr, port, host_url, type_proxy, country, ping_check, date_check, item_datetime FROM proxy_api_proxy LIMIT $1 OFFSET $2",
		limit,
		offset,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}
		return nil, err
	}
	for rows.Next() {
		var p store.Proxy
		err := rows.Scan(
			&p.ID,
			&p.IPAddr,
			&p.Port,
			&p.HostURL,
			&p.TypeProxy,
			&p.Country,
			&p.PingCheck,
			&p.DateCheck,
			&p.ItemDatetime,
		)
		if err != nil {
			return nil, err
		}
		proxies = append(proxies, p)
	}
	return proxies, nil
}

// // Create ...
// func (r *UserRepository) Create(u *models.User) error {
// 	if err := u.Validate(); err != nil {
// 		return err
// 	}

// 	if err := u.BeforeCreate(); err != nil {
// 		return err
// 	}

// 	return r.store.conn.QueryRow(
// 		context.Background(),
// 		"INSERT INTO users (username, encrypted_password) VALUES ($1, $2) RETURNING id",
// 		u.Username,
// 		u.EncryptedPassword,
// 	).Scan(&u.ID)
// }

// // Find ...
// func (r *UserRepository) Find(id int) (*models.User, error) {
// 	u := &models.User{}
// 	if err := r.store.conn.QueryRow(
// 		context.Background(),
// 		"SELECT id, username, encrypted_password FROM users WHERE id = $1",
// 		id,
// 	).Scan(
// 		&u.ID,
// 		&u.Username,
// 		&u.EncryptedPassword,
// 	); err != nil {
// 		if err == pgx.ErrNoRows {
// 			return nil, store.ErrRecordNotFound
// 		}
// 		return nil, err
// 	}
// 	return u, nil
// }

// // FindByUsername ...
// func (r *UserRepository) FindByUsername(username string) (*models.User, error) {
// 	u := &models.User{}
// 	if err := r.store.conn.QueryRow(
// 		context.Background(),
// 		"SELECT id, username, encrypted_password FROM users WHERE username = $1",
// 		username,
// 	).Scan(
// 		&u.ID,
// 		&u.Username,
// 		&u.EncryptedPassword,
// 	); err != nil {
// 		if err == sql.ErrNoRows {
// 			return nil, store.ErrRecordNotFound
// 		}
// 		return nil, err
// 	}
// 	return u, nil
// }
