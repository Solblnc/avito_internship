package database

import (
	"context"
	"encoding/csv"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
	"strconv"
	"strings"
	"time"
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

	return nil
}

func (d *DataBase) AddUser(segmentsAdd []string, segmentsDelete []string, userId int) error {
	ctx := context.Background()
	conn, err := d.db.Acquire(ctx)
	if err != nil {
		log.Fatalf("Unable to acquire a database connection: %w", err)
		return err
	}
	defer conn.Release()
	tx, err := conn.Begin(ctx)
	if err != nil {
		log.Fatalf("Unable to start a transaction: %w", err)
		return err
	}

	query, values := d.addSegments(segmentsAdd, userId)

	if len(segmentsDelete) > 0 {
		query, values = d.deleteSegments(segmentsDelete, userId)
	}

	_, err = tx.Exec(ctx, query, values...)
	if err != nil {
		tx.Rollback(ctx)
		log.Fatalf("Unable to insert segments to users: %w", err)
		return err
	}
	err = tx.Commit(ctx)
	if err != nil {
		log.Fatalf("Unable to commit a transaction: %w", err)
		return err
	}

	return nil
}

type Segment struct {
	Id   int
	Name string
}

func (d *DataBase) GetActualSegments(userId int) []string {
	ctx := context.Background()
	conn, err := d.db.Acquire(ctx)
	if err != nil {
		log.Fatalf("Unable to acquire a database connection: %w", err)
	}
	defer conn.Release()

	row, err := conn.Query(ctx, "SELECT segments.segment_id, segments.segment_name FROM segments_user JOIN segments ON segments_user.segment_id = segments.segment_id  WHERE segments_user.user_id = $1", userId)
	if err != nil {
		log.Fatalf("Unable to get actual segments: %w", err)
	}

	var segments []string

	for row.Next() {
		var segment Segment
		err = row.Scan(&segment.Id, &segment.Name)
		if err != nil {
			log.Fatalf("cannot scan segment from query: %w", err)
		}
		segments = append(segments, segment.Name)
	}

	return segments

}

func (d *DataBase) addSegments(segmentsAdd []string, userId int) (string, []interface{}) {
	values := make([]interface{}, 0)
	placeholders := make([]string, 0)
	i := 1

	for _, segment := range segmentsAdd {
		segmentId, err := d.getSegmentID(segment)
		if err != nil {
			log.Fatal(err)
		}

		values = append(values, userId, segmentId)
		placeholders = append(placeholders, fmt.Sprintf("($%d, $%d, NOW())", i, i+1))
		i += 2
	}

	query := fmt.Sprintf("INSERT INTO segments_user (user_id, segment_id, time) VALUES %s", strings.Join(placeholders, ", "))

	return query, values
}
func (d *DataBase) deleteSegments(segmentsDelete []string, userId int) (string, []interface{}) {
	values := make([]interface{}, 0)
	placeholders := make([]string, 0)
	i := 1

	for _, segment := range segmentsDelete {
		segmentId, err := d.getSegmentID(segment)
		if err != nil {
			log.Fatal(err)
		}

		values = append(values, userId, segmentId)
		placeholders = append(placeholders, fmt.Sprintf("($%d, $%d)", i, i+1))

		i += 2
	}

	query := fmt.Sprintf("DELETE FROM segments_user WHERE (user_id, segment_id) IN (%s)", strings.Join(placeholders, ", "))

	return query, values
}

func (d *DataBase) CreateUser() error {
	ctx := context.Background()
	conn, err := d.db.Acquire(ctx)
	if err != nil {
		log.Fatalf("Unable to acquire a database connection: %w", err)
	}
	defer conn.Release()

	for i := 0; i <= 20; i++ {
		row := conn.QueryRow(ctx, "INSERT INTO users (user_id) values ($1) RETURNING user_id", i)
		var id int
		err = row.Scan(&id)
		if err != nil {
			log.Fatalf("Unable to scan id (creating users): %w", err)
		}
	}
	return nil
}

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

func (d *DataBase) GetHistory(year, month int) (string, error) {
	ctx := context.Background()
	conn, err := d.db.Acquire(ctx)
	if err != nil {
		log.Fatalf("Unable to acquire a database connection: %w", err)
	}
	defer conn.Release()

	row, err := conn.Query(ctx, "SELECT DISTINCT user_id, segments.segment_name, segments_user.time "+
		"FROM segments_user "+
		"JOIN segments ON segments.segment_id = segments_user.segment_id "+
		"WHERE EXTRACT(year FROM segments_user.time) = $1 "+
		"AND EXTRACT(month FROM segments_user.time) = $2;", year, month)
	if err != nil {
		log.Fatalf("Unable to get history data from database: %w", err)
	}

	var buf strings.Builder
	csvBuilder := csv.NewWriter(&buf)
	for row.Next() {
		var (
			userId      int
			segmentName string
			time        time.Time
		)

		err = row.Scan(&userId, &segmentName, &time)
		if err != nil {
			return "", err
		}

		err = csvBuilder.Write([]string{
			strconv.Itoa(userId),
			segmentName,
			time.String(),
		})
		if err != nil {
			return "", err
		}
	}
	csvBuilder.Flush()
	return buf.String(), nil

}
