package transaction

type IncomeReport struct {
	Period            string       `json:"period"`
	TotalIncome       int64        `json:"total_income"`
	TotalFines        int64        `json:"total_fines"`
	MembershipRevenue int64        `json:"membership_revenue"`
	Transactions      int          `json:"transactions"`
	Breakdown         []IncomeLine `json:"breakdown"`
}

type IncomeLine struct {
	Date         string `json:"date"`
	Income       int64  `json:"income"`
	Fines        int64  `json:"fines"`
	Transactions int    `json:"transactions"`
}
