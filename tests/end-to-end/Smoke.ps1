. .\..\Requests.ps1
. .\..\Tools.ps1

# TODO:
#   check with various csv fields len
#   check for existing email collision


$global:Hostname = "localhost"
$global:Port = 8080

$global:MerchantOneId= ""

function TestCreateMerchant {
    $Name = RandomString
    $Email = "$name@email.com"

    $res = CreateMerchant -Hostname $global:Hostname -Port $global:Port -Name $Name -Description "" -Email $Email -Status "active"

    $global:MerchantOneId = $res.merchants[0].uuid

    $res.merchants[0] | ConvertTo-Json
}

$global:MerchantTwoId = ""

function TestUpdateMerchant {
    $Name = RandomString
    $Email = "$name@email.com"

    $res = CreateMerchant -Hostname $global:Hostname -Port $global:Port -Name $Name -Description "" -Email $Email -Status "inactive"

    $global:MerchantTwoId = $res.merchants[0].uuid

    $res = PostMerchant -Id $global:MerchantTwoId -Status "active"

    If ($res.merchants[0].status -ne "active") {
        Write-Output "Error changeding merchant status"
        return
    }

    $res.merchants[0] | ConvertTo-Json
}

$global:AuthorizeTransactionOneId = ""

function TestCreateAuthorizeTransactionOne {
    $res = PostTransaction -Hostname $global:Hostname -Port $global:Port -MerchantId $global:MerchantOneId -Type "authorize" -Amount 100 -CustomerMail "customer@email.com"

    $global:AuthorizeTransactionOneId = $res.transactions[0].uuid

    $res.transactions[0] | ConvertTo-Json
}

$ChargeTransactionId = ""

function TestCreateChargeTransaction {
    $res = PostTransaction -Hostname $global:Hostname -Port $global:Port -ParentId $global:AuthorizeTransactionOneId -MerchantId $global:MerchantOneId -Type "charge" -Amount 100 -CustomerMail "customer@email.com"

    $global:ChargeTransactionId = $res.transactions[0].uuid

    $res.transactions[0] | ConvertTo-Json
}

$global:RefundTransactionId = ""

function TestCreateRefundTransaction {
    $res = PostTransaction -Hostname $global:Hostname -Port $global:Port -ParentId $global:ChargeTransactionId -MerchantId $global:MerchantOneId -Type "refund" -Amount 100 -CustomerMail "customer@email.com"

    $global:RefundTransactionId = $res.transactions[0].uuid

    $res.transactions[0] | ConvertTo-Json
}

$global:AuthorizeTransactionTwoId = ""

function TestCreateAuthorizeTransactionTwo {
    $res = PostTransaction -Hostname $global:Hostname -Port $global:Port -MerchantId $global:MerchantTwoId -Type "authorize" -Amount 100 -CustomerMail "customer@email.com"

    $global:AuthorizeTransactionTwoId = $res.transactions[0].uuid

    $res.transactions[0] | ConvertTo-Json
}

$global:ReversalTransactionId = ""

function TestCreateReversalTransaction {
    $res = PostTransaction -Hostname $global:Hostname -Port $global:Port -ParentId $global:AuthorizeTransactionTwoId -MerchantId $global:MerchantTwoId -Type "reversal" -Amount 0 -CustomerMail "customer@email.com"

    $global:ReversalTransactionId = $res.transactions[0].uuid

    $res.transactions[0] | ConvertTo-Json
}


TestCreateMerchant
TestUpdateMerchant

TestCreateAuthorizeTransactionOne
TestCreateChargeTransaction
TestCreateRefundTransaction

TestCreateAuthorizeTransactionTwo
TestCreateReversalTransaction

TestCreateAuthorizeTransactionOne
TestCreateChargeTransaction

DisplayMerchants
DisplayTransactions
