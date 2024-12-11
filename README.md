# Tarefas para Implementar o Fechamento Automático de Leilões

## 1. Configuração do Ambiente
- [ ] Clone o repositório fornecido.
- [ ] Configure o ambiente de desenvolvimento utilizando Docker/Docker Compose.
- [ ] Certifique-se de que o projeto está rodando corretamente.

---

## 2. Definir Variáveis de Ambiente
- [ ] Adicione variáveis de ambiente para configurar o tempo de duração do leilão (por exemplo, `AUCTION_DURATION_SECONDS`).
- [ ] Atualize os arquivos necessários para carregar essas variáveis, como arquivos `.env` e configurações no Docker Compose.

---

## 3. Implementação do Fechamento Automático

### Função para calcular o tempo do leilão
- [ ] Implemente uma função que calcula a duração de cada leilão com base nas variáveis de ambiente.
- [ ] Certifique-se de retornar o tempo no formato adequado (ex.: `time.Duration`).

### Goroutine para fechamento de leilões
- [ ] Crie uma nova goroutine que:
	- Executa periodicamente (ex.: usando `time.Ticker`).
	- Verifica no banco de dados quais leilões expiraram.
	- Atualiza o estado do leilão para "fechado".
	- Trate concorrência usando locks ou transações no banco, se necessário.

### Atualização no banco de dados
- [ ] Adicione a lógica no arquivo `internal/infra/database/auction/create_auction.go` para suportar o fechamento de leilões expirados.

---

## 4. Testes Automatizados
- [ ] Crie testes para verificar o comportamento de fechamento automático:
	- Simule a criação de um leilão.
	- Aguarde o tempo definido e valide se o leilão foi fechado automaticamente.
- [ ] Utilize ferramentas como `time.Sleep` ou mocks para testar o comportamento dependente de tempo.

---

## 5. Documentação
- [ ] Escreva instruções claras no `README.md`:
	- Como configurar as variáveis de ambiente.
	- Como rodar o projeto com Docker/Docker Compose.
	- Como rodar os testes.

---

## 6. Validação Final
- [ ] Teste manualmente a aplicação para garantir que:
	- Leilões criados expiram automaticamente após o tempo configurado.
	- A aplicação funciona corretamente com múltiplos leilões simultâneos.
- [ ] Revise o código para assegurar que está thread-safe e escalável.
