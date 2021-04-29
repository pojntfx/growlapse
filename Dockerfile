# Build container
FROM golang:1.16.3 AS build

# Setup environment
RUN mkdir -p /data
WORKDIR /data

# Build the release
COPY . .
RUN ./Hydrunfile

# Extract the release
RUN mkdir -p /out
RUN cp out/release/growlapse-agent/growlapse-agent.linux-$(uname -m) /out/growlapse-agent

# Release container
FROM debian

# Add the release
COPY --from=build /out/growlapse-agent /usr/local/bin/growlapse-agent

CMD /usr/local/bin/growlapse-agent
