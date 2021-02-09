# CHANGELOG

## 0.3.0 (February 9, 2021)

#### Additions
- Unit tests and integration tests have now been separated. To run integration tests `go test -v --tags=integration ./...`
- Dockerized testing. Testing is now easier thanks to PR #6. Simply run `make test` to run all unit and integration tests
- Add SQL column notation for simplified aliases. Thanks to PR #7 you can now have pgxscan determine embedded struct notation by using leveraging sql columns. [See documentation for more details.](https://github.com/randallmlough/pgxscan#user-content-example-with-aliasing)
  
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