package database

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
	"strings"
)

type Config struct {
	Host     string `mapstructure:"DB_HOST"`
	Port     string `mapstructure:"DB_PORT"`
	User     string `mapstructure:"DB_USER"`
	DBName   string `mapstructure:"DB_NAME"`
	Password string `mapstructure:"DB_PASSWORD"`
	SSLMode  string `mapstructure:"SSL_MODE"`
}

type Repository interface {
	Create(segment string) (int, error)
	Delete(segment string) error
	AddUser(segmentsAdd []string, segmentsDelete []string, userId int) error
	GetSegments(userId int) []string
}

type DataBase struct {
	db *pgxpool.Pool
}

func NewDataBase(cfg Config) (*DataBase, error) {
	connectionString := fmt.Sprintf(
		"host=%s port=%s user=%s dbname=%s password=%s sslmode=%s", cfg.Host, cfg.Port, cfg.User, cfg.DBName, cfg.Password, cfg.SSLMode,
	)

	db, err := pgxpool.Connect(context.Background(), connectionString)
	if err != nil {
		return &DataBase{}, fmt.Errorf("could not connect to database: %w", err)
	}

	err = db.Ping(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	conn, err := db.Acquire(context.Background())
	if err != nil {
		log.Fatalf("Unable to acquire a database connection: %v\n", err)
	}
	Migrate(conn.Conn(), context.Background())
	conn.Release()

	return &DataBase{db: db}, nil
}

func (d *DataBase) Create(segment string) (int, error) {
	ctx := context.Background()
	conn, err := d.db.Acquire(ctx)
	if err != nil {
		log.Fatalf("Unable to acquire a database connection: %w", err)
	}
	defer conn.Release()

	row := conn.QueryRow(ctx, "INSERT INTO segments (segment_name) values ($1) RETURNING segment_id", segment)

	var id int
	err = row.Scan(&id)
	if err != nil {
		log.Fatalf("Unable to scan id (segments): %w", err)
	}

	return id, nil

}

func (d *DataBase) Delete(segment string) error {
	ctx := context.Background()
	conn, err := d.db.Acquire(ctx)
	if err != nil {
		log.Fatalf("Unable to acquire a database connection: %w", err)
	}
	defer conn.Release()

	_, err = conn.Exec(ctx, "DELETE FROM segments_user WHERE segment_id = (SELECT segment_id FROM segments WHERE segment_name = $1 LIMIT 1)", segment)
	if err != nil {
		if err != nil {
			log.Fatalf("Unable to delete segment from segments_user table: %w", err)
		}
	}

	_, err = conn.Exec(ctx, "  DELETE FROM segments WHERE segment_name = $1", segment)
	if err != nil {
		if err != nil {
			log.Fatalf("Unable to delete segment segments table: %w", err)
		}
	}
	//if ct.RowsAffected() == 0 {
	//	log.Fatal(err)
	//}

	return nil
}

func (d *DataBase) AddUser(segmentsAdd []string, segmentsDelete []string, userId int) error {
	ctx := context.Background()
	conn, err := d.db.Acquire(ctx)
	if err != nil {
		log.Fatalf("Unable to acquire a database connection: %w", err)
	}
	defer conn.Release()

	tx, err := conn.Begin(ctx)
	if err != nil {
		log.Fatalf("Unable to start a transaction: %w", err)
	}

	values := make([]interface{}, 0)
	placeholders := make([]string, 0)
	i := 1

	for _, segment := range segmentsAdd {
		segmentId, err := d.getSegmentID(segment)
		if err != nil {
			tx.Rollback(ctx)
			return err
		}

		values = append(values, userId, segmentId)
		placeholders = append(placeholders, fmt.Sprintf("($%d, $%d, NOW())", i, i+1))
		i += 2
	}

	query := fmt.Sprintf("INSERT INTO segments_user (user_id, segment_id, time) VALUES %s", strings.Join(placeholders, ", "))

	_, err = conn.Exec(ctx, query, values...)
	if err != nil {
		log.Fatalf("Unable to insert segments to users: %w", err)
	}
	err = tx.Commit(ctx)
	if err != nil {
		log.Fatalf("Unable to commit a transaction: %w", err)
	}

	return nil
}

func (d *DataBase) CreateUser() error {
	ctx := context.Background()
	conn, err := d.db.Acquire(ctx)
	if err != nil {
		log.Fatalf("Unable to acquire a database connection: %w", err)
	}
	defer conn.Release()

	for i := 0; i <= 10; i++ {
		row := conn.QueryRow(ctx, "INSERT INTO users (user_id) values ($1) RETURNING user_id", i)
		var id int
		err = row.Scan(&id)
		if err != nil {
			log.Fatalf("Unable to scan id (creating users): %w", err)
		}
	}
	return nil
}

// Функция для получения ID сегмента по его названию
func (d *DataBase) getSegmentID(segmentName string) (int, error) {
	var segmentID int
	ctx := context.Background()
	conn, err := d.db.Acquire(ctx)
	if err != nil {
		log.Fatalf("Unable to acquire a database connection: %w", err)
	}
	defer conn.Release()

	conn.QueryRow(ctx, "SELECT segment_id FROM segments WHERE segment_name = $1", segmentName).Scan(&segmentID)

	return segmentID, nil
}
