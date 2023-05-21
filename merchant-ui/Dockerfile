# Stage 1: Build the React app
FROM node:14.17-alpine as build

WORKDIR /app

# Copy package.json and package-lock.json
COPY package.json package-lock.json ./

# Install dependencies
RUN npm install --silent

# Copy the app source code
COPY . .

# # !!!WORKAROUND!!!
RUN sed -i "s/process\.env\.REACT_APP_PAYMENT_GATEWAY_HOST/\"REACT_APP_PAYMENT_GATEWAY_HOST_WORKAROUND\"/g" src/Login.js
RUN sed -i "s/process\.env\.REACT_APP_PAYMENT_GATEWAY_PORT/\"REACT_APP_PAYMENT_GATEWAY_PORT_WORKAROUND\"/g" src/Login.js
RUN sed -i "s/process\.env\.REACT_APP_PAYMENT_GATEWAY_HOST/\"REACT_APP_PAYMENT_GATEWAY_HOST_WORKAROUND\"/g" src/Dashboard.js
RUN sed -i "s/process\.env\.REACT_APP_PAYMENT_GATEWAY_PORT/\"REACT_APP_PAYMENT_GATEWAY_PORT_WORKAROUND\"/g" src/Dashboard.js

# Build the React app
RUN npm run build

# Stage 2: Serve the built React app
FROM nginx:1.21-alpine

# Copy the built app from the previous stage
COPY --from=build /app/build /usr/share/nginx/html

# Replace default Nginx configuration with custom configuration
COPY nginx.conf /etc/nginx/conf.d/default.conf

# Expose port 3000
EXPOSE 3000

# Set the environment variables
ENV REACT_APP_PAYMENT_GATEWAY_HOST=0.0.0.0
ENV REACT_APP_PAYMENT_GATEWAY_PORT=8080

# !!!WORKAROUND!!!
RUN find /usr/share/nginx/html/static/js/ -type f -name "main.*.js" -exec sed -i "s/\"REACT_APP_PAYMENT_GATEWAY_HOST_WORKAROUND\"/\"${REACT_APP_PAYMENT_GATEWAY_HOST}\"/g" {} +
RUN find /usr/share/nginx/html/static/js/ -type f -name "main.*.js" -exec sed -i "s/\"REACT_APP_PAYMENT_GATEWAY_PORT_WORKAROUND\"/\"${REACT_APP_PAYMENT_GATEWAY_PORT}\"/g" {} +

# Start Nginx server
CMD ["nginx", "-g", "daemon off;"]