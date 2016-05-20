FROM alpine:3.3
ARG NODE_ENV=production
RUN echo $NODE_ENV > /etc/node_env
CMD ["sh", "-c", "cat /etc/node_env"]
