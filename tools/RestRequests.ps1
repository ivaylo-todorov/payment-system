
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
    merchant, first merchant, merchant@email.com
"@

    $res = Invoke-RestMethod -Method 'Post' -Uri $Api -Body $Body -Headers $Headers

    Write-Output $res | ConvertTo-Json
}

function PostMerchant {
    param (
    )

    $MerchantId = "6806cc4c-bf02-4669-aac2-91ad9f13f6fb"

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
        uuid = "170e46f9-b11b-4bd5-be14-59c3ddb2afde"
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
    name, description, email
    root, root user, root@email.com
"@

    $res = Invoke-RestMethod -Method 'Post' -Uri $Api -Body $Body -Headers $Headers

    Write-Output $res | ConvertTo-Json
}


DisplayTransactions
