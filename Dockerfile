FROM scratch
COPY goaccess-pixel /
ENV PORT 80
ENTRYPOINT ["/goaccess-pixel"]
