package ledger

import "fmt"

// Часть сложной подсистемы

type Ledger struct {
}

func (s *Ledger) MakeEntry(accountID, txnType string, amount int) {
	fmt.Printf("Make ledger entry for accountId %s with txnType %s for amount %d\n", accountID, txnType, amount)
	return
}
