package gormagent

type DBType string

const (
	UNKNOWN DBType = "unknown"
	MYSQL   DBType = "mysql"
)

const (
	componentIDUnknown = 0
	componentIDMysql   = 5012
)

type Option func(*options)

type options struct {
	dbType      DBType
	peer        string
	componentID int32

	reportQuery bool
	reportParam bool
}

// WithSqlDBType set dbType option,
// dbType is used for parsing dsn string to peer address
// and setting componentID, if DB type is not support in DBType
// list, please use WithPeerAddr to set peer address manually
func WithSqlDBType(t DBType) Option {
	return func(o *options) {
		o.dbType = t
		o.setComponentID()
	}
}

// WithPeerAddr set the peer address to report
func WithPeerAddr(addr string) Option {
	return func(o *options) {
		o.peer = addr
	}
}

func WithQueryReport() Option {
	return func(o *options) {
		o.reportQuery = true
	}
}

func WithParamReport() Option {
	return func(o *options) {
		o.reportParam = true
	}
}

func (o *options) setComponentID() {
	switch o.dbType {
	case MYSQL:
		o.componentID = componentIDMysql
	default:
		o.componentID = componentIDUnknown
	}
}
