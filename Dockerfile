FROM gcr.io/distroless/static:nonroot

WORKDIR /app

COPY register_deployed_artifact_test_app /app

ENTRYPOINT ["/app/register_deployed_artifact_test_app"]
