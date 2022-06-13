package itests

import (
	"context"
	"encoding/hex"
	"testing"
	"time"

	"github.com/filecoin-project/go-state-types/builtin"
	"github.com/filecoin-project/go-state-types/network"
	"github.com/filecoin-project/lotus/chain/types"
	"github.com/filecoin-project/lotus/itests/kit"
	"github.com/stretchr/testify/require"
)

// The following function will generate and print out the hex encoding of PublishStorageDealsParams with a dealLabel with invalid utf8 encoding.
// This is not possible with our current code structure, so this was run with a modified NewLabelFromString() function in go-state-types
// The result is recorded in the function below for submission.
//func generateInvalidDealProposal(t *testing.T) {
//	ctx := context.Background()
//
//	kit.QuietMiningLogs()
//
//	client16, _, ens := kit.EnsembleMinimal(t, kit.MockProofs(), kit.GenesisNetworkVersion(network.Version16))
//	ens.InterconnectAll().BeginMining(10 * time.Millisecond)
//
//	dealLabel, err := market.NewLabelFromString(string([]byte{0xde, 0xad, 0xbe, 0xef}))
//	require.NoError(t, err)
//
//	dummyCid, err := cid.Parse("bafkqaaa")
//	require.NoError(t, err)
//
//	proposal := market.DealProposal{
//		PieceCID: dummyCid,
//		Client:   client16.DefaultKey.Address,
//		Provider: client16.DefaultKey.Address,
//		Label:    dealLabel,
//	}
//
//	proposalBytes := new(bytes.Buffer)
//	err = proposal.MarshalCBOR(proposalBytes)
//	require.NoError(t, err)
//
//	signature, err := client16.WalletSign(ctx, client16.DefaultKey.Address, proposalBytes.Bytes())
//	require.NoError(t, err)
//
//	params, err := actors.SerializeParams(&market.PublishStorageDealsParams{
//		Deals: []market.ClientDealProposal{{
//			Proposal:        proposal,
//			ClientSignature: *signature,
//		}},
//	})
//	require.NoError(t, err)
//
//	fmt.Println(hex.EncodeToString(params))
//}

func TestSubmitInvalidDealLabel(t *testing.T) {
	ctx := context.Background()

	kit.QuietMiningLogs()

	client16, _, ens := kit.EnsembleMinimal(t, kit.MockProofs(), kit.GenesisNetworkVersion(network.Version16))
	ens.InterconnectAll().BeginMining(10 * time.Millisecond)

	// This hex string is generated by using the commented function above.
	serializedParams, err := hex.DecodeString("8181828bd82a45000155000000f458310394345a1393a5665772ae97456fdc00cf8f537960eeb3838c297da0d60546252d98469cc68b56eb45081e5c959c8cf4fa58310394345a1393a5665772ae97456fdc00cf8f537960eeb3838c297da0d60546252d98469cc68b56eb45081e5c959c8cf4fa64deadbeef00004040405861029927585cf250656cb3ddf337aed516e6393829389e52c45a4fbbe297b8ad3b7a008b3675a5b5c4df43657ff53bcb6c670b4bd83911867ba90b2b14da856a0ecb36782a7d87232272e3b9cd98dd959b2851602fea30c414a4dc1ea12057d31f81")
	require.NoError(t, err)

	_, err = client16.MpoolPushMessage(ctx, &types.Message{
		To:     builtin.StorageMarketActorAddr,
		From:   client16.DefaultKey.Address,
		Value:  types.NewInt(0),
		Method: builtin.MethodsMarket.PublishStorageDeals,
		Params: serializedParams,
	}, nil)
	require.Contains(t, err.Error(), "Serialization error for Cbor protocol: InvalidUtf8")
}