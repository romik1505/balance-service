package balance

import (
	"context"
	"fmt"

	"github.com/romik1505/balance-service/internal/mapper"
	"github.com/romik1505/balance-service/internal/store"
)

const AvitoID = "00000000-0000-0000-0000-000000000000"

func (b *BalanceService) Transfer(ctx context.Context, req mapper.TransferRequest) (mapper.Transfer, error) {
	if err := req.Bind(); err != nil {
		return mapper.Transfer{}, err
	}

	transfer := mapper.Transfer{
		SenderID:   req.SenderID,
		ReceiverID: req.ReceiverID,
		Amount:     req.Amount,
	}

	if req.ReceiverID == "" || req.ReceiverID == AvitoID {
		transfer.Type = mapper.TransferTypeDebit
		transfer.Description = fmt.Sprintf("Оплата услуг пользователя %s на Авито", req.SenderID)
		transfer.ReceiverID = AvitoID
	} else if req.SenderID == "" || req.SenderID == AvitoID {
		transfer.Type = mapper.TransferTypeCredit
		transfer.Description = fmt.Sprintf("Пополнение счета пользователя %s на Авито", req.ReceiverID)
		transfer.SenderID = AvitoID
	} else {
		transfer.Type = mapper.TransferTypeTransfer
		transfer.Description = fmt.Sprintf("Перевод денежных средств от пользователя %s пользователю %s", req.SenderID, req.ReceiverID)
	}

	mTransfer := mapper.ConvertTransferToModel(transfer)
	entryParts := transfer.EntryParts()
	mEntryPart := mapper.ConvertEntryPartToModel(entryParts)

	mTransfer, _, err := b.Storage.InsertTransferWithEntryParts(ctx, mTransfer, mEntryPart)
	if err != nil {
		return mapper.Transfer{}, err
	}
	return mapper.ConvertTransfer(mTransfer), nil
}

func (b *BalanceService) ListTransfers(ctx context.Context, filter store.ListTransfersFilter) (mapper.TransfersResponse, error) {
	transfers, totalItems, err := b.Storage.ListTransfers(ctx, filter)
	if err != nil {
		return mapper.TransfersResponse{}, err
	}
	return mapper.TransfersResponse{
		Items:      mapper.ConvertTransfers(transfers),
		TotalItems: totalItems,
	}, nil
}
