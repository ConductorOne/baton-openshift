FROM gcr.io/distroless/static-debian11:nonroot
ENTRYPOINT ["/baton-openshift"]
COPY baton-openshift /