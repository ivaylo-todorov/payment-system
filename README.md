# payment-system
Payment  System Task

## Usage
There is no UI yet

Requests to the service can be made with powershell cmdlets. For example:

To Create Merchant

```
CreateMerchant -Hostname "localhost" -Port 8080 -Name "merchant name" -Description "" -Email "merchant@email.com" -Status "active"
```

To Create Authorize Transaction

```
PostTransaction -Hostname -MerchantId [from CreateMerchant] -Type "authorize" -Amount 100 -CustomerMail "customer@email.com"
```

To Create Charge Transaction

```
PostTransaction -Hostname -ParentId [authorize transaction id] -MerchantId <from CreateMerchant> -Type "charge" -Amount 100 -CustomerMail "customer@email.com"
```

To show all merchants or all transactions

```
DisplayMerchants
DisplayTransactions
```

Default hostname and port are 'localhost' and '8080' and can be omitted 

Transactions and Merchants can be displayed in the browser also:

```
http://localhost:8080/transactions?render
http://localhost:8080/merchants?render
```

## To run server

From main dir:

```
go run .
```

starts the server listening on http://localhost:8080/


## Tests

There are very basic smoke tests for powershell in \tests\end-to-end\

## TODO

- Databese transaction table should be changed to use 'Self-Referential Has-oe' to take advantage of foreign key contraints
- Add authentication layer

