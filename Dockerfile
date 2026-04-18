# =============================================================================
# STAGE 1 Builder
# Compila o binário
# =============================================================================
FROM golang:1.26-alpine AS builder 

# Define o diretório de trabalho
WORKDIR /app

# Copia o arquivo go.mod e baixa as dependências
COPY go.mod ./
RUN go mod download

# Copia o codigo fonte
COPY . .

# Compila o binário
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-s -w" \
    -o app-demo \
    .


# =============================================================================
# STAGE 2 Runtime
# Imagem final só o binário compilado vai pra produção
# scratch = imagem vazia, zero dependências
# =============================================================================
FROM scratch

# Copia o binário compilado
COPY --from=builder /app/app-demo /app-demo

# Copia os arquivos estáticos
COPY --from=builder /app/static /static

# Executa como usuário não-root (segurança)
USER 65532:65532

# Expõe a porta 8080
EXPOSE 8080

# Comando de inicialização
CMD ["/app-demo"]