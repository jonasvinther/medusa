# Start with the Rocky Linux base image
FROM rockylinux:8

# Set environment variables for Go
ENV GO_VERSION=1.21.0
ENV GOPATH=/go
ENV PATH=$PATH:/usr/local/go/bin:$GOPATH/bin

# Set environment variables for Python
ENV PYTHON_VERSION=3.9

# Update the package list and install dependencies
RUN dnf -y update && \
    dnf -y install epel-release && \
    dnf -y install wget git gcc make && \
    dnf clean all

RUN dnf -y groupinstall "Development Tools" && \
    dnf -y install openssl-devel libffi-devel bzip2-devel sqlite-devel ncurses-devel readline-devel wget curl

RUN wget https://www.python.org/ftp/python/3.9.10/Python-3.9.10.tgz && \
    tar xvf Python-3.9.10.tgz && \
    cd Python-3.9.10 && \
    ./configure --enable-optimizations --enable-loadable-sqlite-extensions && \
    make altinstall -j && \
    cd .. && \
    rm -rf Python-3.9.10*

# Install Go
RUN wget https://go.dev/dl/go${GO_VERSION}.linux-amd64.tar.gz && \
    tar -C /usr/local -xzf go${GO_VERSION}.linux-amd64.tar.gz && \
    rm go${GO_VERSION}.linux-amd64.tar.gz

# Set the working directory for the application
WORKDIR /Medusa

# Copy the local application files into the container
COPY . /Medusa/

# Install Poetry and set up virtual environment
RUN curl -sSL https://install.python-poetry.org | python3.9 - && \
    echo 'export PATH="/home/root/.dev/bin:$PATH"'

# Install Python dependencies using Poetry
RUN python3.9 -m venv /Medusa/venv && \
    /Medusa/venv/bin/pip install --upgrade pip && \
    /Medusa/venv/bin/pip install poetry && \
    /Medusa/venv/bin/poetry install


# Install FastAPI, Uvicorn, and python-multipart
RUN /Medusa/venv/bin/pip install fastapi uvicorn python-multipart pyyaml

# Download and tidy Go modules
RUN go mod tidy

# Build the Medusa binary
RUN go build -o medusa .

# Expose the necessary ports for FastAPI and Medusa
EXPOSE 8080 8000

# Set entrypoint to run Uvicorn with FastAPI app
ENTRYPOINT ["/Medusa/venv/bin/uvicorn", "api.main:app", "--host", "0.0.0.0", "--port", "8000", "--reload"]
