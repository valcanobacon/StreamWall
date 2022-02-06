package money

import (
	"context"
	"log"

	"github.com/google/uuid"
)

type Bank struct {
	sessions     map[SessionID]*Session
	transactions chan Transaction
}

func NewBank() *Bank {
	return &Bank{
		sessions:     make(map[SessionID]*Session),
		transactions: make(chan Transaction),
	}
}

func (b *Bank) GetSession(sid SessionID) *Session {
	s, ok := b.sessions[sid]
	if !ok {
		return b.SetSession(sid, 0)
	}
	return s
}

func (b *Bank) NewSession() *Session {
	id := uuid.New()
	s := b.SetSession(id, 0)
	return s
}

func (b *Bank) SetSession(sid SessionID, amount int64) *Session {
	s := &Session{Id: sid, Credits: amount}
	b.sessions[s.Id] = s
	return s
}

func (b *Bank) NewTransaction(sid SessionID, amount int64) {
	b.transactions <- Transaction{SID: sid, Amount: amount}
}

func (b *Bank) ProcessTransactions(ctx context.Context) {

	for {
		select {
		case <-ctx.Done():
			return
		case t := <-b.transactions:
			s := b.GetSession((t.SID))
			if s == nil {
				continue
			}

			s.Credits += t.Amount

			if t.Amount > 0 {
				log.Printf("Bank %s + %d = %d", s.Id.String(), t.Amount, s.Credits)
			} else if t.Amount < 0 {
				log.Printf("Bank %s - %d = %d", s.Id.String(), -1*t.Amount, s.Credits)
			}
		}
	}
}
