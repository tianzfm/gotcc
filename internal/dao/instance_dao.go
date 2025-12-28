package dao

import (
    "database/sql"
    "errors"
    "time"

    _ "github.com/go-sql-driver/mysql" // MySQL driver
)

type InstanceDAO struct {
    db *sql.DB
}

func NewInstanceDAO(db *sql.DB) *InstanceDAO {
    return &InstanceDAO{db: db}
}

func (dao *InstanceDAO) CreateInstance(instance *Instance) error {
    query := "INSERT INTO task_group_instance (id, flow_id, flow_type, status, task_type, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?)"
    _, err := dao.db.Exec(query, instance.ID, instance.FlowID, instance.FlowType, instance.Status, instance.TaskType, time.Now(), time.Now())
    return err
}

func (dao *InstanceDAO) GetInstanceByID(id string) (*Instance, error) {
    query := "SELECT id, flow_id, flow_type, status, task_type, created_at, updated_at, completed_at FROM task_group_instance WHERE id = ?"
    row := dao.db.QueryRow(query, id)

    var instance Instance
    err := row.Scan(&instance.ID, &instance.FlowID, &instance.FlowType, &instance.Status, &instance.TaskType, &instance.CreatedAt, &instance.UpdatedAt, &instance.CompletedAt)
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, errors.New("instance not found")
        }
        return nil, err
    }
    return &instance, nil
}

func (dao *InstanceDAO) UpdateInstance(instance *Instance) error {
    query := "UPDATE task_group_instance SET flow_type = ?, status = ?, task_type = ?, updated_at = ?, completed_at = ? WHERE id = ?"
    _, err := dao.db.Exec(query, instance.FlowType, instance.Status, instance.TaskType, time.Now(), instance.CompletedAt, instance.ID)
    return err
}

func (dao *InstanceDAO) DeleteInstance(id string) error {
    query := "DELETE FROM task_group_instance WHERE id = ?"
    _, err := dao.db.Exec(query, id)
    return err
}

type Instance struct {
    ID          string
    FlowID      string
    FlowType    string
    Status      string
    TaskType    string
    CreatedAt   time.Time
    UpdatedAt   time.Time
    CompletedAt *time.Time
}