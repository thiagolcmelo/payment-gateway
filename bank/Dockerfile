# Use a Python base image
FROM python:3.9-slim

# Set the working directory inside the container
WORKDIR /app

# Copy the requirements file
COPY requirements.txt .

# Install the dependencies
RUN pip install --no-cache-dir -r requirements.txt

# Copy the app source code to the container
COPY . .

# Set the environment variables
ENV PAYMENT_GATEWAY_HOST=0.0.0.0
ENV PAYMENT_GATEWAY_PORT=8080

# Expose the port on which the FastAPI server listens
EXPOSE 8000

# Start the FastAPI server when the container is run
CMD ["uvicorn", "main:app", "--host", "0.0.0.0", "--port", "8000"]