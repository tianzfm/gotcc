-- filepath: /gotcc/gotcc/internal/persistence/migrations/def.sql
-- 事务流程定义
CREATE TABLE task_group_flow (
    id VARCHAR(64) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT DEFAULT NULL,
    flow_type VARCHAR(50) NOT NULL,
    `version` INT NOT NULL DEFAULT 1 COMMENT '版本号',
    defination JSON,
    is_active TINYINT NOT NULL DEFAULT 1 COMMENT '是否启用',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    create_user VARCHAR(100) NOT NULL,
    updated_user VARCHAR(100) NOT NULL,
    UNIQUE KEY `uk_name_ver` (`name`, `version`)
);

-- 事务流程实例表
CREATE TABLE task_group_instance (
    id VARCHAR(64) PRIMARY KEY,
    flow_id VARCHAR(64) NOT NULL,
    flow_type VARCHAR(50) NOT NULL,
    status ENUM('pending', 'running', 'success', 'failed', 'cancelled', 'rolling_back') DEFAULT 'pending',
    task_type VARCHAR(50) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    completed_at TIMESTAMP NULL,
    INDEX idx_type_status (task_type, status)
);

-- 任务表 - 增加类型相关配置
CREATE TABLE dist_task (
    id VARCHAR(64) PRIMARY KEY,
    group_id VARCHAR(64) NOT NULL,
    name VARCHAR(255) NOT NULL,
    type ENUM('rpc', 'local', 'mq', 'http', 'db', 'file') NOT NULL,
    subtype VARCHAR(50),
    status ENUM('pending', 'running', 'success', 'failed', 'cancelled', 'rollback_success', 'rollback_failed') DEFAULT 'pending',
    priority INT DEFAULT 0,
    config JSON,
    execution_context JSON,
    input_data JSON,
    output_data JSON,
    started_at TIMESTAMP NULL,
    completed_at TIMESTAMP NULL,
    error_message TEXT,
    error_stack TEXT,
    INDEX idx_group_status (group_id, status),
    INDEX idx_type_status (type, status),
    INDEX idx_status_execute (status, execute_at)
);

-- 重试策略配置表
CREATE TABLE retry_policy_config (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    policy_type ENUM('fixed', 'exponential', 'randomized', 'adaptive') DEFAULT 'fixed',
    base_interval INT DEFAULT 60,
    max_interval INT DEFAULT 3600,
    max_attempts INT DEFAULT 3,
    multiplier FLOAT DEFAULT 2.0,
    randomization_factor FLOAT DEFAULT 0.1,
    enabled BOOLEAN DEFAULT TRUE,
    config JSON,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 异常记录表
CREATE TABLE exception_record (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    group_id VARCHAR(64) NOT NULL,
    group_name VARCHAR(255) NOT NULL,
    task_id VARCHAR(64) NOT NULL,
    task_name VARCHAR(255) NOT NULL,
    error_type INT NOT NULL,
    error_code VARCHAR(100),
    error_message TEXT,
    stack_trace TEXT,
    occurred_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    handled BOOLEAN DEFAULT FALSE,
    retry_times INT DEFAULT 0,
    last_retry_at TIMESTAMP NULL,
    INDEX idx_task_group (task_id, group_id),
    INDEX idx_occurred_handled (occurred_at, handled)
);

-- 执行日志表
CREATE TABLE execution_log (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    task_id VARCHAR(64) NOT NULL,
    group_id VARCHAR(64) NOT NULL,
    action ENUM('execute', 'retry', 'rollback', 'cancel', 'timeout') NOT NULL,
    old_status VARCHAR(50),
    new_status VARCHAR(50),
    message TEXT,
    details JSON,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_task_action (task_id, action)
);