FROM nginx:1.25-alpine
COPY ["config/", "/etc/nginx"]
COPY ["tiger/public/", "/tiger/new"]

HEALTHCHECK --interval=30s --timeout=60s --retries=5 --start-period=20s \
  CMD wget -SqO - http://127.0.0.1:${tiger_REACT_PORT:-3000}/new || exit 1