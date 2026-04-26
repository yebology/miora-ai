FROM python:3.10-slim

WORKDIR /app

RUN apt-get update && apt-get install -y --no-install-recommends \
    build-essential \
    libffi-dev \
    libsodium-dev \
    && rm -rf /var/lib/apt/lists/*

RUN pip install --upgrade pip setuptools wheel

COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

# Fix: coinbase-agentkit pins bcl==2.3.1 which fails to build _sodium on aarch64.
# bcl 3.0.0 ships a pre-built wheel that works correctly.
RUN pip install --force-reinstall --no-deps 'bcl>=3.0.0'

COPY . .

EXPOSE 8090

CMD ["python", "main.py"]
