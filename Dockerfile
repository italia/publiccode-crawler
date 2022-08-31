FROM golang:1.17

# Env variables definition
ENV USER developers
ENV HOME /go/developers-italia-backend/crawler

ENV BASE /var/crawler
ENV DATA ${BASE}/data

# Set the work directory
WORKDIR ${HOME}

# Copy crawler files inside the workdir
COPY .git .git
COPY common common
COPY cmd cmd
COPY crawler crawler
COPY git git
COPY internal internal
COPY scanner scanner
COPY metrics metrics
COPY version version
COPY publishers.thirdparty.yml publishers.thirdparty.yml
COPY publishers.yml publishers.yml
COPY config.toml.example config.toml
COPY domains.yml.example domains.yml
COPY go.mod .
COPY go.sum .
COPY main.go .
COPY Makefile .
COPY start.sh .
COPY vitality-ranges.yml .
COPY wait-for-it.sh .

# Run as unprivileged user
RUN adduser --home ${HOME} --shell /bin/sh --disabled-password ${USER}

# Set user ownership on workdir and subdirectories
RUN chown -R ${USER}.${USER} ${HOME}

# Create the crawler output directory structure and set user ownership
# Must match what's written in config.toml.example
RUN mkdir -p ${DATA}
RUN chown -R ${USER}.${USER} ${BASE}

# Set running user
USER ${USER}

# Compile a new crawler
RUN make

# Remove unsed .git
RUN rm -rf .git

CMD ["start.sh"]
