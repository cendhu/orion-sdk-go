package bcdb

import (
	"errors"
	"fmt"

	"github.ibm.com/blockchaindb/server/pkg/constants"
	"github.ibm.com/blockchaindb/server/pkg/types"
)

type provenance struct {
	*commonTxContext
}

func (p *provenance) GetHistoricalData(dbName, key string) ([]*types.ValueWithMetadata, error) {
	path := constants.URLForGetHistoricalData(dbName, key)
	res := &types.GetHistoricalDataResponseEnvelope{}
	err := p.handleRequest(path, &types.GetHistoricalDataQuery{
		UserID: p.userID,
		DBName: dbName,
		Key:    key,
	}, res)
	if err != nil {
		p.logger.Errorf("failed to execute historical data query %s, due to %s", path, err)
		return nil, err
	}
	return res.GetPayload().GetValues(), nil
}

func (p *provenance) GetHistoricalDataAt(dbName, key string, version *types.Version) (*types.ValueWithMetadata, error) {
	path := constants.URLForGetHistoricalDataAt(dbName, key, version)
	res := &types.GetHistoricalDataResponseEnvelope{}
	err := p.handleRequest(path, &types.GetHistoricalDataQuery{
		UserID:  p.userID,
		DBName:  dbName,
		Key:     key,
		Version: version,
	}, res)
	if err != nil {
		p.logger.Errorf("failed to parse execute data query %s, due to %s", path, err)
		return nil, err
	}

	values := res.GetPayload().GetValues()
	if len(values) == 0 {
		return nil, nil
	}
	if len(values) > 1 {
		return nil, errors.New(fmt.Sprintf("error getting historical data fro specific version, more that one record returned"))
	}
	return values[0], nil
}

func (p *provenance) GetPreviousHistoricalData(dbName, key string, version *types.Version) ([]*types.ValueWithMetadata, error) {
	path := constants.URLForGetPreviousHistoricalData(dbName, key, version)
	res := &types.GetHistoricalDataResponseEnvelope{}
	err := p.handleRequest(path, &types.GetHistoricalDataQuery{
		UserID:    p.userID,
		DBName:    dbName,
		Key:       key,
		Version:   version,
		Direction: "previous",
	}, res)
	if err != nil {
		p.logger.Errorf("failed to execute previous historical data query %s, due to %s", path, err)
		return nil, err
	}
	return res.GetPayload().GetValues(), nil
}

func (p *provenance) GetNextHistoricalData(dbName, key string, version *types.Version) ([]*types.ValueWithMetadata, error) {
	path := constants.URLForGetNextHistoricalData(dbName, key, version)
	res := &types.GetHistoricalDataResponseEnvelope{}
	err := p.handleRequest(path, &types.GetHistoricalDataQuery{
		UserID:    p.userID,
		DBName:    dbName,
		Key:       key,
		Version:   version,
		Direction: "next",
	}, res)
	if err != nil {
		p.logger.Errorf("failed to execute next historical data query %s, due to %s", path, err)
		return nil, err
	}
	return res.GetPayload().GetValues(), nil
}

func (p *provenance) GetDataReadByUser(userID string) ([]*types.KVWithMetadata, error) {
	path := constants.URLForGetDataReadBy(userID)
	res := &types.GetDataReadByResponseEnvelope{}
	err := p.handleRequest(path, &types.GetDataReadByQuery{
		UserID:       p.userID,
		TargetUserID: userID,
	}, res)
	if err != nil {
		p.logger.Errorf("failed to execute data read by user query %s, due to %s", path, err)
		return nil, err
	}
	return res.GetPayload().GetKVs(), nil
}

func (p *provenance) GetDataWrittenByUser(userID string) ([]*types.KVWithMetadata, error) {
	path := constants.URLForGetDataWrittenBy(userID)
	res := &types.GetDataWrittenByResponseEnvelope{}
	err := p.handleRequest(path, &types.GetDataWrittenByQuery{
		UserID:       p.userID,
		TargetUserID: userID,
	}, res)
	if err != nil {
		p.logger.Errorf("failed to execute data written by user query %s, due to %s", path, err)
		return nil, err
	}
	return res.GetPayload().GetKVs(), nil
}

func (p *provenance) GetReaders(dbName, key string) ([]string, error) {
	path := constants.URLForGetDataReaders(dbName, key)
	res := &types.GetDataReadersResponseEnvelope{}
	err := p.handleRequest(path, &types.GetDataReadersQuery{
		UserID: p.userID,
		DBName: dbName,
		Key:    key,
	}, res)
	if err != nil {
		p.logger.Errorf("failed to execute data readers query %s, due to %s", path, err)
		return nil, err
	}

	if res.GetPayload().GetReadBy() == nil {
		return nil, nil
	}
	readers := make([]string, 0)
	for k := range res.GetPayload().GetReadBy() {
		readers = append(readers, k)
	}
	return readers, nil
}

func (p *provenance) GetWriters(dbName, key string) ([]string, error) {
	path := constants.URLForGetDataWriters(dbName, key)
	res := &types.GetDataWritersResponseEnvelope{}
	err := p.handleRequest(path, &types.GetDataWritersQuery{
		UserID: p.userID,
		DBName: dbName,
		Key:    key,
	}, res)
	if err != nil {
		p.logger.Errorf("failed to execute data writers query %s, due to %s", path, err)
		return nil, err
	}

	if res.GetPayload().GetWrittenBy() == nil {
		return nil, nil
	}
	writers := make([]string, 0)
	for k := range res.GetPayload().GetWrittenBy() {
		writers = append(writers, k)
	}
	return writers, nil
}

func (p *provenance) GetTxIDsSubmittedByUser(userID string) ([]string, error) {
	path := constants.URLForGetTxIDsSubmittedBy(userID)
	res := &types.GetTxIDsSubmittedByResponseEnvelope{}
	err := p.handleRequest(path, &types.GetTxIDsSubmittedByQuery{
		UserID:       p.userID,
		TargetUserID: userID,
	}, res)
	if err != nil {
		p.logger.Errorf("failed to execute tx id by user query %s, due to %s", path, err)
		return nil, err
	}
	return res.GetPayload().GetTxIDs(), nil
}
