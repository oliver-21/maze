
$Env:GOOS = 'js'
$Env:GOARCH = 'wasm'
go build -o web/yourgame.wasm .
Remove-Item Env:GOOS
Remove-Item Env:GOARCH