# Arquitetura e decisões

Aqui vou documentar e justificar algumas das decisões tomadas ao longo do desenvolvimento do projeto.

## Domain model

Os principais modelos de domínio são Wallet, WagerTransaction e WalletLedgerEntry.

Money é um value object utilizado para representar valores monetários com precisão e sua respectiva moeda.

OutboxEvent e InboxMessage são modelos de suporte à confiabilidade do processamento assíncrono. O outbox permite persistir eventos para publicação após o commit da transação no banco de dados, enquanto o inbox ajuda a evitar o processamento duplicado de mensagens.

![Diagrama de domínio](/assets/modelagemDeDominio.png)

## Fluxo inicial de requisição

Eu tentei mapear o fluxo da request ate o banco de dados, a fim de ter um overview geral da API antes de add camadas como uberFX e SQS que podem dificultar a vizualização do fluxo principal.

![Diagrama de fluxo HTTP](/assets/fluxoInicial.png)
