-- 事务流程定义
CREATE TABLE task_group_flow (
    -- 事务流定义ID
    id VARCHAR(64) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT DEFAULT NULL
    -- 事务流类型 1. 处置任务 2. 处置业务
    flow_type VARCHAR(50) NOT NULL, 
    `version` INT NOT NULL DEFAULT 1 COMMENT '版本号',
    -- 元数据 事务流定义
    defination JSON, 
    is_active TINYINT NOT NULL DEFAULT 1 COMMENT '是否启用',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    create_user VARCHAR(100) NOT NULL,
    updated_user VARCHAR(100) NOT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_name_ver` (`name`, `version`)
);

-- 事务流程实例表
CREATE TABLE task_group_instance (
    id VARCHAR(64) PRIMARY KEY,
    -- 关联的事务流定义ID
    flow_id VARCHAR(64) NOT NULL,
    flow_type VARCHAR(50) NOT NULL, 
    status ENUM('pending', 'running', 'success', 'failed', 'cancelled', 'rolling_back') DEFAULT 'pending',
    task_type VARCHAR(50) NOT NULL, -- 任务类型
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    completed_at TIMESTAMP NULL,
    INDEX idx_type_status (task_type, status)
);

-- 任务表 - 增加类型相关配置
CREATE TABLE dist_task (
    id VARCHAR(64) PRIMARY KEY,
    -- 关联的task_group_instance ID
    group_id VARCHAR(64) NOT NULL,
    name VARCHAR(255) NOT NULL,
    type ENUM('rpc', 'local', 'mq', 'http', 'db', 'file') NOT NULL,
    subtype VARCHAR(50), -- 子类型，如: kafka/rabbitmq, grpc/http
    status ENUM('pending', 'running', 'success', 'failed', 'cancelled', 'rollback_success', 'rollback_failed') DEFAULT 'pending',
    -- 任务优先级, 默认按json定义中的顺序执行
    priority INT DEFAULT 0,
    -- 重试次数记录在异常表中即可？
    -- max_retry INT DEFAULT 3,
    -- retry_count INT DEFAULT 0,
    -- 从json中解析出来的任务配置
    config JSON, -- 任务配置
    execution_context JSON, -- 执行上下文
    input_data JSON, -- 输入参数
    output_data JSON, -- 输出结果
    started_at TIMESTAMP NULL,
    completed_at TIMESTAMP NULL,
    error_message TEXT,
    error_stack TEXT, -- 错误堆栈
    INDEX idx_group_status (group_id, status),
    INDEX idx_type_status (type, status),
    INDEX idx_status_execute (status, execute_at)
);

-- 重试策略配置表
CREATE TABLE retry_policy_config (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    policy_type ENUM('fixed', 'exponential', 'randomized', 'adaptive') DEFAULT 'fixed',
    base_interval INT DEFAULT 60, -- 基础间隔(秒)
    max_interval INT DEFAULT 3600, -- 最大间隔
    max_attempts INT DEFAULT 3,
    multiplier FLOAT DEFAULT 2.0, -- 指数因子
    randomization_factor FLOAT DEFAULT 0.1, -- 随机因子
    enabled BOOLEAN DEFAULT TRUE,
    config JSON, -- 额外配置
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- -- 回滚策略配置表
-- CREATE TABLE rollback_policy_config (
--     id INT AUTO_INCREMENT PRIMARY KEY,
--     name VARCHAR(100) NOT NULL,
--     policy_type ENUM('none', 'auto', 'manual', 'conditional') DEFAULT 'none',
--     conditions JSON, -- 回滚条件
--     rollback_handler VARCHAR(255), -- 回滚处理器
--     enabled BOOLEAN DEFAULT TRUE,
--     config JSON,
--     created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
-- );

-- 异常记录表
CREATE TABLE exception_record (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    group_id VARCHAR(64) NOT NULL,
    group_name VARCHAR(255) NOT NULL,
    task_id VARCHAR(64) NOT NULL,
    task_name VARCHAR(255) NOT NULL,
    -- 1. http_error 2. rpc_error 3. db_error 4. file_error 5. mq_error 6. system_error 7. business_error
    error_type INT NOT NULL,
    error_code VARCHAR(100),
    error_message TEXT,
    stack_trace TEXT,
    occurred_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    handled BOOLEAN DEFAULT FALSE,
    retry_times INT DEFAULT 0,
    last_retry_at TIMESTAMP NULL,
    -- retryable BOOLEAN DEFAULT TRUE,
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