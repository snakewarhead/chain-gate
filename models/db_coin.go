package models

type Coin struct {
	ID            int
	Name          string
	Enable        bool
	MainAddress   string
	PublicKey     string
	Password      string
	APIURL        string
	APIWalletURL  string
	ConfirmNum    int64
	PositionQuery int64
}

func GetCoinEnabled() (*Coin, error) {
	c := &Coin{}

	const sql = "SELECT * FROM coin WHERE enable = 1"
	row := db.QueryRow(sql)
	err := row.Scan(dbColumns(c)...)
	return c, err
}

func UpdateCoin(c *Coin) error {
	stmt, err := db.Prepare("UPDATE coin SET position_query = ? WHERE id = ?")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(
		c.PositionQuery,
		c.ID,
	)
	if err != nil {
		return err
	}

	return nil
}
