# Ruler

A very serious, real, and durable database

## Rules

1. There is only 1 Ruler
2. Every node can only have 1 child node. This forms a linear hierarchy of nodes
3. Only the Ruler can handle reads and writes, a non-Ruler node has to relay the request to the Ruler

## Planned Feautres
### Node States

- **Normal**: The node will respond to read and write requests 
- **Succession**: The hierarchy of nodes will readjust their ranking in the event that any of the nodes becomes unavailable
- **Regicide**: A non-Ruler node will attempt to usurp the Ruler, in order to become the next Ruler. This will require consensus from the other, non-Ruler nodes
- **Resolving**: A node is working towards knowing the current hierarchy. All messages will be cached until the hierarchy is known 

### Messages

- **Read/Write**: The client wanting read or write data to the Ruler hierarchy
- **Health**: A 'heartbeat' message that nodes will send to a subset of all nodes. This lets nodes know which other nodes are available
- **Regicide**: A node will send this message if it wants to become the new Ruler. The recieving node can either 

## Endpoints

### External

- **GET /read/{key}**: Read a value from a key
- **POST /write**: Write a value from a write payload
- **Get /dump**: Read all values stored in a node

## Payload Examples

**Write**:
```
{
    key: 'some-key',
    value: 'some-value'
}
```