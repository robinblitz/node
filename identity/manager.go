// Maps Ethereum account to dto.Identity.
// Currently creates a new eth account with password on CreateNewIdentity().

package identity

import (
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
)

type identityManager struct {
	keystoreManager keystoreInterface
}

func NewIdentityManager(keystore keystoreInterface) *identityManager {
	return &identityManager{
		keystoreManager: keystore,
	}
}

func accountToIdentity(account accounts.Account) Identity {
	identity := FromAddress(account.Address.Hex())
	return identity
}

func identityToAccount(identity Identity) accounts.Account {
	return accounts.Account{
		Address: common.HexToAddress(identity.Address),
	}
}

func (idm *identityManager) CreateNewIdentity(passphrase string) (identity Identity, err error) {
	account, err := idm.keystoreManager.NewAccount(passphrase)
	if err != nil {
		return identity, err
	}

	return accountToIdentity(account), nil
}

func (idm *identityManager) GetIdentities() []Identity {
	accountList := idm.keystoreManager.Accounts()

	var ids = make([]Identity, len(accountList))
	for i, account := range accountList {
		ids[i] = accountToIdentity(account)
	}

	return ids
}

func (idm *identityManager) GetSigner(identity Identity) Signer {
	return NewSigner(idm.keystoreManager, identity)
}

func (idm *identityManager) GetIdentity(identityString string) (identity Identity, err error) {
	account, err := idm.keystoreManager.Find(identityToAccount(FromAddress(identityString)))
	if err != nil {
		return identity, errors.New("identity not found")
	}

	identity = accountToIdentity(account)

	return identity, nil
}

func (idm *identityManager) HasIdentity(identityString string) bool {
	_, err := idm.GetIdentity(identityString)

	return err == nil
}
