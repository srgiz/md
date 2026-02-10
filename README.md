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
flowchart LR
    Request -- controller create --> Cmd -- send to --> Bus(Command Bus) -- exec --> Handler -- return --> Reply -- controller create --> Response
```

```mermaid
flowchart TD
    subgraph UseCaseA1
        CmdA1 --> HandlerA1 --> ReplyA1
    end

    subgraph HandlerA1
        EntityA1
        RepositoryInterfaceA1
    end

    classDef bgA stroke:#85E89D
    DomainContextA:::bgA
    subgraph DomainContextA
        UseCaseA1
        InfrastructureA
    end

    subgraph InfrastructureA
        RepositoryA1(SqlRepositoryA) -.-> RepositoryInterfaceA1
    end

    subgraph UseCaseB1
        CmdB1 --> HandlerB1 --> ReplyB1
    end

    subgraph HandlerB1
        EntityB1
        RepositoryInterfaceB1
    end

    classDef bgB stroke:#79B8FF
    DomainContextB:::bgB
    subgraph DomainContextB
        UseCaseB1
        InfrastructureB
        UseCaseB2
    end

    subgraph UseCaseB2
        CmdB2 --> HandlerB2 --> ReplyB2
    end

    subgraph HandlerB2
        EntityB2
        RepositoryInterfaceB2
    end

    subgraph InfrastructureB
        RepositoryB1(HttpClientB) -.-> RepositoryInterfaceB1
        RepositoryB1 -.-> RepositoryInterfaceB2
    end

    RequestA1 -- "ControllerA1 -> Bus" --> CmdA1
    ReplyA1 -- "Bus -> ControllerA1" --> ResponseA1

    RequestB1 -- "ControllerB1 -> Bus"  --> CmdB1
    ReplyB1 -- "Bus -> ControllerB1" --> ResponseB1

    RequestB2 -- "ControllerB2 -> Bus"  --> CmdB2
    ReplyB2 -- "Bus -> ControllerB2" --> ResponseB2
```
