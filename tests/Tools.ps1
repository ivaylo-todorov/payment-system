
function RandomString {
    param (
        [int]$Length = 6
    )

    $set = "abcdefghijklmnopqrstuvwxyz0123456789".ToCharArray()

    $result = ""
    for ($x = 0; $x -lt $Length; $x++) {
        $result += $set | Get-Random
    }

    return $result
}
