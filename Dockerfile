# Step 1: Use a Golang base image
FROM golang:1.21-alpine AS build-stage 

# Step 2: Set the working directory
WORKDIR /app

# Step 3: Copy Go module files (if applicable)
COPY go.mod go.sum ./

# Step 4: Download dependencies 
RUN go mod download

# Step 5: Copy the rest of your source code 
COPY . .

# Step 6: Build the application
RUN go build -o main

# Step 7: Smaller, final image for production
FROM alpine:latest

# Step 8: Set working directory in the runtime image
WORKDIR /yeet

# Step 9: Copy only the built binary from the build stage
COPY --from=build-stage /app/main /yeet

# Step 10: Define the command to start your application
CMD ["./yeet"] 
