package bcdb

import (
	"github.com/golang/protobuf/proto"
	"github.ibm.com/blockchaindb/server/pkg/constants"
	"github.ibm.com/blockchaindb/server/pkg/cryptoservice"
	"github.ibm.com/blockchaindb/server/pkg/types"
)

// DBsTxContext abstraction for database management transaction context
type DBsTxContext interface {
	TxContext
	// CreateDB creates new database
	CreateDB(dbName string) error
	// DeleteDB deletes database
	DeleteDB(dbName string) error
	// Exists checks whenever database is already created
	Exists(dbName string) (bool, error)
}

type dbsTxContext struct {
	*commonTxContext
	createdDBs map[string]bool
	deletedDBs map[string]bool
}

func (d *dbsTxContext) Commit() (string, error) {
	return d.commit(d, constants.PostDBTx)
}

func (d *dbsTxContext) Abort() error {
	return d.commonTxContext.abort(d)
}

func (d *dbsTxContext) CreateDB(dbName string) error {
	d.createdDBs[dbName] = true
	return nil
}

func (d *dbsTxContext) DeleteDB(dbName string) error {
	d.deletedDBs[dbName] = true
	return nil
}

func (d *dbsTxContext) Exists(dbName string) (bool, error) {
	path := constants.URLForGetDBStatus(dbName)
	res := &types.GetDBStatusResponseEnvelope{}
	err := d.handleRequest(path, &types.GetDBStatusQuery{
		UserID: d.userID,
		DBName: dbName,
	}, res)
	if err != nil {
		d.logger.Errorf("failed to execute database status query, path = %s, due to %s", path, err)
		return false, err
	}
	return res.GetPayload().GetExist(), nil
}

func (d *dbsTxContext) cleanCtx() {
	d.createdDBs = map[string]bool{}
	d.deletedDBs = map[string]bool{}
}

func (d *dbsTxContext) composeEnvelope(txID string) (proto.Message, error) {
	payload := &types.DBAdministrationTx{
		UserID: d.userID,
		TxID:   txID,
	}

	for db := range d.createdDBs {
		payload.CreateDBs = append(payload.CreateDBs, db)
	}

	for db := range d.deletedDBs {
		payload.DeleteDBs = append(payload.DeleteDBs, db)
	}

	signature, err := cryptoservice.SignTx(d.signer, payload)
	if err != nil {
		return nil, err
	}

	return &types.DBAdministrationTxEnvelope{
		Payload:   payload,
		Signature: signature,
	}, nil
}
