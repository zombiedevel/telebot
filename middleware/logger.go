package middleware

import (
	"log"
	"os"

	tele "gopkg.in/tucnak/telebot.v3"
)

type LoggerConfig interface {
	Context(c tele.Context)
	Message(m *tele.Message)
	Callback(c *tele.Callback)
	Query(q *tele.Query)
	ChosenInlineResult(r *tele.ChosenInlineResult)
	ShippingQuery(s *tele.ShippingQuery)
	PreCheckoutQuery(pre *tele.PreCheckoutQuery)
	Poll(poll *tele.Poll)
	PollAnswer(pa *tele.PollAnswer)
}

func Logger(v ...LoggerConfig) tele.MiddlewareFunc {
	var logger LoggerConfig

	if len(v) == 0 {
		logger = DefaultLogger{Logger: log.New(os.Stdout, log.Prefix(), log.Flags())}
	} else {
		logger = v[0]
	}

	return func(next tele.HandlerFunc) tele.HandlerFunc {
		return func(c tele.Context) error {
			logger.Context(c)

			switch {
			// Check callback first to avoid fetching its actual message.
			case c.Callback() != nil:
				logger.Callback(c.Callback())
			case c.Message() != nil:
				logger.Message(c.Message())
			case c.Query() != nil:
				logger.Query(c.Query())
			case c.ChosenInlineResult() != nil:
				logger.ChosenInlineResult(c.ChosenInlineResult())
			case c.ShippingQuery() != nil:
				logger.ShippingQuery(c.ShippingQuery())
			case c.PreCheckoutQuery() != nil:
				logger.PreCheckoutQuery(c.PreCheckoutQuery())
			case c.Poll() != nil:
				logger.Poll(c.Poll())
			case c.PollAnswer() != nil:
				logger.PollAnswer(c.PollAnswer())
			}

			return next(c)
		}
	}
}

type DefaultLogger struct {
	*log.Logger
}

func (l DefaultLogger) Context(c tele.Context) {
	return
}

func (l DefaultLogger) Message(m *tele.Message) {
	l.Println("Message", m.Sender.ID, m.Text)
}

func (l DefaultLogger) Callback(c *tele.Callback) {
	l.Println("Callback", c.Sender.ID, c.Data)
}

func (l DefaultLogger) Query(q *tele.Query) {
	l.Println("Query", q.Sender.ID, q.Text)
}

func (l DefaultLogger) ChosenInlineResult(r *tele.ChosenInlineResult) {
	l.Println("ChosenInlineResult", r.Sender.ID, r.Query)
}

func (l DefaultLogger) ShippingQuery(s *tele.ShippingQuery) {
	l.Println("ShippingQuery", s.Sender.ID, s.Payload)
}

func (l DefaultLogger) PreCheckoutQuery(pre *tele.PreCheckoutQuery) {
	l.Println("PreCheckoutQuery", pre.Sender.ID, pre.Payload)
}

func (l DefaultLogger) Poll(poll *tele.Poll) {
	l.Println("Poll", poll.ID)
}

func (l DefaultLogger) PollAnswer(pa *tele.PollAnswer) {
	l.Println("PollAnswer", pa.Sender.ID, pa.PollID, pa.Options)
}
