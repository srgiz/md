# Доменные контексты и свой фреймворк с CommandBus.

## Запуск

* `docker compose up -d` (create `compose.override.yaml`)
* `docker exec -i -t md-app sh`
* `make migrate`

## Тесты

`make test` or `make coverage-html`

## Остальное команды

`make help`

## Мультики

```mermaid
flowchart TD
    subgraph UseCase1
        Cmd1 --> Handler1 --> Reply1
    end

    subgraph Handler1
        Entity1
        RepositoryInterface1
    end
    
    subgraph DomainContext1
        UseCase1
        Infrastructure1
    end

    subgraph Infrastructure1
        Repository1(SqlRepository) -.-> RepositoryInterface1
    end

    subgraph UseCase2
        Cmd2 --> Handler2 --> Reply2
    end

    subgraph Handler2
        Entity2
        RepositoryInterface2
    end

    subgraph DomainContext2
        UseCase2
        Infrastructure2
    end

    subgraph Infrastructure2
        Repository2(HttpClient) -.-> RepositoryInterface2
    end

    subgraph CommandBus
        DomainContext1
        DomainContext2
    end

    Request1 -- Deserialize request  --> Cmd1
    Reply1 -- Serialize response --> Response1

    Request2 -- Deserialize request  --> Cmd2
    Reply2 -- Serialize response --> Response2
```
