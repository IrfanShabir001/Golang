package shared

import (
	"context"
	"fmt"
	uuid "github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
	"strings"
	"time"
)

type Log struct{}

func Fatal(ctx context.Context, args ...interface{}) context.Context {
	loggedString := trim(args)
	if v := ctx.Value(RequestId{}); v != nil {
		sugarPrint("fatal", v.(string), loggedString)
		return ctx
	}
	uid := uuid.NewV4()
	ctx = context.WithValue(ctx, RequestId{}, uid.String())
	ctx = context.WithValue(ctx, RequestTimestamp{}, time.Now().Unix())

	sugarPrint("fatal", uid.String(), loggedString)
	return ctx
}

func Warn(ctx context.Context, args ...interface{}) context.Context {
	loggedString := trim(args)
	if v := ctx.Value(RequestId{}); v != nil {
		sugarPrint("warn", v.(string), loggedString)
		return ctx
	}
	uid := uuid.NewV4()
	ctx = context.WithValue(ctx, RequestId{}, uid.String())
	ctx = context.WithValue(ctx, RequestTimestamp{}, time.Now().Unix())

	sugarPrint("warn", uid.String(), loggedString)
	return ctx
}

func Info(ctx context.Context, args ...interface{}) context.Context {
	loggedString := trim(args)
	if v := ctx.Value(RequestId{}); v != nil {
		sugarPrint("info", v.(string), loggedString)
		return ctx
	}
	uid := uuid.NewV4()
	ctx = context.WithValue(ctx, RequestId{}, uid.String())
	ctx = context.WithValue(ctx, RequestTimestamp{}, time.Now().Unix())

	sugarPrint("info", uid.String(), loggedString)
	return ctx
}

func Debug(ctx context.Context, args ...interface{}) context.Context {
	loggedString := trim(args)
	if v := ctx.Value(RequestId{}); v != nil {
		sugarPrint("debug", v.(string), loggedString)
		return ctx
	}
	uid := uuid.NewV4()
	ctx = context.WithValue(ctx, RequestId{}, uid.String())
	ctx = context.WithValue(ctx, RequestTimestamp{}, time.Now().Unix())

	sugarPrint("debug", uid.String(), loggedString)
	return ctx
}

func sugarPrint(level, traceId string, s string) {
	log.SetFormatter(&log.JSONFormatter{})
	switch level {
	case "info":
		log.WithFields(
			log.Fields{
				"traceId": traceId,
			},
		).Info(s)
	case "fatal":
		log.WithFields(
			log.Fields{
				"traceId": traceId,
			},
		).Error(s)
	case "warn":
		log.WithFields(
			log.Fields{
				"traceId": traceId,
			},
		).Warn(s)
	case "debug":
		log.WithFields(
			log.Fields{
				"traceId": traceId,
			},
		).Debug(s)
	}
	return
}

func trim(args ...interface{}) string {
	a := fmt.Sprint(args)
	a = strings.TrimPrefix(a, "[[")
	a = strings.TrimSuffix(a, "]]")
	return a
}
