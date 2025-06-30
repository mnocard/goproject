package storage

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5"
)

type Task struct {
	Id           int  `json:"id"`
	UserId       int  `json:"user_id"`
	Points       int  `json:"points"`
	ParentTaskId int  `json:"parent_task_id"`
	IsCompleted  bool `json:"is_completed"`
}

func (s *storage) CreateTask(ctx context.Context, t *Task) (int, error) {
	sql := `
        INSERT INTO tasks (user_id, points, parent_task_id, is_completed)
        VALUES ($1, $2, $3, $4)
        RETURNING id;
	`

	var id int
	err := s.pool.QueryRow(ctx, sql, t.UserId, t.Points, t.ParentTaskId, t.IsCompleted).Scan(&id)
	if err != nil {
		s.logError(err, "storage) CreateTask")
		return 0, err
	}

	return id, nil
}

func (s *storage) FindTaskByUserId(ctx context.Context, tUserId int) (*Task, error) {
	var task Task

	sql := `
    	SELECT id, user_id, points, parent_task_id, is_completed FROM tasks
    	WHERE user_id = $1;
	`

	err := s.pool.QueryRow(ctx, sql, tUserId).Scan(&task.Id, &task.UserId, &task.Points, &task.ParentTaskId, &task.IsCompleted)
	if err != nil {
		s.logError(err, "storage) FindTaskByUserId")
		return nil, err
	}

	log.Println("storage) FindTaskByUserId task:", task)
	return &task, nil
}

func (s *storage) FindTaskById(ctx context.Context, id int) (*Task, error) {
	var task Task

	sql := `
    	SELECT id, user_id, points, parent_task_id, is_completed FROM tasks
    	WHERE id = $1;
	`

	err := s.pool.QueryRow(ctx, sql, id).Scan(&task.Id, &task.UserId, &task.Points, &task.ParentTaskId, &task.IsCompleted)
	if err != nil {
		s.logError(err, "storage) FindTaskById")
		return nil, err
	}

	log.Println("storage) FindTaskById user:", task)
	return &task, nil
}

func (s *storage) UpdateTask(ctx context.Context, t Task) (*Task, error) {
	var task *Task
	var err error

	if t.Id != 0 {
		task, err = s.FindTaskById(ctx, t.Id)
	} else {
		task, err = s.FindTaskByUserId(ctx, t.UserId)
	}

	if err != nil && err != pgx.ErrNoRows {
		s.logError(err, "storage) UpdateTask error")
		return nil, err
	}

	if task == nil {
		return nil, &NotFoundError{}
	}

	sql := `
    	UPDATE tasks
		SET user_id = $2, points = $3, parent_task_id = $4, is_completed = $5
    	WHERE id = $1
		RETURNING id;
	`

	commandTag, err := s.pool.Exec(ctx, sql, task.Id, t.UserId, t.Points, t.ParentTaskId, t.IsCompleted)
	if err != nil {
		s.logError(err, "storage) UpdateTask error")
		return nil, err
	}

	if commandTag.RowsAffected() == 0 {
		err = &NotFoundError{}
		s.logError(err, "storage) UpdateTask not found")
		return nil, err
	}

	log.Println("storage) UpdateTask user:", task)
	return task, nil
}

func (s *storage) DeleteTask(ctx context.Context, id int) (bool, error) {
	sql := `DELETE FROM tasks WHERE id = $1`

	commandTag, err := s.pool.Exec(ctx, sql, id)
	if err != nil {
		s.logError(err, "storage) DeleteTask error")
		return false, err
	}

	if commandTag.RowsAffected() == 0 {
		err = &NotFoundError{}
		s.logError(err, "storage) DeleteTask not found")
		return false, err
	}

	return true, nil
}
