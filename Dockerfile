# Estágio de build
FROM golang:1.22.5-bullseye as build
# Define o diretório de trabalho dentro do container
WORKDIR /app
# Copia todos os arquivos onde o Dockerfile está localizado para dentro do container
COPY . .
# Compila a aplicação
RUN go build -o upload /app/main.go


# Estágio de execução
FROM golang:1.22.5-bullseye
# Define o diretório de trabalho dentro do container
WORKDIR /app
# Copia todos os arquivos onde o Dockerfile está localizado para dentro do container
COPY . ./
# Copia o arquivo compilado do estágio build para a pasta bin
COPY --from=build /app/upload ./bin
# Expoe a aplicação na porta 4000
EXPOSE 4000
# Executa o binário da aplicação
CMD ["./bin/upload"]