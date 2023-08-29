package database

import (
	"context"
	"fmt"
	"log"
	"time"
)

type segmentUser struct {
	ttl         int
	userID      int
	segmentID   int
	segmentName string
	deadline    time.Time
}

var timer = make(chan segmentUser)

func (d *DataBase) Schedule() {
	go func() {
		for task := range timer {
			if err := d.DeleteSegmentDeadline(task.userID, task.segmentID); err != nil {
				log.Println(err)
			}
		}
	}()
	d.CheckDeadlines()

}

func (d *DataBase) CheckDeadlines() {
	ctx := context.Background()
	conn, err := d.db.Acquire(ctx)
	if err != nil {
		log.Fatalf("Unable to acquire a database connection: %w", err)
	}
	defer conn.Release()

	row, err := conn.Query(ctx, "SELECT user_id, segment_id,deadline FROM segments_user;")
	if err != nil {
		log.Fatalf("Unable to get actual segments: %w", err)
	}

	var deadlines []segmentUser

	now := time.Now()

	for row.Next() {
		var deadline segmentUser
		err = row.Scan(&deadline.userID, &deadline.segmentID, &deadline.deadline)
		if err != nil {
			log.Fatalf("cannot scan segment from query: %w", err)
		}
		if now.After(deadline.deadline) {
			deadline.ttl = -1
		} else {
			deadline.ttl = int(deadline.deadline.Sub(now).Seconds())
		}
		deadlines = append(deadlines, deadline)
	}

	for _, task := range deadlines {
		if task.ttl <= 0 {
			if err = d.DeleteSegmentDeadline(task.userID, task.segmentID); err != nil {
				log.Fatal(err)
			} else {
				if err = d.AddSegmentDeadline(task.ttl, task.segmentName, task.userID); err != nil {
					log.Fatal(err)
				}
			}
		}
	}

}

func (d *DataBase) AddSegmentDeadline(ttl int, segmentName string, userId int) error {
	ctx := context.Background()
	conn, err := d.db.Acquire(ctx)
	if err != nil {
		log.Fatalf("Unable to acquire a database connection: %w", err)
		return err
	}
	defer conn.Release()

	segmentId, err := d.getSegmentID(segmentName)
	if err != nil {
		log.Fatal(err)
		return err
	}

	//deadline := time.Now().AddDate(0, 0, ttl)
	deadline := time.Now().Add(time.Duration(ttl) * time.Minute)

	query := fmt.Sprintf("INSERT INTO segments_user (user_id, segment_id, deadline) VALUES ($1,$2,$3)")

	_, err = conn.Exec(ctx, query, userId, segmentId, deadline)
	if err != nil {
		log.Fatalf("Unable to insert data to segements_user_deadline: %w", err)
		return err
	}

	timer <- segmentUser{
		ttl:       ttl,
		userID:    userId,
		segmentID: segmentId,
	}

	return nil
}

func (d *DataBase) DeleteSegmentDeadline(userId, segmentId int) error {
	ctx := context.Background()
	conn, err := d.db.Acquire(ctx)
	if err != nil {
		log.Fatalf("Unable to acquire a database connection: %w", err)
		return err
	}
	defer conn.Release()

	_, err = conn.Exec(ctx, "DELETE FROM segments_user WHERE user_id = $1 AND segment_id = $2", userId, segmentId)
	if err != nil {
		log.Fatalf("Unable to delete segment from segments_user table: %w", err)
		return err
	}
	return nil
}
