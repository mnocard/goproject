package task

import (
	"context"
	"errors"

	pg "github.com/mnocard/go-project/internal/storage"
)

type TaskStorage interface {
	CreateTask(context.Context, *pg.Task) (int, error)
	FindTaskByUserId(context.Context, int) (*pg.Task, error)
	FindTaskById(context.Context, int) (*pg.Task, error)
	UpdateTask(context.Context, pg.Task) (*pg.Task, error)
	DeleteTask(context.Context, int) (bool, error)
}

type Task struct {
	Id           int  `json:"-"`
	UserId       int  `json:"user_id" binding:"required"`
	Points       int  `json:"points"`
	ParentTaskId int  `json:"parent_task_id"`
	IsCompleted  bool `json:"is_completed" default:"false"`
}

type taskService struct {
	tStorage TaskStorage
}

func New(s TaskStorage) *taskService {
	return &taskService{tStorage: s}
}

func (tService *taskService) Create(ctx context.Context, t *Task) (int, error) {
	return tService.tStorage.CreateTask(ctx, &pg.Task{
		UserId:       t.UserId,
		Points:       t.Points,
		ParentTaskId: t.ParentTaskId,
		IsCompleted:  t.IsCompleted,
	})
}

func (tService *taskService) Get(ctx context.Context, id int) (*Task, error) {
	task, err := tService.tStorage.FindTaskById(ctx, id)
	if err != nil {
		return nil, err
	}

	return &Task{
		Id:           task.Id,
		UserId:       task.UserId,
		Points:       task.Points,
		ParentTaskId: task.ParentTaskId,
		IsCompleted:  task.IsCompleted,
	}, err
}

func (tService *taskService) FindTaskByUserId(ctx context.Context, userId int) (*Task, error) {
	task, err := tService.tStorage.FindTaskByUserId(ctx, userId)
	if err != nil {
		return nil, err
	}

	return &Task{
		Id:           task.Id,
		UserId:       task.UserId,
		Points:       task.Points,
		ParentTaskId: task.ParentTaskId,
		IsCompleted:  task.IsCompleted,
	}, err
}

func (tService *taskService) Update(ctx context.Context, t *Task) (int, error) {
	task, err := tService.tStorage.UpdateTask(ctx, pg.Task{
		UserId:       t.UserId,
		Points:       t.Points,
		ParentTaskId: t.ParentTaskId,
		IsCompleted:  t.IsCompleted,
	})

	if err != nil {
		return 0, err
	}

	return task.Id, err
}

func (tService *taskService) Delete(ctx context.Context, id int) (bool, error) {
	task, err := tService.Get(ctx, id)
	if err != nil {
		return false, err
	}

	return tService.tStorage.DeleteTask(ctx, task.Id)
}

func (tService *taskService) CompleteTask(ctx context.Context, userId int, taskId int) (int, error) {
	task, err := tService.tStorage.FindTaskById(ctx, taskId)
	if err != nil {
		return 0, err
	}

	if task.UserId != userId {
		return 0, errors.New("it is not your task!")
	}

	if task.IsCompleted {
		return 0, errors.New("task is already completed!")
	}

	task.IsCompleted = true
	tService.tStorage.UpdateTask(ctx, *task)

	return task.Points, nil
}
