FROM node:16-alpine as build-stage

WORKDIR /app
COPY frontend/. .
RUN cp .env.prod .env
RUN yarn
RUN yarn build

FROM nginx:stable-alpine as production-stage
COPY --from=build-stage /app/dist /usr/share/nginx/html
COPY .docker/frontend/default.conf /etc/nginx/conf.d/default.conf
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
