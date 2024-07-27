package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/tursodatabase/libsql-client-go/libsql"

	"tender-scraper/config"
	"tender-scraper/types"
)

func ConnectDatabase(cfg *config.Config) (*sql.DB, error) {
	connectionString := fmt.Sprintf("%s?authToken=%s", cfg.DatabaseURL, cfg.AuthToken)
	db, err := sql.Open("libsql", connectionString)

	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	log.Println("Successfully connected to the database.")
	return db, nil
}

func AddUser(db *sql.DB, user types.UserInfo) error {

	query := `SELECT email FROM users WHERE email = ?`
	row := db.QueryRow(query, user.Email)

	var email string

	err := row.Scan(&email)
	if err == nil {
		fmt.Println("User with email %s already exists", user.Email)
		return fmt.Errorf("user with email %s already exists", user.Email)
	}

	query = `INSERT INTO users (email, name, company, position, phone, g7Company, existingG7Project, range, upcomingG7Project, description) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	_, err = db.Exec(query, user.Email, user.Name, user.Company, user.Position, user.Phone, user.G7Company, user.ExistingG7Project, user.Range, user.UpcomingG7Project, user.Description)
	if err != nil {
		return err
	}
	return nil
}

func BatchInsertTenders(db *sql.DB, tenders []types.Tender) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	sqlCmd := "INSERT INTO tenders (name, link, is_notified, kod_bidang, kebenaran_khas, hari_lawat_tapak, tarikh_iklan, taraf) VALUES (?, ?, ?, ?, ?, ?, ?, ?) ON CONFLICT(name) DO NOTHING"
	fmt.Println("SQL Command:", sqlCmd)

	stmt, err := tx.Prepare(sqlCmd)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, tender := range tenders {
		_, err = stmt.Exec(tender.Name, tender.Link, tender.IsNotified, tender.KodBidang, tender.KebenaranKhas, tender.TarikhLawatan, tender.TarikhIklan, tender.Taraf)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

func GetTenders(db *sql.DB) ([]types.Tender, error) {
	query := `SELECT id, name, link, is_notified, kod_bidang, kebenaran_khas, taraf FROM tenders WHERE is_notified = 0 LIMIT 10`
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tenders []types.Tender
	index := 0
	for rows.Next() {
		var tender types.Tender
		err := rows.Scan(&tender.ID, &tender.Name, &tender.Link, &tender.IsNotified, &tender.KodBidang, &tender.KebenaranKhas, &tender.Taraf)
		if err != nil {
			return nil, err
		}
		tender.Index = index + 1
		tenders = append(tenders, tender)
		index++
	}

	return tenders, nil
}

func GetClosestTarikhIklan(db *sql.DB) (string, error) {
	var closestTarikhIklan string

	query := `
	SELECT tarikh_iklan
	FROM tenders
	ORDER BY ABS(DATE(tarikh_iklan) - CURRENT_DATE)
	LIMIT 1;
	`

	db.QueryRow(query).Scan(&closestTarikhIklan)
	// if err != nil {
	// 	log.Printf("Error fetching closest tarikh iklan: %v", err)
	// 	return "", err
	// }

	return closestTarikhIklan, nil
}

// get users
func GetUsers(db *sql.DB) ([]types.UserInfo, error) {
	query := `SELECT name, email FROM users`
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []types.UserInfo
	for rows.Next() {
		var user types.UserInfo
		// Scan both name and email columns
		err := rows.Scan(&user.Name, &user.Email)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func MarkTendersAsNotified(db *sql.DB, tenderIDs []int) error {
	tx, err := db.BeginTx(context.Background(), &sql.TxOptions{})
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare("UPDATE tenders SET is_notified = 1 WHERE id = ?")
	if err != nil {
		tx.Rollback()
		return err
	}
	defer stmt.Close()

	for _, id := range tenderIDs {
		if _, err := stmt.Exec(id); err != nil {
			tx.Rollback()
			return err
		}
	}

	// Commit the transaction
	return tx.Commit()
}
