FROM golang:1.17

# Env variables definition
ENV USER developers
ENV HOME /go/developers-italia-backend/crawler

ENV DEFAULT_TIMEOUT 300

ENV BASE /var/crawler
ENV DATA ${BASE}/data
ENV OUTPUT ${BASE}/output

# Set the work directory
WORKDIR ${HOME}

# Copy crawler files inside the workdir
COPY .git .git
COPY cmd cmd
COPY crawler crawler
COPY elastic elastic
COPY ipa ipa
COPY jekyll jekyll
COPY metrics metrics
COPY version version
COPY publishers.thirparty.yml publishers.thirparty.yml
COPY publishers.yml publishers.yml
COPY blacklist blacklist
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
RUN mkdir -p ${OUTPUT}
RUN chown -R ${USER}.${USER} ${BASE}

# Set running user
USER ${USER}

# Compile a new crawler
RUN make

# Remove unsed .git
RUN rm -rf .git

# By default, wait until Elasticsearch is not ready. Then start the crawler
CMD ["bash", "-c", "./wait-for-it.sh ${ELASTIC_URL} -t ${DEFAULT_TIMEOUT} -- ./start.sh"]
