package gormagent

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/SkyAPM/go2sky"
	"github.com/phprao/go-skywalking.git/tracerhelper"
	"github.com/phprao/go-skywalking.git/tracerhelper/util"
	"gorm.io/gorm"
	agentV3 "skywalking.apache.org/repo/goapi/collect/language/agent/v3"
)

var (
	_ gorm.Plugin = &SkyWalking{}
)

const spanKey = "spanKey"

type SkyWalking struct {
	tracer *go2sky.Tracer
	gcm    *util.GoroutineContextManager
	opts   *options
}

func SetGormPlugin(peerServiceName string) *SkyWalking {
	return NewPlugin(
		WithSqlDBType(MYSQL),
		WithQueryReport(),
		WithParamReport(),
		WithPeerAddr(peerServiceName),
	)
}

func NewPlugin(opts ...Option) *SkyWalking {
	options := &options{
		dbType:      UNKNOWN,
		componentID: componentIDUnknown,
		peer:        "unknown",
		reportQuery: false,
		reportParam: false,
	}

	for _, o := range opts {
		o(options)
	}

	return &SkyWalking{
		tracer: tracerhelper.GetTracer(),
		gcm:    tracerhelper.GetGcm(),
		opts:   options,
	}
}

func (s *SkyWalking) Name() string {
	return "gorm:skyWalking"
}

func (s *SkyWalking) Initialize(db *gorm.DB) (err error) {
	// before database operation
	db.Callback().Create().Before("gorm:create").Register("sky_create_span", s.BeforeCallback("create"))
	db.Callback().Query().Before("gorm:query").Register("sky_create_span", s.BeforeCallback("query"))
	db.Callback().Update().Before("gorm:update").Register("sky_create_span", s.BeforeCallback("update"))
	db.Callback().Delete().Before("gorm:delete").Register("sky_create_span", s.BeforeCallback("delete"))
	db.Callback().Row().Before("gorm:row").Register("sky_create_span", s.BeforeCallback("row"))
	db.Callback().Raw().Before("gorm:raw").Register("sky_create_span", s.BeforeCallback("raw"))

	// after database operation
	db.Callback().Create().After("gorm:create").Register("sky_end_span", s.AfterCallback())
	db.Callback().Query().After("gorm:query").Register("sky_end_span", s.AfterCallback())
	db.Callback().Update().After("gorm:update").Register("sky_end_span", s.AfterCallback())
	db.Callback().Delete().After("gorm:delete").Register("sky_end_span", s.AfterCallback())
	db.Callback().Row().After("gorm:row").Register("sky_end_span", s.AfterCallback())
	db.Callback().Raw().After("gorm:raw").Register("sky_end_span", s.AfterCallback())

	return
}

func (s *SkyWalking) BeforeCallback(operation string) func(db *gorm.DB) {
	tracer := s.tracer
	peer := s.opts.peer

	if tracer == nil {
		return func(db *gorm.DB) {}
	}

	return func(db *gorm.DB) {
		tableName := db.Statement.Table
		operation := fmt.Sprintf("%s/%s", tableName, operation)

		ctx, ok := s.gcm.GetContext()
		if !ok {
			return
		}
		span, err := tracer.CreateExitSpan(*ctx, operation, peer, func(key, value string) error {
			return nil
		})
		if err != nil {
			db.Logger.Error(db.Statement.Context, "gorm:skyWalking failed to create exit span, got error: %v", err)
			return
		}

		// set span from db instance's context to pass span
		db.Set(spanKey, span)
	}
}

func (s *SkyWalking) AfterCallback() func(db *gorm.DB) {
	tracer := s.tracer
	if tracer == nil {
		return func(db *gorm.DB) {}
	}

	return func(db *gorm.DB) {
		// get span from db instance's context
		spanInterface, _ := db.Get(spanKey)
		span, ok := spanInterface.(go2sky.Span)
		if !ok {
			return
		}

		defer span.End()

		sql := db.Statement.SQL.String()
		vars := db.Statement.Vars
		err := db.Statement.Error

		span.SetComponent(s.opts.componentID)
		span.SetSpanLayer(agentV3.SpanLayer_Database)
		span.Tag(go2sky.TagDBType, string(s.opts.dbType))
		span.Tag(go2sky.TagDBInstance, s.opts.peer)
		span.Tag("component", "mysql")

		if s.opts.reportQuery {
			span.Tag(go2sky.TagDBStatement, sql)
		}
		if s.opts.reportParam && len(vars) != 0 {
			span.Tag(go2sky.TagDBSqlParameters, argsToString(vars))
		}

		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			span.Error(time.Now(), err.Error())
		}
	}
}

func argsToString(args []interface{}) string {
	sb := strings.Builder{}

	switch len(args) {
	case 0:
		return ""
	case 1:
		return fmt.Sprintf("%v", args[0])
	}

	sb.WriteString(fmt.Sprintf("%v", args[0]))
	for _, arg := range args[1:] {
		sb.WriteString(fmt.Sprintf(", %v", arg))
	}
	return sb.String()
}
