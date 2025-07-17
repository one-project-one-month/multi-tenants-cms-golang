package utils

import (
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/sirupsen/logrus"
	"math"
	"math/big"
)

func NumericToInt64(n pgtype.Numeric) int64 {
	return n.Int.Int64() * int64(math.Pow10(int(n.Exp)))
}

func NumericToUint64(n pgtype.Numeric) uint64 {
	if !n.Valid {
		return 0
	}

	if n.Exp == 0 {
		return n.Int.Uint64()
	}

	// If there's an exponent, handle it carefully
	f := new(big.Float).SetInt(n.Int)
	if n.Exp > 0 {
		m := new(big.Float).SetInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(n.Exp)), nil))
		f.Mul(f, m)
	} else {
		d := new(big.Float).SetInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(-n.Exp)), nil))
		f.Quo(f, d)
	}

	val, _ := f.Uint64()
	return val
}

func NumericToFloat64(n pgtype.Numeric) float64 {
	if !n.Valid {
		return 0
	}

	intValue := n.Int.Int64()
	scale := int(n.Exp)

	if scale >= 0 {
		return float64(intValue * int64(math.Pow10(scale)))
	} else {
		return float64(intValue) / math.Pow10(-scale)
	}
}

func Float64ToBigInt(f float64) *big.Int {
	return big.NewInt(int64(f))
}

func Uint64ToBigInt(u uint64) *big.Int {
	return big.NewInt(int64(u))
}

func Uint64ToFloat64(u uint64) float64 {
	return float64(u)
}

func Uint64ToBigIntWithPrecision(u uint64, precision int) *big.Int {
	if precision < 0 {
		logrus.WithFields(logrus.Fields{
			"precision": precision,
		}).Error("Precision must be non-negative")
		return big.NewInt(0)
	}

	f := float64(u) / math.Pow10(precision)
	return big.NewInt(int64(f))
}

func Int32ToBigInt(i int32) *big.Int {
	return big.NewInt(int64(i))
}

func Float64ToNumeric(f float64) pgtype.Numeric {
	n := pgtype.Numeric{}
	err := n.Scan(f)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"value": f,
		}).WithError(err).Error("Failed to scan float64 into pgtype.Numeric")
		return pgtype.Numeric{}
	}
	return n
}

func Int64ToNumeric(i int64) pgtype.Numeric {
	n := pgtype.Numeric{}
	err := n.Scan(i)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"value": i,
		}).WithError(err).Error("Failed to scan int64 into pgtype.Numeric")
		return pgtype.Numeric{}
	}
	return n
}
