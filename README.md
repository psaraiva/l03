# Sistema de Leilões em Go

Este projeto é um microsserviço para um sistema de leilões, desenvolvido em Go e utilizando MongoDB como banco de dados e utiliza Docker.

## Pré-requisitos

Antes de começar, garanta que você tenha as seguintes ferramentas instaladas:
*   **Git:** Para clonar o repositório.
*   **Docker e Docker Compose:** Para containerizar e orquestrar a aplicação e o banco de dados.
*   **Visual Studio Code:** Editor de código recomendado.
*   **Extensão [REST Client](https://marketplace.visualstudio.com/items?itemName=humao.rest-client)** para o VS Code: Para interagir facilmente com a API.

---

## 🚀 Executando a Aplicação

Siga os passos abaixo para ter o ambiente rodando localmente.

### Passo 1: Clonar o Repositório

Abra seu terminal e clone o projeto para sua máquina local.

```bash
git clone https://github.com/psaraiva/l03
cd l03
```

### Passo 2: Configurar o Ambiente

A aplicação precisa de algumas variáveis de ambiente para se conectar ao banco de dados. Aproveite o arquivo de exemplo.

No terminal, dentro da raiz do projeto (`l03`), execute o comando abaixo para criar o seu arquivo `.env` local:

```bash
cp cmd/auction/.env-example cmd/auction/.env
```
O arquivo `docker-compose.yml` está configurado para injetar automaticamente essas variáveis no contêiner da aplicação.

### 🔧 Variáveis de Ambiente

O arquivo `.env` controla as configurações da aplicação. Abaixo está a descrição de cada variável:

| Variável                    | Descrição                                                                                                  | Exemplo/Padrão |
| --------------------------- | ---------------------------------------------------------------------------------------------------------- | -------------- |
| `MONGODB_URL`               | A URL de conexão para a instância do MongoDB.                                                              | `mongodb://root:password@mongo:27017` |
| `MONGODB_DATABASE`          | O nome do banco de dados que a aplicação utilizará.                                                        | `l03`          |
| `AUCTION_DURATION_DEFAULT`  | A duração padrão de um leilão após ser iniciado. O formato deve ser compatível com `time.ParseDuration`.    | `10m` (10 minutos) |
| `BID_BATCH_BUFFER_SIZE`     | O número máximo de lances a serem agrupados em memória (buffer) antes de serem persistidos no banco de dados. | `5`            |
| `BID_BATCH_BUFFER_DURATION` | O tempo máximo de espera para persistir o buffer de lances, mesmo que o tamanho máximo não tenha sido atingido. | `1m` (1 minutos) |

### Passo 3: Executar com Docker Compose

Com tudo configurado, o próximo passo é iniciar os contêineres. O Docker Compose irá ler o arquivo `docker-compose.yml`, construir a imagem da sua aplicação Go, baixar a imagem do MongoDB e iniciar ambos os serviços.

Execute o seguinte comando na raiz do projeto:

```bash
docker-compose up -d --build
```

Você verá os logs de ambos os serviços no output padrão. Quando a mensagem `Starting http server on port :8080` aparecer, significa que a aplicação está pronta para receber requisições em `http://localhost:8080`.

```bash
docker logs -f l03_app
```

---

## ⚡ Interagindo com a API

A maneira mais fácil de testar a API é usando o arquivo `cli_rest/requests.http` que está na raiz do projeto.

1.  Abra o projeto no VS Code.
2.  Instale a extensão **REST Client**, caso ainda não a tenha.
3.  Abra o arquivo `cli_rest/requests.http`.
4.  Você verá as requisições formatadas. Acima de cada uma, haverá um link `Send Request`.
5.  Clique em `Send Request` para executar cada chamada na ordem:
    *   Primeiro, crie um usuário.
    *   Depois, crie um leilão.
    *   Consulte o leilão pelo ID.
    *   Por fim, faça um lance no leilão.
    *   Acompanhe o fechamento do leilão via terminal ou consultando seu status.
    *   Consulta o vencedor do leilão.

Uma nova aba será aberta no VS Code a cada requisição, mostrando a resposta da API.

---

## 🏛️ Arquitetura e Estratégia de Leilão

A aplicação é composta por dois componentes principais: um servidor web para a API e um processo *worker* em segundo plano, responsável por gerenciar o ciclo de vida dos leilões.

### Ciclo de Vida do Leilão

O leilão progride através de uma máquina de estados simples:
* **Ativo** (aguardando início);
* **Em Execução** (aceitando lances);
* **Concluído** (finalizado).

Obs: Lances são permitidos apenas durante o estado **Em Execução**.

### Funcionamento do Worker

O *worker* automatiza o processo:

-   Periodicamente, ele busca por leilões com status **Ativo**.
-   Ao encontrar um, altera seu status para **Em Execução** e inicia uma goroutine com um temporizador.
-   Ao final da duração do leilão, uma goroutine o finaliza, alterando seu status para **Concluído**.

### Processamento de Lances (original)

Para otimizar a performance, os lances não são escritos no banco de dados individualmente. Eles são recebidos e agrupados em um buffer. O conteúdo do buffer é persistido em lote quando atinge sua capacidade máxima ou um tempo limite é excedido.
