# CHANGELOG

## 0.1.0 (April 27, 2020)

#### Additions
- `NewScanner` now accepts an `Option` variadic argument.
- A `pgx.ErrNoRows` error is returned on a `Query` that returns a length of zero.
    - Identical functionality to the `Row` interface. Allows for app conditional logic if query returns nothing.
    - Can be turned off by passing in `ErrNoRowsQuery(false)` into the `NewScanner` function.  

#### Breaking Changes
- Scanner no longer accepts a non initialized value for the Destination. 
  - An error will be returned during validation if a non initialized destination is passed into the scanner.
  - Reason for breaking change: No longer need to do a recursive call, which makes closing the scanner easier and allow for additional defered logic.