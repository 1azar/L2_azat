package walletFacade

import ( // фасад реализует работу с множеством иных зависимостей
	account2 "L2_azat/patterns/facade/account"
	ledger2 "L2_azat/patterns/facade/ledger"
	notification2 "L2_azat/patterns/facade/notification"
	securityCode2 "L2_azat/patterns/facade/securityCode"
	wallet2 "L2_azat/patterns/facade/wallet"
	"fmt"
)

// Нащ фасад, который знает как работать с компонентами сложной подсистемы

type WalletFacade struct {
	account      *account2.Account
	wallet       *wallet2.Wallet
	securityCode *securityCode2.SecurityCode
	notification *notification2.Notification
	ledger       *ledger2.Ledger
}

func NewWalletFacade(accountID string, code int) *WalletFacade {
	fmt.Println("Starting createEvent account")
	walletFacacde := &WalletFacade{
		account:      account2.NewAccount(accountID),
		securityCode: securityCode2.NewSecurityCode(code),
		wallet:       wallet2.NewWallet(),
		notification: &notification2.Notification{},
		ledger:       &ledger2.Ledger{},
	}
	fmt.Println("Account created")
	return walletFacacde
}

func (w *WalletFacade) AddMoneyToWallet(accountID string, securityCode int, amount int) error {
	fmt.Println("Starting add money to wallet")
	err := w.account.CheckAccount(accountID)
	if err != nil {
		return err
	}
	err = w.securityCode.CheckCode(securityCode)
	if err != nil {
		return err
	}
	w.wallet.CreditBalance(amount)
	w.notification.SendWalletCreditNotification()
	w.ledger.MakeEntry(accountID, "credit", amount)
	return nil
}

func (w *WalletFacade) DeductMoneyFromWallet(accountID string, securityCode int, amount int) error {
	fmt.Println("Starting debit money from wallet")
	err := w.account.CheckAccount(accountID)
	if err != nil {
		return err
	}

	err = w.securityCode.CheckCode(securityCode)
	if err != nil {
		return err
	}
	err = w.wallet.DebitBalance(amount)
	if err != nil {
		return err
	}
	w.notification.SendWalletDebitNotification()
	w.ledger.MakeEntry(accountID, "debit", amount)
	return nil
}
