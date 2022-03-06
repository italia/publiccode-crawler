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
COPY crawler/cmd cmd
COPY crawler/crawler crawler
COPY crawler/elastic elastic
COPY crawler/ipa ipa
COPY crawler/jekyll jekyll
COPY crawler/metrics metrics
COPY crawler/version version
COPY crawler/whitelist whitelist
COPY crawler/blacklist blacklist
COPY crawler/config.toml.example config.toml
COPY crawler/domains.yml.example domains.yml
COPY crawler/go.mod .
COPY crawler/go.sum .
COPY crawler/main.go .
COPY crawler/Makefile .
COPY crawler/start.sh .
COPY crawler/vitality-ranges.yml .
COPY crawler/wait-for-it.sh .

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
