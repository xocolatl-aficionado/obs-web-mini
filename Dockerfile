# Use a base image with Node.js for the frontend
FROM node:18-alpine AS build

# Set the working directory
WORKDIR /app

# Install Git and Node.js dependencies
RUN apk add --no-cache git
COPY ./ ./
RUN npm ci && npm run build

# Use a smaller base image to run the app
FROM node:18-alpine

# Set the working directory
WORKDIR /app

# Copy the built files from the build stage
COPY --from=build /app /app

# Install production dependencies (if needed)
RUN npm ci 

# Expose ports
EXPOSE 5000 3001

# Command to start the frontend
CMD ["npm", "run", "start:all"]

# Include a label for the repository source
LABEL org.opencontainers.image.source="https://github.com/xocolatl-aficionado/obs-web-mini"
