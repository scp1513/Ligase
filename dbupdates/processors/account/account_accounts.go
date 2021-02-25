package processors

import (
	"context"

	"github.com/finogeeks/ligase/common"
	"github.com/finogeeks/ligase/dbupdates/dbregistry"
	"github.com/finogeeks/ligase/dbupdates/dbupdatetypes"
	"github.com/finogeeks/ligase/storage/model"
	"github.com/finogeeks/ligase/common/config"
	"github.com/finogeeks/ligase/model/dbtypes"
	"github.com/finogeeks/ligase/skunkworks/log"
)

func init() {
	dbregistry.Register("account_accounts", NewDBAccountAccountsProcessor, nil)
}

type DBAccountAccountsProcessor struct {
	name string
	cfg  *config.Dendrite
	db   model.AccountsDatabase
}

func NewDBAccountAccountsProcessor(
	name string,
	cfg *config.Dendrite,
) dbupdatetypes.DBEventSeqProcessor {
	p := new(DBAccountAccountsProcessor)
	p.name = name
	p.cfg = cfg

	return p
}

func (p *DBAccountAccountsProcessor) Start() {
	db, err := common.GetDBInstance("accounts", p.cfg)
	if err != nil {
		log.Panicf("failed to connect to accounts db")
	}
	p.db = db.(model.AccountsDatabase)
}

func (p *DBAccountAccountsProcessor) Process(ctx context.Context, inputs []dbupdatetypes.DBEventDataInput) error {
	if len(inputs) == 0 {
		return nil
	}

	switch inputs[0].Event.Key {
	case dbtypes.AccountInsertKey:
		p.processInsert(ctx, inputs)
	default:
		log.Errorf("invalid %s event key %d", p.name, inputs[0].Event.Key)
	}

	return nil
}

func (p *DBAccountAccountsProcessor) processInsert(ctx context.Context, inputs []dbupdatetypes.DBEventDataInput) error {
	for _, v := range inputs {
		msg := v.Event.AccountDBEvents.AccountInsert
		err := p.db.OnInsertAccount(ctx, msg.UserID, msg.PassWordHash, msg.AppServiceID, msg.CreatedTs)
		if err != nil {
			log.Error(p.name, "insert err", err, msg.UserID, msg.PassWordHash, msg.AppServiceID, msg.CreatedTs)
		}
	}
	return nil
}
