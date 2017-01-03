# cache
Golang Simple Cache.  Supports S3 and File System.

# Install

`go get "github.com/rameshvk/cache"`

# Usage

The basic interface for all caches and high order components is `ICache`:

```go
type ICache interface {
  Add(key string, value interface{}, expiration *time.Time) error
  Get(key string) (interface{}, *time.Time, error)
}
```

## S3 Object Cache

Creating an S3 object cache is as follows:

```go
  import "github.com/rameshvk/cache"
  s3Cache := NewS3ObjectCache("s3://your_bucket_here/prefix", "us-west-2")
```

Some caveats:

1. There is no explicit credentials setup -- you can set them up via global AWS config or environment for instance
2. The S3 object cache expects values to be of type `[]byte`.  If you would like to marshal/unmarshal other types, use `NewMarshaler`
3. Adding a nil value or fetching an expired object causes it to be removed from S3.
4. No explicit support to purge older S3 objects periodically (yet).
5. No explicit support for extra meta data or other custom fields (yet).

## Marshaler

It is useful to provide a mechanism to marshal/unmarshal other types than `[]byte` that the S3 Object Cache supports.  This can be done by using the `NewMarshaler` function which converts any `ICache` interface (with the help of the provided IEncoderDecoder) into one that supports arbitrary interface types.

The `NewJSONMarshaler` provides a simpler version of this which simply uses `encoding/json` to encode the values.  Note that it does not fully decode the return values since it has no knowledge of expected types.

```go
   s3CacheWithAnyTypeValue := NewJSONMarshaler(
     NewS3ObjectCache("s3://your_bucket_here/prefix", "us-west-2")
   )
```
 
## LRU
 
The LRU wrapper provides an in-memory LRU read-through cache.
 
```go

   ttl := time.Minute * 10; // refetch from S3 after 10 min
   maxCount := 1000; // do not cache more than 1000 items
   fetcher := nil; // do not do any side-loading
   fastCache := NewLRUCache(
     NewS3ObjectCache("s3://your_bucket_here/prefix", "us-west-2"),
     ttl,
     maxCount,
     fetcher
   )
```

The LRU cache also supports side-loading -- i.e. if the data is cached in S3 but fetched from another source.  This is done by providing a fetcher function to do the side-loading.

## File System cache

The File System cache provides a local file system based cache.  Note that the actual key values refer to file names and so the regular OS restrictions on names apply.  Also note that the provided data is fully serialized and deserialized (using JSON encoding).

```go
  fsCache := NewFileCache(directory)
```

Note that if the directory is specified as an empty string, a temp directory would be chosen.  The files are not naturally cleaned up on their own (but similar to the S3 cache, they get cleaned up if a stale file is accessed).


# Developing locally

Only dependency is AWS (for S3)

`go get -u github.com/aws/aws-sdk-go`
