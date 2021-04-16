package redisobs

import (
	"context"
	"fmt"
	"strings"

	"github.com/go-redis/redis/v7"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

// NewHook returns redis.Hook
func NewHook(cl *redis.Client) redis.Hook {
	return &redisHook{redisOpts: cl.Options()}
}

type redisHook struct {
	redisOpts *redis.Options
}

func (h *redisHook) BeforeProcess(ctx context.Context, cmd redis.Cmder) (context.Context, error) {
	nCtx := h.log(ctx, "redis-command", []redis.Cmder{cmd})
	return nCtx, nil
}

func (h *redisHook) AfterProcess(ctx context.Context, cmd redis.Cmder) error {
	opentracing.SpanFromContext(ctx).Finish()
	return nil
}

func (h *redisHook) BeforeProcessPipeline(ctx context.Context, cmds []redis.Cmder) (context.Context, error) {
	nCtx := h.log(ctx, "redis-pipeline", cmds)
	return nCtx, nil
}

func (h *redisHook) AfterProcessPipeline(ctx context.Context, cmds []redis.Cmder) error {
	opentracing.SpanFromContext(ctx).Finish()
	return nil
}

func (h *redisHook) log(ctx context.Context, operationName string, cmds []redis.Cmder) context.Context {
	var name = cmds[0].Name()
	if operationName == "redis-pipeline" {
		name = "pipeline"
	}
	span, newCtx := opentracing.StartSpanFromContext(ctx, "redis "+name)
	span.LogFields(
		log.String("db.system", "redis"),
		log.String("db.host", h.redisOpts.Addr),
		log.Int("db.name", h.redisOpts.DB),
		log.String("db.operation", formatMethods(cmds)),
		log.String("db.statement", formatCmds(cmds)),
	)
	return newCtx
}

func formatCmds(cmds []redis.Cmder) string {
	s := ""
	for _, cmd := range cmds {
		for _, v := range cmd.Args() {
			s += fmt.Sprint(v) + " "
		}
		s = strings.TrimRight(s, " "+"\n")
	}
	return strings.Trim(s, "\n")
}

func formatMethods(cmds []redis.Cmder) string {
	cmdsAsDbMethods := make([]string, len(cmds))
	for i, cmd := range cmds {
		dbMethod := cmd.Name()
		cmdsAsDbMethods[i] = dbMethod
	}
	return strings.Join(cmdsAsDbMethods, " -> ")
}
