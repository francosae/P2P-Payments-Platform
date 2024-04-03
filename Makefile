# Services
SERVICES = api-gateway auth-service pool-transactions

# Default target
server: 

# Individual targets for each service
api-gateway:
	$(MAKE) -C api-gateway server

auth-service:
	$(MAKE) -C auth-service server

pool-transactions:
	$(MAKE) -C pool-transactions server

.PHONY: $(SERVICES) server
