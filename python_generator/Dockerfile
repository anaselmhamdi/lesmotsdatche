FROM python:3.12-slim

WORKDIR /app

# Install uv and curl (keep curl for healthcheck)
ENV UV_INSTALL_DIR="/usr/local/bin"
RUN apt-get update && apt-get install -y --no-install-recommends curl && \
    curl -LsSf https://astral.sh/uv/install.sh | sh && \
    rm -rf /var/lib/apt/lists/*

# Copy dependency files
COPY pyproject.toml uv.lock ./

# Copy source code (needed for editable install)
COPY domain/ domain/
COPY generator/ generator/
COPY language/ language/
COPY llm/ llm/
COPY api/ api/
COPY data/ data/
COPY seed/ seed/
COPY cli.py .

# Install dependencies
RUN uv sync --frozen --no-dev

# Create non-root user
RUN useradd -m -u 1000 appuser && chown -R appuser:appuser /app
USER appuser

# Default command runs the API
EXPOSE 8080
CMD ["uv", "run", "uvicorn", "api.routes:app", "--host", "0.0.0.0", "--port", "8080"]
