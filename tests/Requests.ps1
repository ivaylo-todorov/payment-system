
# Methods

function GetTransactions {
    param (
        $Hostname = "localhost",
        $Port = 8080
    )

    $Api = "http://$($Hostname):$Port/" + "transactions"
    
    Invoke-RestMethod -Method 'Get' -Uri $Api
}

function DisplayTransactions {
    param (
        $Hostname = "localhost",
        $Port = 8080
    )

    $res = GetTransactions -Hostname $Hostname -Port $Port

    $res.transactions | ForEach-Object {[PSCustomObject]$_} | Format-Table
}

function PostTransaction {
    param (
        $Hostname = "localhost",
        $Port = 8080,

        $MerchantId,
        $ParentId,
        $Type,
        $Amount,
        $CustomerMail,
        $CustomerPhone
    )

    $Api = "http://$($Hostname):$Port/" + "transactions"

    $Headers = @{
        'Content-Type'='application/json'
    }

    $Transaction = @{
        merchant_uuid = $MerchantId
        parent_uuid = $ParentId
        type = $Type
        amount = $Amount
        customer_email = $CustomerMail
        customer_phone = $CustomerPhone
    }
    
    $Body = @{
        transaction = $Transaction
    }

    Invoke-RestMethod -Method 'Post' -Uri $Api -Body ($Body|ConvertTo-Json) -Headers $Headers
}

function GetMerchants {
    param (
        $Hostname = "localhost",
        $Port = 8080
    )

    $Api = "http://$($Hostname):$Port/" + "merchants"
    
    Invoke-RestMethod -Method 'Get' -Uri $Api
}

function DisplayMerchants {
    param (
        $Hostname = "localhost",
        $Port = 8080
    )

    $res = GetMerchants -Hostname $Hostname -Port $Port

    $res.merchants | ForEach-Object {[PSCustomObject]$_} | Format-Table
}

function CreateMerchant {
    param (
        $Hostname = "localhost",
        $Port = 8080,

        $Name,
        $Description,
        $Email,
        $Status
    )

    $Api = "http://$($Hostname):$Port/" + "merchants"

    $Headers = @{
        'Content-Type'='text/csv'
    }

    $Body = @"
    $Name, $Description, $Email, $Status
"@

    Invoke-RestMethod -Method 'Post' -Uri $Api -Body $Body -Headers $Headers
}

function PostMerchant {
    param (
        $Hostname = "localhost",
        $Port = 8080,
    
        $Id,
        $Name,
        $Description,
        $Email,
        $Status
    )

    $Api = "http://$($Hostname):$Port/" + "merchants/" + $Id

    $Headers = @{
        'Content-Type'='application/json'
    }

    $Merchant = @{
        uuid = $Id
        name = $Name
        description = $Description
        email = $Email
        status = $Status
    }

    $Body = @{
        merchant = $Merchant
    }

    Invoke-RestMethod -Method 'Post' -Uri $Api -Body ($Body|ConvertTo-Json) -Headers $Headers
}

function DeleteMerchant {
    param (
        $Hostname = "localhost",
        $Port = 8080,

        $Id
    )

    $Api = "http://$($Hostname):$Port/" + "merchants"

    $Headers = @{
        'Content-Type'='application/json'
    }

    $Merchant = @{
        uuid = $Id
    }

    $Body = @{
        merchant = $Merchant
    }

    Invoke-RestMethod -Method 'Delete' -Uri $Api -Body ($Body|ConvertTo-Json) -Headers $Headers
}

function CreateAdmin {
    param (
        $Hostname = "localhost",
        $Port = 8080,

        $Name,
        $Description,
        $Email
    )

    $Api = "http://$($Hostname):$Port/" + "admins"

    $Headers = @{
        'Content-Type'='text/csv'
    }

    $Body = @"
    $Name, $Description, $Email
"@

    Invoke-RestMethod -Method 'Post' -Uri $Api -Body $Body -Headers $Headers
}
