
$Remote = "http://localhost:8080/"


CreateAdmin


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
    name, description, email
"@

    $res = Invoke-RestMethod -Method 'Post' -Uri $Api -Body $Body -Headers $Headers

    Write-Output $res
}

function PostMerchant {
    param (
    )

    $MerchantId = "1"

    $Api = $Remote + "merchants/" + $MerchantId

    $Headers = @{
        'Content-Type'='application/json'
    }

    $Merchant = @{
        name = "Merchant One"
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
        uuid = "2"
        name = "Merchant Two"
    }

    $Body = @{
        merchant = $Merchant
    }

    $res = Invoke-RestMethod -Method 'Delete' -Uri $Api -Body ($Body|ConvertTo-Json) -Headers $Headers

    Write-Output $res | ConvertTo-Json
}

function CreateAdmin {
    param (
    )

    $Api = $Remote + "admins"

    $Headers = @{
        'Content-Type'='text/csv'
    }

    $Body = @"
    name, description, email
"@

    $res = Invoke-RestMethod -Method 'Post' -Uri $Api -Body $Body -Headers $Headers

    Write-Output $res
}
