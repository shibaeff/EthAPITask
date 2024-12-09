# Ethereum Validator API

This is a RESTful API application that provides information about Ethereum block rewards and validator sync committee duties.
## Some assumptions made

0. All slotno requests about "rewards" do not make sense for slots <  4700013.
Info about the sync committees and the state of validator indexes also seems to lag for slots before the Merge (however, beaconcha serves them well).
I think that's something to do with Altair and go out of the scope of this assignment.
1. All the parameters are for mainnet (it would be nicer to make config with custom rpc and network params);
it's imaginable to spin up a local Kurtosis env with some modified eth consensus
2. The service is designed in a manner that it's used locally by company, so no fancy async patterns and 
a global eth client object
3. Fee sums in gwei are *statistically correct*, I tested the values against block explorers;
no blazing precision to the 9th decimal or something.
4. Reward for the MEV block is extracted from the MEV reward transfer tx and added to the overall reward for block
5. int64 is used a base type for the endpoints (current block height and rewards kind of fit this data type); 
uint64 could be better but `big.NewInt` constructor in golang works with signed version and i prefer to use it for readability
6. the design is intentionally procedural since there's no complex business logic or databases involved (that's a simple statictical service, not an online shop)
7. There're some beacon API libraries from protolambda and ethpandaops; but they do not seem to implement 100% of spec;
that's why I decided to hardcode types and endpoints.
8. Some of the code is deprecated, but it's there to show the way I thought.
9. There's a method to convert val indices into public keys using native beacon node API, but it's too slow to make 
10. Ethereum reward structure is complex. It's impacted by the sync committee rewards, attestations rewards, fees rewards and potential mev rewards.
That's outlined here: https://eth2book.info/capella/part2/incentives/rewards/#rewards. 
I try to calculate the total sum of rewards for proposing validator.

## Definition of MEV blocks

To identify MEV block, we should look at the last tx in the block in 99% of cases. 
Then it tx recv address should be checked against some Flashbots data lists.
I derived some simpler probablistic way: 
1) The last transaction's recv adress is taken
2) Ethscan indexer is queried about the last n incoming txs with this account
3) If the last n txs incoming to this addr were always the last in the corresponding blocks, with 99% it's a MEV address

False positives are hard to imagine, while there should be some false negatives, but Simon confirmed that it was acceptable.

## Modes of rewards

The app has two modes: `beast` and `light` (see config). I would encourage to do mass testing in `light` and try `beast` just out interest.
In the light mode, only EL rewards are taken into account: `tx fees - burnt fees + mev (if block is mev)`.
In the beast mode, I also try to calc the attestation and sync committee rewards, which block proposer can also get. However, they appear 
to be statistically insignificant compared to EL rewards. That's why I'd rater just monitor attestations and sync duties
rather than hunt down precise formulas for these rewards (also, one should take into account that there's a lag in these rewards distribution).


## Prerequisites

- Go 1.19 or higher
- Internet connection (for Quicknode RPC access)

## Building and Running
### From binaries
1. Clone the repository
2. Install dependencies:
   ```bash
   go mod tidy
   ```
3. Take values from `config.yaml.example` to `config.yaml`; you can choose between two modes `light` and `beast`.
4. Run the server:
   ```bash
   make run
   ```
Or:
   ```bash
   docker compose up
   ```


## API Endpoints

The api Swagger and API itself are deployed here:
http://206.81.25.233:8000/swagger/index.html

http://206.81.25.233:8000
### Get Block Reward
```bash
curl http://localhost:8000/blockreward/{slot}
```

### Get Sync Duties
```bash
curl http://localhost:8000/syncduties/{slot}
```

## Testing

For some fuzzy-style tests run this script: 
```
bash test_endpoints.sh
```

If you're running a server, be sure to check the local swagger: http://localhost:8000/swagger/index.html
Some particular test cases:
```bash
echo normal slot
curl -X 'GET' \
  'http://localhost:8000/blockreward/10579354' \
  -H 'accept: application/json'

echo slot in the future
curl -X 'GET' \
  'http://localhost:8000/blockreward/105793540' \
  -H 'accept: application/json' \
  -w "\nHTTP Status: %{http_code}\n"
# {"error":"Slot is in the future"}, 400

echo request with incorrect parameters
curl -X 'GET' \
  'http://localhost:8000/blockreward/tokyohotel' \
  -H 'accept: application/json' \
  -w "\nHTTP Status: %{http_code}\n"
# {"error":"Invalid slot number"} HTTP Status: 400


echo missed slot 
curl -X 'GET' \
  'http://localhost:8000/blockreward/10579330' \
  -H 'accept: application/json' \
  -w "\nHTTP Status: %{http_code}\n"
# {"error":"block not found for slot"} HTTP Status: 404

echo normal slot
curl -X 'GET' \
  'http://localhost:8000/syncduties/10579354' \
  -H 'accept: application/json'

echo slot in the future
curl -X 'GET' \
  'http://localhost:8000/syncduties/105793540' \
  -H 'accept: application/json' \
  -w "\nHTTP Status: %{http_code}\n"
# {"error":"Slot is in the future"}, 400

echo request with incorrect parameters
curl -X 'GET' \
  'http://localhost:8000/syncduties/tokyohotel' \
  -H 'accept: application/json' \
  -w "\nHTTP Status: %{http_code}\n"
# {"error":"Invalid slot number"} HTTP Status: 400


echo missed slot 
curl -X 'GET' \
  'http://localhost:8000/syncduties/10579330' \
  -H 'accept: application/json' \
  -w "\nHTTP Status: %{http_code}\n"
# {"error":"block not found for slot"} HTTP Status: 404
```
## Design Choices

- **Gin Framework**: Used for its performance, middleware support, and ease of use
- **go-ethereum**: Official Go Ethereum client implementation for blockchain interaction
- **Error Handling**: Comprehensive error handling with appropriate HTTP status codes
- **Modular Structure**: Code organized into handlers, services, and models packages
- **Several clients** there're adapaters for EL and CL APIs; for EL I'm using the client from geth's codebase
- **Common models** I define them in the sep module and propagate throught the application for convenience
- **Libs** I use viper for file configs, cobra for flags, logrus for logging; a small middleware struct is used to pass 
configs from the entrypoint to handlers and further to modules

## Further improvements

- add more contexts to handle connections and graceful shutdown
- more precise formulas for sync duties / attestation rewards
- robust infra, so the two above make some sense
- API security checks
