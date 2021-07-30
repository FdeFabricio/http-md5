# http-md5
This project makes HTTP requests and prints the address of the request along with the MD5 hash of the response.
In order to perform multiple requests concurrently, it implements a pool of workers that perform the task for each URL independently and in parallel.
The maximum number of concurrent requests is defined by a flag `-parallel` and its default value is `10`.
As soon as each worker perform the HTTP request, reads the response body and compute its checksum, it prints the URL followed by its MD5.
Therefore, the order in which the URLs are provided as arguments might not be followed when outputting the result. 

## How to run
This project was tested on a machine with go version `go1.16.2 darwin/amd64`.
To run the project, follow the instructions below:

1. Build the project with the following command:
```shell
go build .
```

2. Execute the binary passing URLs as arguments:
```shell
./http-md5 google.com http://bbc.com facebook.com yahoo.com
```

You can alternatively inform the limit of concurrent workers with the flag `parallel` as such:
```shell
./http-md5 -parallel 3 google.com http://bbc.com facebook.com yahoo.com
```


## Considerations

* **URL validation:** the solution does not perform any robust URL validation,
  besides including the prefix `"http://"` in case the URL does not have it already.
  One could use a regular expression to validate them, but since the description of the task did not mention
  what URL formats are valid I decided to not include.
  The solution can make HTTP requests on invalid URLs which will eventually fail,
  just like requesting any valid URL that does not exist, for example, `"http://silverpotatoes.com"`.
  
* **Error handling:** since the description of the task did not mention how to handle errors,
  this solution also prints the errors, but with a `[ERROR]` prefix.
  This way the user can know which URLs failed to have their checksum computed.
  Requests in which their responses have a status code different from HTTP 200 OK are considered as an error case.
  