package utils

import (
	galileo "github.com/Sharefunds/galileo-client"
	"github.com/Sharefunds/pool-transactions/pkg/models"
	"github.com/Sharefunds/pool-transactions/pkg/pb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func ConvertToPBParticipant(modelParticipant models.Participant) *pb.Participant {
	return &pb.Participant{
		ParticipantId: modelParticipant.ParticipantID.String(),
		PoolId:        modelParticipant.PoolID.String(),
		UserId:        modelParticipant.UserID,
		Role:          modelParticipant.Role,
		JoinedAt:      timestamppb.New(modelParticipant.JoinedAt),
	}
}

func ConvertToPBTransaction(modelTransaction models.Transaction) *pb.Transaction {
	return &pb.Transaction{
		TransactionId: modelTransaction.TransactionID.String(),
		FromAccountId: modelTransaction.FromAccountID,
		ToAccountId:   modelTransaction.ToAccountID,
		Amount:        float32(modelTransaction.Amount),
		Description:   modelTransaction.Description,
		Status:        modelTransaction.Status,
		CreatedAt:     timestamppb.New(modelTransaction.CreatedAt),
	}
}

func ConvertToPBPool(modelPool models.Pool) *pb.Pool {
	return &pb.Pool{
		PoolId:      modelPool.PoolID.String(),
		Poolname:    modelPool.PoolName,
		Description: modelPool.Description,
		BalanceGoal: int32(modelPool.BalanceGoal),
		CreatedAt:   timestamppb.New(modelPool.CreatedAt).String(),
	}
}

func ConvertTransactions(transactions []galileo.ResponseData3Transactions) []*pb.GalileoTransactions {
	var pbTransactions []*pb.GalileoTransactions
	for _, t := range transactions {
		pbTrans := &pb.GalileoTransactions{
			PmtRefNo:           t.PmtRefNo,
			ActId:              t.ActId,
			ActType:            t.ActType,
			Mcc:                t.Mcc,
			PostTs:             t.PostTs.String(),
			Amt:                t.Amt,
			Details:            t.Details,
			Description:        t.Description,
			SourceId:           t.SourceId,
			BalId:              t.BalId,
			ProdId:             t.ProdId,
			AuthTs:             t.AuthTs.String(),
			TransCode:          t.TransCode,
			AchTransactionId:   t.AchTransactionId,
			ExternalTransId:    t.ExternalTransId,
			OriginalAuthId:     t.OriginalAuthId,
			NetworkId:          t.NetworkId,
			LocalAmt:           t.LocalAmt,
			LocalCurrCode:      t.LocalCurrCode,
			SettleAmt:          t.SettleAmt,
			SettleCurrCode:     t.SettleCurrCode,
			BillingAmt:         t.BillingAmt,
			BillingCurrCode:    t.BillingCurrCode,
			IacTax:             t.IacTax,
			IvaTax:             t.IvaTax,
			FundingAccountPrn:  t.FundingAccountPrn,
			SpendingAccountPrn: t.SpendingAccountPrn,
		}
		pbTransactions = append(pbTransactions, pbTrans)
	}
	return pbTransactions
}
