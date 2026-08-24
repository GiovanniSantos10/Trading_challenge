# Cotador BRL → MXN

API em Go + front em React. O cliente informa quanto MXN quer enviar, o sistema cota em BRL (via USDT, usando os melhores preços entre Binance e OKX) e o cliente confirma dentro de um prazo de 10 segundos.

Preço: `mxn × ask(USDT/BRL, menor entre Binance e OKX) × 1 / bid(USDT/MXN, Binance) × 1,006`, arredondado para cima, 2 casas, uma única vez.

## Stack

- **Backend**: Go (stdlib + [`modernc.org/sqlite`](https://pkg.go.dev/modernc.org/sqlite), o único driver de banco usado).
- **Front**: React (Vite), sem axios/UI kit/lib de estado — `fetch` nativo e `useState`/`useEffect`.

## Como rodar

### Backend

```bash
go run .
```

Variáveis de ambiente (opcionais):

| Variável  | Default      | Uso                          |
|-----------|--------------|-------------------------------|
| `PORT`    | `8080`       | Porta HTTP da API             |
| `DB_PATH` | `quotes.db`  | Caminho do arquivo SQLite     |

O arquivo do banco é criado sozinho na primeira execução (schema incluso).

### Front-end

```bash
cd frontend
npm install
npm run dev
```

Abre em `http://localhost:5173`. A URL da API vem de `frontend/.env` (`VITE_API_URL`, já configurado pra `http://localhost:8080`).

### Testes

```bash
go test ./...              # backend
cd frontend && npm run lint && npm run build   # front (lint + build de produção)
```

## Endpoints

| Método | Rota                     | Descrição                                  |
|--------|--------------------------|---------------------------------------------|
| GET    | `/health`                | Healthcheck                                 |
| POST   | `/quotes`                | Cria cotação — body `{"mxn": number}`       |
| POST   | `/quotes/{id}/confirm`   | Confirma cotação pendente (dentro de 10s)   |
| GET    | `/history`               | Lista cotações confirmadas                  |

## Estrutura

```
main.go              — bootstrap: HTTP server, LivePriceProvider (busca preços em paralelo)
exchanges/            — clientes HTTP da Binance e da OKX
handlers/             — regra de negócio: cotação, expiração, persistência SQLite, rotas HTTP
frontend/src/         — App (estado central) + páginas (Cotação, Histórico) + componentes
```

## Limitações conhecidas / o que faria com mais tempo

- **Só o par BRL/MXN**: a fórmula e os pares buscados nas exchanges são fixos no código. Generalizar pra outras moedas exigiria trocar `CalculateBRL` por algo parametrizado por par, e uma tabela de pares suportados em vez de constantes soltas.
- **Sem teste automatizado provando o paralelismo** das chamadas às exchanges — só validei manualmente (tempo de resposta de uma cotação real). Escrever esse teste exigiria expor as variáveis de base URL do pacote `exchanges` (hoje não-exportadas) ou injetar dependências, e achei invasivo demais pro ganho nesse prazo.
- **Sem autenticação/rate limiting** na API — qualquer um pode criar e confirmar cotações. Faria sentido pelo menos limitar taxa de criação de cotação por IP antes de considerar isso pronto pra produção.
- **Sem retry/backoff** nas chamadas às exchanges — uma falha de rede pontual em qualquer uma das três chamadas derruba a cotação inteira (`502`). Um retry curto (1-2 tentativas) reduziria falso-negativo por instabilidade passageira.
- **SQLite sem WAL mode configurado** — não é um problema no volume de uma cotadora de teste, mas sob escrita concorrente real valeria habilitar `journal_mode=WAL` explicitamente.
- **Sem CI** — não há GitHub Actions rodando `go test`/`npm run build` a cada push. Seria o próximo passo natural pra manter a régua de qualidade sem depender de rodar tudo manualmente.
- **Front sem testes automatizados de componente** — só lint (`oxlint`) e build. A lógica mais delicada (o contador de expiração, especialmente) se beneficiaria de testes com Vitest + Testing Library.
- **Cotações pendentes não expiram por retenção** — ficam em memória até serem confirmadas ou até alguém tentar confirmá-las depois do prazo; não há uma goroutine de limpeza. Sem impacto funcional (o `Confirm` já rejeita cotação expirada), mas cotações pendentes nunca confirmadas ficam ocupando memória indefinidamente num processo de vida longa.

Decisões técnicas mais detalhadas (com as alternativas rejeitadas e o porquê) estão em [`DECISIONS.md`](./DECISIONS.md).
