# Distributed Systems Examples

## Run the Demo

- `./start.sh`

## Test Scenarios

### Success

- ./start.sh
- Submit an Order

```bash
curl --request POST \
  --url http://localhost:3000/order \
  --data '{
  "customer": {
    "firstName": "Matt",
    "lastName": "Chaffe",
    "city": "UK"
  },
  "products": [{
    "productCode": 1,
    "quantity": 1,
    "total": 100
  }],
  "total": 100
}'
```

### Failure

- Test a failure case `docker-compose stop transmission-endpoint`
- Submit the order

```bash
curl --request POST \
  --url http://localhost:3000/order \
  --data '{
  "customer": {
    "firstName": "Matt",
    "lastName": "Chaffe",
    "city": "UK"
  },
  "products": [{
    "productCode": 1,
    "quantity": 1,
    "total": 100
  }],
  "total": 100
}'
```

## View Workflow History

- visit [http://localhost:8233/namespaces/distributed-systems-examples/workflows](http://localhost:8233/namespaces/distributed-systems-examples/workflows)
