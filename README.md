# cloud-run-auth

This package implements [credentials.PerRPCCredentials](https://pkg.go.dev/google.golang.org/grpc@v1.48.0/credentials#PerRPCCredentials) interface. This implementation allows you to invoke Cloud Run gRPC from another Cloud Run gRPC service, once the client service account has the role `roles/run.invoker`.

## Example

```go
    addr := os.Getenv("EMAIL_SERVICE_ADDR")
	systemRoots, err := x509.SystemCertPool()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	creds := credentials.NewTLS(&tls.Config{
		RootCAs: systemRoots,
	})
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(creds),
		grpc.WithPerRPCCredentials(cloudrunauth.MetadataServerToken{ServiceURL: os.Getenv("AUDIENCE")}),
		grpc.WithUnaryInterceptor(grpc_middleware.ChainUnaryClient(otelgrpc.UnaryClientInterceptor())),
		grpc.WithStreamInterceptor(grpc_middleware.ChainStreamClient(otelgrpc.StreamClientInterceptor())),
		grpc.WithAuthority(addr),
	}
	conn, err := grpc.DialContext(ctx, addr, opts...)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	c := notificationv1.NewEmailClient(conn)
	_, err = c.Send(ctx, &nr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
```