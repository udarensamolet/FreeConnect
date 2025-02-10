# Stage 1: Build the Angular application
FROM node:20 as build

# Set the working directory
WORKDIR /app

# Copy package.json and package-lock.json to leverage Docker cache
COPY package*.json ./

# Install Angular CLI globally (optional if you want to use the global CLI)
RUN npm install -g @angular/cli

# Install project dependencies (including devDependencies)
RUN npm install

# Copy the rest of the application source code
COPY . .

# Build the Angular application for production.
# The double-dash (-- --prod) passes the --prod flag to the underlying ng build command.
RUN npm run build -- --prod

# Stage 2: Serve the application with Nginx
FROM nginx:alpine

# Remove default Nginx static assets
RUN rm -rf /usr/share/nginx/html/*

# Copy the production build from the previous stage to Nginx's public directory.
# Update the path if your Angular output directory is different.
COPY --from=build /app/dist/freeconnect-frontend /usr/share/nginx/html

# Expose port 80 for Nginx
EXPOSE 80

# Start Nginx in the foreground
CMD ["nginx", "-g", "daemon off;"]
