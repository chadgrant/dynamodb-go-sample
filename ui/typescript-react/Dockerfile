FROM node:12.13.1-alpine as builder

WORKDIR /app

ENV PATH /app/node_modules/.bin:$PATH

COPY package.json .
RUN npm install
RUN npm install react-scripts@3.3.0 -g

COPY public public
COPY src src
RUN npm run build

FROM nginx:1.16.0-alpine
COPY --from=builder /app/build /usr/share/nginx/html
COPY nginx.conf /etc/nginx/conf.d
COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh
EXPOSE 80
ENTRYPOINT ["/entrypoint.sh"]
CMD ["nginx", "-g", "daemon off;"]