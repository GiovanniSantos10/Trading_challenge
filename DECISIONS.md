# Decisões técnicas

Este projeto foi construído numa sessão de pareamento com o Claude Code (Anthropic) — eu descrevia o requisito, discutíamos as alternativas e eu aprovava cada mudança antes de aplicar. As decisões abaixo são as mais relevantes dessa sessão: o que ficou, o que foi descartado e por quê.

## 1. Driver SQLite: `modernc.org/sqlite` em vez de `mattn/go-sqlite3`

O enunciado permite só stdlib + um driver de SQLite, então essa escolha era inevitável — a dúvida era qual driver. `mattn/go-sqlite3` é o mais usado e maduro, mas é um wrapper cgo em cima da libsqlite3 em C: exige `CGO_ENABLED=1` e um compilador C instalado (gcc via MSYS2/TDM-GCC no Windows). Como estava desenvolvendo no Windows sem esse toolchain configurado, optei pelo `modernc.org/sqlite`, que é SQLite reimplementado em Go puro — compila sem depender de nada externo, ao custo de ser um pouco mais lento que a versão em C. Irrelevante no volume de dados de uma cotadora.

## 2. Chamadas às exchanges em paralelo, sem lib de concorrência

O enunciado pede explicitamente que as requisições à Binance e à OKX sejam feitas em paralelo — na primeira versão eu tinha implementado sequencialmente e só percebi o erro relendo o enunciado com calma. A correção óbvia seria usar `golang.org/x/sync/errgroup`, que deixa o código mais enxuto, mas é uma dependência fora da stdlib, e a regra do projeto é não adicionar nada sem necessidade real. Optei por `sync.WaitGroup` puro: três goroutines, cada uma escrevendo só na sua própria variável de resultado (sem dado compartilhado entre elas, então sem necessidade de mutex também). Mais verboso que `errgroup`, mas zero dependência nova pra resolver um problema que a stdlib já resolve sozinha.

## 3. Contador de expiração: estado derivado, não um `secondsLeft` que soma sozinho

A primeira versão do componente de cotação guardava `secondsLeft` como estado e decrementava a cada segundo. Migrei pra guardar só um timestamp (`now`, atualizado por um `setInterval`) e calcular `secondsLeft`/`expirado` a cada render a partir dele e do `expires_at` que a API devolve — o servidor é quem decide o prazo, o front só espelha. Essa mudança quase introduziu um bug: numa tentativa de silenciar um aviso de lint (que reclamava de um `setState` síncrono dentro de um efeito), removi a sincronização do relógio no momento em que uma cotação nova chegava. Resultado: se o usuário demorasse pra pedir a cotação, o contador nascia inflado (cheguei a ver 28s em vez de 10s) e "saltava" pro valor certo um segundo depois. A correção foi reintroduzir aquele `setState` síncrono — nesse caso específico ele não é redundância, é a sincronização em si. Ficou registrado em comentário no código pra não repetir o erro.

## 4. Front-end: rotas com um hook de 20 linhas, não `react-router-dom`

O enunciado pede duas páginas (Cotação e Histórico). A escolha natural seria `react-router-dom`, mas o projeto define explicitamente "sem lib de estado, componentes simples" pro front, e duas rotas estáticas não justificam uma dependência inteira de roteamento. Escrevi um hook próprio (`useRoute`) que lê `window.location.pathname`, escuta o evento `popstate` do navegador e navega com `history.pushState` — é o mesmo padrão que a documentação do React usa como exemplo de quando um `useEffect` faz sentido (sincronizar com um sistema externo, aqui a URL). Se o projeto crescesse pra muitas rotas, aninhamento ou parâmetros complexos, `react-router-dom` passaria a valer a pena — pra duas páginas fixas, não.

## 5. CORS liberado (`Access-Control-Allow-Origin: *`), não restrito por variável de ambiente

A alternativa mais "correta" na aparência seria liberar só a origem do front via variável de ambiente (`CORS_ORIGIN=http://localhost:5173`). Rejeitei porque a API não usa cookies nem sessão — não há credencial nenhuma que uma origem maliciosa pudesse roubar liberando geral — e porque isso criaria acoplamento desnecessário entre a porta que o Vite escolhe e uma variável que alguém precisaria lembrar de configurar. Pra uma API sem autenticação, liberar geral é mais simples e não abre risco real.

---

## Sobre o processo

Boa parte do valor dessa sessão não foi só "pedir e receber código pronto": em vários pontos o Claude sinalizou problemas que eu não tinha visto — a fórmula usando o par errado de moedas numa versão inicial do `GetPrices`, o rótulo de "você paga/recebe" invertido em relação ao fluxo real de dinheiro, e o próprio bug do contador acima. Preferi manter esse registro porque acho que decisão técnica sem o "por que rejeitei a alternativa" vira só uma lista de escolhas arbitrárias — e o objetivo aqui era documentar raciocínio, não só resultado.
