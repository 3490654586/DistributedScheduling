package common

import "errors"

const (
	JOB_SAVE_DIR = "/cron/jobs/"
	JOB_KILL_DIR = "/cron/kill/"
	JOB_LOCK_DIR =  "/cron/lock/"
	// 保存任务事件
	JOB_EVENT_SAVE = 1

	// 删除任务事件
	JOB_EVENT_DELETE = 2

	// 强杀任务事件
	JOB_EVENT_KILL = 3

)

var ERR_LOCK error = errors.New("锁被占用")
