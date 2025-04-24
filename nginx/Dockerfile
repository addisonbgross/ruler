FROM nginx:alpine

WORKDIR /app

COPY nginx.conf /etc/nginx/nginx.conf
COPY client.pem /etc/nginx/client.pem
COPY client.key /etc/nginx/client.key

EXPOSE 80 443

CMD ["nginx", "-g", "daemon off;"]
