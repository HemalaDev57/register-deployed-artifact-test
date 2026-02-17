FROM gcr.io/distroless/static:nonroot

WORKDIR /app

COPY gha_register_deployed_artifact_app /app

ENTRYPOINT ["/app/gha_register_deployed_artifact_app"]
