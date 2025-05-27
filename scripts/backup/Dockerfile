FROM alpine:latest

# Install required packages
RUN apk add --no-cache \
    postgresql-client \
    curl \
    bash \
    coreutils \
    tar \
    gzip

# Create app directory
WORKDIR /app

# Copy backup scripts
COPY backup.sh /app/backup.sh
COPY restore.sh /app/restore.sh

# Make scripts executable
RUN chmod +x /app/backup.sh /app/restore.sh

# Create backup directory
RUN mkdir -p /app/backups

# Set default command
CMD ["/app/backup.sh"] 