# Distributed Systems Examples

## To run demo

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

- Check the Transmission
```bash
curl --request GET \
  --url http://localhost:3001/
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
- Check the logs and the DB after 5 mins to check for timeout poller