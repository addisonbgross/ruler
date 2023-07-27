# Ruler

A very serious, real, and durable database

## Planned Rules

1. There is only 1 Ruler
2. Every node can only have 1 child node. This forms a linear hierarchy of nodes
3. Only the Ruler can handle reads and writes, a non-Ruler node has to relay the request to the Ruler

## Planned Features

- Key space division between nodes, with adjustments when nodes leave or are added
- Multi-node replication based on key spaces
- [Dead-letter-queue](https://en.wikipedia.org/wiki/Dead_letter_queue) for failed replications, with retires

## Endpoints

### External

- **GET /read/{key}**: Read a value from a key (will replicate to all other nodes)
- **POST /write**: Write a value from a write payload
- **GET /dump**: Read all values stored in a node

## Payload Examples

**Write**:
```
{
    key: 'some-key',
    value: 'some-value'
    isreplicate: false
}
```