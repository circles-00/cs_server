FROM ubuntu:focal

ARG ADMIN_NICKNAME
ARG ADMIN_PASSWORD

RUN \
  apt update -qq && \
  apt install -qq -y lib32gcc1 && \
  apt install -qq -y screen

RUN useradd cstrike

RUN mkdir -p /home/cstrike

USER cstrike

WORKDIR /home/cstrike

COPY --chown=cstrike:cstrike ./server .

RUN echo "$ADMIN_NICKNAME $ADMIN_PASSWORD \"abcdefghijklmnopqrstu\" \"ab\"" >> /home/cstrike/cstrike/addons/amxmodx/configs/users.ini

ENTRYPOINT ./hlds_run -game cstrike -strictportbind -ip 0.0.0.0 -port $PORT +map $START_MAP -maxplayers $MAX_PLAYERS -pingboost 3 -sys_ticrate 1000
