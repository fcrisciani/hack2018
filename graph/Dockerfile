FROM node:9-alpine

COPY . /usr/app
WORKDIR /usr/app
RUN yarn && yarn build

FROM nginx:alpine

COPY --from=0 /usr/app/build /usr/share/nginx/html

ENTRYPOINT ["nginx", "-g", "daemon off;"]
