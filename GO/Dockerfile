FROM golang:1.8

# Set Working directory
WORKDIR /AWBeta

# Set GOPath
ENV GOPATH /AWBeta

# Expose port
EXPOSE 3000

# Copy all the content to /AWBeta
COPY . /AWBeta

# Compile project
RUN go build -o exec


ENTRYPOINT ./exec --port 3000 --host="0.0.0.0"