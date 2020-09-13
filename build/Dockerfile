# Dockerfile for the library server application

FROM golang:1.14.6

# Install API application and configuration
WORKDIR /opt/library
COPY library . 
COPY conf/library_prod.toml conf/library.toml
COPY internal/app/views internal/app/views
COPY docs docs 

# Default HTTP listen port for the library 
ARG http_port=8080

# Expose the library HTTP port
EXPOSE ${http_port}

CMD ["/opt/library/library"] 
