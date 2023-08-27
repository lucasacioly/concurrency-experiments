$goPath = "C:\Program Files\Go\bin\go.exe"  # Caminho para o executável 'go.exe'
$serverPath = Join-Path $PSScriptRoot "server\server.go"  # Caminho completo para o arquivo 'server.go'
$clientPath = Join-Path $PSScriptRoot "client\client.go"  # Caminho completo para o arquivo 'client.go'

$numberOfClients = Read-Host "Insira o numero de clientes: "   # Quantidade de clientes a serem iniciados

# Verifica se o executável 'go.exe' existe no caminho especificado
if (Test-Path $goPath) {
    # Inicia o servidor usando o comando 'go run'
    Start-Process -FilePath $goPath -ArgumentList "run", $serverPath -NoNewWindow

    # Inicia a quantidade 'x' de clientes usando o comando 'go run'
    for ($i = 1; $i -le $numberOfClients; $i++) {
        $clientID = $i
        Start-Process -FilePath $goPath -ArgumentList "run", $clientPath, "-clients=$numberOfClients", "-id=$clientID" -NoNewWindow
    }
} else {
    Write-Host "O caminho para o executável 'go.exe' está incorreto."
}
