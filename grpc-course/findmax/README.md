# FindMaximum API

## Bi-Directional Streaming API

The Server takes a stream of Request message that has one integer, and returns a stream of Responses that represent the current maximum between all these integers

Ex: 

The client will send a stream of number (1,5,3,6,2,20) and the server will respond with a stream of (1,5,6,20)