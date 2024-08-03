# Use the official Rust image as a base image
FROM rust:1.70-slim-buster as builder

# Set the working directory inside the container
WORKDIR /app

# Copy the Cargo.toml and Cargo.lock files
COPY Cargo.toml Cargo.lock ./

# Copy the source code
COPY src ./src

# Copy the .env file
COPY .env .env

# Build the application
RUN cargo build --release

# Create a smaller final image
FROM debian:buster-slim

# Set the working directory inside the container
WORKDIR /app

# Copy the .env file
COPY .env .env

# Copy the compiled binary from the builder stage
COPY --from=builder /app/target/release/your_package_name .

# Set the binary as the entry point
ENTRYPOINT ["./your_package_name"]

# Specify environment variables
ENV RUST_LOG=info
