package type_

import (
	bin "github.com/gagliardetto/binary"
	"github.com/gagliardetto/solana-go"
	"github.com/pkg/errors"
)

type SwapV2Keys struct {
	AmmConfig        solana.PublicKey
	PairAddress      solana.PublicKey
	Vaults           map[solana.PublicKey]solana.PublicKey
	ObservationState solana.PublicKey
	ExBitmapAccount  solana.PublicKey
	RemainAccounts   []solana.PublicKey
}

func (t *SwapV2Keys) ToAccounts(
	userAddress solana.PublicKey,
	inputToken solana.PublicKey,
	outputToken solana.PublicKey,
) ([]*solana.AccountMeta, error) {
	userInputAssociatedAccount, _, err := solana.FindAssociatedTokenAddress(
		userAddress,
		inputToken,
	)
	if err != nil {
		return nil, errors.Wrapf(err, "<userAddress: %s> <tokenAddress: %s>", userAddress, inputToken)
	}

	userOutputAssociatedAccount, _, err := solana.FindAssociatedTokenAddress(
		userAddress,
		outputToken,
	)
	if err != nil {
		return nil, errors.Wrapf(err, "<userAddress: %s> <tokenAddress: %s>", userAddress, outputToken)
	}

	accounts := []*solana.AccountMeta{
		{
			PublicKey:  userAddress,
			IsSigner:   true,
			IsWritable: true,
		},
		{
			PublicKey:  t.AmmConfig,
			IsSigner:   false,
			IsWritable: false,
		},
		{
			PublicKey:  t.PairAddress,
			IsSigner:   false,
			IsWritable: true,
		},
		{
			PublicKey:  userInputAssociatedAccount,
			IsSigner:   false,
			IsWritable: true,
		},
		{
			PublicKey:  userOutputAssociatedAccount,
			IsSigner:   false,
			IsWritable: true,
		},
		{
			PublicKey:  t.Vaults[inputToken],
			IsSigner:   false,
			IsWritable: true,
		},
		{
			PublicKey:  t.Vaults[outputToken],
			IsSigner:   false,
			IsWritable: true,
		},
		{
			PublicKey:  t.ObservationState,
			IsSigner:   false,
			IsWritable: true,
		},
		{
			PublicKey:  solana.TokenProgramID,
			IsSigner:   false,
			IsWritable: false,
		},
		{
			PublicKey:  solana.Token2022ProgramID,
			IsSigner:   false,
			IsWritable: false,
		},
		{
			PublicKey:  solana.MemoProgramID,
			IsSigner:   false,
			IsWritable: false,
		},
		{
			PublicKey:  inputToken,
			IsSigner:   false,
			IsWritable: false,
		},
		{
			PublicKey:  outputToken,
			IsSigner:   false,
			IsWritable: false,
		},
		{
			PublicKey:  t.ExBitmapAccount,
			IsSigner:   false,
			IsWritable: true,
		},
	}

	for _, remainAccount := range t.RemainAccounts {
		accounts = append(accounts, &solana.AccountMeta{
			PublicKey:  remainAccount,
			IsSigner:   false,
			IsWritable: true,
		})
	}

	return accounts, nil
}

type SwapKeys struct {
	AmmConfig        solana.PublicKey
	PairAddress      solana.PublicKey
	Vaults           map[solana.PublicKey]solana.PublicKey
	ObservationState solana.PublicKey
	TickArrayAccount solana.PublicKey
	RemainAccounts   []solana.PublicKey
}

func (t *SwapKeys) ToAccounts(
	userAddress solana.PublicKey,
	inputToken solana.PublicKey,
	outputToken solana.PublicKey,
) ([]*solana.AccountMeta, error) {
	userInputAssociatedAccount, _, err := solana.FindAssociatedTokenAddress(
		userAddress,
		inputToken,
	)
	if err != nil {
		return nil, errors.Wrapf(err, "<userAddress: %s> <tokenAddress: %s>", userAddress, inputToken)
	}

	userOutputAssociatedAccount, _, err := solana.FindAssociatedTokenAddress(
		userAddress,
		outputToken,
	)
	if err != nil {
		return nil, errors.Wrapf(err, "<userAddress: %s> <tokenAddress: %s>", userAddress, outputToken)
	}

	accounts := []*solana.AccountMeta{
		{
			PublicKey:  userAddress,
			IsSigner:   true,
			IsWritable: true,
		},
		{
			PublicKey:  t.AmmConfig,
			IsSigner:   false,
			IsWritable: false,
		},
		{
			PublicKey:  t.PairAddress,
			IsSigner:   false,
			IsWritable: true,
		},
		{
			PublicKey:  userInputAssociatedAccount,
			IsSigner:   false,
			IsWritable: true,
		},
		{
			PublicKey:  userOutputAssociatedAccount,
			IsSigner:   false,
			IsWritable: true,
		},
		{
			PublicKey:  t.Vaults[inputToken],
			IsSigner:   false,
			IsWritable: true,
		},
		{
			PublicKey:  t.Vaults[outputToken],
			IsSigner:   false,
			IsWritable: true,
		},
		{
			PublicKey:  t.ObservationState,
			IsSigner:   false,
			IsWritable: true,
		},
		{
			PublicKey:  solana.TokenProgramID,
			IsSigner:   false,
			IsWritable: false,
		},
		{
			PublicKey:  t.TickArrayAccount,
			IsSigner:   false,
			IsWritable: true,
		},
	}

	for _, remainAccount := range t.RemainAccounts {
		accounts = append(accounts, &solana.AccountMeta{
			PublicKey:  remainAccount,
			IsSigner:   false,
			IsWritable: true,
		})
	}

	return accounts, nil
}

type PoolInfo struct {
	Id                  uint64
	Bump                [1]uint8
	AmmConfig           solana.PublicKey
	Owner               solana.PublicKey
	TokenMint0          solana.PublicKey
	TokenMint1          solana.PublicKey
	TokenVault0         solana.PublicKey
	TokenVault1         solana.PublicKey
	ObservationKey      solana.PublicKey
	MintDecimals0       uint8
	MintDecimals1       uint8
	TickSpacing         uint16
	Liquidity           bin.Uint128
	SqrtPriceX64        bin.Uint128
	TickCurrent         int32
	Padding3            uint16
	Padding4            uint16
	FeeGrowthGlobal0X64 bin.Uint128
	FeeGrowthGlobal1X64 bin.Uint128
	ProtocolFeesToken0  uint64
	ProtocolFeesToken1  uint64
	SwapInAmountToken0  bin.Uint128
	SwapOutAmountToken1 bin.Uint128
	SwapInAmountToken1  bin.Uint128
	SwapOutAmountToken0 bin.Uint128
}

type ExtraDatasType struct {
	ReserveInputWithDecimals  uint64 `json:"reserve_input_with_decimals"`
	ReserveOutputWithDecimals uint64 `json:"reserve_output_with_decimals"`
}
