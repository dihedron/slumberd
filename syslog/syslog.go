package syslog

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"os"
	"strings"
	"time"

	"github.com/juju/rfc/v2/rfc5424"
)

const (
	DefaultEnterprise  = "dihedron"
	DefaultSendTimeout = time.Duration(10 * time.Second)
	DefaultSendMaxSize = 0
)

// Option is a functional option type that allows us to configure the Syslog.
type Option func(*Syslog)

// WithApplication allows to specify the name of the application that is
// going to send events to the syslog.
func WithApplication(application string) Option {
	return func(sl *Syslog) {
		if application != "" {
			sl.application = application
		}
	}
}

// WithEnterprise allows to specify the name of the enterprise that is
// responsible for the application sending events to the syslog.
func WithEnterprise(enterprise string) Option {
	return func(sl *Syslog) {
		if enterprise != "" {
			sl.enterprise = enterprise
		}
	}
}

// WithProcess allows to specify the name of the process that is
// sending events to the syslog.
func WithProcess(process string) Option {
	return func(sl *Syslog) {
		if process != "" {
			sl.process = process
		}
	}
}

// Syslog wraps a syslog connection and stores all common
// configuration elements.
type Syslog struct {
	application string
	hostname    string
	enterprise  string
	process     string
	client      *rfc5424.Client
}

// New creates a new Syslog, initialising all relevant fields
// to the defaults unless options are provided; if application is
// not provided, it defaults to os.Args[0]; if enterprise is not
// provided, it defaults to DefaultEnterprise, which is set to
// "dihedron"; if process is not provided, it defaults to the
// current PID.
func New(options ...Option) (syslog *Syslog, err error) {
	hostname := ""
	if hostname, err = os.Hostname(); err != nil {
		slog.Error("error retrieving hostname", "error", err)
		return nil, err
	}

	syslog = &Syslog{
		application: os.Args[0],
		enterprise:  DefaultEnterprise,
		process:     fmt.Sprintf("%d", os.Getpid()),
		hostname:    hostname,
	}
	// apply functional options
	for _, option := range options {
		option(syslog)
	}

	syslog.client, err = rfc5424.Open("/dev/log", rfc5424.ClientConfig{
		MaxSize:     DefaultSendMaxSize,
		SendTimeout: DefaultSendTimeout,
	}, func(network, address string) (rfc5424.Conn, error) {
		return net.Dial("unixgram", address)
	})
	if err != nil {
		slog.Error("error opening syslog client", "error", err)
		return nil, err
	}
	return
}

// Send prepares a message in RFC5424-compliant format and
// sends it through the client.
func (s *Syslog) Send(message *Message) error {
	msg := rfc5424.Message{
		Header: rfc5424.Header{
			Priority: rfc5424.Priority{
				Severity: message.Severity,
				Facility: message.Facility,
			},
			Timestamp: rfc5424.Timestamp{Time: time.Now()},
			Hostname:  rfc5424.Hostname{FQDN: s.hostname},
			AppName:   rfc5424.AppName(s.application),
			ProcID:    rfc5424.ProcID(s.process),
			MsgID:     rfc5424.MsgID(message.ID),
		},
		StructuredData: rfc5424.StructuredData{},
	}
	if s, ok := (message.Content).(string); ok {
		msg.Msg = s
	} else if message.Content != nil {
		if s, err := json.Marshal(message.Content); err != nil {
			return err
		} else {
			msg.Msg = string(s)
		}
	} else {
		return errors.New("no text in message")
	}
	for k, v := range message.Data {
		msg.StructuredData = append(msg.StructuredData, s.newStructuredDataElement(k, v...))
	}
	return s.client.Send(msg)
}

// Message contains the set of information that is specific
// to a given message: its priority (in terms of facility and
// severity), the message type identifier (to group similar
// messages, the message text and optionally a set of
// parameters.
type Message struct {
	Facility rfc5424.Facility
	Severity rfc5424.Severity
	ID       string
	Content  any // either a string or an object that will be marshalled to JSON
	Data     map[string][]string
}

// Element is an implementation of the StructuredDataElement interface.
type Element struct {
	id     rfc5424.StructuredDataName
	params []rfc5424.StructuredDataParam
}

func (s *Syslog) newStructuredDataElement(id string, parameters ...string) *Element {
	params := make([]rfc5424.StructuredDataParam, len(parameters))
	for i, str := range parameters {
		parts := strings.SplitN(str, "=", 2)
		params[i].Name = rfc5424.StructuredDataName(parts[0])
		params[i].Value = rfc5424.StructuredDataParamValue(parts[1])
	}
	return &Element{
		id:     rfc5424.StructuredDataName(fmt.Sprintf("%s@%s", id, s.enterprise)),
		params: params,
	}
}

func (e *Element) ID() rfc5424.StructuredDataName {
	return e.id
}

func (e *Element) Params() []rfc5424.StructuredDataParam {
	return e.params
}

func (e *Element) Validate() error {
	return nil
}
