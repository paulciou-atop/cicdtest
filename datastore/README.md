## Data store service

Provides API to allow storing general data such as configuration data, metadata about devices, specific information about users, etc.

Uses SQL databases like PostgreSQL or MySQL for relational data that require OLTP support.  
For scalability we also need to consider distributed SQL like TiDB, Cockroach DB, YugabyteDB, etc.


Can also use NoSQL databases for storing less schema rigid data like document data in JSON format.  
Some NoSQL databases that can be used may be MongoDB, Redis, Cassandra, CouchDB, etc.

Either way, there needs to be a API layer that makes this service independent of the specific backend database being used.  Do not be tied to only one database.

