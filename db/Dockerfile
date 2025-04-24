FROM postgres:latest

EXPOSE 5432

COPY schema.sql /docker-entrypoint-initdb.d/

CMD ["postgres"]