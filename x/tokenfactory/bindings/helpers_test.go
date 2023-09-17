package bindings_test

import (
	"os"
	"testing"
	"time"

	"github.com/CosmWasm/wasmd/x/wasm/keeper"
	"github.com/stretchr/testify/require"

	"github.com/cometbft/cometbft/crypto"
	"github.com/cometbft/cometbft/crypto/ed25519"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	banktestutil "github.com/cosmos/cosmos-sdk/x/bank/testutil"

	"github.com/terpnetwork/terp-core/v2/app"
)

func CreateTestInput(t *testing.T) (*app.TerpApp, sdk.Context) {
	osmosis := app.Setup(t)
	ctx := osmosis.BaseApp.NewContext(false, tmproto.Header{Height: 1, ChainID: "testing", Time: time.Now().UTC()})
	return osmosis, ctx
}

func FundAccount(t *testing.T, ctx sdk.Context, terpapp *app.TerpApp, acct sdk.AccAddress) {
	err := banktestutil.FundAccount(terpapp.BankKeeper, ctx, acct, sdk.NewCoins(
		sdk.NewCoin("uosmo", sdk.NewInt(10000000000)),
	))
	require.NoError(t, err)
}

// we need to make this deterministic (same every test run), as content might affect gas costs
func keyPubAddr() (crypto.PrivKey, crypto.PubKey, sdk.AccAddress) {
	key := ed25519.GenPrivKey()
	pub := key.PubKey()
	addr := sdk.AccAddress(pub.Address())
	return key, pub, addr
}

func RandomAccountAddress() sdk.AccAddress {
	_, _, addr := keyPubAddr()
	return addr
}

func RandomBech32AccountAddress() string {
	return RandomAccountAddress().String()
}

func storeReflectCode(t *testing.T, ctx sdk.Context, terpapp *app.TerpApp, addr sdk.AccAddress) uint64 {
	wasmCode, err := os.ReadFile("./testdata/token_reflect.wasm")
	require.NoError(t, err)

	contractKeeper := keeper.NewDefaultPermissionKeeper(terpapp.WasmKeeper)
	codeID, _, err := contractKeeper.Create(ctx, addr, wasmCode, nil)
	require.NoError(t, err)

	return codeID
}

func instantiateReflectContract(t *testing.T, ctx sdk.Context, terpapp *app.TerpApp, funder sdk.AccAddress) sdk.AccAddress {
	initMsgBz := []byte("{}")
	contractKeeper := keeper.NewDefaultPermissionKeeper(terpapp.WasmKeeper)
	codeID := uint64(1)
	addr, _, err := contractKeeper.Instantiate(ctx, codeID, funder, funder, initMsgBz, "demo contract", nil)
	require.NoError(t, err)

	return addr
}

func fundAccount(t *testing.T, ctx sdk.Context, terpapp *app.TerpApp, addr sdk.AccAddress, coins sdk.Coins) {
	err := banktestutil.FundAccount(
		terpapp.BankKeeper,
		ctx,
		addr,
		coins,
	)
	require.NoError(t, err)
}

func SetupCustomApp(t *testing.T, addr sdk.AccAddress) (*app.TerpApp, sdk.Context) {
	terpapp, ctx := CreateTestInput(t)
	wasmKeeper := terpapp.WasmKeeper

	storeReflectCode(t, ctx, terpapp, addr)

	cInfo := wasmKeeper.GetCodeInfo(ctx, 1)
	require.NotNil(t, cInfo)

	return terpapp, ctx
}
