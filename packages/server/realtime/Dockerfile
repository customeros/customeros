ARG ELIXIR_VERSION=1.16.1
ARG OTP_VERSION=26.2.1
ARG DEBIAN_VERSION=bullseye-20240130-slim

ARG BUILDER_IMAGE="hexpm/elixir:${ELIXIR_VERSION}-erlang-${OTP_VERSION}-debian-${DEBIAN_VERSION}"
ARG RUNNER_IMAGE="debian:${DEBIAN_VERSION}"

FROM ${BUILDER_IMAGE} AS builder

# install build dependencies
RUN apt-get update -y && \
  apt-get install -y build-essential git unzip curl && \
  apt-get clean && \
  rm -f /var/lib/apt/lists/*_*

ENV BUN_INSTALL="/root/.bun"
# install bun
RUN curl -fsSL https://bun.sh/install | bash
# Add bun to PATH
ENV PATH="${BUN_INSTALL}/bin:$PATH"

# prepare build dir
WORKDIR /app
COPY priv priv
COPY scripts scripts
RUN mkdir -p priv/scripts

WORKDIR /app/scripts
RUN bun install
RUN bun build --compile convert_lexical_to_yjs.ts --output convert_lexical_to_yjs
RUN ls -la
RUN ls -la /app/priv/
RUN ls -la /app/priv/scripts
RUN mv convert_lexical_to_yjs /app/priv/scripts/

WORKDIR /app

# install hex + rebar
RUN mix local.hex --force && \
  mix local.rebar --force

# set build ENV
ENV MIX_ENV="prod"

# install mix dependencies
COPY mix.exs mix.lock ./
RUN mix deps.get --only $MIX_ENV
RUN mkdir config

# copy compile-time config files before we compile dependencies
# to ensure any relevant config change will trigger the dependencies
# to be re-compiled.
COPY config/config.exs config/${MIX_ENV}.exs config/
RUN mix deps.compile

# COPY priv priv
COPY lib lib

# COPY assets assets

# compile assets
# RUN mix assets.deploy # we have no static assets

# Compile the release
RUN mix compile

# Changes to config/runtime.exs don't require recompiling the code
COPY config/runtime.exs config/

COPY rel rel
RUN mix release

# start a components build stage so that the final image will only contain
# the compiled release and other runtime necessities
FROM ${RUNNER_IMAGE}

RUN apt-get update -y && \
  apt-get install -y libstdc++6 openssl libncurses5 locales ca-certificates \
  && apt-get clean && rm -f /var/lib/apt/lists/*_*

# Set the locale
RUN sed -i '/en_US.UTF-8/s/^# //g' /etc/locale.gen && locale-gen

ENV LANG=en_US.UTF-8
ENV LANGUAGE=en_US:en
ENV LC_ALL=en_US.UTF-8

WORKDIR "/app"
RUN chown nobody /app

# set runner ENV
ENV MIX_ENV="prod"

# Only copy the final release from the build stage
COPY --from=builder --chown=nobody:root /app/_build/${MIX_ENV}/rel/realtime ./
# Copy the entrypoint script
COPY entrypoint.sh /app/entrypoint.sh
RUN chmod +x /app/entrypoint.sh

USER nobody

ENTRYPOINT ["/app/entrypoint.sh"]
