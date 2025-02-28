package components

import (
	"bytes"
	"fmt"

	"github.com/Rhymond/go-money"
	"github.com/olekukonko/tablewriter"
	"github.com/samber/lo"
	"github.com/yohamta/donburi"
)

type WalletTX struct {
	FromID    string `json:"from_id"`
	From      string `json:"from"`
	ToID      string `json:"to_id"`
	To        string `json:"to"`
	Price     int64  `json:"price"`
	CreatedAt int64  `json:"created_at"`
	Due       int64  `json:"due"`
	Note      string `json:"note"`
}

type WalletData struct {
	Txs       []WalletTX `json:"Txs"`
	Pending   []WalletTX `json:"Pending"`
	balance   int64      `json:"balance"`
	mn        *money.Money
	txtblbf   *bytes.Buffer
	tblwriter *tablewriter.Table
}

func NewWalletData() *WalletData {
	buf := bytes.NewBuffer([]byte{})
	writer := tablewriter.NewWriter(buf)

	return &WalletData{
		Txs:       []WalletTX{},
		Pending:   []WalletTX{},
		mn:        money.New(0, money.USD),
		txtblbf:   buf,
		tblwriter: writer,
	}
}

func (w *WalletData) Init() {
	buf := bytes.NewBuffer([]byte{})
	writer := tablewriter.NewWriter(buf)

	w.txtblbf = buf
	w.tblwriter = writer
	w.SetBalance()
	w.mn = money.New(w.balance, money.USD)
}

func (w *WalletData) AddPendingTX(tx WalletTX) {
	w.Pending = append(w.Pending, tx)
}

func (w *WalletData) ReconcilePendingTxs() {
	for _, tx := range w.Pending {
		w.mn.Add(money.New(tx.Price, money.USD))
		w.Txs = append(w.Txs, tx)
	}
	w.Pending = []WalletTX{}
}

func (w *WalletData) PendingTxCnt() int {
	return len(w.Pending)
}

func (w *WalletData) TxCnt() int {
	return len(w.Txs)
}

func (w *WalletData) LatestTransactionsTable(limit int) string {
	if len(w.Txs) == 0 {
		return ""
	}
	tmp := [][]string{}
	for i := len(w.Txs) - 1; i >= 0 || len(tmp) == limit; i-- {
		if i < 0 {
			break
		}
		money := money.New(w.Txs[i].Price, money.USD)
		tmp = append(tmp, []string{
			fmt.Sprintf("%v", w.Txs[i].CreatedAt),
			trimTXFrom(w.Txs[i].From),
			w.Txs[i].Note,
			money.Display(),
		})
	}
	if len(tmp) != len(w.Txs) {
		tmp = append(tmp, []string{"", "", "", "..."})
	}

	w.txtblbf.Reset()
	w.tblwriter.SetHeader([]string{"Date", "Origin", "Note", "Amount"})
	w.tblwriter.ClearRows()
	w.tblwriter.AppendBulk(tmp)
	w.tblwriter.Render()
	return w.txtblbf.String()
}

func trimTXFrom(from string) string {
	if len(from) <= 12 {
		return from
	} else {
		return from[:12]
	}
}

func (w *WalletData) AddTX(tx WalletTX) (int64, error) {
	var err error
	w.Txs = append(w.Txs, tx)
	w.mn, err = w.mn.Add(money.New(tx.Price, money.USD))
	return w.mn.Amount(), err
}

func (w *WalletData) SetBalance() {
	w.balance = lo.SumBy(w.Txs, func(tx WalletTX) int64 {
		return tx.Price
	})
}

func (w *WalletData) Balance() int64 {
	w.balance = 0
	for _, tx := range w.Txs {
		w.balance += tx.Price
	}
	return w.balance
}

func (w *WalletData) BalanceDisplay() string {
	return w.mn.Display()
}

var Wallet = donburi.NewComponentType[WalletData]()
