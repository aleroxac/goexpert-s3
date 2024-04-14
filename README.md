# goexpert-s3


## Modo de uso
``` shell
### 1. Crie os arquivos
timeout 0.5s go run cmd/generator/main.go

### 2. Faça o upload dos arquivos para o s3
go run cmd/uploader/main.go

### 3. Execute o uploader
bash -c "time go run cmd/uploader/main.go"
```
