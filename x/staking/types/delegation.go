package types

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/staking/exported"
	yaml "gopkg.in/yaml.v2"
)

// Implements Delegation interface
var _ exported.DelegationI = Delegation{}

// String implements the Stringer interface for a DVPair object.
func (dv DVPair) String() string {
	out, _ := yaml.Marshal(dv)
	return string(out)
}

// String implements the Stringer interface for a DVVTriplet object.
func (dvv DVVTriplet) String() string {
	out, _ := yaml.Marshal(dvv)
	return string(out)
}

// NewDelegation creates a new delegation object
func NewDelegation(delegatorAddr sdk.AccAddress, validatorAddr sdk.ValAddress, shares sdk.Dec) Delegation {
	return Delegation{
		DelegatorAddress: delegatorAddr,
		ValidatorAddress: validatorAddr,
		Shares:           shares,
	}
}

// MustMarshalDelegation returns the delegation bytes. Panics if fails
func MustMarshalDelegation(cdc codec.Marshaler, delegation Delegation) []byte {
	return cdc.MustMarshalBinaryLengthPrefixed(&delegation)
}

// MustUnmarshalDelegation return the unmarshaled delegation from bytes.
// Panics if fails.
func MustUnmarshalDelegation(cdc codec.Marshaler, value []byte) Delegation {
	delegation, err := UnmarshalDelegation(cdc, value)
	if err != nil {
		panic(err)
	}
	return delegation
}

// return the delegation
func UnmarshalDelegation(cdc codec.Marshaler, value []byte) (delegation Delegation, err error) {
	err = cdc.UnmarshalBinaryLengthPrefixed(value, &delegation)
	return delegation, err
}

// nolint - for Delegation
func (d Delegation) GetDelegatorAddr() sdk.AccAddress { return d.DelegatorAddress }
func (d Delegation) GetValidatorAddr() sdk.ValAddress { return d.ValidatorAddress }
func (d Delegation) GetShares() sdk.Dec               { return d.Shares }

// String returns a human readable string representation of a Delegation.
func (d Delegation) String() string {
	out, _ := yaml.Marshal(d)
	return string(out)
}

// Delegations is a collection of delegations
type Delegations []Delegation

func (d Delegations) String() (out string) {
	for _, del := range d {
		out += del.String() + "\n"
	}
	return strings.TrimSpace(out)
}

func NewUnbondingDelegationEntry(creationHeight int64, completionTime time.Time, balance sdk.Int) UnbondingDelegationEntry {
	return UnbondingDelegationEntry{
		CreationHeight: creationHeight,
		CompletionTime: completionTime,
		InitialBalance: balance,
		Balance:        balance,
	}
}

// String implements the stringer interface for a UnbondingDelegationEntry.
func (ubd UnbondingDelegationEntry) String() string {
	out, _ := yaml.Marshal(ubd)
	return string(out)
}

// IsMature - is the current entry mature
func (e UnbondingDelegationEntry) IsMature(currentTime time.Time) bool {
	return !e.CompletionTime.After(currentTime)
}

// NewUnbondingDelegation - create a new unbonding delegation object
func NewUnbondingDelegation(
	delegatorAddr sdk.AccAddress, validatorAddr sdk.ValAddress,
	creationHeight int64, minTime time.Time, balance sdk.Int,
) UnbondingDelegation {

	return UnbondingDelegation{
		DelegatorAddress: delegatorAddr,
		ValidatorAddress: validatorAddr,
		Entries: []UnbondingDelegationEntry{
			NewUnbondingDelegationEntry(creationHeight, minTime, balance),
		},
	}
}

// AddEntry - append entry to the unbonding delegation
func (d *UnbondingDelegation) AddEntry(creationHeight int64, minTime time.Time, balance sdk.Int) {
	entry := NewUnbondingDelegationEntry(creationHeight, minTime, balance)
	d.Entries = append(d.Entries, entry)
}

// RemoveEntry - remove entry at index i to the unbonding delegation
func (d *UnbondingDelegation) RemoveEntry(i int64) {
	d.Entries = append(d.Entries[:i], d.Entries[i+1:]...)
}

// return the unbonding delegation
func MustMarshalUBD(cdc codec.Marshaler, ubd UnbondingDelegation) []byte {
	return cdc.MustMarshalBinaryLengthPrefixed(&ubd)
}

// unmarshal a unbonding delegation from a store value
func MustUnmarshalUBD(cdc codec.Marshaler, value []byte) UnbondingDelegation {
	ubd, err := UnmarshalUBD(cdc, value)
	if err != nil {
		panic(err)
	}
	return ubd
}

// unmarshal a unbonding delegation from a store value
func UnmarshalUBD(cdc codec.Marshaler, value []byte) (ubd UnbondingDelegation, err error) {
	err = cdc.UnmarshalBinaryLengthPrefixed(value, &ubd)
	return ubd, err
}

// String returns a human readable string representation of an UnbondingDelegation.
func (d UnbondingDelegation) String() string {
	out := fmt.Sprintf(`Unbonding Delegations between:
  Delegator:                 %s
  Validator:                 %s
	Entries:`, d.DelegatorAddress, d.ValidatorAddress)
	for i, entry := range d.Entries {
		out += fmt.Sprintf(`    Unbonding Delegation %d:
      Creation Height:           %v
      Min time to unbond (unix): %v
      Expected balance:          %s`, i, entry.CreationHeight,
			entry.CompletionTime, entry.Balance)
	}
	return out
}

// UnbondingDelegations is a collection of UnbondingDelegation
type UnbondingDelegations []UnbondingDelegation

func (ubds UnbondingDelegations) String() (out string) {
	for _, u := range ubds {
		out += u.String() + "\n"
	}
	return strings.TrimSpace(out)
}

func NewRedelegationEntry(creationHeight int64, completionTime time.Time, balance sdk.Int, sharesDst sdk.Dec) RedelegationEntry {
	return RedelegationEntry{
		CreationHeight: creationHeight,
		CompletionTime: completionTime,
		InitialBalance: balance,
		SharesDst:      sharesDst,
	}
}

// String implements the Stringer interface for a RedelegationEntry object.
func (red RedelegationEntry) String() string {
	out, _ := yaml.Marshal(red)
	return string(out)
}

// IsMature - is the current entry mature
func (e RedelegationEntry) IsMature(currentTime time.Time) bool {
	return !e.CompletionTime.After(currentTime)
}

func NewRedelegation(
	delegatorAddr sdk.AccAddress, validatorSrcAddr, validatorDstAddr sdk.ValAddress,
	creationHeight int64, minTime time.Time, balance sdk.Int, sharesDst sdk.Dec,
) Redelegation {

	return Redelegation{
		DelegatorAddress:    delegatorAddr,
		ValidatorSrcAddress: validatorSrcAddr,
		ValidatorDstAddress: validatorDstAddr,
		Entries: []RedelegationEntry{
			NewRedelegationEntry(creationHeight, minTime, balance, sharesDst),
		},
	}
}

// AddEntry - append entry to the unbonding delegation
func (d *Redelegation) AddEntry(creationHeight int64, minTime time.Time, balance sdk.Int, sharesDst sdk.Dec) {
	entry := NewRedelegationEntry(creationHeight, minTime, balance, sharesDst)
	d.Entries = append(d.Entries, entry)
}

// RemoveEntry - remove entry at index i to the unbonding delegation
func (d *Redelegation) RemoveEntry(i int64) {
	d.Entries = append(d.Entries[:i], d.Entries[i+1:]...)
}

// MustMarshalRED returns the Redelegation bytes. Panics if fails.
func MustMarshalRED(cdc codec.Marshaler, red Redelegation) []byte {
	return cdc.MustMarshalBinaryLengthPrefixed(&red)
}

// MustUnmarshalRED unmarshals a redelegation from a store value. Panics if fails.
func MustUnmarshalRED(cdc codec.Marshaler, value []byte) Redelegation {
	red, err := UnmarshalRED(cdc, value)
	if err != nil {
		panic(err)
	}
	return red
}

// UnmarshalRED unmarshals a redelegation from a store value
func UnmarshalRED(cdc codec.Marshaler, value []byte) (red Redelegation, err error) {
	err = cdc.UnmarshalBinaryLengthPrefixed(value, &red)
	return red, err
}

// String returns a human readable string representation of a Redelegation.
func (d Redelegation) String() string {
	out := fmt.Sprintf(`Redelegations between:
  Delegator:                 %s
  Source Validator:          %s
  Destination Validator:     %s
  Entries:
`,
		d.DelegatorAddress, d.ValidatorSrcAddress, d.ValidatorDstAddress,
	)

	for i, entry := range d.Entries {
		out += fmt.Sprintf(`    Redelegation Entry #%d:
      Creation height:           %v
      Min time to unbond (unix): %v
      Dest Shares:               %s
`,
			i, entry.CreationHeight, entry.CompletionTime, entry.SharesDst,
		)
	}

	return strings.TrimRight(out, "\n")
}

// Redelegations are a collection of Redelegation
type Redelegations []Redelegation

func (d Redelegations) String() (out string) {
	for _, red := range d {
		out += red.String() + "\n"
	}
	return strings.TrimSpace(out)
}

// ----------------------------------------------------------------------------
// Client Types

// DelegationResponse is equivalent to Delegation except that it contains a balance
// in addition to shares which is more suitable for client responses.
type DelegationResponse struct {
	Delegation
	Balance sdk.Coin `json:"balance" yaml:"balance"`
}

// NewDelegationResp creates a new DelegationResponse instance
func NewDelegationResp(
	delegatorAddr sdk.AccAddress, validatorAddr sdk.ValAddress, shares sdk.Dec, balance sdk.Coin,
) DelegationResponse {
	return DelegationResponse{
		Delegation: NewDelegation(delegatorAddr, validatorAddr, shares),
		Balance:    balance,
	}
}

// String implements the Stringer interface for DelegationResponse.
func (d DelegationResponse) String() string {
	return fmt.Sprintf("%s\n  Balance:   %s", d.Delegation.String(), d.Balance)
}

type delegationRespAlias DelegationResponse

// MarshalJSON implements the json.Marshaler interface. This is so we can
// achieve a flattened structure while embedding other types.
func (d DelegationResponse) MarshalJSON() ([]byte, error) {
	return json.Marshal((delegationRespAlias)(d))
}

// UnmarshalJSON implements the json.Unmarshaler interface. This is so we can
// achieve a flattened structure while embedding other types.
func (d *DelegationResponse) UnmarshalJSON(bz []byte) error {
	return json.Unmarshal(bz, (*delegationRespAlias)(d))
}

// DelegationResponses is a collection of DelegationResp
type DelegationResponses []DelegationResponse

// String implements the Stringer interface for DelegationResponses.
func (d DelegationResponses) String() (out string) {
	for _, del := range d {
		out += del.String() + "\n"
	}
	return strings.TrimSpace(out)
}

// RedelegationResponse is equivalent to a Redelegation except that its entries
// contain a balance in addition to shares which is more suitable for client
// responses.
type RedelegationResponse struct {
	Redelegation
	Entries []RedelegationEntryResponse `json:"entries" yaml:"entries"`
}

// NewRedelegationResponse crates a new RedelegationEntryResponse instance.
func NewRedelegationResponse(
	delegatorAddr sdk.AccAddress, validatorSrc, validatorDst sdk.ValAddress, entries []RedelegationEntryResponse,
) RedelegationResponse {
	return RedelegationResponse{
		Redelegation: Redelegation{
			DelegatorAddress:    delegatorAddr,
			ValidatorSrcAddress: validatorSrc,
			ValidatorDstAddress: validatorDst,
		},
		Entries: entries,
	}
}

// RedelegationEntryResponse is equivalent to a RedelegationEntry except that it
// contains a balance in addition to shares which is more suitable for client
// responses.
type RedelegationEntryResponse struct {
	RedelegationEntry
	Balance sdk.Int `json:"balance"`
}

// NewRedelegationEntryResponse creates a new RedelegationEntryResponse instance.
func NewRedelegationEntryResponse(
	creationHeight int64, completionTime time.Time, sharesDst sdk.Dec, initialBalance, balance sdk.Int) RedelegationEntryResponse {
	return RedelegationEntryResponse{
		RedelegationEntry: NewRedelegationEntry(creationHeight, completionTime, initialBalance, sharesDst),
		Balance:           balance,
	}
}

// String implements the Stringer interface for RedelegationResp.
func (r RedelegationResponse) String() string {
	out := fmt.Sprintf(`Redelegations between:
  Delegator:                 %s
  Source Validator:          %s
  Destination Validator:     %s
  Entries:
`,
		r.DelegatorAddress, r.ValidatorSrcAddress, r.ValidatorDstAddress,
	)

	for i, entry := range r.Entries {
		out += fmt.Sprintf(`    Redelegation Entry #%d:
      Creation height:           %v
      Min time to unbond (unix): %v
      Initial Balance:           %s
      Shares:                    %s
      Balance:                   %s
`,
			i, entry.CreationHeight, entry.CompletionTime, entry.InitialBalance, entry.SharesDst, entry.Balance,
		)
	}

	return strings.TrimRight(out, "\n")
}

type redelegationRespAlias RedelegationResponse

// MarshalJSON implements the json.Marshaler interface. This is so we can
// achieve a flattened structure while embedding other types.
func (r RedelegationResponse) MarshalJSON() ([]byte, error) {
	return json.Marshal((redelegationRespAlias)(r))
}

// UnmarshalJSON implements the json.Unmarshaler interface. This is so we can
// achieve a flattened structure while embedding other types.
func (r *RedelegationResponse) UnmarshalJSON(bz []byte) error {
	return json.Unmarshal(bz, (*redelegationRespAlias)(r))
}

// RedelegationResponses are a collection of RedelegationResp
type RedelegationResponses []RedelegationResponse

func (r RedelegationResponses) String() (out string) {
	for _, red := range r {
		out += red.String() + "\n"
	}
	return strings.TrimSpace(out)
}
