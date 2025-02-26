package util

import (
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/shopspring/decimal"
)

// DecimalToNumeric конвертирует decimal.Decimal в pgtype.Numeric
func DecimalToNumeric(d decimal.Decimal) (pgtype.Numeric, error) {
	numStr := d.String()
	var numeric pgtype.Numeric
	err := numeric.Scan(numStr)
	if err != nil {
		return pgtype.Numeric{}, fmt.Errorf("failed to convert decimal to numeric: %v", err)
	}
	return numeric, nil
}

// NumericToDecimal конвертирует pgtype.Numeric в decimal.Decimal
func NumericToDecimal(n pgtype.Numeric) (decimal.Decimal, error) {
	numStr, err := n.Value()
	if err != nil {
		return decimal.Decimal{}, fmt.Errorf("failed to convert numeric to string: %v", err)
	}
	return decimal.NewFromString(numStr.(string))
}
