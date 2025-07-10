FROM gcr.io/distroless/static-debian11:nonroot
# FROM public.ecr.aws/debian/debian:11
ENTRYPOINT ["/baton-openshift"]
COPY baton-openshift /