package cron

import (
	"context"
	"fmt"
	"time"

	"github.com/robfig/cron/v3"
)

type chronos struct {
	*cron.Cron
}

func NewCron() *chronos {
	return &chronos{
		Cron: cron.New(cron.WithSeconds()),
	}
}

func (c *chronos) Name() string {
	return "cron"
}

func (c *chronos) Run(ctx context.Context) error {
	// 添加一个任务，每10分钟执行一次
	_, err := c.AddFunc("0 */10 * * * *", func() {
		fmt.Println("Cron job executed!", time.Now().Format(time.DateTime))
	})
	if err != nil {
		return err
	}

	c.Cron.Run()
	return nil
}

func (c *chronos) Close(ctx context.Context) error {
	c.Stop()
	return nil
}
