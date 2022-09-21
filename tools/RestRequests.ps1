
$Remote = "http://localhost:8080/"


# Methods

function DisplayTransactions {
    param (
    )

    $Api = $Remote + "transactions"
    
    $res = Invoke-RestMethod -Method 'Get' -Uri $Api

    Write-Output $res | ConvertTo-Json
}

function PostTransaction {
    param (
    )

    $Api = $Remote + "transactions"

    $Headers = @{
        'Content-Type'='application/json'
    }

    $Transaction = @{
        merchant_uuid = "d68e2827-4d7c-4c84-823c-4e63572732dd"
        type = "charge"
        amount = 123
        customer_email = "customer@email.com"
    }
    
    $Body = @{
        transaction = $Transaction
    }

    $res = Invoke-RestMethod -Method 'Post' -Uri $Api -Body ($Body|ConvertTo-Json) -Headers $Headers

    Write-Output $res | ConvertTo-Json
}

function DisplayMerchants {
    param (
    )

    $Api = $Remote + "merchants"
    
    $res = Invoke-RestMethod -Method 'Get' -Uri $Api

    Write-Output $res | ConvertTo-Json
}

function CreateMerchant {
    param (
    )

    $Api = $Remote + "merchants"

    $Headers = @{
        'Content-Type'='text/csv'
    }

    $Body = @"
    merchant2, second merchant, merchant2@email.com
"@

    $res = Invoke-RestMethod -Method 'Post' -Uri $Api -Body $Body -Headers $Headers

    Write-Output $res | ConvertTo-Json
}

function PostMerchant {
    param (
    )

    $MerchantId = "d68e2827-4d7c-4c84-823c-4e63572732dd"

    $Api = $Remote + "merchants/" + $MerchantId

    $Headers = @{
        'Content-Type'='application/json'
    }

    $Merchant = @{
        description = "Merchant Two"
    }

    $Body = @{
        merchant = $Merchant
    }

    $res = Invoke-RestMethod -Method 'Post' -Uri $Api -Body ($Body|ConvertTo-Json) -Headers $Headers

    Write-Output $res | ConvertTo-Json
}

function DeleteMerchant {
    param (
    )

    $Api = $Remote + "merchants"

    $Headers = @{
        'Content-Type'='application/json'
    }

    $Merchant = @{
        uuid = "9addb6a6-be2d-4fff-a883-a42ab761dfa5"
        name = "Merchant Two"
    }

    $Body = @{
        merchant = $Merchant
    }

    Invoke-RestMethod -Method 'Delete' -Uri $Api -Body ($Body|ConvertTo-Json) -Headers $Headers
}

function CreateAdmin {
    param (
    )

    $Api = $Remote + "admins"

    $Headers = @{
        'Content-Type'='text/csv'
    }

    $Body = @"
    root, root user, root@email.com
"@

    $res = Invoke-RestMethod -Method 'Post' -Uri $Api -Body $Body -Headers $Headers

    Write-Output $res | ConvertTo-Json
}


# PostMerchant

# DisplayMerchants

DisplayTransactions

# PostTransaction