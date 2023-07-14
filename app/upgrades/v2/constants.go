package v2

import (
	store "github.com/cosmos/cosmos-sdk/store/types"
	"github.com/hard-nett/terp-node/v2/app/upgrades"
	feesharetypes "github.com/hard-nett/terp-node/v2/x/feeshare/types"
	"github.com/hard-nett/terp-node/v2/x/globalfee"
	ibchookstypes "github.com/hard-nett/terp-node/v2/x/ibchooks/types"
	packetforwardtypes "github.com/strangelove-ventures/packet-forward-middleware/v7/router/types"
)

// UpgradeName defines the on-chain upgrade name for the upgrade.
const UpgradeName = "v2"

var Upgrade = upgrades.Upgrade{
	UpgradeName:          UpgradeName,
	CreateUpgradeHandler: CreateV2UpgradeHandler,
	StoreUpgrades: store.StoreUpgrades{
		Added: []string{
			globalfee.ModuleName,
			ibchookstypes.StoreKey,
			packetforwardtypes.StoreKey,
			feesharetypes.ModuleName,
		},
	},
}
