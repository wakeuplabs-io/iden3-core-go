package backupsrv

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/iden3/go-iden3/services/claimsrv"
)

type SaveBackupMsg struct {
	IdAddrHex       string
	Data            string
	DataSignature   string
	KSign           common.Address
	ProofOfKSignHex claimsrv.ProofOfClaimHex
	RelayAddr       common.Address
	Timestamp       uint64
}
