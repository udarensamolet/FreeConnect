# Use a newer Node image (Node 16 or Node 18 recommended)
FROM node:20

# Set the working directory inside the container
WORKDIR /app

# Copy the package files to leverage Docker cache for dependency installation
COPY package*.json ./

# (Optional) Install Angular CLI globally if needed
RUN npm install -g @angular/cli

# Install the dependencies
RUN npm install

# Copy the rest of the application code
COPY . .

# Expose the port on which Angular’s development server runs
EXPOSE 4200

# Start the Angular development server, ensuring it listens on all network interfaces.
CMD ["ng", "serve", "--host", "0.0.0.0"]