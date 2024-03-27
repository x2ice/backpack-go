package backpack

import "errors"

var ErrCredentialsRequired = errors.New("api key and secret are required")

var ErrInvalidOrderSide = errors.New("order side must be Ask or Bid")

var ErrInsufficientFunds = errors.New("insufficient funds")
