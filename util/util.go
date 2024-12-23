package util

import (
	"encoding/hex"
	"fmt"
	"time"

	bin "github.com/gagliardetto/binary"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	constant "github.com/pefish/go-coin-sol/constant"
	type_ "github.com/pefish/go-coin-sol/type"
	go_decimal "github.com/pefish/go-decimal"
	go_http "github.com/pefish/go-http"
	i_logger "github.com/pefish/go-interface/i-logger"
)

func FindInnerInstructions(meta *rpc.TransactionMeta, index uint64) []solana.CompiledInstruction {
	for _, innerInstruction := range meta.InnerInstructions {
		if innerInstruction.Index == uint16(index) {
			return innerInstruction.Instructions
		}
	}
	return nil
}

func GetComputeUnitPriceFromHelius(
	logger i_logger.ILogger,
	key string,
	accountKeys []string,
) (uint64, error) {
	var httpResult struct {
		Result struct {
			PriorityFeeEstimate float64 `json:"priorityFeeEstimate"`
		} `json:"result"`
	}
	_, _, err := go_http.NewHttpRequester(
		go_http.WithLogger(logger),
		go_http.WithTimeout(10*time.Second),
	).PostForStruct(
		&go_http.RequestParams{
			Url: fmt.Sprintf("https://mainnet.helius-rpc.com/?api-key=%s", key),
			Params: map[string]interface{}{
				"jsonrpc": "2.0",
				"id":      "helius-example",
				"method":  "getPriorityFeeEstimate",
				"params": []map[string]interface{}{
					{
						"accountKeys": accountKeys,
						"options": map[string]interface{}{
							"recommended": true,
						},
					},
				},
			},
		},
		&httpResult,
	)
	if err != nil {
		return 0, err
	}
	return go_decimal.Decimal.MustStart(httpResult.Result.PriorityFeeEstimate).RoundDown(0).MustEndForUint64(), nil
}

func GetFeeInfoFromTx(meta *rpc.TransactionMeta, transaction *solana.Transaction) (*type_.FeeInfo, error) {
	accountKeys := transaction.Message.AccountKeys
	if meta.LoadedAddresses.Writable != nil {
		accountKeys = append(accountKeys, meta.LoadedAddresses.Writable...)
	}
	if meta.LoadedAddresses.ReadOnly != nil {
		accountKeys = append(accountKeys, meta.LoadedAddresses.ReadOnly...)
	}

	totalFee := go_decimal.Decimal.MustStart(meta.Fee).MustUnShiftedBy(constant.SOL_Decimals).EndForString()
	priorityFee := "0"
	computeUnitPrice := 0

	var setComputeUnitLimitInstru solana.CompiledInstruction
	var setComputeUnitPriceInstru solana.CompiledInstruction
	for _, instruction := range transaction.Message.Instructions {
		programPKey := accountKeys[instruction.ProgramIDIndex]
		if !programPKey.Equals(constant.Compute_Budget) {
			continue
		}
		methodId := hex.EncodeToString(instruction.Data)[:2]
		if methodId == "02" {
			setComputeUnitLimitInstru = instruction
		}
		if methodId == "03" {
			setComputeUnitPriceInstru = instruction
		}
	}
	computeUnitLimit := 200000
	if setComputeUnitLimitInstru.ProgramIDIndex != 0 {
		var params struct {
			Id    uint8  `json:"id"`
			Units uint32 `json:"units"`
		}
		err := bin.NewBorshDecoder(setComputeUnitLimitInstru.Data).Decode(&params)
		if err != nil {
			return nil, err
		}
		computeUnitLimit = int(params.Units)
	}

	if setComputeUnitPriceInstru.ProgramIDIndex != 0 {
		var params struct {
			Id            uint8  `json:"id"`
			MicroLamports uint64 `json:"microLamports"`
		}
		err := bin.NewBorshDecoder(setComputeUnitPriceInstru.Data).Decode(&params)
		if err != nil {
			return nil, err
		}
		computeUnitPrice = int(params.MicroLamports)

		priorityFee = go_decimal.Decimal.MustStart(computeUnitPrice).MustMulti(computeUnitLimit).MustUnShiftedBy(constant.SOL_Decimals + 6).EndForString()
	}

	return &type_.FeeInfo{
		BaseFee:          go_decimal.Decimal.MustStart(totalFee).MustSubForString(priorityFee),
		PriorityFee:      priorityFee,
		TotalFee:         totalFee,
		ComputeUnitPrice: uint64(computeUnitPrice),
	}, nil
}
