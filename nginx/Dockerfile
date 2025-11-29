FROM nginx:alpine

COPY nginx/nginx.conf /etc/nginx/templates/default.conf.template

CMD ["/bin/sh", "-c", "envsubst '$$AUTH_REST_PORT $$PARKING_REST_PORT $$BOOKING_REST_PORT $$PAYMENT_REST_PORT $$KEYCLOAK_INNER_PORT $$NGINX_PORT $$GRAFANA_PORT $$PROMETHEUS_PORT' < /etc/nginx/templates/default.conf.template > /etc/nginx/conf.d/default.conf && nginx -g 'daemon off;'"]

