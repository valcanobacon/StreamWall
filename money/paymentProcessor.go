package money

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/lightningnetwork/lnd/lnrpc"
)

type PaymentProcessor struct {
	Client lnrpc.LightningClient
	Bank   *Bank
	Index  uint64
}

func (p *PaymentProcessor) Run(ctx context.Context) {

	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		invoices, err := p.nextBatch()
		if err == nil {
			p.processBatch(invoices)
			p.Index += uint64(len(invoices))
		}

		time.Sleep(3 * time.Second)
	}

}

func (p *PaymentProcessor) nextBatch() ([]*lnrpc.Invoice, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	req := &lnrpc.ListInvoiceRequest{IndexOffset: p.Index}
	res, err := p.Client.ListInvoices(ctx, req)
	if err != nil {
		return nil, err
	}

	return res.Invoices, nil
}

func (p *PaymentProcessor) processBatch(invoices []*lnrpc.Invoice) {
	for i := 0; i < len(invoices); i++ {
		invoice := invoices[i]

		p.processInvoice(invoice)
	}
}

func (p *PaymentProcessor) processInvoice(invoice *lnrpc.Invoice) {
	htlcs := invoice.GetHtlcs()

	for j := 0; j < len(htlcs); j++ {

		htlc := htlcs[j]

		customRecords := htlc.GetCustomRecords()

		record, ok := customRecords[123123]

		if !ok {
			continue
		}

		value := invoice.Value
		sid, err := uuid.Parse(strings.Trim(string(record[:]), " \n\r"))
		if err != nil {
			continue
		}

		p.Bank.NewTransaction(sid, value)
	}
}
