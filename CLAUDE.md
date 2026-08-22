# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Cotadora BRL→MXN — teste técnico de fim de semana

### O que é
API em Go + front em React. Cliente informa MXN; sistema cota em BRL via USDT.

Preço: `mxn × ask(USDT/BRL, menor entre Binance e OKX) × 1/bid(USDT/MXN, Binance) × 1,006`,
arredondado PARA CIMA (teto) com 2 casas, uma única vez, no final.

Cotação expira em 10s; só as confirmadas vão para o SQLite. Histórico lista só confirmadas.

### Regras de trabalho
- Go: apenas stdlib + driver SQLite. NÃO adicione dependências sem perguntar antes.
- React: sem axios, sem UI kit, sem lib de estado; fetch nativo e componentes simples.
- Sempre propor um plano antes de editar; esperar ok antes de aplicar.
- Mudanças pequenas e incrementais; ao final de cada uma, rodar os testes.
- Explicar em 2–3 linhas o porquê de cada escolha não óbvia 
- Commits pequenos com mensagem imperativa curta; só commitar quando pedido.
