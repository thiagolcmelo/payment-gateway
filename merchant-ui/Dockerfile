# Base image
FROM node:14-alpine

# Set the working directory inside the container
WORKDIR /app

# Copy package.json and package-lock.json to the working directory
COPY package*.json ./

# Install dependencies
RUN npm install --production

# Copy the entire app to the working directory
COPY . .

# Set the environment variables
ENV REACT_APP_PAYMENT_GATEWAY_HOST=0.0.0.0
ENV REACT_APP_PAYMENT_GATEWAY_PORT=8080

# Build the React app
RUN npm run build

# Expose the port that the app will run on
EXPOSE 3000

# Start the app
CMD ["npm", "start"]

# # Stage 1: Build the React app
# FROM node:14.17-alpine as build

# WORKDIR /app

# # Copy package.json and package-lock.json
# COPY package.json package-lock.json ./

# # Install dependencies
# RUN npm install --silent

# # Copy the app source code
# COPY . .

# # Build the React app
# RUN npm run build

# # Stage 2: Serve the built React app
# FROM nginx:1.21-alpine

# # Copy the built app from the previous stage
# COPY --from=build /app/build /usr/share/nginx/html

# # Replace default Nginx configuration with custom configuration
# COPY nginx.conf /etc/nginx/conf.d/default.conf

# # Expose port 3000
# EXPOSE 3000

# # Set the environment variables
# ENV REACT_APP_PAYMENT_GATEWAY_HOST=0.0.0.0
# ENV REACT_APP_PAYMENT_GATEWAY_PORT=8080

# # Start Nginx server
# CMD ["nginx", "-g", "daemon off;"]