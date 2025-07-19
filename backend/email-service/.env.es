# SMTP Configuration
ES_SMTP_HOST=smtp.gmail.com
ES_SMTP_PORT=587
ES_SMTP_USER=your_email@gmail.com
SMTP_PASSWORD=your_email_password
ES_FROM_ADDR=your_email@gmail.com
ES_SMTP_MAX_CONNECTIONS=10
ES_SMTP_CONNECTION_TIMEOUT=30       # in seconds
ES_SMTP_SEND_TIMEOUT=60             # in seconds
ES_SMTP_KEEP_ALIVE=true

# NATS Configuration
ES_NATS_URL=nats://localhost:4222
ES_NATS_STREAM=EMAILS
ES_NATS_SUBJECTS=email.verification,email.notification,page.approval
ES_NATS_MAX_DELIVER=3
ES_NATS_ACK_WAIT=300                # in seconds
ES_NATS_MAX_ACK_PENDING=100

# Template Configuration
ES_TEMPLATES_DIR=templates
ES_TEMPLATES_EXT=.html
ES_TEMPLATES_CACHE=true
ES_TEMPLATES_WATCH=false
ES_TEMPLATES_RELOAD=false

# Email Processing Workers
ES_PROCESSING_WORKERS=4             # Default to runtime.NumCPU() if unset
ES_PROCESSING_BATCH_SIZE=10
ES_PROCESSING_RETRY_DELAY=60        # in seconds
ES_PROCESSING_MAX_RETRIES=3

# Monitoring & Logging
ES_MONITORING_METRICS=true
ES_MONITORING_LOG_LEVEL=INFO
