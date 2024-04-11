package cron

import (
	"context"
	"fmt"
	"time"

	"github.com/robfig/cron/v3"
)

func Run(ctx context.Context) error {
	c := cron.New(cron.WithSeconds())

	// 添加一个任务，每10分钟执行一次
	_, err := c.AddFunc("0 */10 * * * *", func() {
		fmt.Println("Cron job executed!", time.Now().Format(time.DateTime))
	})
	if err != nil {
		return err
	}

	c.Start()
	return nil
}
