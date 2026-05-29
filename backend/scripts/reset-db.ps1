# reset-db.ps1
# Resetea la BD (down -> up -> seed) usando un usuario admin.
# Uso:  .\reset-db.ps1
# Opcional: .\reset-db.ps1 -BaseUrl "http://localhost:8080" -Email "admin@correo.cl"

param(
    [string]$BaseUrl = "http://localhost:8080",
    [string]$Email,
    [string]$Password
)

$ErrorActionPreference = "Stop"

# 1) Pedir credenciales si no vinieron por parametro
if (-not $Email)    { $Email = Read-Host "Email admin" }
if (-not $Password) {
    $secure = Read-Host "Password" -AsSecureString
    $bstr   = [Runtime.InteropServices.Marshal]::SecureStringToBSTR($secure)
    $Password = [Runtime.InteropServices.Marshal]::PtrToStringAuto($bstr)
}

# 2) Login -> token
$body = @{ email = $Email; password = $Password } | ConvertTo-Json
try {
    $login = Invoke-RestMethod -Uri "$BaseUrl/api/v1/auth/login" -Method Post `
        -ContentType "application/json" -Body $body
} catch {
    Write-Host "Login fallo: $($_.Exception.Message)" -ForegroundColor Red
    exit 1
}

# 3) Caso TOTP: no se puede automatizar desde aca
if ($login.data.totp_required) {
    Write-Host "Este usuario tiene 2FA (TOTP) activado. Usa un admin sin 2FA o desactivalo." -ForegroundColor Yellow
    exit 1
}

$token = $login.data.token
if (-not $token) {
    Write-Host "No vino token en la respuesta. Revisa credenciales/rol." -ForegroundColor Red
    exit 1
}
Write-Host "Token OK: $($token.Substring(0,20))..." -ForegroundColor Green

# 4) Reset BD
try {
    $reset = Invoke-RestMethod -Uri "$BaseUrl/dev/reset-db" -Method Post `
        -Headers @{ Authorization = "Bearer $token" }
    Write-Host "BD reseteada:" -ForegroundColor Green
    $reset | ConvertTo-Json -Depth 5
} catch {
    $code = $_.Exception.Response.StatusCode.value__
    switch ($code) {
        403 { Write-Host "403: el usuario NO es admin." -ForegroundColor Red }
        401 { Write-Host "401: token invalido/vencido." -ForegroundColor Red }
        default { Write-Host "Error ${code}: $($_.Exception.Message)" -ForegroundColor Red }
    }
    if ($_.ErrorDetails.Message) {
        Write-Host "--- Detalle del server ---" -ForegroundColor Yellow
        Write-Host $_.ErrorDetails.Message
    }
    exit 1
}
