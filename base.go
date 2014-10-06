// Use of this source code is governed by the CC0 1.0
// license that can be found in the LICENSE file or here:
// http://creativecommons.org/publicdomain/zero/1.0/

// Package base implements encoding and decoding of data in to specified base encoded data.
//
// Package base handles any base between base2 and base62 eg. base6, base32, base36, base60 etc.
//
// The base ASCII representation is probably not compatible with any other implementation of the corresponding base.
// For a standard implementations of base32 please see encoding/base32 in the standard lib.
package base

import (
	"log"
	"math/big"
)

const digits = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

// Precalculated buffer size multipliers for all valid bases
var bufferSizeMultiplier = [...]float64{0, 0, 8, 5.052, 4, 3.452, 3.097, 2.852, 2.671, 2.53, 2.413, 2.317, 2.233, 2.168, 2.104, 2.052, 2, 1.962, 1.923, 1.884, 1.852, 1.826, 1.8, 1.775, 1.749, 1.73, 1.704, 1.684, 1.665, 1.652, 1.633, 1.62, 1.6, 1.588, 1.575, 1.562, 1.549, 1.542, 1.53, 1.517, 1.504, 1.497, 1.483, 1.478, 1.471, 1.459, 1.452, 1.446, 1.439, 1.426, 1.42, 1.413, 1.407, 1.4, 1.394, 1.388, 1.381, 1.375, 1.368, 1.362, 1.355, 1.355, 1.349}

// Encode takes an []byte b containing byte data and returns []byte r containing base b encoded data.
//
// b can not be grater than 62 or less than 2. If b is over 36 then r is case sensitive.
func Encode(u []byte, b int) (r []byte) {

	if b < 2 || b > len(digits) {
		log.Fatalln("Illegal Encode base")
	}

	a := big.NewInt(0).SetBytes(u)
	base := big.NewInt(int64(b))
	rem := big.NewInt(0)

	// Calculate the necessary buffer size for the defined base
	i := int((float64(len(a.Bytes())))*bufferSizeMultiplier[b]) + 1
	d := make([]byte, i)

	for a.Cmp(base) >= 0 {
		i--
		a.QuoRem(a, base, rem)
		d[i] = digits[int(rem.Int64())]
	}

	// last character when a < base
	i--
	d[i] = digits[int(a.Int64())]

	return d[i:]
}

// Decode takes an []byte u containing base b encoded data and returns []byte r containing byte data.
//
// b can not be grater than 62 or less than 2. If b is over 36 then u is case sensitive.
//
// u may not contain characters outside of the base character representation, e.g. base 2 can only contain "0" and "1" while base62 can only contain 0-9a-zA-Z.
func Decode(u []byte, b int) (r []byte) {

	if b < 2 || b > len(digits) {
		log.Fatalln("Illegal Decode base")
	}

	base := big.NewInt(int64(b))
	v := big.NewInt(0)
	n := big.NewInt(0)

	for i := 0; i < len(u); i++ {
		u := u[i]
		switch {
		case '0' <= u && u <= '9':
			v.SetInt64(int64(u - '0'))
		case 'a' <= u && u <= 'z':
			v.SetInt64(int64(u - 'a' + 10))
		case 'A' <= u && u <= 'Z':
			if b <= 36 {
				v.SetInt64(int64(u - 'A' + 10))
			} else {
				v.SetInt64(int64(u - 'A' + 36))
			}
		}
		if v.Int64() >= base.Int64() {
			// If a parser uses base.Decode he can use recover() to handle this otherwise fatal error.
			// See http://blog.golang.org/defer-panic-and-recover for more information.
			panic("Illegal characters in u []byte")
		}
		n.Mul(n, base)
		n.Add(n, v)
	}
	return n.Bytes()
}
