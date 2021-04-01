# scratch-test

A simple api endpoint with custom in-memory database implementation.

### Assumptions:
* A custom database is to be implemented and other existing database implementation
  can't be used.
  
* The data stored in the database is in the format as given in the `challenge` document

## Usage
1. Just running `go run main.go` from the `scratch-test` directory should get the
   server up and running.
    
2. A docker image can be built and run from the Dockerfile.
```bash
# Build the app image
docker build -t scratch-app .

# Run the app image in interactive mode to access logs
docker run -it -p 8080:8080 scratch-app
```

3. Only `/` endpoint is available, but different filters can be used.
```bash
curl -X GET http://localhost:8080/?state=florida\&state=ak\&limit=4
```
Returns four max entries that have state florida or alaska.
